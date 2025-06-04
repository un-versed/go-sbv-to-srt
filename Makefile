# Variables
BINARY_NAME=go-sbv-to-srt
MAIN_PACKAGE=.
BUILD_DIR=dist
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building ${BINARY_NAME}..."
	go build ${LDFLAGS} -o ${BINARY_NAME} ${MAIN_PACKAGE}

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	go test -race ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f ${BINARY_NAME}
	@rm -rf ${BUILD_DIR}
	@rm -f coverage.out coverage.html

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
.PHONY: vet
vet:
	@echo "Vetting code..."
	go vet ./...

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p ${BUILD_DIR}
	@for os in linux windows darwin; do \
		for arch in amd64 arm64; do \
			ext=""; \
			if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
			echo "Building $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-$$os-$$arch$$ext ${MAIN_PACKAGE}; \
		done; \
	done

# Run the application with sample data
.PHONY: run-sample
run-sample: build
	@echo "Running with sample data..."
	./${BINARY_NAME} -i ./testdata/sample.sbv

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  test        - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-race   - Run tests with race detection"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Install dependencies"
	@echo "  fmt         - Format code"
	@echo "  vet         - Vet code"
	@echo "  lint        - Lint code (requires golangci-lint)"
	@echo "  build-all   - Build for multiple platforms"
	@echo "  run-sample  - Run with sample data"
	@echo "  help        - Show this help"

# Development workflow
.PHONY: dev
dev: deps fmt vet test build
	@echo "Development workflow complete!"
