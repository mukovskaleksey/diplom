from __future__ import annotations

import sys
from pathlib import Path

import grpc

BASE_DIR = Path(__file__).resolve().parent.parent
GEN_DIR = BASE_DIR / "gen"

# добавляем пути
sys.path.insert(0, str(BASE_DIR))
sys.path.insert(0, str(GEN_DIR))

from analysis import analysis_pb2, analysis_pb2_grpc


def main() -> None:
    tests = [
        "I forgot my password and cannot log in",
        "I want to cancel my order",
        "Where is my refund?",
        "My card was charged twice",
        "I need to change my shipping address",

        "Я забыл пароль и не могу войти",
        "Хочу отменить заказ",
        "Где мой возврат?",
        "С меня списали деньги два раза",
        "Хочу изменить адрес доставки",

        "Не получается зарегистрировать аккаунт",
        "Как связаться с оператором поддержки?",
        "Когда привезут мой заказ?",
        "Мне нужен счет за покупку",
        "Хочу оформить новый заказ",
    ]

    with grpc.insecure_channel("localhost:50052") as channel:
        stub = analysis_pb2_grpc.AnalysisServiceStub(channel)

        for text in tests:
            response = stub.ClassifyMessage(
                analysis_pb2.ClassifyMessageRequest(message=text)
            )
            print("-" * 80)
            print("text:", text)
            print("translated:", response.translated_text)
            print("intent:", response.intent)
            print("raw_category:", response.raw_category)
            print("category:", response.category)
            print("confidence:", response.confidence)


if __name__ == "__main__":
    main()