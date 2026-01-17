# [IMPL:TOKEN_VALIDATION_SCRIPT] Token Validation Script

**Cross-References**: [ARCH:TOKEN_VALIDATION_AUTOMATION] [REQ:STDD_SETUP]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

Automate `[PROC:TOKEN_VALIDATION]` via `scripts/validate_tokens.sh`.

## Rationale

- Contributors need a single command to prove token references are registered
- Satisfies modernization tasks blocked on running `[PROC:TOKEN_VALIDATION]`

## Implementation Approach

- Added `scripts/validate_tokens.sh` (Bash, `set -euo pipefail`)
- Script requires `git` and `rg`, builds the token registry from `stdd/semantic-tokens.md`, and scans tracked source files (`*.go`, module files, shell scripts, workflows, Makefile) unless custom paths are supplied
- Emits diagnostic output and fails if tokens found in source are missing from the registry
- Produces success message with counts for audit trails

## Code Markers

- Script header includes `[IMPL:TOKEN_VALIDATION_SCRIPT] [ARCH:TOKEN_VALIDATION_AUTOMATION] [REQ:STDD_SETUP]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `scripts/validate_tokens.sh` - script header

Tests that must reference `[REQ:STDD_SETUP]`:
- [ ] N/A - script is self-validating

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | `./scripts/validate_tokens.sh` (default globs) → verified 12 token references across 35 files |

## Related Decisions

- Depends on: —
- See also: [ARCH:TOKEN_VALIDATION_AUTOMATION], [REQ:STDD_SETUP], [PROC:TOKEN_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
