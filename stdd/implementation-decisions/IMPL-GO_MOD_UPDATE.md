# [IMPL:GO_MOD_UPDATE] Go Mod Update

**Cross-References**: [ARCH:GO_RUNTIME_STRATEGY] [REQ:GO_TOOLCHAIN_LTS]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

Update `go.mod` to current LTS and align tooling.

## Rationale

- Required to use supported compiler features and security fixes
- Aligns local/CI builds and reduces drift

## Implementation Approach

- Set `go 1.24.0` plus `toolchain go1.24.3` in `go.mod` with `[IMPL:GO_MOD_UPDATE] [ARCH:GO_RUNTIME_STRATEGY] [REQ:GO_TOOLCHAIN_LTS]` comments
- Recomputed the module graph with `go mod tidy` under `go1.24.3`
- Updated `message.log` to use a constant format string so `go vet` passes on the modern toolchain
- Verified the upgraded toolchain by running `go fmt ./...`, `go vet ./...`, and `go test ./...` on darwin/arm64 (go1.24.3)

## Code Markers

- `go.mod` header and CI workflow `setup-go` steps

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `go.mod` - version comments
- [ ] CI workflow - setup-go version

Tests that must reference `[REQ:GO_TOOLCHAIN_LTS]`:
- [ ] N/A - verified via toolchain execution

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | go1.24.3 darwin/arm64 |

Validation commands executed:
- `go fmt ./...` (touched `info/info_unix.go`, `info/info_windows.go`)
- `go vet ./...`
- `go test ./...`
- `./scripts/validate_tokens.sh` → verified 12 token references across 35 files

## Related Decisions

- Depends on: —
- See also: [ARCH:GO_RUNTIME_STRATEGY], [REQ:GO_TOOLCHAIN_LTS]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
