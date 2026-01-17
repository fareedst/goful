# [IMPL:DOC_CONTRIBUTING] CONTRIBUTING Guide

**Cross-References**: [ARCH:CONTRIBUTION_PROCESS] [REQ:CONTRIBUTING_GUIDE]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

Add contributor standards document.

## Rationale

- Aligns development workflow and review expectations
- Documents Go / Makefile targets, CI steps, and STDD-specific requirements (semantic tokens, module validation, debug logging)

## Implementation Approach

Sections:

- **Tooling & Setup** (Go LTS, `make` targets, local environment variables)
- **Workflow Checklist** enumerating fmt → vet → test → race/staticcheck (via CI) plus manual `./scripts/validate_tokens.sh`
- **Semantic Token Discipline** linking to registry updates and `[PROC:TOKEN_AUDIT]`
- **Module Validation Expectations** referencing `[REQ:MODULE_VALIDATION]` and `KeymapBaselineSuite`
- **Debug Logging Policy** (retain `DEBUG:`/`DIAGNOSTIC:` outputs)
- **Review Gate** referencing required doc/test updates before opening PRs

Provide copy/paste friendly command blocks (e.g., `make fmt`, `go test ./...`).

Link to `ARCHITECTURE.md`, `README.md`, and STDD docs for quick navigation.

## Code Markers

- `CONTRIBUTING.md` includes `[IMPL:DOC_CONTRIBUTING] [ARCH:CONTRIBUTION_PROCESS] [REQ:CONTRIBUTING_GUIDE] [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION] [REQ:MODULE_VALIDATION]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `CONTRIBUTING.md` - document header

Tests that must reference `[REQ:CONTRIBUTING_GUIDE]`:
- [ ] Document is cross-linked from README

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | Includes instructions to run validation script |

## Related Decisions

- Depends on: [IMPL:DOC_ARCH_GUIDE]
- See also: [ARCH:CONTRIBUTION_PROCESS], [REQ:CONTRIBUTING_GUIDE], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
