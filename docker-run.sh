#!/bin/bash
# Helper script for running Goful in Docker
# [IMPL:DOCKER_COMPOSE_CONFIG] [ARCH:DOCKER_BUILD_STRATEGY] [REQ:DOCKER_INTERACTIVE_SETUP]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
IMAGE_NAME="goful:latest"

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed or not in PATH" >&2
    exit 1
fi

# Check if docker-compose is available (optional, fallback to docker run)
USE_COMPOSE=false
if command -v docker-compose &> /dev/null || docker compose version &> /dev/null 2>&1; then
    USE_COMPOSE=true
fi

# Build image if it doesn't exist or if --build flag is provided
if [[ "${1:-}" == "--build" ]] || ! docker image inspect "$IMAGE_NAME" &> /dev/null; then
    echo "Building Docker image: $IMAGE_NAME"
    docker build -t "$IMAGE_NAME" "$SCRIPT_DIR"
    if [[ "${1:-}" == "--build" ]]; then
        shift  # Remove --build flag from arguments
    fi
fi

# Run container with docker-compose if available, otherwise use docker run
if [ "$USE_COMPOSE" = true ]; then
    # Use docker-compose for easier volume and environment management
    cd "$SCRIPT_DIR"
    if docker compose version &> /dev/null 2>&1; then
        docker compose run --rm goful "$@"
    else
        docker-compose run --rm goful "$@"
    fi
else
    # Fallback to docker run
    docker run -it --rm \
        -v "$SCRIPT_DIR:/workspace" \
        -v goful-config:/root/.goful \
        -w /workspace \
        -e TERM=xterm-256color \
        -e COLORTERM=truecolor \
        "$IMAGE_NAME" \
        "$@"
fi
