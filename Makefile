APP_NAME ?= shanraq
GO ?= go
MAIN_PKG ?= ./cmd/app
CLI_MIGRATE ?= ./cmd/cli/migrate
BIN_DIR ?= bin
BIN_APP ?= $(BIN_DIR)/$(APP_NAME)
BIN_MIGRATE ?= $(BIN_DIR)/migrate

.PHONY: all build run clean fmt lint test tidy deps migrate-up migrate-down migrate-steps seed watch

all: build

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

build: $(BIN_DIR)
	$(GO) build -o $(BIN_APP) $(MAIN_PKG)

run:
	$(GO) run $(MAIN_PKG)

fmt:
	$(GO) fmt ./...

lint:
	$(GO) vet ./...

test:
	$(GO) test ./...

tidy:
	$(GO) mod tidy

deps: tidy
	$(GO) mod download

$(BIN_MIGRATE): $(BIN_DIR)
	$(GO) build -o $(BIN_MIGRATE) $(CLI_MIGRATE)

MIGRATE ?= $(BIN_MIGRATE)
MIGRATIONS_DIR ?= migrations
DATABASE_URL ?=

migrate-up: $(BIN_MIGRATE)
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	$(MIGRATE) -database "$(DATABASE_URL)" -dir $(MIGRATIONS_DIR) -command up

migrate-down: $(BIN_MIGRATE)
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	$(MIGRATE) -database "$(DATABASE_URL)" -dir $(MIGRATIONS_DIR) -command down -steps 1

migrate-steps: $(BIN_MIGRATE)
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	@if [ -z "$(steps)" ]; then echo "steps=<N> is required"; exit 1; fi
	$(MIGRATE) -database "$(DATABASE_URL)" -dir $(MIGRATIONS_DIR) -command steps -steps $(steps)

seed: migrate-up
	@echo "Seed data applied via migrations"

clean:
	rm -rf $(BIN_DIR)

watch:
	@command -v air >/dev/null 2>&1 && air || reflex -r '\.go$$' -- sh -c '$(GO) run $(MAIN_PKG)'
