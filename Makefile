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

# Build targets
build: ## Build the project
	@echo "🔨 Building project..."
	@$(GO) build ./...

install: ## Install the CLI tool
	@echo "📦 Installing CLI tool..."
	@$(GO) install cmd/acousticalc/main.go

test: ## Run all tests
	@echo "🧪 Running all tests..."
	@$(MAKE) -C $(SCRIPTS_DIR) test-all

# Linting and formatting targets
lint: ## Run all linters
	@echo "🔍 Running linters..."
	@$(MAKE) -C $(SCRIPTS_DIR) lint

lint-fix: ## Fix linting issues automatically
	@echo "🔧 Fixing linting issues..."
	@$(MAKE) -C $(SCRIPTS_DIR) lint-fix

format: ## Format code
	@echo "📝 Formatting code..."
	@$(MAKE) -C $(SCRIPTS_DIR) format

format-check: ## Check if code is formatted
	@echo "🔍 Checking code formatting..."
	@$(MAKE) -C $(SCRIPTS_DIR) format-check

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	@$(MAKE) -C $(SCRIPTS_DIR) vet

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