# =============================================================================
# orion-v2 — Root Makefile
# =============================================================================
# Targets are namespaced by service: product/*, db/*, migrate/*, etc.
# Run `make help` to see all available targets.
# =============================================================================

.DEFAULT_GOAL := help

# ── Variables ─────────────────────────────────────────────────────────────────
PRODUCT_SVC    := ./services/product-service
MIGRATIONS_DIR := $(PRODUCT_SVC)/migrations

# Load root .env so DATABASE_DSN is available for migrate/_ targets.
ifneq (,$(wildcard .env))
  include .env
  export
endif

# ── Help ──────────────────────────────────────────────────────────────────────
.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_/%-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'

# ── Tools ─────────────────────────────────────────────────────────────────────
.PHONY: install/tools
install/tools: ## Install all required CLI tools (run once per machine)
	@echo "→ Installing golang-migrate CLI..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "✓ Tools installed."

# ── Docker ────────────────────────────────────────────────────────────────────
.PHONY: up down down/v logs/db
up: ## Start all docker compose services
	docker compose up -d

down: ## Stop and remove all containers (keeps volumes)
	docker compose down

down/v: ## Stop and remove containers + volumes (WARNING: destroys all data)
	docker compose down -v

logs/db: ## Follow postgres container logs
	docker compose logs -f db

# ── Build ─────────────────────────────────────────────────────────────────────
.PHONY: build/product build/all
build/product: ## Build the product-service binary
	go build -o bin/product-service $(PRODUCT_SVC)/cmd/main.go

build/all: ## Build all service binaries
	go build ./...

# ── Run ───────────────────────────────────────────────────────────────────────
.PHONY: run/product
run/product: ## Run product-service locally (requires .env)
	go run $(PRODUCT_SVC)/cmd/main.go

# ── Migrations (product-service) ──────────────────────────────────────────────
# Requires: make install/tools  (installs the `migrate` CLI once)
# Requires: DATABASE_DSN set in .env at the repo root
.PHONY: migrate/up migrate/down migrate/down/1 migrate/version migrate/force

migrate/up: ## Apply all pending migrations (product-service)
	@echo "→ migrate up (product-service)..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_DSN)" up

migrate/down: ## Rollback ALL migrations (product-service) ⚠️ destructive
	@echo "→ migrate down ALL (product-service)..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_DSN)" down

migrate/down/1: ## Rollback the last migration (product-service)
	@echo "→ migrate down 1 (product-service)..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_DSN)" down 1

migrate/version: ## Show current migration version (product-service)
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_DSN)" version

migrate/force: ## Force-set migration version — usage: make migrate/force V=2
	@test -n "$(V)" || (echo "Usage: make migrate/force V=<version>"; exit 1)
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_DSN)" force $(V)

# ── Tidy ──────────────────────────────────────────────────────────────────────
.PHONY: tidy
tidy: ## Run go mod tidy across all modules
	cd pkg && go mod tidy
	cd $(PRODUCT_SVC) && GOWORK=off go mod tidy
	go work sync
