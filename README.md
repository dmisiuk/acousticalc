# acousticalc

A TUI-based application for calculating room acoustics.

## Build

To build the application, run:

```bash
go build -o acousticalc ./cmd/acousticalc
```

## Usage

After building, you can run the application with:

```bash
./acousticalc "2 + 3 * 4"
```

## Installation

Download the latest release from the [GitHub Releases](https://github.com/dmisiuk/acousticalc/releases) page.

## GitHub Actions

This project uses GitHub Actions for CI/CD:

- CI workflow: Runs on every push and pull request to test the application
- Release workflow: Runs when a new tag is pushed to create a new release

## Development

To run tests:

```bash
go test ./...
```

To run tests with coverage:

```bash
go test -cover ./...
```
