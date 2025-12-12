.PHONY: help setup run build test clean deps docker-up docker-down migrate-up migrate-down migrate-status migrate-create

# Default target
help:
	@echo "Available commands:"
	@echo "  make setup           - Full setup (deps, docker, migrations)"
	@echo "  make run             - Run the application"
	@echo "  make build           - Build the application"
	@echo "  make test            - Run tests"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make docker-up       - Start development services (DB, Redis, MinIO)"
	@echo "  make docker-down     - Stop development services"
	@echo "  make migrate-up      - Run all pending migrations"
	@echo "  make migrate-down    - Rollback last migration"
	@echo "  make migrate-create  - Create new migration (usage: make migrate-create NAME=x)"

# Full Setup: Deps -> Docker -> Migrate
setup: deps docker-up
	@echo "Waiting for services to be ready..."
	@sleep 5
	@make migrate-up
	@echo "Setup complete! Run 'make run' to start."

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Dependencies installed"

# Docker
docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down

# Run the application
run:
	@go run main.go

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/app main.go
	@echo "Build complete: bin/app"

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f example-app
	@echo "Clean complete"

# Migrations
migrate-up:
	@echo "Running migrations..."
	@go run cmd/migrate/main.go -direction=up

migrate-down:
	@echo "Rolling back last migration..."
	@go run cmd/migrate/main.go -direction=down -steps=1

migrate-status:
	@echo "Checking migration status..."
	@go run cmd/migrate/main.go -direction=version

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	@migrate create -ext sql -dir database/migrations -seq $(NAME)
	@echo "Migration files created in database/migrations/"
