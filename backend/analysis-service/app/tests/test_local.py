from __future__ import annotations

import sys
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent.parent
sys.path.insert(0, str(BASE_DIR))

from model_service import TextClassifier


def main() -> None:
    clf = TextClassifier()

    tests = [
        "Я забыл пароль",
        "Хочу отменить заказ",
        "Где мой возврат?",
        "С меня списали деньги два раза",
        "Хочу изменить адрес доставки",
    ]

    for text in tests:
        intent, raw_category, category, confidence, translated = clf.predict(text)

        print("-" * 80)
        print("text:", text)
        print("translated:", translated)
        print("intent:", intent)
        print("raw_category:", raw_category)
        print("category:", category)
        print("confidence:", confidence)


if __name__ == "__main__":
    main()