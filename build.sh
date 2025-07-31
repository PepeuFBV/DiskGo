#!/bin/bash

# DiskGo Build Script
# This script builds the DiskGo application for Linux distribution

set -e  # Exit on any error

echo "Building DiskGo for Linux..."

# Clean previous builds
if [ -f "diskgo" ]; then
    echo "Removing previous build..."
    rm diskgo
fi

# Build with optimizations
echo "Compiling with optimizations..."
CGO_ENABLED=1 go build \
    -ldflags="-s -w -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
    -o diskgo \
    src/main.go

# Check if build was successful
if [ -f "diskgo" ]; then
    echo "✅ Build successful!"
    echo "Binary size: $(du -h diskgo | cut -f1)"
    echo "Binary location: $(pwd)/diskgo"
    echo ""
    echo "To run: ./diskgo"
    echo "To install system-wide: sudo cp diskgo /usr/local/bin/"
else
    echo "❌ Build failed!"
    exit 1
fi
