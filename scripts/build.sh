#!/bin/bash

set -e

APP_NAME="Goofed"
DIST_DIR="./.build"
ROOT="./cmd/goofed/main.go"

mkdir -p "$DIST_DIR"

echo "🔧 Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o "$DIST_DIR/$APP_NAME.exe" -ldflags="-H=windowsgui" "$ROOT"

echo "🍎 Building for macOS..."
echo "Apple Silicon"
GOOS=darwin GOARCH=arm64 go build -o "$DIST_DIR/$APP_NAME-macos" "$ROOT"

# echo "Apple Intel"
# export SDKROOT="$(xcrun --sdk macosx --show-sdk-path)"
# GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 \
#   CC="clang -arch x86_64 -isysroot $SDKROOT" \
#   go build -o "$DIST_DIR/$APP_NAME-macos" "$ROOT"

echo "🐧 Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o "$DIST_DIR/$APP_NAME-linux" "$ROOT"

# echo "✅ All builds completed. Files in $DIST_DIR:" "$ROOT"
# ls -lh "$DIST_DIR"
