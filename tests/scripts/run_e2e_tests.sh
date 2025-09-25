#!/bin/bash
set -e

if [[ "$RUNNER_OS" == "Linux" ]]; then
    sudo apt-get update
    sudo apt-get install -y asciinema
elif [[ "$RUNNER_OS" == "macOS" ]]; then
    brew install asciinema
elif [[ "$RUNNER_OS" == "Windows" ]]; then
    pwsh -Command "choco install asciinema -y"
    export PATH="$PATH:/c/Program Files/Asciinema/bin"
fi

if [[ "$RUNNER_OS" == "Windows" ]]; then
    go build -o cmd/acousticalc/acousticalc.exe ./cmd/acousticalc
else
    go build -o cmd/acousticalc/acousticalc ./cmd/acousticalc
fi

go test -v -timeout=120s ./tests/e2e/...