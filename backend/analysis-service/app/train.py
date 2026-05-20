from __future__ import annotations

import json
import re
from pathlib import Path
from sklearn.metrics import (
    accuracy_score,
    classification_report,
    confusion_matrix,
    precision_score,
    recall_score,
    f1_score,
)
import joblib
import pandas as pd
from sklearn.feature_extraction.text import ENGLISH_STOP_WORDS, TfidfVectorizer
from sklearn.model_selection import train_test_split
from sklearn.pipeline import Pipeline
from sklearn.preprocessing import LabelEncoder
from sklearn.svm import LinearSVC


BASE_DIR = Path(__file__).resolve().parent.parent
DATA_PATH = BASE_DIR / "data" / "tickets.csv"
MODELS_DIR = BASE_DIR / "models"

TEXT_COLUMN = "instruction"
TARGET_COLUMN = "intent"
CATEGORY_COLUMN = "category"


def clean_text(text: str) -> str:
    text = str(text).lower()
    text = re.sub(r"http\S+|www\.\S+", " ", text)
    text = re.sub(r"\S+@\S+", " ", text)
    text = re.sub(r"[^a-z0-9\s_]", " ", text, flags=re.IGNORECASE)
    text = re.sub(r"\s+", " ", text).strip()
    return text


def load_dataset(csv_path: Path) -> pd.DataFrame:
    if not csv_path.exists():
        raise FileNotFoundError(f"Dataset not found: {csv_path}")

    df = pd.read_csv(csv_path)

    print("Dataset loaded")
    print("Shape:", df.shape)
    print("Columns:", list(df.columns))

    required_columns = {TEXT_COLUMN, TARGET_COLUMN, CATEGORY_COLUMN}
    missing = required_columns - set(df.columns)
    if missing:
        raise ValueError(f"Missing required columns: {missing}")

    work_df = df[[TEXT_COLUMN, TARGET_COLUMN, CATEGORY_COLUMN]].copy()
    work_df = work_df.dropna()

    work_df["message"] = work_df[TEXT_COLUMN].astype(str).str.strip()
    work_df["intent"] = work_df[TARGET_COLUMN].astype(str).str.strip()
    work_df["category"] = work_df[CATEGORY_COLUMN].astype(str).str.strip()

    work_df = work_df[
        (work_df["message"] != "") &
        (work_df["intent"] != "") &
        (work_df["category"] != "")
    ]

    work_df["message"] = work_df["message"].apply(clean_text)
    work_df = work_df[work_df["message"] != ""]

    print("After cleaning:", work_df.shape)
    print("\nSample messages:")
    for text in work_df["message"].head(10).tolist():
        print("-", text)

    return work_df[["message", "intent", "category"]]


def build_pipeline() -> Pipeline:
    return Pipeline(
        steps=[
            (
                "tfidf",
                TfidfVectorizer(
                    lowercase=False,
                    stop_words=list(ENGLISH_STOP_WORDS),
                    ngram_range=(1, 2),
                    max_features=50000,
                    min_df=2,
                    max_df=0.95,
                    sublinear_tf=True,
                ),
            ),
            (
                "clf",
                LinearSVC(C=1.0, random_state=42),
            ),
        ]
    )


def build_intent_category_map(df: pd.DataFrame) -> dict[str, str]:
    mapping = (
        df[["intent", "category"]]
        .drop_duplicates()
        .set_index("intent")["category"]
        .to_dict()
    )
    return mapping


def main() -> None:
    print("Training started...")
    print("BASE_DIR =", BASE_DIR)
    print("DATA_PATH =", DATA_PATH)
    print("MODELS_DIR =", MODELS_DIR)

    MODELS_DIR.mkdir(parents=True, exist_ok=True)

    df = load_dataset(DATA_PATH)

    print("\nIntent distribution:")
    print(df["intent"].value_counts())

    print("\nCategory distribution:")
    print(df["category"].value_counts())

    X = df["message"]
    y = df["intent"]

    label_encoder = LabelEncoder()
    y_encoded = label_encoder.fit_transform(y)

    print("\nIntents:")
    for idx, class_name in enumerate(label_encoder.classes_):
        print(f"{idx}: {class_name}")

    X_train, X_test, y_train, y_test = train_test_split(
        X,
        y_encoded,
        test_size=0.2,
        random_state=42,
        stratify=y_encoded,
    )

    print("\nTrain size:", len(X_train))
    print("Test size:", len(X_test))

    pipeline = build_pipeline()

    print("\nFitting model...")
    pipeline.fit(X_train, y_train)
    print("Training completed")

    y_pred = pipeline.predict(X_test)
    accuracy = accuracy_score(y_test, y_pred)
    precision = precision_score(y_test, y_pred, average="macro", zero_division=0)
    recall = recall_score(y_test, y_pred, average="macro", zero_division=0)
    f1 = f1_score(y_test, y_pred, average="macro", zero_division=0)

    print("\nAccuracy:", round(accuracy, 4))
    print("Precision:", round(precision, 4))
    print("Recall:", round(recall, 4))
    print("F1:", round(f1, 4))

    print("\nAccuracy:", round(accuracy, 4))
    print("\nClassification report:")
    print(
        classification_report(
            y_test,
            y_pred,
            target_names=label_encoder.classes_,
            digits=4,
            zero_division=0,
        )
    )

    print("Confusion matrix:")
    print(confusion_matrix(y_test, y_pred))

    classifier_path = MODELS_DIR / "classifier.joblib"
    encoder_path = MODELS_DIR / "label_encoder.joblib"
    mapping_path = MODELS_DIR / "intent_to_category.json"

    joblib.dump(pipeline, classifier_path)
    joblib.dump(label_encoder, encoder_path)

    intent_to_category = build_intent_category_map(df)
    mapping_path.write_text(
        json.dumps(intent_to_category, ensure_ascii=False, indent=2),
        encoding="utf-8",
    )

    print("\nModel saved:")
    print(classifier_path)
    print(encoder_path)
    print(mapping_path)


if __name__ == "__main__":
    main()