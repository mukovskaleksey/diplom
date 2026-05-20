from __future__ import annotations

import sys
from pathlib import Path
from concurrent import futures

import grpc


BASE_DIR = Path(__file__).resolve().parent
GEN_DIR = BASE_DIR / "gen"

if str(GEN_DIR) not in sys.path:
    sys.path.insert(0, str(GEN_DIR))

from analysis import analysis_pb2, analysis_pb2_grpc
from app.model_service import TextClassifier


class AnalysisServiceServicer(analysis_pb2_grpc.AnalysisServiceServicer):
    def __init__(self) -> None:
        self.classifier = TextClassifier()

    def ClassifyMessage(self, request, context):
        message = request.message.strip()

        if not message:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details("message is required")
            return analysis_pb2.ClassifyMessageResponse()

        try:
            intent, raw_category, short_category, confidence, translated = self.classifier.predict(message)

            return analysis_pb2.ClassifyMessageResponse(
                intent=intent,
                raw_category=raw_category,
                category=short_category,
                confidence=confidence,
                translated_text=translated,
            )
        except Exception as exc:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(exc))
            return analysis_pb2.ClassifyMessageResponse()


def serve() -> None:
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    analysis_pb2_grpc.add_AnalysisServiceServicer_to_server(
        AnalysisServiceServicer(),
        server,
    )

    server.add_insecure_port("[::]:50052")
    server.start()

    print("analysis-service started on :50052")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()