#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"

echo "📂 SCRIPT_DIR=$SCRIPT_DIR"
echo "📂 BACKEND_DIR=$BACKEND_DIR"

cd "$BACKEND_DIR"

echo "🔄 Cleaning old generated files..."

rm -rf analysis
rm -rf core
rm -rf backend

rm -f proto/core/*.pb.go
rm -f proto/core/*_grpc.pb.go
rm -f proto/chat/*.pb.go
rm -f proto/chat/*_grpc.pb.go
rm -f proto/analysis/*.pb.go
rm -f proto/analysis/*_grpc.pb.go

rm -rf analysis-service/app/gen
mkdir -p analysis-service/app/gen
mkdir -p analysis-service/app/gen/analysis
touch analysis-service/app/gen/__init__.py
touch analysis-service/app/gen/analysis/__init__.py

echo "⚙️ Generating Go protobuf..."

protoc \
  -I ./proto \
  --go_out=./proto \
  --go_opt=paths=source_relative \
  --go-grpc_out=./proto \
  --go-grpc_opt=paths=source_relative \
  ./proto/core/core.proto

protoc \
  -I ./proto \
  --go_out=./proto \
  --go_opt=paths=source_relative \
  --go-grpc_out=./proto \
  --go-grpc_opt=paths=source_relative \
  ./proto/chat/chat.proto

protoc \
  -I ./proto \
  --go_out=./proto \
  --go_opt=paths=source_relative \
  --go-grpc_out=./proto \
  --go-grpc_opt=paths=source_relative \
  ./proto/analysis/analysis.proto

echo "⚙️ Generating Python protobuf..."

PYTHON_BIN="$BACKEND_DIR/analysis-service/venv/bin/python"

if [ ! -x "$PYTHON_BIN" ]; then
  echo "❌ Python venv not found: $PYTHON_BIN"
  exit 1
fi

"$PYTHON_BIN" -m grpc_tools.protoc \
  -I ./proto \
  --python_out=./analysis-service/app/gen \
  --grpc_python_out=./analysis-service/app/gen \
  ./proto/analysis/analysis.proto

touch analysis-service/app/gen/__init__.py
touch analysis-service/app/gen/analysis/__init__.py

echo "✅ Proto generation complete"