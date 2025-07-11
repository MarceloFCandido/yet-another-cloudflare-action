#!/bin/bash

# Local Docker build test script with optimizations

echo "🚀 Testing optimized Docker build locally..."

# Enable BuildKit
export DOCKER_BUILDKIT=1

# Check if we have a buildx builder that supports cache
if ! docker buildx ls | grep -q "docker-container"; then
    echo "⚠️  No buildx builder with cache support found."
    echo "📝 Creating a new builder instance..."
    docker buildx create --name mybuilder --driver docker-container --use
    docker buildx inspect --bootstrap
fi

# Get current builder
CURRENT_BUILDER=$(docker buildx ls | grep -E '^\*' | awk '{print $1}' | sed 's/\*//')
echo "🔧 Using builder: $CURRENT_BUILDER"

# Detect platform
if [[ "$(uname -m)" == "arm64" ]] || [[ "$(uname -m)" == "aarch64" ]]; then
    PLATFORM="linux/arm64"
else
    PLATFORM="linux/amd64"
fi

echo "📦 Building for $PLATFORM..."

# Build with caching (if supported)
if docker buildx ls | grep -q "docker-container"; then
    # Full cache support
    docker buildx build \
        --platform $PLATFORM \
        --tag yaca:local-test \
        --cache-from type=local,src=/tmp/.buildx-cache \
        --cache-to type=local,dest=/tmp/.buildx-cache,mode=max \
        --build-arg BUILDKIT_INLINE_CACHE=1 \
        --load \
        .
else
    # Fallback for default driver
    echo "⚠️  Using default driver (limited cache support)"
    docker buildx build \
        --platform $PLATFORM \
        --tag yaca:local-test \
        --build-arg BUILDKIT_INLINE_CACHE=1 \
        --load \
        .
fi

if [ $? -eq 0 ]; then
    echo "✅ Build complete!"
    echo "🐳 Image tagged as: yaca:local-test"
    echo ""
    echo "To run the image:"
    echo "  docker run --rm yaca:local-test"
    echo ""
    echo "To test multi-platform build:"
    echo "  docker buildx build --platform linux/amd64,linux/arm64 ."
    echo ""
    echo "To switch back to default builder:"
    echo "  docker buildx use default"
else
    echo "❌ Build failed!"
    echo ""
    echo "Troubleshooting tips:"
    echo "1. Ensure Docker Desktop is running"
    echo "2. Try: docker buildx use default"
    echo "3. Or create a new builder: docker buildx create --use"
fi