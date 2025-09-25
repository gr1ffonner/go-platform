POSTGRES_DSN = "postgres://admin:admin@localhost:5432/go_platform?sslmode=disable"
MIGRATION_DIR_PG = ./migrations/postgres/
APP_DIR= ./cmd/app


# Run application with PostgreSQL
run:
	@export $$(grep -v '^#' .env | xargs) >/dev/null 2>&1; \
	go run $(APP_DIR)/main.go

# Start all services
up:
	COMPOSE_PROJECT_NAME=go-platform docker compose -f docker-compose.yml --profile=test --env-file=.env-docker up -d --build

# Start only infrastructure
up-infra:
	COMPOSE_PROJECT_NAME=go-platform docker compose -f docker-compose.yml --env-file=.env up -d 

# Stop all services
down:
	COMPOSE_PROJECT_NAME=go-platform docker compose -f docker-compose.yml --profile=test --env-file=.env-docker down --remove-orphans

# Stop and remove volumes
down-infra:
	COMPOSE_PROJECT_NAME=go-platform docker compose -f docker-compose.yml --env-file=.env down --remove-orphans

clean:
	COMPOSE_PROJECT_NAME=go-platform docker compose -f docker-compose.yml --profile=test --env-file=.env-docker down -v --remove-orphans

clean-infra:
	COMPOSE_PROJECT_NAME=go-platform docker compose -f docker-compose.yml --env-file=.env down -v --remove-orphans



# Run all tests
test:
	@echo "Running all tests..."
	go test ./... -v

# Run tests with verbose output and coverage
test-verbose:
	@echo "Running tests with verbose output and coverage..."
	go test ./... -v -cover

# Database migrations - PostgreSQL
migrate-up:
	@echo "Applying PostgreSQL migrations..."
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) up

migrate-rollback:
	@echo "Rolling back all PostgreSQL migrations..."
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) reset

migrate-down:
	@echo "Rolling back last PostgreSQL migration..."
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) down

migrate-status:
	@echo "PostgreSQL migration status:"
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) status

migrate-create:
	@read -p "Enter PostgreSQL migration name: " NAME; \
	goose -dir $(MIGRATION_DIR_PG) create $$NAME sql

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

