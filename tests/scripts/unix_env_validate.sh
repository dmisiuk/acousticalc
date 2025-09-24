#!/bin/bash

# Unix Environment Validation Script for AcoustiCalc Testing Framework
# Validates and optimizes Unix environment for testing

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
MIN_GO_VERSION="1.21"
MIN_PARALLEL_JOBS=2
OPTIMAL_PARALLEL_JOBS=$(nproc)
REQUIRED_DISK_SPACE=100 # MB

# Logging
log_info() { echo -e "${BLUE}[INFO]${NC} $*"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $*"; }
log_warning() { echo -e "${YELLOW}[WARNING]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

# Main validation function
main() {
    log_info "Starting Unix environment validation..."
    echo "=============================================="

    local failed=0
    local warnings=0

    # System validation
    validate_system || ((failed++))
    validate_go_environment || ((failed++))
    validate_dependencies || ((failed++))
    validate_performance_settings || ((warnings++))
    validate_file_system || ((warnings++))
    validate_network_connectivity || ((warnings++))
    validate_security_settings || ((warnings++))
    validate_development_tools || ((warnings++))

    echo "=============================================="

    # Summary
    if [[ $failed -gt 0 ]]; then
        log_error "Validation failed with $failed error(s) and $warnings warning(s)"
        exit 1
    else
        log_success "Validation passed with $warnings warning(s)"
        echo ""
        print_recommendations
        exit 0
    fi
}

# System validation
validate_system() {
    log_info "Validating system environment..."

    # Check Unix-like system
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        OS="Linux"
        log_success "Linux detected"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macOS"
        log_success "macOS detected"
    else
        log_error "Unsupported operating system: $OSTYPE"
        return 1
    fi

    # Check architecture
    ARCH=$(uname -m)
    if [[ "$ARCH" == "x86_64" || "$ARCH" == "arm64" ]]; then
        log_success "Architecture: $ARCH"
    else
        log_warning "Uncommon architecture: $ARCH"
    fi

    # Check available memory
    if [[ "$OS" == "Linux" ]]; then
        TOTAL_MEM=$(free -m | awk 'NR==2{printf "%.0f", $2}')
        AVAILABLE_MEM=$(free -m | awk 'NR==2{printf "%.0f", $7}')
    elif [[ "$OS" == "macOS" ]]; then
        TOTAL_MEM=$(sysctl -n hw.memsize | awk '{printf "%.0f", $1/1024/1024}')
        AVAILABLE_MEM=$(vm_stat | grep "Pages free" | awk '{printf "%.0f", $3 * 4 / 1024}')
    fi

    if [[ ${TOTAL_MEM:-0} -lt 1024 ]]; then
        log_warning "Low memory: ${TOTAL_MEM}MB total, ${AVAILABLE_MEM}MB available"
    else
        log_success "Memory: ${TOTAL_MEM}MB total, ${AVAILABLE_MEM}MB available"
    fi

    # Check CPU cores
    CPU_CORES=$(nproc)
    if [[ $CPU_CORES -lt $MIN_PARALLEL_JOBS ]]; then
        log_warning "Low CPU cores: $CPU_CORES (minimum recommended: $MIN_PARALLEL_JOBS)"
    else
        log_success "CPU cores: $CPU_CORES"
    fi

    # Check disk space
    AVAILABLE_DISK=$(df -m "$PROJECT_ROOT" | awk 'NR==2{print $4}')
    if [[ ${AVAILABLE_DISK:-0} -lt $REQUIRED_DISK_SPACE ]]; then
        log_warning "Low disk space: ${AVAILABLE_DISK}MB available (minimum recommended: ${REQUIRED_DISK_SPACE}MB)"
    else
        log_success "Disk space: ${AVAILABLE_DISK}MB available"
    fi

    return 0
}

# Go environment validation
validate_go_environment() {
    log_info "Validating Go environment..."

    # Check Go installation
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed"
        return 1
    fi

    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if ! validate_version "$GO_VERSION" "$MIN_GO_VERSION"; then
        log_error "Go version $GO_VERSION is below minimum required $MIN_GO_VERSION"
        return 1
    fi
    log_success "Go version: $GO_VERSION"

    # Check Go environment variables
    check_go_env_vars

    # Check Go modules
    if [[ ! -f "$PROJECT_ROOT/go.mod" ]]; then
        log_error "go.mod not found in project root"
        return 1
    fi
    log_success "Go modules detected"

    # Check if we can build the project
    if ! go build ./... >/dev/null 2>&1; then
        log_error "Failed to build project"
        return 1
    fi
    log_success "Project builds successfully"

    return 0
}

# Dependency validation
validate_dependencies() {
    log_info "Validating dependencies..."

    local required_deps=("git" "make" "bc" "find" "xargs" "mktemp" "timeout")
    local missing_deps=()

    for dep in "${required_deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            missing_deps+=("$dep")
        else
            log_success "$dep: $(command -v "$dep")"
        fi
    done

    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        log_error "Missing dependencies: ${missing_deps[*]}"
        return 1
    fi

    # Check optional but recommended dependencies
    local optional_deps=("inotifywait" "fswatch" "entr" "watchexec")
    local missing_optional=()

    for dep in "${optional_deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            missing_optional+=("$dep")
        else
            log_success "$dep: $(command -v "$dep") (optional)"
        fi
    done

    if [[ ${#missing_optional[@]} -gt 0 ]]; then
        log_warning "Optional dependencies not found: ${missing_optional[*]}"
        log_info "These are recommended for better development experience"
    fi

    return 0
}

# Performance settings validation
validate_performance_settings() {
    log_info "Validating performance settings..."

    # Check file descriptor limits
    local fd_limit=$(ulimit -n)
    if [[ $fd_limit -lt 1024 ]]; then
        log_warning "File descriptor limit is low: $fd_limit (recommended: 1024+)"
    else
        log_success "File descriptor limit: $fd_limit"
    fi

    # Check process limits
    local process_limit=$(ulimit -u)
    if [[ $process_limit -lt 1024 ]]; then
        log_warning "Process limit is low: $process_limit (recommended: 1024+)"
    else
        log_success "Process limit: $process_limit"
    fi

    # Check GOMAXPROCS
    local gomaxprocs=${GOMAXPROCS:-}
    if [[ -z "$gomaxprocs" ]]; then
        log_info "GOMAXPROCS not set, will use default"
    else
        log_success "GOMAXPROCS: $gomaxprocs"
    fi

    # Check swap settings (Linux only)
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        local swap_usage=$(free -m | awk 'NR==3{printf "%.1f", $3/$2*100}')
        if (( $(echo "$swap_usage > 50" | bc -l) )); then
            log_warning "High swap usage: ${swap_usage}%"
        else
            log_success "Swap usage: ${swap_usage}%"
        fi
    fi

    return 0
}

# File system validation
validate_file_system() {
    log_info "Validating file system..."

    # Check file system type
    local fs_type=$(df -T "$PROJECT_ROOT" | awk 'NR==2{print $2}')
    log_info "File system type: $fs_type"

    # Check permissions
    if [[ ! -r "$PROJECT_ROOT" || ! -w "$PROJECT_ROOT" ]]; then
        log_error "Insufficient permissions on project directory"
        return 1
    fi
    log_success "Project directory permissions OK"

    # Check for case sensitivity (important for Go)
    local test_file="$PROJECT_ROOT/CaSeSeNsItIvItYtEsT.tmp"
    touch "$test_file"
    if [[ ! -f "$test_file" ]]; then
        log_warning "File system may be case-insensitive (not recommended for Go development)"
    else
        log_success "File system appears to be case-sensitive"
    fi
    rm -f "$test_file"

    # Check available inodes
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        local inode_usage=$(df -i "$PROJECT_ROOT" | awk 'NR==2{printf "%.1f", $3/$2*100}')
        if (( $(echo "$inode_usage > 90" | bc -l) )); then
            log_warning "High inode usage: ${inode_usage}%"
        else
            log_success "Inode usage: ${inode_usage}%"
        fi
    fi

    return 0
}

# Network connectivity validation
validate_network_connectivity() {
    log_info "Validating network connectivity..."

    # Check basic internet connectivity
    if ping -c 1 -W 5 8.8.8.8 >/dev/null 2>&1; then
        log_success "Internet connectivity OK"
    else
        log_warning "No internet connectivity (may affect module downloads)"
    fi

    # Check GitHub connectivity
    if curl -s --connect-timeout 5 https://github.com >/dev/null 2>&1; then
        log_success "GitHub connectivity OK"
    else
        log_warning "Cannot reach GitHub (may affect CI operations)"
    fi

    # Check proxy settings
    if [[ -n "${HTTP_PROXY:-}" || -n "${HTTPS_PROXY:-}" ]]; then
        log_info "Proxy settings detected"
        log_info "HTTP_PROXY: ${HTTP_PROXY:-not set}"
        log_info "HTTPS_PROXY: ${HTTPS_PROXY:-not set}"
    fi

    return 0
}

# Security settings validation
validate_security_settings() {
    log_info "Validating security settings..."

    # Check if running as root
    if [[ $EUID -eq 0 ]]; then
        log_warning "Running as root (not recommended for development)"
    else
        log_success "Running as normal user"
    fi

    # Check Go module proxy settings
    local goproxy=${GOPROXY:-}
    if [[ -z "$goproxy" ]]; then
        log_info "GOPROXY not set (using default)"
    else
        log_success "GOPROXY: $goproxy"
    fi

    # Check GONOPROXY and GONOSUMDB
    if [[ -n "${GONOPROXY:-}" || -n "${GONOSUMDB:-}" ]]; then
        log_info "Private module settings detected"
    fi

    return 0
}

# Development tools validation
validate_development_tools() {
    log_info "Validating development tools..."

    # Check for IDE integration
    if [[ -n "${VSCODE_PID:-}" ]]; then
        log_success "VS Code detected"
    elif [[ -n "${TMUX:-}" ]]; then
        log_success "tmux session detected"
    elif [[ -n "${STY:-}" ]]; then
        log_success "screen session detected"
    fi

    # Check for shell enhancements
    if [[ -n "${ZSH_VERSION:-}" ]]; then
        log_success "Zsh shell detected"
        # Check for Oh My Zsh
        if [[ -n "${ZSH:-}" ]]; then
            log_success "Oh My Zsh detected"
        fi
    elif [[ -n "${BASH_VERSION:-}" ]]; then
        log_success "Bash shell detected"
    fi

    # Check for Go development tools
    local go_tools=("gofmt" "golint" "staticcheck")
    for tool in "${go_tools[@]}"; do
        if command -v "$tool" &> /dev/null; then
            log_success "$tool: $(command -v "$tool")"
        else
            log_info "$tool: not installed (optional)"
        fi
    done

    return 0
}

# Helper functions
validate_version() {
    local current=$1
    local minimum=$2

    # Simple version comparison (works for semantic versions)
    if [[ "$current" == "$minimum" ]]; then
        return 0
    fi

    local current_parts=(${current//./ })
    local minimum_parts=(${minimum//./ })

    for i in "${!minimum_parts[@]}"; do
        if [[ ${current_parts[i]:-0} -lt ${minimum_parts[i]} ]]; then
            return 1
        elif [[ ${current_parts[i]:-0} -gt ${minimum_parts[i]} ]]; then
            return 0
        fi
    done

    return 0
}

check_go_env_vars() {
    local env_vars=("GOPATH" "GOROOT" "GOBIN" "GOPROXY" "GONOPROXY" "GONOSUMDB" "GOOS" "GOARCH")

    for var in "${env_vars[@]}"; do
        local value=${!var:-}
        if [[ -n "$value" ]]; then
            log_success "$var: $value"
        fi
    done
}

print_recommendations() {
    log_info "Recommendations for optimal testing experience:"
    echo ""
    echo "1. Set GOMAXPROCS to match CPU cores: export GOMAXPROCS=$(nproc)"
    echo "2. Increase file descriptor limit: ulimit -n 65536"
    echo "3. Install file watcher for development: sudo apt install inotify-tools (Linux) or brew install fswatch (macOS)"
    echo "4. Configure Go module proxy for faster downloads: export GOPROXY=https://proxy.golang.org,direct"
    echo "5. Install Go development tools: go install golang.org/x/tools/cmd/...@latest"
    echo "6. Use Unix Makefile for test operations: cd tests/scripts && make help"
    echo ""
    log_success "Environment validation complete!"
}

# Run main function
main "$@"