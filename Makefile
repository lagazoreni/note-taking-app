.PHONY: help dev-server dev-web build-web build-server fmt lint audit test-server test-web test-contract test-e2e test-offline test-import-export test-performance check image smoke-image db-status db-check db-reset

SHELL := /bin/sh
GO := go
NPM := npm
WEB_DIR := web
SERVER_DIR := server

help:
	@printf '%s\n' 'Use make build-web build-server check image'

dev-server:
	mkdir -p .local/data
	cd $(SERVER_DIR) && NOTED_DATA_DIR=../.local/data $(GO) run ./cmd/noted

dev-web:
	cd $(WEB_DIR) && $(NPM) run dev -- --host 127.0.0.1

build-web:
	cd $(WEB_DIR) && $(NPM) ci && $(NPM) run check && $(NPM) run build

build-server:
	cd $(SERVER_DIR) && $(GO) mod download && $(GO) build -o noted ./cmd/noted

fmt:
	cd $(SERVER_DIR) && $(GO) fmt ./...
	cd $(WEB_DIR) && $(NPM) run format:check

lint:
	cd $(SERVER_DIR) && $(GO) vet ./...
	cd $(WEB_DIR) && $(NPM) run lint

audit:
	cd $(SERVER_DIR) && $(GO) list -m all >/dev/null && $(GO) vet ./...
	cd $(WEB_DIR) && $(NPM) audit --audit-level=high

test-server:
	cd $(SERVER_DIR) && $(GO) test ./...

test-web:
	cd $(WEB_DIR) && $(NPM) test

test-contract:
	cd $(SERVER_DIR) && $(GO) test ./internal/api/contracttest/... ./internal/api/handlers/...
	cd $(WEB_DIR) && $(NPM) test -- --run

test-e2e:
	$(NPM) exec playwright test --config=playwright.config.ts tests/e2e

test-offline:
	$(NPM) exec playwright test --config=playwright.config.ts tests/e2e/no-internet.spec.ts

test-import-export:
	$(NPM) exec playwright test --config=playwright.config.ts tests/e2e/us7-import-export.spec.ts

test-performance:
	cd $(SERVER_DIR) && $(GO) test ./internal/store/sqlite -run 'Performance|Query' -count=1

check: fmt lint test-server test-web test-contract

image:
	docker build -f containers/Containerfile -t noted:local .

smoke-image:
	bash tests/smoke/container.sh noted:local

db-status:
	cd $(SERVER_DIR) && $(GO) run ./cmd/noted db-status

db-check:
	cd $(SERVER_DIR) && $(GO) run ./cmd/noted db-check

db-reset:
	cd $(SERVER_DIR) && $(GO) run ./cmd/noted db-reset
