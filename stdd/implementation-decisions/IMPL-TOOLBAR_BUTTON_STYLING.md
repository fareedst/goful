# [IMPL:TOOLBAR_BUTTON_STYLING] Toolbar Button Styling

**Cross-References**: [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_BUTTON_STYLING]  
**Status**: Active  
**Created**: 2026-01-18  
**Last Updated**: 2026-01-18

---

## Decision

Implement toolbar button styling that distinguishes action buttons (plain text) from toggle buttons (color-coded states). The parent button changes from `[^]` with reverse highlight to plain `^`, and toggle buttons use Lime+Bold foreground when ON instead of reverse video.

## Rationale

- Reverse style was used for all buttons, but this conflates action buttons with toggle state indicators
- Action buttons (one-time operations) should appear as plain text to indicate they have no persistent state
- Toggle buttons should use color to indicate their ON/OFF state, providing consistent visual feedback
- Lime+Bold is consistent with existing "active" indicators in the codebase (messageInfo, cmdlineCommand, executable)
- Reducing the parent button from `[^]` to `^` saves horizontal space and clarifies it is an action

## Implementation Approach

### Changes to `filer/filer.go` `drawHeader()`

**Parent Button** (action):
```go
// Before
parentBtn := "[^]"
x = widget.SetCells(x, y, parentBtn, look.Default().Reverse(true))

// After
// [IMPL:TOOLBAR_BUTTON_STYLING] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_BUTTON_STYLING]
// Action button: plain text, no highlight
parentBtn := "^"
x = widget.SetCells(x, y, parentBtn, look.Default())
```

**Linked Toggle Button**:
```go
// Before
linkedStyle := look.Default()
if linkedNavIndicator != nil && linkedNavIndicator() {
    linkedStyle = linkedStyle.Reverse(true)
}

// After
// [IMPL:TOOLBAR_BUTTON_STYLING] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_BUTTON_STYLING]
// Toggle button: Lime+Bold when ON, Yellow when OFF
linkedStyle := look.Default().Foreground(tcell.ColorYellow)
if linkedNavIndicator != nil && linkedNavIndicator() {
    linkedStyle = look.Default().Foreground(tcell.ColorLime).Bold(true)
}
```

**Ignore-Failures Toggle Button**:
```go
// Before
ignoreStyle := look.Default()
if syncIgnoreFailuresIndicator != nil && syncIgnoreFailuresIndicator() {
    ignoreStyle = ignoreStyle.Reverse(true)
}

// After
// [IMPL:TOOLBAR_BUTTON_STYLING] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_BUTTON_STYLING]
// Toggle button: Lime+Bold when ON, Yellow when OFF
ignoreStyle := look.Default().Foreground(tcell.ColorYellow)
if syncIgnoreFailuresIndicator != nil && syncIgnoreFailuresIndicator() {
    ignoreStyle = look.Default().Foreground(tcell.ColorLime).Bold(true)
}
```

### Test Updates in `filer/toolbar_test.go`

Update parent button bounds from width 3 (`[^]`) to width 1 (`^`):
```go
// Before
toolbarButtons["parent"] = toolbarBounds{x1: 0, y: 0, x2: 2} // "[^]"

// After
toolbarButtons["parent"] = toolbarBounds{x1: 0, y: 0, x2: 0} // "^"
```

## Code Markers

Files that must carry `[IMPL:TOOLBAR_BUTTON_STYLING]` annotations:
- `filer/filer.go` - `drawHeader()` button styling changes

Tests that must reference `[REQ:TOOLBAR_BUTTON_STYLING]`:
- `filer/toolbar_test.go` - Updated bounds tests

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [x] `filer/filer.go` `drawHeader()` - `[IMPL:TOOLBAR_BUTTON_STYLING] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_BUTTON_STYLING]`

Tests that reference requirements:
- [x] Test names/comments include `[REQ:TOOLBAR_BUTTON_STYLING]`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-18 | — | ✅ Pass | Token validation: 1784 refs across 85 files |

## Related Decisions

- Updates: [IMPL:TOOLBAR_PARENT_BUTTON] (button text and style)
- Updates: [IMPL:TOOLBAR_LINKED_TOGGLE] (toggle styling)
- Updates: [IMPL:TOOLBAR_IGNORE_FAILURES] (toggle styling)
- See also: [ARCH:TOOLBAR_LAYOUT]

---

*Created on 2026-01-18*
