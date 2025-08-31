POSTGRES_DSN = "postgres://admin:admin@localhost:5432/go_platform?sslmode=disable"
MYSQL_DSN = "admin:admin@tcp(localhost:3306)/go_platform?parseTime=true&loc=UTC"
CLICKHOUSE_DSN = "clickhouse://admin:admin@localhost:9000/go_platform"
MIGRATION_DIR_PG = ./migrations/postgres/
MIGRATION_DIR_MYSQL = ./migrations/mysql/
MIGRATION_DIR_CLICKHOUSE = ./migrations/clickhouse/
APP_DIR= ./cmd/app

.PHONY: run run-pg run-mysql run-ch up-full-pg up-full-mysql up-full-ch up-pg up-mysql up-ch down migrate-up migrate-down migrate-status migrate-create test test-verbose swagger-init dev-pg dev-mysql dev-ch full-pg full-mysql full-ch


# Run application with PostgreSQL
run-pg:
	@export $$(grep -v '^#' ./docker/.env-pg | xargs) >/dev/null 2>&1; \
	go run $(APP_DIR)/main.go

# Run application with MySQL
run-mysql:
	@export $$(grep -v '^#' ./docker/.env-mysql | xargs) >/dev/null 2>&1; \
	go run $(APP_DIR)/main.go

# Run application with ClickHouse
run-ch:
	@export $$(grep -v '^#' ./docker/.env-ch | xargs) >/dev/null 2>&1; \
	go run $(APP_DIR)/main.go

# Start full stack with PostgreSQL (infrastructure + app)
full-pg:
	COMPOSE_PROJECT_NAME=go-platform docker compose -f docker/docker-compose.postgres.yml --env-file=docker/.env-docker-pg --profile=test up -d --build 

# Start full stack with MySQL (infrastructure + app)
full-mysql:
	COMPOSE_PROJECT_NAME=go-platform docker compose -f docker/docker-compose.mysql.yml --env-file=docker/.env-docker-mysql --profile=test up -d --build 

# Start full stack with ClickHouse (infrastructure + app)
full-ch:
	COMPOSE_PROJECT_NAME=go-platform docker compose -f docker/docker-compose.clickhouse.yml --env-file=docker/.env-docker-ch --profile=test up -d --build 

# Start PostgreSQL infrastructure
up-pg:
	COMPOSE_PROJECT_NAME=go-platform docker compose -f docker/docker-compose.postgres.yml --env-file=docker/.env-docker-pg up -d --build

# Start MySQL infrastructure
up-mysql:
	COMPOSE_PROJECT_NAME=go-platform docker compose -f docker/docker-compose.mysql.yml --env-file=docker/.env-docker-mysql up -d --build

# Start ClickHouse infrastructure
up-ch:
	COMPOSE_PROJECT_NAME=go-platform docker compose -f docker/docker-compose.clickhouse.yml --env-file=docker/.env-docker-ch up -d --build

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

# Run quick load test
test-quick:
	@echo "Running quick test..."
	./tests/quick_test.sh

# Run aggressive load test
test-load:
	@echo "Running aggressive load test..."
	./tests/load_test.sh

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
	docker exec platform_minio mc alias set myminio http://minio:9000 minioadmin minioadmin
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

