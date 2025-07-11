#!/bin/bash

# Simple local Docker build test (no external cache)

echo "🚀 Testing Docker build locally (simple mode)..."

# Enable BuildKit
export DOCKER_BUILDKIT=1

# Detect platform
if [[ "$(uname -m)" == "arm64" ]] || [[ "$(uname -m)" == "aarch64" ]]; then
    PLATFORM="linux/arm64"
else
    PLATFORM="linux/amd64"
fi

echo "📦 Building for $PLATFORM..."
echo "⚡ Using Docker's internal build cache"

# Simple build with internal cache only
docker build \
    --platform $PLATFORM \
    --tag yaca:local-test \
    --build-arg BUILDKIT_INLINE_CACHE=1 \
    .

if [ $? -eq 0 ]; then
    echo "✅ Build complete!"
    echo "🐳 Image tagged as: yaca:local-test"
    echo ""
    echo "To run the image:"
    echo "  docker run --rm yaca:local-test"
    echo ""
    echo "Note: This uses Docker's internal cache only."
    echo "Subsequent builds will be faster if layers haven't changed."
else
    echo "❌ Build failed!"
fi