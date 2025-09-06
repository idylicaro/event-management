# Event Management Makefile
# This Makefile provides commands for development, database migrations, and Docker operations

# Load environment variables from .env file
include .env
export

# Default target
.DEFAULT_GOAL := help

# Variables
MIGRATIONS_PATH := migrations
DOCKER_COMPOSE_FILE := docker-compose.yml
CONTAINER_NAME := postgres
DOCKER_COMPOSE := docker compose

# Local development database URL (for when running outside Docker)
LOCAL_DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

# Docker development database URL (for migrations inside Docker network)
DOCKER_DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

##@ Help
help: ## Display this help
	@echo "Event Management Development Commands"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
dev: ## Start development environment with Docker Compose
	@echo "$(GREEN)Starting development environment...$(NC)"
	$(DOCKER_COMPOSE) up --build

dev-detached: ## Start development environment in background
	@echo "$(GREEN)Starting development environment in background...$(NC)"
	$(DOCKER_COMPOSE) up --build -d

stop: ## Stop development environment
	@echo "$(YELLOW)Stopping development environment...$(NC)"
	$(DOCKER_COMPOSE) down

clean: ## Stop and remove containers, networks, and volumes
	@echo "$(RED)Cleaning up Docker environment...$(NC)"
	$(DOCKER_COMPOSE) down --volumes --remove-orphans

logs: ## Show application logs
	$(DOCKER_COMPOSE) logs -f app

logs-db: ## Show database logs
	$(DOCKER_COMPOSE) logs -f db

##@ Build
build: ## Build the Go application
	@echo "$(GREEN)Building application...$(NC)"
	go build -o event-management .

run: ## Run the application locally (requires local database)
	@echo "$(GREEN)Running application locally...$(NC)"
	go run .

test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	go test ./...

##@ Database Migrations (Local Development)
migrate-up: ## Run all up migrations (local database)
	@echo "$(GREEN)Running migrations up (local)...$(NC)"
	migrate -path $(MIGRATIONS_PATH) -database "$(LOCAL_DATABASE_URL)" up

migrate-down: ## Run one down migration (local database)
	@echo "$(YELLOW)Running migration down (local)...$(NC)"
	migrate -path $(MIGRATIONS_PATH) -database "$(LOCAL_DATABASE_URL)" down 1

migrate-version: ## Show current migration version (local database)
	@echo "$(GREEN)Current migration version (local):$(NC)"
	migrate -path $(MIGRATIONS_PATH) -database "$(LOCAL_DATABASE_URL)" version

migrate-force: ## Force migration to specific version (local database) - Usage: make migrate-force VERSION=1
	@echo "$(RED)Forcing migration to version $(VERSION) (local)...$(NC)"
	migrate -path $(MIGRATIONS_PATH) -database "$(LOCAL_DATABASE_URL)" force $(VERSION)

##@ Database Migrations (Docker Environment)
docker-migrate-up: ## Run all up migrations in Docker environment
	@echo "$(GREEN)Running migrations up (Docker)...$(NC)"
	$(DOCKER_COMPOSE) exec app migrate -path $(MIGRATIONS_PATH) -database "$(DOCKER_DATABASE_URL)" up

docker-migrate-down: ## Run one down migration in Docker environment
	@echo "$(YELLOW)Running migration down (Docker)...$(NC)"
	$(DOCKER_COMPOSE) exec app migrate -path $(MIGRATIONS_PATH) -database "$(DOCKER_DATABASE_URL)" down 1

docker-migrate-version: ## Show current migration version in Docker environment
	@echo "$(GREEN)Current migration version (Docker):$(NC)"
	$(DOCKER_COMPOSE) exec app migrate -path $(MIGRATIONS_PATH) -database "$(DOCKER_DATABASE_URL)" version

docker-migrate-force: ## Force migration to specific version in Docker - Usage: make docker-migrate-force VERSION=1
	@echo "$(RED)Forcing migration to version $(VERSION) (Docker)...$(NC)"
	$(DOCKER_COMPOSE) exec app migrate -path $(MIGRATIONS_PATH) -database "$(DOCKER_DATABASE_URL)" force $(VERSION)

##@ Database Operations
db-shell: ## Connect to PostgreSQL shell in Docker
	@echo "$(GREEN)Connecting to database shell...$(NC)"
	$(DOCKER_COMPOSE) exec -it db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

db-tables: ## Show database tables
	@echo "$(GREEN)Database tables:$(NC)"
	$(DOCKER_COMPOSE) exec db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "\dt"

db-reset: ## Reset database (drop and recreate)
	@echo "$(RED)Resetting database...$(NC)"
	$(DOCKER_COMPOSE) exec db psql -U $(POSTGRES_USER) -c "DROP DATABASE IF EXISTS $(POSTGRES_DB);"
	$(DOCKER_COMPOSE) exec db psql -U $(POSTGRES_USER) -c "CREATE DATABASE $(POSTGRES_DB);"

##@ Installation
install-migrate: ## Install golang-migrate tool
	@echo "$(GREEN)Installing golang-migrate...$(NC)"
	@if command -v migrate >/dev/null 2>&1; then \
		echo "$(YELLOW)golang-migrate is already installed$(NC)"; \
	else \
		echo "$(GREEN)Installing golang-migrate...$(NC)"; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
		echo "$(GREEN)golang-migrate installed successfully$(NC)"; \
	fi

check-tools: ## Check if required tools are installed
	@echo "$(GREEN)Checking required tools...$(NC)"
	@command -v docker >/dev/null 2>&1 || (echo "$(RED)Docker is not installed$(NC)" && exit 1)
	@docker compose version >/dev/null 2>&1 || (echo "$(RED)Docker Compose is not installed$(NC)" && exit 1)
	@command -v go >/dev/null 2>&1 || (echo "$(RED)Go is not installed$(NC)" && exit 1)
	@command -v migrate >/dev/null 2>&1 || (echo "$(YELLOW)golang-migrate is not installed. Run 'make install-migrate'$(NC)")
	@echo "$(GREEN)All required tools are available$(NC)"

##@ Utilities
health: ## Check application health
	@echo "$(GREEN)Checking application health...$(NC)"
	curl -f http://localhost:8080/health || echo "$(RED)Application is not responding$(NC)"

status: ## Show Docker containers status
	@echo "$(GREEN)Docker containers status:$(NC)"
	$(DOCKER_COMPOSE) ps

##@ Quick Setup
setup: install-migrate dev-detached docker-migrate-up ## Quick setup: install tools, start environment, run migrations
	@echo "$(GREEN)Development environment is ready!$(NC)"
	@echo "$(YELLOW)Application: http://localhost:8080$(NC)"
	@echo "$(YELLOW)Swagger API: http://localhost:8080/swagger/index.html$(NC)"
	@echo "$(YELLOW)Database: localhost:5432$(NC)"

.PHONY: help dev dev-detached stop clean logs logs-db build run test \
        migrate-up migrate-down migrate-version migrate-force \
        docker-migrate-up docker-migrate-down docker-migrate-version docker-migrate-force \
        db-shell db-tables db-reset install-migrate check-tools health status setup
