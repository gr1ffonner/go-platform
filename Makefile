POSTGRES_DSN = "postgres://admin:admin@localhost:5432/go_platform?sslmode=disable"
MYSQL_DSN = "admin:admin@tcp(localhost:3306)/go_platform?parseTime=true&loc=UTC"
CLICKHOUSE_DSN = "clickhouse://admin:admin@localhost:9000/go_platform"
MIGRATION_DIR_PG = ./migrations/postgres/
MIGRATION_DIR_MYSQL = ./migrations/mysql/
MIGRATION_DIR_CLICKHOUSE = ./migrations/clickhouse/
APP_DIR= ./cmd/app

.PHONY: run up up-dev down migrate-up migrate-down migrate-status migrate-create test test-verbose swagger-init

# Run application locally
run:
	@export $$(grep -v '^#' ./.env | xargs) >/dev/null 2>&1; \
	go run $(APP_DIR)/main.go

# Start the services with Docker Compose (infrastructure + app)
up: 
	COMPOSE_PROJECT_NAME=go-platform docker compose -f docker-compose.yml --env-file=.env-docker --profile=test up -d --build 

# Start the infrastructure with Docker Compose
up-dev: 
	COMPOSE_PROJECT_NAME=go-platform docker compose -f docker-compose.yml --env-file=.env up -d --build

# Stop all services
down:
	COMPOSE_PROJECT_NAME=go-platform docker compose down

# Stop and remove volumes
down-clean:
	COMPOSE_PROJECT_NAME=go-platform docker compose down -v

# Run all tests
test:
	@echo "Running all tests..."
	go test ./... -v

# Run tests with verbose output and coverage
test-verbose:
	@echo "Running tests with verbose output and coverage..."
	go test ./... -v -cover

# Database migrations - PostgreSQL
migrate-up-pg:
	@echo "Applying PostgreSQL migrations..."
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) up

migrate-down-pg:
	@echo "Rolling back last PostgreSQL migration..."
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) down

migrate-status-pg:
	@echo "PostgreSQL migration status:"
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) status

migrate-create-pg:
	@read -p "Enter PostgreSQL migration name: " NAME; \
	goose -dir $(MIGRATION_DIR_PG) create $$NAME sql

# Database migrations - MySQL
migrate-up-mysql:
	@echo "Applying MySQL migrations..."
	goose -dir $(MIGRATION_DIR_MYSQL) mysql $(MYSQL_DSN) up

migrate-down-mysql:
	@echo "Rolling back last MySQL migration..."
	goose -dir $(MIGRATION_DIR_MYSQL) mysql $(MYSQL_DSN) down

migrate-status-mysql:
	@echo "MySQL migration status:"
	goose -dir $(MIGRATION_DIR_MYSQL) mysql $(MYSQL_DSN) status

migrate-create-mysql:
	@read -p "Enter MySQL migration name: " NAME; \
	goose -dir $(MIGRATION_DIR_MYSQL) create $$NAME sql

# Database migrations - ClickHouse
migrate-up-clickhouse:
	@echo "Applying ClickHouse migrations..."
	goose -dir $(MIGRATION_DIR_CLICKHOUSE) clickhouse $(CLICKHOUSE_DSN) up

migrate-down-clickhouse:
	@echo "Rolling back last ClickHouse migration..."
	goose -dir $(MIGRATION_DIR_CLICKHOUSE) clickhouse $(CLICKHOUSE_DSN) down

migrate-status-clickhouse:
	@echo "ClickHouse migration status:"
	goose -dir $(MIGRATION_DIR_CLICKHOUSE) clickhouse $(CLICKHOUSE_DSN) status

migrate-create-clickhouse:
	@read -p "Enter ClickHouse migration name: " NAME; \
	goose -dir $(MIGRATION_DIR_CLICKHOUSE) create $$NAME sql

# All database migrations
migrate-up: migrate-up-pg migrate-up-mysql migrate-up-clickhouse
migrate-down: migrate-down-pg migrate-down-mysql migrate-down-clickhouse
migrate-status: migrate-status-pg migrate-status-mysql migrate-status-clickhouse

bucket-create:
	docker exec platform_minio mc alias set myminio http://localhost:9000 minioadmin minioadmin
	docker exec platform_minio mc mb myminio/dogs

# Swagger documentation
check-swagger: 
	@command -v which swag >/dev/null 2>&1 || { \
		echo "swaggo not found, installing..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	}

swagger-init: check-swagger 
	swag fmt && swag init --pdl=1 -g cmd/app/main.go -o api/

# Protobuf
proto-all:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
    --go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
    api/protobuf/*.proto

