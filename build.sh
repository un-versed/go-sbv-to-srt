#!/bin/bash

# Build script for go-sbv-to-srt

set -e

echo "Building go-sbv-to-srt..."

# Get version info
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# Build for current platform
echo "Building for current platform..."
go build -ldflags "${LDFLAGS}" -o go-sbv-to-srt .

echo "Build complete!"
echo "Version: ${VERSION}"
echo "Build time: ${BUILD_TIME}"
echo "Git commit: ${GIT_COMMIT}"

# Test the build
echo "Testing the build..."
./go-sbv-to-srt --help
