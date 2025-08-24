POSTGRES_DSN = "postgres://admin:admin@localhost:5432/go_platform?sslmode=disable"
MIGRATION_DIR_PG = ./migrations/postgres/
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

# Database migrations
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

proto-all:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
    --go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
    api/protobuf/*.proto

