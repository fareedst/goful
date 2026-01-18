# [IMPL:HELP_STYLING] Help Popup Styling and Mouse Scroll

**Cross-References**: [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]  
**Status**: Active  
**Created**: 2026-01-18  
**Last Updated**: 2026-01-18

---

## Decision

Implement help popup visual styling with theme-aware colors and mouse wheel scroll support by extending the look package and creating a custom content drawer.

## Rationale

- Centralizing styles in `look` package follows existing patterns (`look.Directory()`, `look.Marked()`, etc.)
- Custom drawer allows per-entry styling without modifying `widget.ListBox` internals
- Mouse scroll forwarding reuses existing wheel handling patterns from `mouseHandler()`
- Provides visual hierarchy distinguishing section headers from key bindings

## Implementation Approach

### 1. Style Variables (`look/look.go`)

Add style variables and accessors following existing patterns:

```go
var (
    helpBorder tcell.Style
    helpHeader tcell.Style
    helpKey    tcell.Style
    helpDesc   tcell.Style
)

func HelpBorder() tcell.Style { return helpBorder }
func HelpHeader() tcell.Style { return helpHeader }
func HelpKey() tcell.Style { return helpKey }
func HelpDesc() tcell.Style { return helpDesc }

func SetHelpBorder(s tcell.Style) { helpBorder = s }
func SetHelpHeader(s tcell.Style) { helpHeader = s }
func SetHelpKey(s tcell.Style) { helpKey = s }
func SetHelpDesc(s tcell.Style) { helpDesc = s }
```

### 2. Theme Configuration

Configure in each theme function (`setDefault`, `setMidnight`, `setBlack`, `setWhite`):

| Theme    | Border | Header       | Key          | Description |
|----------|--------|--------------|--------------|-------------|
| Default  | Cyan   | Yellow+Bold  | Lime+Bold    | Default     |
| Midnight | Aqua   | Yellow+Bold  | Yellow+Bold  | White       |
| Black    | Aqua   | Yellow+Bold  | Lime+Bold    | White       |
| White    | Navy   | Olive+Bold   | Green+Bold   | Black       |

### 3. helpEntry Drawer (`help/help.go`)

Custom drawer implementing `widget.Drawer`:

```go
type helpEntry struct {
    text     string
    isHeader bool  // Detected from "===" prefix
}

func (e *helpEntry) Name() string { return e.text }

func (e *helpEntry) Draw(x, y, width int, focus bool) {
    if e.isHeader {
        // Render with look.HelpHeader()
    } else if e.text == "" {
        // Blank line
    } else {
        // Parse key (first 21 chars) and description (rest)
        // Render key with look.HelpKey(), description with look.HelpDesc()
    }
    if focus {
        style = style.Reverse(true)
    }
}
```

### 4. Colored Border Drawing (`help/help.go`)

Override `Draw()` method to render borders with `look.HelpBorder()`:

```go
func (h *Help) Draw() {
    h.Clear()
    h.drawColoredBorder()
    h.drawHeader()
    h.drawScrollbar()
    h.drawContents()
}

func (h *Help) drawColoredBorder() {
    // Use widget.SetCells with look.HelpBorder() for border characters
}
```

### 5. Mouse Wheel Forwarding (`app/goful.go`)

Update `mouseHandler()` to forward wheel events to modal widgets:

```go
func (g *Goful) mouseHandler(ev *tcell.EventMouse) {
    x, y := ev.Position()
    buttons := ev.Buttons()

    // [IMPL:HELP_STYLING] Forward wheel events to modal widgets
    if !widget.IsNil(g.Next()) {
        if buttons&tcell.WheelUp != 0 {
            g.Next().Input("M-p") // Scroll up
            return
        }
        if buttons&tcell.WheelDown != 0 {
            g.Next().Input("M-n") // Scroll down
            return
        }
        return // Ignore other mouse events for modals
    }
    // ... existing code ...
}
```

## Code Markers

- `look/look.go`: `// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]`
- `help/help.go`: `// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]`
- `app/goful.go`: `// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [x] `look/look.go` - help style variables and accessors
- [x] `help/help.go` - helpEntry.Draw(), Help.Draw(), drawColoredBorder()
- [x] `app/goful.go` - mouseHandler() wheel forwarding section

Tests that must reference `[REQ:HELP_POPUP_STYLING]`:
- [ ] `help/help_test.go` - TestHelpEntry_Draw_REQ_HELP_POPUP_STYLING (future enhancement)

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-18 | — | ✅ Pass | `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1708 token references across 79 files.` |

## Related Decisions

- Depends on: [ARCH:HELP_WIDGET], [IMPL:HELP_POPUP]
- See also: [REQ:HELP_POPUP_STYLING], [REQ:MOUSE_FILE_SELECT]

---

*Created on 2026-01-18*
