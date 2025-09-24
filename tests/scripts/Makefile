# Unix Makefile for AcoustiCalc Testing Framework
# Provides Unix-specific test operations and optimizations

.PHONY: help test test-unit test-integration test-benchmark test-coverage \
        test-watch test-clean test-validate test-all test-parallel \
        coverage-html coverage-summary benchmark-detailed \
        setup-unix setup-completion unix-report clean-artifacts \
        lint lint-fix format format-check vet staticcheck security-scan \
        pre-commit install-hooks

# Default target
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Configuration
GO := go
PROJECT_ROOT := $(shell cd $(dir $(firstword $(MAKEFILE_LIST)))../../ && pwd)
TESTS_DIR := $(PROJECT_ROOT)/tests
ARTIFACTS_DIR := $(TESTS_DIR)/artifacts
UNIX_SCRIPT := $(TESTS_DIR)/scripts/unix_test_tools.sh

# Unix-specific detection
UNAME_S := $(shell uname -s)
NPROC := $(shell nproc)

# Environment variables
export PARALLEL_JOBS ?= $(NPROC)
export COVERAGE_ENABLED ?= true
export VERBOSE ?= false

# Test targets
test: test-unit test-integration test-benchmark test-coverage ## Run all tests with Unix optimizations

test-unit: ## Run unit tests with Unix optimizations
	@echo "🔍 Running unit tests..."
	@$(UNIX_SCRIPT) run --parallel $(PARALLEL_JOBS) $(if $(filter true,$(VERBOSE)),--verbose,)

test-integration: ## Run integration tests with Unix optimizations
	@echo "🔗 Running integration tests..."
	@cd $(TESTS_DIR)/integration && $(GO) test -v -timeout=60s \
		$(if $(filter true,$(COVERAGE_ENABLED)),-coverprofile=$(ARTIFACTS_DIR)/coverage/integration_coverage.out -covermode=atomic,) \
		-parallel=$(PARALLEL_JOBS) ./...

test-benchmark: ## Run performance benchmarks
	@echo "⚡ Running performance benchmarks..."
	@cd $(TESTS_DIR)/unit && $(GO) test -bench=. -benchmem -timeout=120s \
		-parallel=$(PARALLEL_JOBS) ./... > $(ARTIFACTS_DIR/reports/benchmark_results.txt 2>&1

test-coverage: coverage-html coverage-summary ## Generate coverage reports

test-watch: ## Watch file changes and run tests automatically
	@echo "👀 Starting file watcher..."
	@$(UNIX_SCRIPT) watch

test-clean: clean-artifacts ## Clean test artifacts

test-validate: ## Validate Unix environment
	@echo "✅ Validating Unix environment..."
	@$(UNIX_SCRIPT) validate

test-all: test-unit test-integration test-benchmark unix-report ## Run complete test suite with reports

test-parallel: ## Run tests with maximum parallelization
	@echo "🚀 Running tests with maximum parallelization..."
	@PARALLEL_JOBS=$(NPROC) $(MAKE) test-all

# Coverage targets
coverage-html: ## Generate HTML coverage report
	@echo "📊 Generating HTML coverage report..."
	@mkdir -p $(ARTIFACTS_DIR)/coverage
	@if [ -f $(ARTIFACTS_DIR)/coverage/unit_coverage.out ]; then \
		$(GO) tool cover -html=$(ARTIFACTS_DIR)/coverage/unit_coverage.out \
			-o $(ARTIFACTS_DIR/coverage/coverage.html; \
		echo "✅ HTML coverage report: $(ARTIFACTS_DIR)/coverage/coverage.html"; \
	else \
		echo "⚠️  No coverage file found, run tests first"; \
	fi

coverage-summary: ## Generate coverage summary
	@echo "📈 Generating coverage summary..."
	@if [ -f $(ARTIFACTS_DIR)/coverage/unit_coverage.out ]; then \
		$(GO) tool cover -func=$(ARTIFACTS_DIR)/coverage/unit_coverage.out \
			> $(ARTIFACTS_DIR/coverage/coverage_summary.txt; \
		echo "✅ Coverage summary: $(ARTIFACTS_DIR)/coverage/coverage_summary.txt"; \
		grep "total:" $(ARTIFACTS_DIR)/coverage/coverage_summary.txt; \
	else \
		echo "⚠️  No coverage file found, run tests first"; \
	fi

# Benchmark targets
benchmark-detailed: ## Run detailed benchmarks with analysis
	@echo "🔬 Running detailed benchmarks..."
	@cd $(TESTS_DIR)/unit && $(GO) test -bench=. -benchmem -benchtime=5s \
		-count=3 -timeout=300s ./... > $(ARTIFACTS_DIR/reports/detailed_benchmarks.txt 2>&1
	@echo "✅ Detailed benchmarks: $(ARTIFACTS_DIR/reports/detailed_benchmarks.txt"

# Setup targets
setup-unix: ## Setup Unix-specific testing environment
	@echo "🔧 Setting up Unix testing environment..."
	@mkdir -p $(ARTIFACTS_DIR)/{coverage,reports,platform_results}
	@chmod +x $(UNIX_SCRIPT)
	@echo "✅ Unix testing environment setup complete"

setup-completion: ## Setup shell completion
	@echo "🔧 Setting up shell completion..."
	@$(UNIX_SCRIPT) completion
	@echo "✅ Shell completion setup complete"

# Report targets
unix-report: ## Generate Unix-specific test report
	@echo "📋 Generating Unix test report..."
	@$(UNIX_SCRIPT) validate && echo "✅ Unix environment validated"
	@echo "System Information:" > $(ARTIFACTS_DIR/reports/unix_system_info.txt
	@uname -a >> $(ARTIFACTS_DIR/reports/unix_system_info.txt
	@echo "Go Version:" >> $(ARTIFACTS_DIR/reports/unix_system_info.txt
	@$(GO) version >> $(ARTIFACTS_DIR/reports/unix_system_info.txt
	@echo "Parallel Jobs: $(PARALLEL_JOBS)" >> $(ARTIFACTS_DIR/reports/unix_system_info.txt
	@echo "✅ Unix system report: $(ARTIFACTS_DIR/reports/unix_system_info.txt"

# Performance monitoring
monitor-performance: ## Monitor system performance during tests
	@echo "📊 Monitoring performance..."
	@if command -v top >/dev/null 2>&1; then \
		echo "CPU Usage:" && top -l 1 -n 0 | grep "CPU usage" || true; \
	fi
	@if command -v vm_stat >/dev/null 2>&1; then \
		echo "Memory Info:" && vm_stat | head -5 || true; \
	fi

# Clean targets
clean-artifacts: ## Clean all test artifacts
	@echo "🧹 Cleaning test artifacts..."
	@rm -rf $(ARTIFACTS_DIR)
	@echo "✅ Test artifacts cleaned"

clean-coverage: ## Clean only coverage artifacts
	@echo "🧹 Cleaning coverage artifacts..."
	@rm -rf $(ARTIFACTS_DIR)/coverage
	@echo "✅ Coverage artifacts cleaned"

# Unix-specific optimizations
unix-optimize: ## Apply Unix-specific optimizations
	@echo "⚡ Applying Unix-specific optimizations..."
	@# Set optimal file descriptor limit
	@ulimit -n 65536 || true
	@# Set optimal process priority
	@renice -n 10 $$ >/dev/null 2>&1 || true
	@echo "✅ Unix optimizations applied"

# Development workflow targets
dev-test: test-unit coverage-html ## Quick development test cycle
	@echo "🚀 Development test cycle complete"

dev-full: test-all unix-report ## Full development test cycle with reports
	@echo "🚀 Full development test cycle complete"

# CI-specific targets
ci-test: test-unit test-integration test-coverage ## CI-optimized test run
	@echo "🤖 CI test run complete"

ci-benchmark: test-benchmark benchmark-detailed ## CI-optimized benchmark run
	@echo "🤖 CI benchmark run complete"

# Platform-specific targets
linux-test: ## Linux-specific test optimizations
	@if [ "$(UNAME_S)" = "Linux" ]; then \
		echo "🐧 Applying Linux-specific optimizations..."; \
		export GOMAXPROCS=$(NPROC); \
		$(MAKE) test-all; \
	else \
		echo "⚠️  Not running on Linux, skipping Linux-specific optimizations"; \
	fi

darwin-test: ## macOS-specific test optimizations
	@if [ "$(UNAME_S)" = "Darwin" ]; then \
		echo "🍎 Applying macOS-specific optimizations..."; \
		export GOMAXPROCS=$(shell sysctl -n hw.ncpu); \
		$(MAKE) test-all; \
	else \
		echo "⚠️  Not running on macOS, skipping macOS-specific optimizations"; \
	fi

# Utility targets
show-config: ## Show current configuration
	@echo "⚙️  Current Configuration:"
	@echo "  Go Version: $$($(GO) version)"
	@echo "  System: $(UNAME_S)"
	@echo "  CPU Cores: $(NPROC)"
	@echo "  Parallel Jobs: $(PARALLEL_JOBS)"
	@echo "  Coverage Enabled: $(COVERAGE_ENABLED)"
	@echo "  Verbose Output: $(VERBOSE)"
	@echo "  Project Root: $(PROJECT_ROOT)"
	@echo "  Tests Directory: $(TESTS_DIR)"
	@echo "  Artifacts Directory: $(ARTIFACTS_DIR)"

check-deps: ## Check Unix dependencies
	@echo "🔍 Checking Unix dependencies..."
	@command -v go >/dev/null 2>&1 || { echo "❌ Go not found"; exit 1; }
	@command -v bc >/dev/null 2>&1 || { echo "❌ bc not found"; exit 1; }
	@command -v find >/dev/null 2>&1 || { echo "❌ find not found"; exit 1; }
	@command -v xargs >/dev/null 2>&1 || { echo "❌ xargs not found"; exit 1; }
	@command -v mktemp >/dev/null 2>&1 || { echo "❌ mktemp not found"; exit 1; }
	@command -v timeout >/dev/null 2>&1 || { echo "❌ timeout not found"; exit 1; }
	@echo "✅ All Unix dependencies found"

# Debug targets
debug-test: ## Run tests with debug output
	@echo "🐛 Running tests with debug output..."
	@VERBOSE=true $(MAKE) test-all

debug-coverage: ## Debug coverage generation
	@echo "🐛 Debugging coverage generation..."
	@ls -la $(ARTIFACTS_DIR)/coverage/ 2>/dev/null || echo "No coverage directory"
	@find $(ARTIFACTS_DIR) -name "*coverage*" 2>/dev/null || echo "No coverage files found"

# Quick targets for common operations
quick: test-unit ## Quick unit test run
quick-coverage: test-unit coverage-html ## Quick test with coverage
quick-bench: test-benchmark ## Quick benchmark run

# Default parallel jobs for different systems
ifeq ($(UNAME_S),Linux)
export PARALLEL_JOBS ?= $(shell nproc)
else ifeq ($(UNAME_S),Darwin)
export PARALLEL_JOBS ?= $(shell sysctl -n hw.ncpu)
else
export PARALLEL_JOBS ?= 2
endif

# Linting and formatting targets
lint: ## Run all linters (golint, go vet, staticcheck)
	@echo "🔍 Running all linters..."
	@$(MAKE) format-check
	@$(MAKE) vet
	@$(MAKE) staticcheck
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Running golangci-lint..."; \
		cd $(PROJECT_ROOT) && golangci-lint run --timeout=5m; \
	else \
		echo "⚠️  golangci-lint not found, skipping"; \
	fi
	@echo "✅ All linters passed"

lint-fix: ## Fix linting issues automatically
	@echo "🔧 Fixing linting issues..."
	@$(MAKE) format
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Running golangci-lint with fixes..."; \
		cd $(PROJECT_ROOT) && golangci-lint run --fix --timeout=5m; \
	fi
	@echo "✅ Linting issues fixed"

format: ## Format Go code using go fmt and goimports
	@echo "📝 Formatting Go code..."
	@cd $(PROJECT_ROOT) && go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		echo "Running goimports..."; \
		cd $(PROJECT_ROOT) && goimports -w .; \
	fi
	@echo "✅ Code formatted"

format-check: ## Check if Go code is properly formatted
	@echo "🔍 Checking code formatting..."
	@cd $(PROJECT_ROOT) && \
		files=$$(go fmt ./...); \
		if [ -n "$$files" ]; then \
			echo "❌ Files are not formatted:"; \
			echo "$$files"; \
			exit 1; \
		fi
	@echo "✅ Code is properly formatted"

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	@cd $(PROJECT_ROOT) && go vet ./...
	@echo "✅ go vet passed"

staticcheck: ## Run staticcheck if available
	@echo "🔍 Running staticcheck..."
	@if command -v staticcheck >/dev/null 2>&1; then \
		cd $(PROJECT_ROOT) && staticcheck ./...; \
		echo "✅ staticcheck passed"; \
	else \
		echo "⚠️  staticcheck not found, skipping"; \
	fi

security-scan: ## Run security scanning tools
	@echo "🔒 Running security scan..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		cd $(PROJECT_ROOT) && govulncheck ./...; \
		echo "✅ Security scan completed"; \
	else \
		echo "⚠️  govulncheck not found, skipping"; \
	fi

pre-commit: ## Run all pre-commit checks
	@echo "🚀 Running pre-commit checks..."
	@$(MAKE) lint
	@$(MAKE) test-unit
	@echo "✅ All pre-commit checks passed"

install-hooks: ## Install git hooks
	@echo "🪝 Installing git hooks..."
	@mkdir -p $(PROJECT_ROOT)/.git/hooks
	@echo '#!/bin/bash' > $(PROJECT_ROOT)/.git/hooks/pre-commit
	@echo 'echo "Running pre-commit checks..."' >> $(PROJECT_ROOT)/.git/hooks/pre-commit
	@echo 'cd $$(git rev-parse --show-toplevel)' >> $(PROJECT_ROOT)/.git/hooks/pre-commit
	@echo 'make -C tests/scripts pre-commit' >> $(PROJECT_ROOT)/.git/hooks/pre-commit
	@chmod +x $(PROJECT_ROOT)/.git/hooks/pre-commit
	@echo "✅ Git hooks installed"

# Development workflow targets with linting
dev-test: lint test-unit coverage-html ## Quick development test cycle with linting
	@echo "🚀 Development test cycle complete"

dev-full: lint test-all unix-report ## Full development test cycle with linting and reports
	@echo "🚀 Full development test cycle complete"

# CI-specific targets with linting
ci-test: lint test-unit test-integration test-coverage ## CI-optimized test run with linting
	@echo "🤖 CI test run complete"

# Quick validation targets
validate: lint test-unit ## Quick validation (lint + unit tests)
	@echo "✅ Quick validation complete"

# Include environment-specific Makefile if it exists
-include $(PROJECT_ROOT)/Makefile.local