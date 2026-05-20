from __future__ import annotations

import json
import re
from pathlib import Path

import joblib
from deep_translator import GoogleTranslator


BASE_DIR = Path(__file__).resolve().parent.parent
MODELS_DIR = BASE_DIR / "models"

INTENT_TO_SHORT_CATEGORY = {
    "create_account": "ACCOUNT",
    "delete_account": "ACCOUNT",
    "edit_account": "ACCOUNT",
    "recover_password": "ACCOUNT",
    "registration_problems": "ACCOUNT",
    "switch_account": "ACCOUNT",

    "place_order": "ORDER",
    "change_order": "ORDER",
    "track_order": "ORDER",
    "cancel_order": "ORDER",

    "get_refund": "REFUND",
    "track_refund": "REFUND",
    "check_refund_policy": "REFUND",

    "payment_issue": "PAYMENT",
    "check_payment_methods": "PAYMENT",
    "check_invoice": "PAYMENT",
    "get_invoice": "PAYMENT",
    "check_cancellation_fee": "PAYMENT",

    "delivery_options": "DELIVERY",
    "delivery_period": "DELIVERY",
    "change_shipping_address": "DELIVERY",
    "set_up_shipping_address": "DELIVERY",

    "contact_customer_service": "SUPPORT",
    "contact_human_agent": "SUPPORT",
    "complaint": "SUPPORT",
    "review": "SUPPORT",
    "newsletter_subscription": "SUPPORT",
}


RULES = [
    {
        "category": "PAYMENT",
        "intent": "payment_issue",
        "raw_category": "PAYMENT",
        "patterns": [
            r"\bсписал[аи]?\b",
            r"\bсписали деньги\b",
            r"\bдв(а|е) раз[аы]\b",
            r"\bденьги\b",
            r"\bоплат[аы]\b",
            r"\bоплатить\b",
            r"\bплатеж\b",
            r"\bплатёж\b",
            r"\bкарт[аы]\b",
            r"\bсчет\b",
            r"\bсч[её]т\b",
            r"\binvoice\b",
            r"\bbill(ing)?\b",
            r"\bcharged\b",
            r"\bcharge\b",
            r"\bpayment\b",
            r"\bcard\b",
            r"\bpay\b",
        ],
    },
    {
        "category": "REFUND",
        "intent": "get_refund",
        "raw_category": "REFUND",
        "patterns": [
            r"\bвозврат\b",
            r"\bвернуть деньги\b",
            r"\bверните деньги\b",
            r"\brefund\b",
            r"\bmoney back\b",
            r"\breimburse\b",
        ],
    },
    {
        "category": "DELIVERY",
        "intent": "change_shipping_address",
        "raw_category": "SHIPPING",
        "patterns": [
            r"\bадрес доставки\b",
            r"\bдоставк[аи]\b",
            r"\bкурьер\b",
            r"\bshipping\b",
            r"\bdelivery\b",
            r"\bchange address\b",
            r"\bshipping address\b",
        ],
    },
    {
        "category": "ORDER",
        "intent": "cancel_order",
        "raw_category": "ORDER",
        "patterns": [
            r"\bзаказ\b",
            r"\bотменить заказ\b",
            r"\bоформить заказ\b",
            r"\bновый заказ\b",
            r"\border\b",
            r"\bcancel order\b",
            r"\bplace order\b",
            r"\btrack order\b",
        ],
    },
    {
        "category": "ACCOUNT",
        "intent": "recover_password",
        "raw_category": "ACCOUNT",
        "patterns": [
            r"\bпарол[ья]\b",
            r"\bне могу войти\b",
            r"\bвойти в аккаунт\b",
            r"\bаккаунт\b",
            r"\bучетн(ая|ый) запись\b",
            r"\bучётн(ая|ый) запись\b",
            r"\blog ?in\b",
            r"\bpassword\b",
            r"\baccount\b",
            r"\bsign ?in\b",
            r"\bregister\b",
            r"\bregistration\b",
        ],
    },
    {
        "category": "SUPPORT",
        "intent": "contact_human_agent",
        "raw_category": "CONTACT",
        "patterns": [
            r"\bоператор\b",
            r"\bподдержк[аи]\b",
            r"\bсвязаться\b",
            r"\bжалоб[аы]\b",
            r"\bagent\b",
            r"\bsupport\b",
            r"\bhuman\b",
            r"\bcontact\b",
            r"\bcomplaint\b",
        ],
    },
]


class TextClassifier:
    def __init__(self) -> None:
        classifier_path = MODELS_DIR / "classifier.joblib"
        encoder_path = MODELS_DIR / "label_encoder.joblib"
        mapping_path = MODELS_DIR / "intent_to_category.json"

        if not classifier_path.exists():
            raise FileNotFoundError(f"Classifier not found: {classifier_path}")
        if not encoder_path.exists():
            raise FileNotFoundError(f"Label encoder not found: {encoder_path}")
        if not mapping_path.exists():
            raise FileNotFoundError(f"Intent map not found: {mapping_path}")

        self.pipeline = joblib.load(classifier_path)
        self.label_encoder = joblib.load(encoder_path)
        self.intent_to_category = json.loads(mapping_path.read_text(encoding="utf-8"))

    def normalize_text(self, text: str) -> str:
        text = text.lower().strip()
        text = text.replace("ё", "е")
        text = re.sub(r"\s+", " ", text)
        return text

    def translate_to_english(self, text: str) -> str:
        try:
            return GoogleTranslator(source="auto", target="en").translate(text)
        except Exception:
            return text

    def predict_by_rules(self, text: str) -> tuple[str, str, str, float, str] | None:
        normalized = self.normalize_text(text)

        scores: dict[str, int] = {}
        matched_meta: dict[str, tuple[str, str]] = {}

        for rule in RULES:
            count = 0
            for pattern in rule["patterns"]:
                if re.search(pattern, normalized, flags=re.IGNORECASE):
                    count += 1

            if count > 0:
                category = rule["category"]
                scores[category] = scores.get(category, 0) + count
                matched_meta[category] = (rule["intent"], rule["raw_category"])

        if not scores:
            return None

        best_category = max(scores, key=scores.get)
        intent, raw_category = matched_meta[best_category]

        return intent, raw_category, best_category, 0.95, text

    def predict_by_model(self, text: str) -> tuple[str, str, str, float, str]:
        translated = self.translate_to_english(text)

        pred_encoded = self.pipeline.predict([translated])[0]
        intent = self.label_encoder.inverse_transform([pred_encoded])[0]

        raw_category = self.intent_to_category.get(intent, "UNKNOWN")
        short_category = INTENT_TO_SHORT_CATEGORY.get(intent, "SUPPORT")

        return intent, raw_category, short_category, 1.0, translated

    def predict(self, text: str) -> tuple[str, str, str, float, str]:
        if not text or not text.strip():
            raise ValueError("empty text cannot be classified")

        rule_result = self.predict_by_rules(text)
        if rule_result is not None:
            return rule_result

        return self.predict_by_model(text)