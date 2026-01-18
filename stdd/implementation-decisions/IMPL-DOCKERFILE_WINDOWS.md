# [IMPL:DOCKERFILE_WINDOWS] Windows Dockerfile

**Cross-References**: [ARCH:DOCKER_WINDOWS_BUILD] [REQ:DOCKER_WINDOWS_CONTAINER]  
**Status**: Active  
**Created**: 2026-01-18  
**Last Updated**: 2026-01-18

---

## Decision

Create a separate Windows-specific Dockerfile (`Dockerfile.windows`) and supporting files to enable Windows container execution of Goful for testing on Windows Server environments.

## Rationale

- Windows containers require different base images than Linux containers
- Enables testing Goful's Windows binary in a containerized environment
- Complements the existing Alpine-based Linux container for cross-platform validation
- The existing codebase already has Windows support via `info/info_windows.go` with proper build tags
- tcell library supports Windows console applications

## Implementation Approach

### Files to Create

1. **`Dockerfile.windows`** - Windows-specific multi-stage build

```dockerfile
# Build stage: Cross-compile Go binary for Windows
# [IMPL:DOCKERFILE_WINDOWS] [ARCH:DOCKER_WINDOWS_BUILD] [REQ:DOCKER_WINDOWS_CONTAINER]
FROM golang:1.23 AS builder

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Enable automatic toolchain download for Go 1.24+
ENV GOTOOLCHAIN=auto
RUN go mod download

# Copy source code
COPY . .

# Build static binary for Windows
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o goful.exe .

# Runtime stage: Windows Server Core for full console API support
FROM mcr.microsoft.com/windows/servercore:ltsc2022

# Copy binary from builder stage
COPY --from=builder /build/goful.exe C:/goful/goful.exe

# Set working directory
WORKDIR C:/workspace

# Default entrypoint: run goful interactively
ENTRYPOINT ["C:\\goful\\goful.exe"]
```

2. **`docker-compose.windows.yml`** - Windows-specific compose file

```yaml
# Docker Compose configuration for Goful on Windows
# [IMPL:DOCKERFILE_WINDOWS] [ARCH:DOCKER_WINDOWS_BUILD] [REQ:DOCKER_WINDOWS_CONTAINER]

services:
  goful:
    build:
      context: .
      dockerfile: Dockerfile.windows
    image: goful:windows
    volumes:
      # Mount current directory for file operations (Windows paths)
      - .:C:/workspace
    # Enable interactive terminal
    stdin_open: true
    tty: true
    working_dir: C:/workspace
    # Pass through any command-line arguments to goful
    command: []
```

3. **`docker-run.ps1`** - PowerShell helper script

```powershell
# Helper script for running Goful in Windows Docker container
# [IMPL:DOCKERFILE_WINDOWS] [ARCH:DOCKER_WINDOWS_BUILD] [REQ:DOCKER_WINDOWS_CONTAINER]

param(
    [switch]$Build,
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$Args
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ImageName = "goful:windows"

# Check if Docker is available
if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Error "Docker is not installed or not in PATH"
    exit 1
}

# Build image if requested or if it doesn't exist
$ImageExists = docker image inspect $ImageName 2>$null
if ($Build -or -not $ImageExists) {
    Write-Host "Building Docker image: $ImageName"
    docker build -f "$ScriptDir\Dockerfile.windows" -t $ImageName $ScriptDir
}

# Run container interactively
docker run -it --rm `
    -v "${ScriptDir}:C:\workspace" `
    -w "C:\workspace" `
    $ImageName @Args
```

### Platform Compatibility

- Windows containers require Windows host (Docker Desktop with Windows containers mode or Windows Server with containers feature)
- Build can be performed from any platform using cross-compilation (`GOOS=windows`)
- Runtime requires Windows Server 2019, 2022, or Windows 10/11 with containers

### Limitations

| Feature | Status | Notes |
|---------|--------|-------|
| Basic keyboard | ✅ Works | tcell supports Windows console |
| Display/colors | ⚠️ Limited | Windows console may not support all ANSI escapes |
| Mouse input | ⚠️ Limited | May be limited in container environment |
| Terminal resize | ⚠️ Limited | Windows console API behavior may differ |
| nsync features | ❌ Stub only | darwin-only implementation |

## Code Markers

- `Dockerfile.windows`: `[IMPL:DOCKERFILE_WINDOWS] [ARCH:DOCKER_WINDOWS_BUILD] [REQ:DOCKER_WINDOWS_CONTAINER]`
- `docker-compose.windows.yml`: `[IMPL:DOCKERFILE_WINDOWS] [ARCH:DOCKER_WINDOWS_BUILD] [REQ:DOCKER_WINDOWS_CONTAINER]`
- `docker-run.ps1`: `[IMPL:DOCKERFILE_WINDOWS] [ARCH:DOCKER_WINDOWS_BUILD] [REQ:DOCKER_WINDOWS_CONTAINER]`
- `Makefile`: Docker targets annotated with `[IMPL:DOCKERFILE_WINDOWS] [ARCH:DOCKER_WINDOWS_BUILD] [REQ:DOCKER_WINDOWS_CONTAINER]`

### Makefile Targets

The following Makefile targets were added for Docker operations:

```makefile
# Linux container targets
docker-build         # Build Alpine image
docker-run           # Run interactively

# Windows container targets  
docker-build-windows # Build ServerCore image
docker-run-windows   # Run interactively

# Cleanup
docker-clean         # Remove both images
```

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that carry annotations:
- [x] `Dockerfile.windows` header comment
- [x] `docker-compose.windows.yml` header comment
- [x] `docker-run.ps1` header comment
- [x] `Makefile` Docker variables and targets

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-18 | — | ✅ Pass | All files created, Makefile targets added, README updated; token validation passes (1760 refs/85 files); Windows host testing deferred |

## Related Decisions

- Depends on: [ARCH:DOCKER_WINDOWS_BUILD]
- See also: [IMPL:DOCKERFILE_MULTISTAGE], [IMPL:DOCKER_COMPOSE_CONFIG], [ARCH:DOCKER_BUILD_STRATEGY]

---

*Created 2026-01-18*
