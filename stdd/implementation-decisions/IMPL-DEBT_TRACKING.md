# [IMPL:DEBT_TRACKING] Debt Tracking

**Cross-References**: [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

Track debt via issues and TODO annotations.

## Rationale

- Makes risk visible and assignable
- Keeps inline breadcrumbs synchronized with the central backlog so `[PROC:TOKEN_AUDIT]` can verify coverage

## Implementation Approach

- Capture the backlog in `stdd/debt-log.md` (D1–D4 for the initial pass) with owner, risk, TODO reference, and next action columns
- Annotate each hotspot (`app/goful.go`, `main.go`, `cmdline/cmdline.go`, `filer/filer.go`) with `TODO(goful-maintainers)` comments that describe the risk and include `[IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]`
- Update `stdd/tasks.md` to link the P1 debt triage task back to the backlog and record `[PROC:TOKEN_VALIDATION]` output after audits

## Code Markers

- TODOs carry `[IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]`
- Documentation links refer to the backlog file so future contributors know where to extend the list

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `app/goful.go` - TODO comments
- [ ] `main.go` - TODO comments
- [ ] `cmdline/cmdline.go` - TODO comments
- [ ] `filer/filer.go` - TODO comments
- [ ] `stdd/debt-log.md` - backlog entries

Tests that must reference `[REQ:DEBT_TRIAGE]`:
- [ ] N/A - debt tracking is documentation, not code

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | `./scripts/validate_tokens.sh` → verified 148 token references across 44 files |

## Related Decisions

- Depends on: —
- See also: [ARCH:DEBT_MANAGEMENT], [REQ:DEBT_TRIAGE]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
