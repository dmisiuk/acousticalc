#!/bin/bash

# Unix Test Tools for AcoustiCalc Testing Framework
# Provides Unix-specific testing utilities and optimizations

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test execution modes
MODE="${1:-run}"
PARALLEL_JOBS="${PARALLEL_JOBS:-$(nproc)}"
COVERAGE_ENABLED="${COVERAGE_ENABLED:-true}"
VERBOSE="${VERBOSE:-false}"

# Project directories
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
TESTS_DIR="$PROJECT_ROOT/tests"
ARTIFACTS_DIR="$TESTS_DIR/artifacts"

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*"
}

# Performance monitoring
start_time=$(date +%s.%N)

elapsed_time() {
    end_time=$(date +%s.%N)
    echo "$end_time - $start_time" | bc
}

# Unix environment validation
validate_unix_environment() {
    log_info "Validating Unix environment..."

    # Check if running on Unix-like system
    if [[ "$OSTYPE" == "linux-gnu"* ]] || [[ "$OSTYPE" == "darwin"* ]]; then
        log_success "Unix-like environment detected: $OSTYPE"
    else
        log_error "This script is designed for Unix-like systems only"
        exit 1
    fi

    # Check required Unix tools
    local required_tools=("go" "bc" "find" "xargs" "mktemp" "timeout")
    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            log_error "Required tool '$tool' not found"
            exit 1
        fi
    done

    log_success "Unix environment validation completed"
}

# Test runner functions
run_unit_tests() {
    log_info "Running unit tests with Unix optimizations..."

    # Change to project root for consistent paths
    cd "$PROJECT_ROOT"

    local test_cmd="go test -v -timeout=30s"

    if [[ "$COVERAGE_ENABLED" == "true" ]]; then
        test_cmd="$test_cmd -coverprofile=$ARTIFACTS_DIR/coverage/unit_coverage.out"
        test_cmd="$test_cmd -covermode=atomic"
    fi

    if [[ "$PARALLEL_JOBS" -gt 1 ]]; then
        test_cmd="$test_cmd -parallel=$PARALLEL_JOBS"
    fi

    test_cmd="$test_cmd ./tests/unit/..."

    if [[ "$VERBOSE" == "true" ]]; then
        test_cmd="$test_cmd -v"
    fi

    log_info "Executing: $test_cmd"

    if eval "$test_cmd"; then
        log_success "Unit tests completed successfully"
        return 0
    else
        log_error "Unit tests failed"
        return 1
    fi
}

run_integration_tests() {
    log_info "Running integration tests with Unix optimizations..."

    # Change to project root for consistent paths
    cd "$PROJECT_ROOT"

    local test_cmd="go test -v -timeout=60s"

    if [[ "$COVERAGE_ENABLED" == "true" ]]; then
        test_cmd="$test_cmd -coverprofile=$ARTIFACTS_DIR/coverage/integration_coverage.out"
        test_cmd="$test_cmd -covermode=atomic"
    fi

    if [[ "$PARALLEL_JOBS" -gt 1 ]]; then
        test_cmd="$test_cmd -parallel=$PARALLEL_JOBS"
    fi

    test_cmd="$test_cmd ./tests/integration/..."

    if [[ "$VERBOSE" == "true" ]]; then
        test_cmd="$test_cmd -v"
    fi

    log_info "Executing: $test_cmd"

    if eval "$test_cmd"; then
        log_success "Integration tests completed successfully"
        return 0
    else
        log_error "Integration tests failed"
        return 1
    fi
}

run_benchmarks() {
    log_info "Running performance benchmarks..."

    # Change to project root for consistent paths
    cd "$PROJECT_ROOT"

    local benchmark_cmd="go test -bench=. -benchmem -timeout=120s"

    if [[ "$PARALLEL_JOBS" -gt 1 ]]; then
        benchmark_cmd="$benchmark_cmd -parallel=$PARALLEL_JOBS"
    fi

    benchmark_cmd="$benchmark_cmd ./tests/unit/..."

    log_info "Executing: $benchmark_cmd"

    if eval "$benchmark_cmd" > "$ARTIFACTS_DIR/reports/benchmark_results.txt" 2>&1; then
        log_success "Benchmarks completed successfully"
        return 0
    else
        log_error "Benchmarks failed"
        return 1
    fi
}

generate_coverage_report() {
    if [[ "$COVERAGE_ENABLED" != "true" ]]; then
        log_info "Coverage generation disabled, skipping..."
        return 0
    fi

    log_info "Generating coverage reports..."

    # Create coverage directory
    mkdir -p "$ARTIFACTS_DIR/coverage"

    # Combine coverage profiles if multiple exist
    local coverage_files=("$ARTIFACTS_DIR/coverage/"*"_coverage.out")
    if [[ ${#coverage_files[@]} -gt 1 ]]; then
        log_info "Combining multiple coverage profiles..."

        # Create combined coverage profile
        echo "mode: atomic" > "$ARTIFACTS_DIR/coverage/combined_coverage.out"
        for file in "${coverage_files[@]}"; do
            if [[ -f "$file" && "$file" != *"combined_coverage.out" ]]; then
                tail -n +2 "$file" >> "$ARTIFACTS_DIR/coverage/combined_coverage.out"
            fi
        done
    fi

    # Generate HTML coverage report
    local coverage_file="$ARTIFACTS_DIR/coverage/combined_coverage.out"
    if [[ ! -f "$coverage_file" ]]; then
        # Use unit coverage if combined doesn't exist
        coverage_file="$ARTIFACTS_DIR/coverage/unit_coverage.out"
    fi

    if [[ -f "$coverage_file" ]]; then
        log_info "Generating HTML coverage report..."
        go tool cover -html="$coverage_file" -o "$ARTIFACTS_DIR/coverage/coverage.html"
        log_success "HTML coverage report generated: $ARTIFACTS_DIR/coverage/coverage.html"

        # Generate coverage summary
        log_info "Generating coverage summary..."
        go tool cover -func="$coverage_file" > "$ARTIFACTS_DIR/coverage/coverage_summary.txt"
        log_success "Coverage summary generated: $ARTIFACTS_DIR/coverage/coverage_summary.txt"

        # Create coverage record for historical analysis
        log_info "Creating coverage record for historical analysis..."
        local timestamp=$(date +%Y%m%d_%H%M%S)
        if [[ -x "$SCRIPT_DIR/coverage.sh" ]]; then
            "$SCRIPT_DIR/coverage.sh" record "$coverage_file" "$timestamp"
        else
            log_warning "Coverage trends script not found or not executable"
        fi
    else
        log_warning "No coverage files found, skipping coverage report generation"
    fi
}

# File watching for development
watch_tests() {
    log_info "Starting file watcher for test development..."

    # Use inotifywait on Linux, fswatch on macOS
    local watch_cmd=""
    if command -v inotifywait &> /dev/null; then
        # Linux
        watch_cmd="inotifywait -r -e modify,create,delete --include '\\.(go)$' ./pkg ./tests"
    elif command -v fswatch &> /dev/null; then
        # macOS
        watch_cmd="fswatch -o -e '.*' --include '\\.go$' ./pkg ./tests"
    else
        log_error "No file watcher found. Please install inotify-tools (Linux) or fswatch (macOS)"
        exit 1
    fi

    log_info "Watching for Go file changes..."
    log_info "Press Ctrl+C to stop watching"

    while true; do
        if eval "$watch_cmd"; then
            log_info "Files changed, running tests..."
            run_all_tests
            echo "----------------------------------------"
        fi
    done
}

# Run all tests
run_all_tests() {
    log_info "Running all tests with Unix optimizations..."

    local start_time=$(date +%s)
    local failed=0

    # Create artifacts directories
    mkdir -p "$ARTIFACTS_DIR"/{coverage,reports,platform_results}

    # Run test suites
    run_unit_tests || ((failed++))
    run_integration_tests || ((failed++))
    run_benchmarks || ((failed++))

    # Generate reports
    generate_coverage_report

    # Run coverage trends analysis if coverage script exists
    if [[ -x "$SCRIPT_DIR/coverage.sh" ]]; then
        log_info "Running coverage trends analysis..."
        "$SCRIPT_DIR/coverage.sh" all
    else
        log_warning "Coverage trends script not found or not executable"
    fi

    local end_time=$(date +%s)
    local duration=$((end_time - start_time))

    if [[ $failed -eq 0 ]]; then
        log_success "All tests completed successfully in ${duration}s"
        return 0
    else
        log_error "$failed test suite(s) failed in ${duration}s"
        return 1
    fi
}

# Unix-specific performance monitoring
monitor_performance() {
    log_info "Monitoring test performance..."

    local cpu_usage="N/A"
    local memory_usage="N/A"
    local memory_usage_label="Memory Free"

    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS-specific commands
        cpu_usage=$(top -l 1 -n 0 | grep "CPU usage" | awk '{print $3}' | sed 's/%//')
        memory_usage=$(vm_stat | grep "Pages free" | awk '{print $3}' | sed 's/\.$//')
        memory_usage_label="Memory Free Pages"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux-specific commands
        cpu_usage=$(top -bn1 | grep '%Cpu(s)' | awk '{print $2+$4}')
        memory_usage=$(free -m | awk '/^Mem:/{print $4}')
        memory_usage_label="Memory Free (MB)"
    fi

    log_info "CPU Usage: ${cpu_usage}%"
    log_info "${memory_usage_label}: ${memory_usage}"

    # Save performance metrics
    {
        echo "Performance Metrics - $(date)"
        echo "CPU Usage: ${cpu_usage}%"
        echo "${memory_usage_label}: ${memory_usage}"
        echo "Parallel Jobs: $PARALLEL_JOBS"
        echo "Coverage Enabled: $COVERAGE_ENABLED"
    } >> "$ARTIFACTS_DIR/reports/performance_metrics.txt"
}

# Generate Unix-specific test report
generate_unix_report() {
    log_info "Generating Unix-specific test report..."

    local report_file="$ARTIFACTS_DIR/reports/unix_test_report.txt"

    {
        echo "Unix Test Framework Report"
        echo "========================="
        echo "Generated: $(date)"
        echo "System: $(uname -a)"
        echo "Go Version: $(go version)"
        echo ""
        echo "Configuration:"
        echo "  Parallel Jobs: $PARALLEL_JOBS"
        echo "  Coverage Enabled: $COVERAGE_ENABLED"
        echo "  Verbose Output: $VERBOSE"
        echo ""
        echo "Environment:"
        echo "  OS: $OSTYPE"
        echo "  Shell: $SHELL"
        echo "  CPU Cores: $(nproc)"
        echo ""

        # Include coverage summary if available
        if [[ -f "$ARTIFACTS_DIR/coverage/coverage_summary.txt" ]]; then
            echo "Coverage Summary:"
            echo "----------------"
            cat "$ARTIFACTS_DIR/coverage/coverage_summary.txt"
        fi

    } > "$report_file"

    log_success "Unix test report generated: $report_file"
}

# Setup shell completion
setup_completion() {
    log_info "Setting up shell completion..."

    local completion_script="$SCRIPT_DIR/unix_test_completion.sh"

    cat > "$completion_script" << 'EOF'
#!/bin/bash
# Shell completion for unix_test_tools.sh

_unix_test_tools_completion() {
    local cur prev commands
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    commands="run watch benchmark coverage validate clean help"

    if [[ ${COMP_CWORD} -eq 1 ]] ; then
        COMPREPLY=( $(compgen -W "${commands}" -- "${cur}") )
        return 0
    fi

    case "${prev}" in
        run|watch|benchmark|coverage|validate|clean)
            COMPREPLY=( $(compgen -W "--verbose --parallel --no-coverage" -- "${cur}") )
            return 0
            ;;
        *)
            ;;
    esac
}

complete -F _unix_test_tools_completion unix_test_tools.sh
EOF

    chmod +x "$completion_script"

    log_info "To enable completion, run: source $completion_script"
    log_success "Shell completion setup completed"
}

# Clean test artifacts
clean_artifacts() {
    log_info "Cleaning test artifacts..."

    if [[ -d "$ARTIFACTS_DIR" ]]; then
        rm -rf "$ARTIFACTS_DIR"
        log_success "Test artifacts cleaned"
    else
        log_info "No artifacts to clean"
    fi
}

# Help function
show_help() {
    cat << EOF
Unix Test Tools for AcoustiCalc

USAGE:
    unix_test_tools.sh [MODE] [OPTIONS]

MODES:
    run         Run all tests (default)
    watch       Watch file changes and run tests automatically
    benchmark   Run performance benchmarks only
    coverage    Generate coverage reports only
    validate    Validate Unix environment
    clean       Clean test artifacts
    help        Show this help message
    completion  Setup shell completion

OPTIONS:
    --verbose           Enable verbose output
    --parallel N        Set number of parallel jobs (default: nproc)
    --no-coverage       Disable coverage generation

ENVIRONMENT VARIABLES:
    PARALLEL_JOBS       Number of parallel jobs (default: nproc)
    COVERAGE_ENABLED    Enable coverage reports (default: true)
    VERBOSE             Enable verbose output (default: false)

EXAMPLES:
    unix_test_tools.sh run --verbose
    unix_test_tools.sh watch --parallel 4
    unix_test_tools.sh benchmark --no-coverage
    PARALLEL_JOBS=8 unix_test_tools.sh run

EOF
}

# Main execution
main() {
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --verbose)
                VERBOSE=true
                shift
                ;;
            --parallel)
                PARALLEL_JOBS="$2"
                shift 2
                ;;
            --no-coverage)
                COVERAGE_ENABLED=false
                shift
                ;;
            *)
                MODE="$1"
                shift
                ;;
        esac
    done

    # Create artifacts directory
    mkdir -p "$ARTIFACTS_DIR"/{coverage,reports,platform_results}

    # Execute based on mode
    case "$MODE" in
        run)
            validate_unix_environment
            monitor_performance
            run_all_tests
            generate_unix_report
            ;;
        watch)
            validate_unix_environment
            watch_tests
            ;;
        benchmark)
            validate_unix_environment
            run_benchmarks
            ;;
        coverage)
            generate_coverage_report
            ;;
        validate)
            validate_unix_environment
            ;;
        clean)
            clean_artifacts
            ;;
        completion)
            setup_completion
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "Unknown mode: $MODE"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"