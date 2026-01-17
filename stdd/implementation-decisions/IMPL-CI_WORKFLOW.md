# [IMPL:CI_WORKFLOW] CI Workflow

**Cross-References**: [ARCH:CI_PIPELINE] [REQ:CI_PIPELINE_CORE]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

GitHub Actions workflow for fmt/vet/tests.

## Rationale

- Enforces consistency and prevents regressions on PRs

## Implementation Approach

- Added `.github/workflows/ci.yml` with a `format-vet-test` job that:
  - Checks out the repo, sets up Go `1.24.3` via `actions/setup-go@v5`, and caches modules (`go.sum`)
  - Runs `gofmt -w $(git ls-files '*.go')` followed by `git status --short` and `git diff --exit-code` to enforce formatting
  - Executes `go vet ./...` and `go test ./...` for regression coverage
  - Runs `./scripts/validate_tokens.sh` so every CI pass records `[PROC:TOKEN_VALIDATION]`
- Each shell block embeds `[IMPL:CI_WORKFLOW] [ARCH:CI_PIPELINE] [REQ:CI_PIPELINE_CORE]` (or the token validation equivalents) to keep traceability in the workflow itself

## Code Markers

- Workflow file comments with `[IMPL:CI_WORKFLOW] [ARCH:CI_PIPELINE] [REQ:CI_PIPELINE_CORE]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `.github/workflows/ci.yml` - all job steps

Tests that must reference `[REQ:CI_PIPELINE_CORE]`:
- [ ] N/A - workflow is the test itself

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | `./scripts/validate_tokens.sh` → verified 19 token references across 36 files |

## Related Decisions

- Depends on: [IMPL:GO_MOD_UPDATE]
- See also: [ARCH:CI_PIPELINE], [REQ:CI_PIPELINE_CORE]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
