# [IMPL:MAKE_RELEASE_TARGETS] Release Targets

**Cross-References**: [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

Automate reproducible releases via Makefile + CI + tag workflows.

## Rationale

- Guarantees that every release binary is built with the same flags (CGO disabled, `-trimpath -ldflags "-s -w"`), stored under predictable filenames, and accompanied by SHA256 digests
- Keeps local, CI, and tag-triggered release flows identical: `make release` locally mirrors both the CI `release-matrix` job and the GitHub Releases workflow

## Implementation Approach

### Makefile Enhancements

- New helpers: `vet`, `lint` (fmt + vet), `clean-release`, `release`, with defaults `DIST_DIR=dist`, `RELEASE_PLATFORMS="linux/amd64 linux/arm64 darwin/arm64"`, and `SHASUM=shasum -a 256`
- `release` target accepts optional `PLATFORM=os/arch`; otherwise iterates `RELEASE_PLATFORMS`. For each platform it emits `DIAGNOSTIC: [IMPL:MAKE_RELEASE_TARGETS] ...`, builds `dist/goful_${os}_${arch}`, and writes `dist/goful_${os}_${arch}.sha256`
- Targets stamp `[IMPL:MAKE_RELEASE_TARGETS] [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX]` in echoes so logs remain traceable

### CI Workflow (`release-matrix` job)

- Adds matrix include entries for linux/amd64, linux/arm64, darwin/arm64
- Step 1: checkout + setup Go 1.24.3
- Step 2: `make release PLATFORM=${{matrix.goos}}/${{matrix.goarch}}`
- Step 3: display checksum file and upload both binary + `.sha256` via `actions/upload-artifact`

### Release Workflow (`.github/workflows/release.yml`)

- Trigger: `push` tags matching `v*`
- Job `matrix-build` mirrors the CI matrix and runs the same `make release` command, uploading artifacts per platform
- Job `publish` downloads all artifacts (merged) and calls `softprops/action-gh-release` to attach binaries + `.sha256` files to the GitHub release while logging checksum contents for `ArtifactDeterminismAudit`

## Code Markers

- Makefile release recipes and workflow shell blocks include `[IMPL:MAKE_RELEASE_TARGETS] [ARCH:BUILD_MATRIX] [REQ:RELEASE_BUILD_MATRIX]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `Makefile` - release targets
- [ ] `.github/workflows/ci.yml` - release-matrix job
- [ ] `.github/workflows/release.yml` - all jobs

Tests that must reference `[REQ:RELEASE_BUILD_MATRIX]`:
- [ ] Local: `make release PLATFORM=$(go env GOOS)/$(go env GOARCH)` followed by `ls dist/`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | Host builds succeed, CI artifacts uploaded |

## Related Decisions

- Depends on: [IMPL:CI_WORKFLOW]
- See also: [ARCH:BUILD_MATRIX], [REQ:RELEASE_BUILD_MATRIX], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
