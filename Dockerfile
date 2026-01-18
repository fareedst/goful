# Multi-stage Dockerfile for Goful
# [IMPL:DOCKERFILE_MULTISTAGE] [ARCH:DOCKER_BUILD_STRATEGY] [REQ:DOCKER_INTERACTIVE_SETUP]

# Build stage: Compile Goful binary
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Enable automatic toolchain download for Go 1.24+
ENV GOTOOLCHAIN=auto
RUN go mod download

# Copy source code
COPY . .

# Build static binary with CGO disabled for portability
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o goful .

# Runtime stage: Minimal image with just the binary
FROM alpine:latest

# Install ca-certificates for HTTPS support
RUN apk add --no-cache ca-certificates

# Copy binary from builder stage
COPY --from=builder /build/goful /usr/local/bin/goful

# Set terminal environment variables required for tcell
ENV TERM=xterm-256color
ENV COLORTERM=truecolor

# Set working directory
WORKDIR /workspace

# Default entrypoint: run goful interactively
ENTRYPOINT ["goful"]
