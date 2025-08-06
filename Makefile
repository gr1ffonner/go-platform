POSTGRES_DSN = "postgres://admin:admin@localhost:5432/go_platform?sslmode=disable"
MIGRATION_DIR_PG = ./migrations/postgres/
APP_DIR= ./cmd/app
PRODUCER_DIR= ./cmd/producer

.PHONY: run up up-dev down migrate-up-pg migrate-down-pg migrate-status-pg migrate-create-pg test test-jwt test-verbose

run:
	@export $$(grep -v '^#' ./.env | xargs) >/dev/null 2>&1; \
	go run $(APP_DIR)/main.go

run-producer:
	@export $$(grep -v '^#' ./.env | xargs) >/dev/null 2>&1; \
	go run $(PRODUCER_DIR)/main.go

# Infrastructure only (postgres, redis, nats)
up-infra:
	COMPOSE_PROJECT_NAME=go-platform docker compose --profile infra up -d

# Development environment (infra + app)
up-dev:
	COMPOSE_PROJECT_NAME=go-platform docker compose --profile dev up -d --build

# Test environment
up-test:
	COMPOSE_PROJECT_NAME=go-platform docker compose --profile test up -d --build

# Full stack (all services)
up-full:
	COMPOSE_PROJECT_NAME=go-platform docker compose --profile full up -d

# Messaging stack (kafka, rabbitmq)
up-messaging:
	COMPOSE_PROJECT_NAME=go-platform docker compose --profile messaging up -d

# Analytics stack (clickhouse)
up-analytics:
	COMPOSE_PROJECT_NAME=go-platform docker compose --profile analytics up -d

# Storage (minio s3)
up-storage:
	COMPOSE_PROJECT_NAME=go-platform docker compose --profile storage up -d

# Monitoring stack
up-monitoring:
	COMPOSE_PROJECT_NAME=go-platform docker compose --profile monitoring up -d

# Production-like environment
up-prod:
	COMPOSE_PROJECT_NAME=go-platform docker compose --env-file=.env-docker --profile prod up -d

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

migrate-up:
	@echo "Applying PostgreSQL migrations..."
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) up

migrate-down:
	@echo "Rolling back last PostgreSQL migration..."
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) down


migrate-status:
	@echo "PostgreSQL migration status:"
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) status


migrate-create:
	@read -p "Enter PostgreSQL migration name: " NAME; \
	goose -dir $(MIGRATION_DIR_PG) create $$NAME sql

check-swagger: 
	@command -v which swag >/dev/null 2>&1 || { \
		echo "swaggo not found, installing..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	}

# Initialize Swagger documentation
swagger-init: check-swagger 
	swag fmt && swag init --pdl=1 -g cmd/app/main.go -o api/
