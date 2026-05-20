PROTO_DIR=api/proto
GEN_DIR=api/proto/gen/pb
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

export GOWORK=off
export GOPROXY=https://proxy.golang.org,direct
export GOSUMDB=sum.golang.org
export GOPRIVATE=
export GONOSUMDB=
export GONOPROXY=

.PHONY: tools deps proto clean run build tidy env

APP_DB_URL=postgres://postgres:postgres@localhost:5431/support?sslmode=disable
MIGRATE_DB_URL=postgres://postgres:postgres@postgres:5432/support?sslmode=disable
DOCKER_COMPOSE=docker compose -p core-service -f deployments/docker-compose.yaml

.PHONY: proto clean tidy run build docker-up docker-down docker-logs db-wait migrate-up migrate-down bootstrap




proto:
	@echo "Generating protobuf..."
	@mkdir -p $(GEN_DIR)
	protoc \
		--go_out=$(GEN_DIR) \
		--go-grpc_out=$(GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "Done"

run:
	@echo "Starting core-service..."
	GOPROXY=$(GOPROXY) GOSUMDB=$(GOSUMDB) GOWORK=$(GOWORK) \
	GOPRIVATE= GONOSUMDB= GONOPROXY= \
	go run main.go

build:
	@echo "Building..."
	@mkdir -p bin
	GOPROXY=$(GOPROXY) GOSUMDB=$(GOSUMDB) GOWORK=$(GOWORK) \
	GOPRIVATE= GONOSUMDB= GONOPROXY= \
	go build -o bin/core-service main.go
	@echo "Done"

tidy:
	@echo "Tidying..."
	GOPROXY=$(GOPROXY) GOSUMDB=$(GOSUMDB) GOWORK=$(GOWORK) \
	GOPRIVATE= GONOSUMDB= GONOPROXY= \
	go mod tidy
	@echo "Done"

clean:
	@echo "Cleaning generated files..."
	rm -rf $(GEN_DIR)
	@echo "Done"

docker-up:
	@echo "Starting postgres container..."
	$(DOCKER_COMPOSE) up -d
	@echo "Done"

docker-down:
	@echo "Stopping containers..."
	$(DOCKER_COMPOSE) down
	@echo "Done"

docker-logs:
	$(DOCKER_COMPOSE) logs -f

db-wait:
	@echo "Waiting for postgres..."
	@until docker exec core-service-postgres pg_isready -U postgres -d support > /dev/null 2>&1; do \
		sleep 2; \
	done
	@echo "Postgres is ready"

migrate-up:
	@echo "Applying migrations..."
	docker run --rm \
		--network core-service_default \
		-v $(PWD)/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$(MIGRATE_DB_URL)" \
		up
	@echo "Migrations applied"

migrate-down:
	@echo "Rolling back migrations..."
	docker run --rm \
		--network core-service_default \
		-v $(PWD)/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$(MIGRATE_DB_URL)" \
		down -all
	@echo "Rollback complete"

bootstrap: docker-up db-wait migrate-up
	@echo "Bootstrap complete"