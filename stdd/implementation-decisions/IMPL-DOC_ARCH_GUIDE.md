# [IMPL:DOC_ARCH_GUIDE] Architecture Guide

**Cross-References**: [ARCH:DOCS_STRUCTURE] [REQ:ARCH_DOCUMENTATION]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

Write `ARCHITECTURE.md`.

## Rationale

- Provides concise understanding of packages and data flow
- Establishes a stable "map" before larger refactors touch keymap/menu wiring

## Implementation Approach

Structure the document into:

1. **Overview & Goals** referencing `[REQ:ARCH_DOCUMENTATION]`
2. **Runtime Flow** describing `main` → `configpaths.Resolver` → `app.Goful` event loop
3. **Module Deep Dives** (`app`, `filer`, `widget`, `cmdline`, `menu`, `look/message/info/progress`, `configpaths`, `util`)
4. **Validation & Testing Surfaces** listing module-level tests (widgets, cmdline, integration flows, keymap baselines)

- Embed ASCII-style flow diagrams or bullet chains to highlight dependencies
- Cross-link to requirements/architecture/implementation tokens inline to preserve STDD traceability

## Code Markers

- Document contains `[IMPL:DOC_ARCH_GUIDE] [ARCH:DOCS_STRUCTURE] [REQ:ARCH_DOCUMENTATION]` plus relevant cross-references for each section (e.g., `[REQ:CONFIGURABLE_STATE_PATHS]`, `[ARCH:STATE_PATH_SELECTION]`)

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `ARCHITECTURE.md` - document header and sections

Tests that must reference `[REQ:ARCH_DOCUMENTATION]`:
- [ ] Document review ensures every section references at least one `[REQ:*]` token

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | Document linked from README.md |

## Related Decisions

- Depends on: —
- See also: [ARCH:DOCS_STRUCTURE], [REQ:ARCH_DOCUMENTATION], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
