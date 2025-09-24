#!/bin/bash

# AcoustiCalc Visual Testing Screenshot Utility
# Cross-platform screenshot capture for Unix systems (macOS/Linux)
# Supports native utilities with PNG/sRGB output

set -euo pipefail

# Configuration
SCREENSHOT_DIR="${SCREENSHOT_DIR:-tests/artifacts/screenshots}"
FORMAT="png"
QUALITY=100

# Detect platform and set screenshot command
detect_platform() {
    case "$(uname -s)" in
        Darwin*)
            PLATFORM="macos"
            SCREENSHOT_CMD="screencapture"
            ;;
        Linux*)
            PLATFORM="linux"
            # Try different Linux screenshot tools in order of preference
            if command -v gnome-screenshot >/dev/null 2>&1; then
                SCREENSHOT_CMD="gnome-screenshot"
            elif command -v scrot >/dev/null 2>&1; then
                SCREENSHOT_CMD="scrot"
            elif command -v import >/dev/null 2>&1; then
                SCREENSHOT_CMD="import"
            else
                echo "Error: No screenshot utility found. Please install gnome-screenshot, scrot, or ImageMagick" >&2
                exit 1
            fi
            ;;
        *)
            echo "Error: Unsupported platform $(uname -s)" >&2
            exit 1
            ;;
    esac
}

# Create screenshot with platform-specific command
capture_screenshot() {
    local output_file="$1"
    local delay="${2:-0}"

    # Ensure output directory exists
    mkdir -p "$(dirname "$output_file")"

    case "$PLATFORM" in
        macos)
            # macOS screencapture with PNG format and optional delay
            if [ "$delay" -gt 0 ]; then
                screencapture -t png -T "$delay" "$output_file"
            else
                screencapture -t png "$output_file"
            fi
            ;;
        linux)
            case "$SCREENSHOT_CMD" in
                gnome-screenshot)
                    if [ "$delay" -gt 0 ]; then
                        gnome-screenshot -f "$output_file" --delay="$delay"
                    else
                        gnome-screenshot -f "$output_file"
                    fi
                    ;;
                scrot)
                    if [ "$delay" -gt 0 ]; then
                        scrot -d "$delay" "$output_file"
                    else
                        scrot "$output_file"
                    fi
                    ;;
                import)
                    if [ "$delay" -gt 0 ]; then
                        sleep "$delay"
                    fi
                    import -window root "$output_file"
                    ;;
            esac
            ;;
    esac
}

# Generate screenshot filename with metadata
generate_filename() {
    local test_name="$1"
    local event_type="$2"
    local timestamp="${3:-$(date +%Y%m%d_%H%M%S)}"

    echo "${test_name}_${event_type}_${timestamp}.png"
}

# Add metadata to screenshot file
add_metadata() {
    local file="$1"
    local test_name="$2"
    local event_type="$3"
    local platform="$4"

    # Use exiftool if available to add metadata
    if command -v exiftool >/dev/null 2>&1; then
        exiftool -overwrite_original \
            -XMP:Subject="AcoustiCalc Visual Test" \
            -XMP:Title="$test_name - $event_type" \
            -XMP:Creator="AcoustiCalc Visual Testing Framework" \
            -XMP:Description="Test: $test_name, Event: $event_type, Platform: $platform" \
            "$file" >/dev/null 2>&1 || true
    fi
}

# Optimize screenshot for demo use
optimize_screenshot() {
    local input_file="$1"
    local output_file="$2"
    local max_width="${3:-1920}"
    local max_height="${4:-1080}"

    # Use ImageMagick convert if available
    if command -v convert >/dev/null 2>&1; then
        convert "$input_file" \
            -resize "${max_width}x${max_height}>" \
            -quality 95 \
            -colorspace sRGB \
            "$output_file"
    else
        # Fallback: just copy the file
        cp "$input_file" "$output_file"
    fi
}

# Main screenshot capture function
capture_test_screenshot() {
    local test_name="$1"
    local event_type="$2"
    local category="${3:-unit}"
    local delay="${4:-0}"

    # Generate paths
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local filename=$(generate_filename "$test_name" "$event_type" "$timestamp")
    local output_path="$SCREENSHOT_DIR/$category/$filename"

    # Capture screenshot
    echo "Capturing screenshot: $output_path"
    capture_screenshot "$output_path" "$delay"

    # Add metadata
    add_metadata "$output_path" "$test_name" "$event_type" "$PLATFORM"

    # Verify file was created and has reasonable size
    if [ -f "$output_path" ] && [ "$(stat -f%z "$output_path" 2>/dev/null || stat -c%s "$output_path" 2>/dev/null)" -gt 1024 ]; then
        echo "Screenshot captured successfully: $output_path"
        echo "$output_path"
    else
        echo "Error: Screenshot capture failed or file is too small" >&2
        exit 1
    fi
}

# Performance test - measure screenshot time
performance_test() {
    local test_name="screenshot_performance"
    local start_time=$(date +%s%N)

    capture_test_screenshot "$test_name" "performance" "unit" 0 >/dev/null

    local end_time=$(date +%s%N)
    local duration_ms=$(( (end_time - start_time) / 1000000 ))

    echo "Screenshot performance: ${duration_ms}ms"

    # Check if under 5 second constraint
    if [ "$duration_ms" -lt 5000 ]; then
        echo "✅ Performance acceptable (<5s constraint)"
    else
        echo "❌ Performance too slow (${duration_ms}ms > 5000ms)"
    fi
}

# Batch screenshot capture for test suite
batch_capture() {
    local test_name="$1"
    local category="${2:-unit}"

    echo "Starting batch screenshot capture for: $test_name"

    # Capture test lifecycle events
    capture_test_screenshot "$test_name" "start" "$category" 0
    sleep 1
    capture_test_screenshot "$test_name" "process" "$category" 0
    sleep 1
    capture_test_screenshot "$test_name" "complete" "$category" 0

    echo "Batch capture completed for: $test_name"
}

# Display help information
show_help() {
    cat << EOF
AcoustiCalc Visual Testing Screenshot Utility

Usage: $0 [COMMAND] [OPTIONS]

Commands:
    capture TEST_NAME EVENT_TYPE [CATEGORY] [DELAY]
        Capture a single screenshot

    batch TEST_NAME [CATEGORY]
        Capture a series of screenshots for test lifecycle

    performance
        Run performance test and measure capture speed

    help
        Show this help message

Examples:
    $0 capture calculator_test start unit 0
    $0 batch integration_test integration
    $0 performance

Environment Variables:
    SCREENSHOT_DIR    Directory for screenshots (default: tests/artifacts/screenshots)

Requirements:
    macOS: screencapture (built-in)
    Linux: gnome-screenshot, scrot, or ImageMagick

EOF
}

# Initialize platform detection
detect_platform
echo "Platform detected: $PLATFORM using $SCREENSHOT_CMD"

# Command dispatch
case "${1:-help}" in
    capture)
        if [ $# -lt 3 ]; then
            echo "Error: capture requires TEST_NAME and EVENT_TYPE" >&2
            exit 1
        fi
        capture_test_screenshot "$2" "$3" "${4:-unit}" "${5:-0}"
        ;;
    batch)
        if [ $# -lt 2 ]; then
            echo "Error: batch requires TEST_NAME" >&2
            exit 1
        fi
        batch_capture "$2" "${3:-unit}"
        ;;
    performance)
        performance_test
        ;;
    help)
        show_help
        ;;
    *)
        echo "Error: Unknown command '$1'" >&2
        show_help
        exit 1
        ;;
esac