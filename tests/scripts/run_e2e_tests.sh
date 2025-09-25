#!/bin/bash
set -e

if [[ "$RUNNER_OS" == "Linux" ]]; then
    sudo apt-get update
    sudo apt-get install -y asciinema
elif [[ "$RUNNER_OS" == "macOS" ]]; then
    brew install asciinema
fi

go build -o cmd/acousticalc/acousticalc ./cmd/acousticalc

go test -v -timeout=120s ./tests/e2e/...