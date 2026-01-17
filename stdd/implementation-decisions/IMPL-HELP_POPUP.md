# [IMPL:HELP_POPUP] Help Popup Implementation

**Cross-References**: [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]  
**Status**: Active  
**Created**: 2026-01-12  
**Last Updated**: 2026-01-17

---

## Decision

Create a `help` package with a scrollable keystroke catalog popup based on `widget.ListBox`.

## Rationale

- Provides users with quick access to keystroke documentation without leaving the application
- Reusing the existing `widget.ListBox` pattern ensures consistent scrolling behavior and styling
- The popup model (connect via `g.next`, disconnect on exit) integrates naturally with goful's widget chain

## Implementation Approach

### Create `help/help.go`

- `Help` struct embeds `*widget.ListBox` and stores a reference to the parent filer widget
- `New(filer widget.Widget)` constructor creates a popup sized to ~80% of screen dimensions
- The keystroke catalog is defined as a Go slice `keystrokeCatalog` containing all key bindings
- Each entry is formatted as "key  description" and added to the ListBox via `AppendString()`
- `Input(key string)` dispatches to navigation methods (`MoveCursor`, `PageDown`, etc.) or `Exit()`
- Exit triggers: `?` (toggle), `q`, `C-g`, `C-[` (Escape)
- `Exit()` calls `filer.Disconnect()` to return to normal filer operation

### Add `Help()` method to `app/goful.go`

- Creates a new `help.New(g)` widget and assigns it to `g.next`

### Add keystroke binding in `main.go`

- `g.AddKeymap("?", func() { g.Help() })`

## Code Markers

- `help/help.go` includes `// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]`
- `app/goful.go` `Help()` method includes `// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]`
- `main.go` keystroke binding includes `// [IMPL:HELP_POPUP] [REQ:HELP_POPUP]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `help/help.go` - all functions
- [ ] `app/goful.go` - Help() method
- [ ] `main.go` - keystroke binding

Tests that must reference `[REQ:HELP_POPUP]`:
- [ ] Manual verification for now; future automated tests can reference `[REQ:HELP_POPUP]`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| — | — | ⏳ Pending | To be captured after implementation |

## Related Decisions

- Depends on: —
- See also: [ARCH:HELP_WIDGET], [REQ:HELP_POPUP], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
