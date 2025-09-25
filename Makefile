# Main Makefile for AcoustiCalc
# Provides convenient access to common development tasks

.PHONY: help build install test lint format clean setup-hooks \
        dev-check ci-check validate install-deps

# Default target
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Configuration
GO := go
PROJECT_ROOT := $(shell pwd)
TESTS_DIR := $(PROJECT_ROOT)/tests
SCRIPTS_DIR := $(TESTS_DIR)/scripts

# Windows compatibility - check if directory exists
ifeq ($(OS),Windows_NT)
    SCRIPTS_EXISTS := $(shell if exist "$(SCRIPTS_DIR)" echo 1)
else
    SCRIPTS_EXISTS := $(shell test -d "$(SCRIPTS_DIR)" && echo 1)
endif

# Build targets
build: ## Build the project
	@echo "🔨 Building project..."
	@$(GO) build ./...

install: ## Install the CLI tool
	@echo "📦 Installing CLI tool..."
	@$(GO) install cmd/acousticalc/main.go

test: ## Run all tests
	@echo "🧪 Running all tests..."
ifeq ($(SCRIPTS_EXISTS),1)
	@$(MAKE) -C $(SCRIPTS_DIR) test-all
else
	@echo "⚠️  Scripts directory not found, running go test directly..."
	@go test -v ./tests/unit/... ./tests/integration/... ./tests/e2e/... ./tests/recording/... ./tests/cross_platform/... ./tests/reporting/...
endif

# Linting and formatting targets
lint: ## Run all linters
	@echo "🔍 Running linters..."
ifeq ($(SCRIPTS_EXISTS),1)
	@$(MAKE) -C $(SCRIPTS_DIR) lint
else
	@echo "⚠️  Scripts directory not found, running linters directly..."
	@go fmt ./...
	@go vet ./...
endif

lint-fix: ## Fix linting issues automatically
	@echo "🔧 Fixing linting issues..."
ifeq ($(SCRIPTS_EXISTS),1)
	@$(MAKE) -C $(SCRIPTS_DIR) lint-fix
else
	@echo "⚠️  Scripts directory not found, running go fmt directly..."
	@go fmt ./...
endif

format: ## Format code
	@echo "📝 Formatting code..."
ifeq ($(SCRIPTS_EXISTS),1)
	@$(MAKE) -C $(SCRIPTS_DIR) format
else
	@echo "⚠️  Scripts directory not found, running go fmt directly..."
	@go fmt ./...
endif

format-check: ## Check if code is formatted
	@echo "🔍 Checking code formatting..."
ifeq ($(SCRIPTS_EXISTS),1)
	@$(MAKE) -C $(SCRIPTS_DIR) format-check
else
	@echo "⚠️  Scripts directory not found, running go fmt directly..."
	@go fmt ./...
endif

vet: ## Run go vet
	@echo "🔍 Running go vet..."
ifeq ($(SCRIPTS_EXISTS),1)
	@$(MAKE) -C $(SCRIPTS_DIR) vet
else
	@echo "⚠️  Scripts directory not found, running go vet directly..."
	@go vet ./...
endif

staticcheck: ## Run staticcheck
	@echo "🔍 Running staticcheck..."
	@$(MAKE) -C $(SCRIPTS_DIR) staticcheck

security-scan: ## Run security scan
	@echo "🔒 Running security scan..."
	@$(MAKE) -C $(SCRIPTS_DIR) security-scan

# Development workflow targets
dev-check: ## Quick development check (lint + unit tests)
	@echo "🚀 Running development checks..."
	@$(MAKE) -C $(SCRIPTS_DIR) validate

ci-check: ## CI-style check (lint + all tests)
	@echo "🤖 Running CI checks..."
	@$(MAKE) -C $(SCRIPTS_DIR) ci-test

validate: ## Quick validation
	@echo "✅ Running validation..."
	@$(MAKE) -C $(SCRIPTS_DIR) validate

# Setup targets
setup-hooks: ## Install git hooks
	@echo "🪝 Installing git hooks..."
	@$(MAKE) -C $(SCRIPTS_DIR) install-hooks

install-deps: ## Install development dependencies
	@echo "📦 Installing development dependencies..."
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Installing staticcheck..."
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "Installing govulncheck..."
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "✅ Development dependencies installed"

# Clean targets
clean: ## Clean build artifacts
	@echo "🧹 Cleaning artifacts..."
	@$(GO) clean
	@$(MAKE) -C $(SCRIPTS_DIR) clean-artifacts

clean-cache: ## Clean Go cache
	@echo "🧹 Cleaning Go cache..."
	@$(GO) clean -cache
	@$(GO) clean -modcache
	@$(GO) clean -testcache

# E2E testing targets
test-e2e: ## Run E2E tests
	@echo "🔗 Running E2E tests..."
	@go test -v -timeout=120s ./tests/e2e/...

test-recording: ## Run recording tests
	@echo "🎥 Running recording tests..."
	@go test -v -timeout=60s ./tests/recording/...

test-cross-platform: ## Run cross-platform tests
	@echo "🌐 Running cross-platform tests..."
	@go test -v -timeout=60s ./tests/cross_platform/...

test-reporting: ## Run reporting tests
	@echo "📊 Running reporting tests..."
	@go test -v -timeout=60s ./tests/reporting/...

test-e2e-all: test-e2e test-recording test-cross-platform test-reporting ## Run all E2E-related tests

# Quick commands
quick: dev-check ## Quick development check
full: ci-check ## Full CI-style check

# Pre-commit hook (called by git)
pre-commit: ## Run pre-commit checks
	@echo "🚀 Running pre-commit checks..."
	@$(MAKE) dev-check
	@echo "✅ Pre-commit checks passed"

# Default to help if no target specified
.DEFAULT_GOAL := help