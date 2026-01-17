# [IMPL:TEST_WIDGETS] Widget Tests

**Cross-References**: [ARCH:TEST_STRATEGY_UI] [REQ:UI_PRIMITIVE_TESTS]  
**Status**: Active  
**Created**: 2026-01-01  
**Last Updated**: 2026-01-17

---

## Decision

Add unit/snapshot tests for widget primitives.

## Rationale

- Protect rendering/event handling behaviors

## Implementation Approach

### Module Identification `[REQ:MODULE_VALIDATION]`

- `widget.ListBox` (cursor + scrolling state machine)
- `widget.Gauge` (progress rendering & normalization)
- `widget.TextBox` (buffer editing helpers such as `InsertBytes`/`DeleteBytes`)
- Supporting `widget.Window` helpers (column calculations, offsets)

### Validation Criteria

- Pure functions (cursor math, offset adjustments, column sizing) validated with table-driven Go tests
- Rendering helpers validated indirectly by inspecting state (e.g., `ScrollRate`, `ColumnAdjustContentsWidth`) to avoid brittle terminal assertions; future work can introduce snapshot harnesses once deterministic screen fakes exist

### Immediate Test Plan

- `TestListBoxCursorClamping_REQ_UI_PRIMITIVE_TESTS`: proves `SetCursor`, `MoveCursor`, and `SetCursorByName` respect `Lower()/Upper()` bounds and fallback semantics
- `TestListBoxScrollRate_REQ_UI_PRIMITIVE_TESTS`: verifies offset math for `ScrollRate` ("Top"/percentage/"Bot")
- `TestListBoxColumnAdjust_REQ_UI_PRIMITIVE_TESTS`: confirms column auto-fit honors widest content within available width
- Follow-on work: add gauge fill-ratio tests and textbox editing regressions; snapshot harness will cover highlight rendering once `SetCells` fakes land

Tests live in `widget/*.go` so they can directly access unexported helpers while retaining `[REQ:UI_PRIMITIVE_TESTS] [ARCH:TEST_STRATEGY_UI] [IMPL:TEST_WIDGETS]` breadcrumbs for auditability.

## Code Markers

- Test files include `[IMPL:TEST_WIDGETS] [ARCH:TEST_STRATEGY_UI] [REQ:UI_PRIMITIVE_TESTS]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `widget/listbox_test.go` - `[IMPL:TEST_WIDGETS]`
- [ ] `widget/widget_test.go` - `[IMPL:TEST_WIDGETS]`

Tests that must reference `[REQ:UI_PRIMITIVE_TESTS]`:
- [ ] `TestListBoxCursorClamping_REQ_UI_PRIMITIVE_TESTS`
- [ ] `TestListBoxScrollRate_REQ_UI_PRIMITIVE_TESTS`
- [ ] `TestListBoxColumnAdjust_REQ_UI_PRIMITIVE_TESTS`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-01 | — | ✅ Pass | Widget tests pass |

## Related Decisions

- Depends on: —
- See also: [ARCH:TEST_STRATEGY_UI], [REQ:UI_PRIMITIVE_TESTS]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
