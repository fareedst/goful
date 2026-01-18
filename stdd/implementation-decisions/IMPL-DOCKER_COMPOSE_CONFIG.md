# [IMPL:DOCKER_COMPOSE_CONFIG] Docker Compose Configuration

**Cross-References**: [ARCH:DOCKER_BUILD_STRATEGY] [REQ:DOCKER_INTERACTIVE_SETUP]  
**Status**: Active  
**Created**: 2026-01-18  
**Last Updated**: 2026-01-18

---

## Decision

Provide docker-compose.yml and helper script for simplified interactive container execution with proper volume mounts and terminal configuration.

## Rationale

- Docker Compose simplifies multi-option docker run commands
- Named volumes enable config persistence across container restarts
- Helper script provides one-command execution for developers
- Terminal environment variables required for tcell color support

## Implementation Approach

### docker-compose.yml

```yaml
services:
  goful:
    build:
      context: .
      dockerfile: Dockerfile
    image: goful:latest
    volumes:
      - .:/workspace
      - goful-config:/root/.goful
    environment:
      - TERM=xterm-256color
      - COLORTERM=truecolor
    stdin_open: true
    tty: true
    working_dir: /workspace

volumes:
  goful-config:
```

### Helper Script (docker-run.sh)

Features:
- Checks Docker availability
- Auto-builds image if not present or with `--build` flag
- Uses docker-compose when available, falls back to docker run
- Passes through CLI arguments to goful

### Usage

```bash
# Build and run interactively
./docker-run.sh

# Run with specific directories
./docker-run.sh /workspace/dir1 /workspace/dir2

# Force rebuild
./docker-run.sh --build

# Or use docker-compose directly
docker compose run --rm goful
```

## Code Markers

- `docker-compose.yml`: `[IMPL:DOCKER_COMPOSE_CONFIG] [ARCH:DOCKER_BUILD_STRATEGY] [REQ:DOCKER_INTERACTIVE_SETUP]`
- `docker-run.sh`: `[IMPL:DOCKER_COMPOSE_CONFIG] [ARCH:DOCKER_BUILD_STRATEGY] [REQ:DOCKER_INTERACTIVE_SETUP]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that carry annotations:
- [x] `docker-compose.yml` header comment
- [x] `docker-run.sh` header comment

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-18 | — | ✅ Pass | Helper script executes correctly, docker-compose runs container |

## Related Decisions

- Depends on: [ARCH:DOCKER_BUILD_STRATEGY], [IMPL:DOCKERFILE_MULTISTAGE]
- See also: [PROC:DOCKER_CONTAINER_SETUP]

---

*Created 2026-01-18*
