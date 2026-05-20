export GOWORK=off
export GOPROXY=https://proxy.golang.org,direct
export GOSUMDB=sum.golang.org
export GOPRIVATE=
export GONOSUMDB=
export GONOPROXY=

APP_DB_URL=postgres://postgres:postgres@localhost:5430/chat?sslmode=disable
MIGRATE_DB_URL=postgres://postgres:postgres@postgres:5432/chat?sslmode=disable
DOCKER_COMPOSE=docker compose -p chat-service -f deployments/docker-compose.yaml

.PHONY: clean run build tidy docker-up docker-down docker-logs db-wait migrate-up migrate-down bootstrap

run:
	@echo "Starting chat-service..."
	GOPROXY=$(GOPROXY) GOSUMDB=$(GOSUMDB) GOWORK=$(GOWORK) \
	GOPRIVATE= GONOSUMDB= GONOPROXY= \
	go run main.go

build:
	@echo "Building..."
	@mkdir -p bin
	GOPROXY=$(GOPROXY) GOSUMDB=$(GOSUMDB) GOWORK=$(GOWORK) \
	GOPRIVATE= GONOSUMDB= GONOPROXY= \
	go build -o bin/chat-service main.go
	@echo "Done"

tidy:
	@echo "Tidying..."
	GOPROXY=$(GOPROXY) GOSUMDB=$(GOSUMDB) GOWORK=$(GOWORK) \
	GOPRIVATE= GONOSUMDB= GONOPROXY= \
	go mod tidy
	@echo "Done"

clean:
	@echo "Cleaning..."
	rm -rf bin
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
	@until docker exec chat-service-postgres pg_isready -U postgres -d chat > /dev/null 2>&1; do \
		sleep 2; \
	done
	@echo "Postgres is ready"

migrate-up:
	@echo "Applying migrations..."
	docker run --rm \
		--network chat-service_default \
		-v $(PWD)/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$(MIGRATE_DB_URL)" \
		up
	@echo "Migrations applied"

migrate-down:
	@echo "Rolling back migrations..."
	docker run --rm \
		--network chat-service_default \
		-v $(PWD)/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$(MIGRATE_DB_URL)" \
		down -all
	@echo "Rollback complete"

bootstrap: docker-up db-wait migrate-up
	@echo "Bootstrap complete"