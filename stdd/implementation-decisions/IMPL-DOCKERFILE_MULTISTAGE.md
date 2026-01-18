# [IMPL:DOCKERFILE_MULTISTAGE] Multi-Stage Dockerfile

**Cross-References**: [ARCH:DOCKER_BUILD_STRATEGY] [REQ:DOCKER_INTERACTIVE_SETUP]  
**Status**: Active  
**Created**: 2026-01-18  
**Last Updated**: 2026-01-18

---

## Decision

Use a multi-stage Dockerfile to build Goful for Linux containers, with platform-specific build tags to handle macOS-only dependencies.

## Rationale

- Multi-stage builds separate build-time dependencies from runtime, minimizing image size
- Static binaries (`CGO_ENABLED=0`) ensure portability across Linux distributions
- `GOTOOLCHAIN=auto` allows the Go 1.23 base image to download Go 1.24.3 as required by go.mod
- Build tags enable graceful degradation of nsync features on Linux

## Implementation Approach

### Dockerfile Structure

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o goful .

# Runtime stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/goful /usr/local/bin/goful
ENV TERM=xterm-256color
ENV COLORTERM=truecolor
WORKDIR /workspace
ENTRYPOINT ["goful"]
```

### Platform Compatibility

Created build-tagged files to handle nsync's macOS-specific syscalls:
- `app/nsync.go` - `//go:build darwin` (full implementation)
- `app/nsync_stub.go` - `//go:build !darwin` (stubs that fall back to regular copy/move)
- `app/nsync_test.go` - `//go:build darwin` (tests only on macOS)

### .dockerignore

Excludes build artifacts, git files, and IDE configuration to minimize build context size.

## Code Markers

- `Dockerfile`: `[IMPL:DOCKERFILE_MULTISTAGE] [ARCH:DOCKER_BUILD_STRATEGY] [REQ:DOCKER_INTERACTIVE_SETUP]`
- `.dockerignore`: `[IMPL:DOCKERFILE_MULTISTAGE] [ARCH:DOCKER_BUILD_STRATEGY] [REQ:DOCKER_INTERACTIVE_SETUP]`
- `app/nsync_stub.go`: `[IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that carry annotations:
- [x] `Dockerfile` header comment
- [x] `.dockerignore` header comment
- [x] `app/nsync_stub.go` package comment and function comments

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-18 | — | ✅ Pass | Docker image builds successfully, tests pass on macOS |

## Related Decisions

- Depends on: [ARCH:DOCKER_BUILD_STRATEGY]
- See also: [IMPL:DOCKER_COMPOSE_CONFIG], [PROC:DOCKER_CONTAINER_SETUP]

---

*Created 2026-01-18*
