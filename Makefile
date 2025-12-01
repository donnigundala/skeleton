.PHONY: help run build test migrate-up migrate-down migrate-status migrate-create clean

# Default target
help:
	@echo "Available commands:"
	@echo "  make run             - Run the application"
	@echo "  make build           - Build the application"
	@echo "  make test            - Run tests"
	@echo "  make migrate-up      - Run all pending migrations"
	@echo "  make migrate-down    - Rollback last migration"
	@echo "  make migrate-status  - Show migration status"
	@echo "  make migrate-create  - Create new migration (use NAME=migration_name)"
	@echo "  make clean           - Clean build artifacts"

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

# Run all pending migrations
migrate-up:
	@echo "Running migrations..."
	@go run cmd/migrate/main.go -direction=up

# Rollback last migration
migrate-down:
	@echo "Rolling back last migration..."
	@go run cmd/migrate/main.go -direction=down -steps=1

# Show migration status
migrate-status:
	@echo "Checking migration status..."
	@go run cmd/migrate/main.go -direction=version

# Create new migration
# Usage: make migrate-create NAME=add_posts_table
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	@migrate create -ext sql -dir database/migrations -seq $(NAME)
	@echo "Migration files created in database/migrations/"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f example-app
	@echo "Clean complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# Install go-migrate CLI
install-migrate:
	@echo "Installing go-migrate CLI..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "go-migrate CLI installed"

# Development database setup
db-setup:
	@echo "Setting up development database..."
	@createdb myapp || echo "Database may already exist"
	@make migrate-up
	@echo "Database setup complete"

# Reset database (WARNING: Deletes all data!)
db-reset:
	@echo "WARNING: This will delete all data!"
	@read -p "Are you sure? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	@make migrate-down
	@dropdb myapp || echo "Database may not exist"
	@make db-setup
	@echo "Database reset complete"
