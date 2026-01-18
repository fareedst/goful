# Linked Cursor Synchronization [IMPL:LINKED_CURSOR_SYNC]

**Status**: Active

## Cross-References

- [ARCH:LINKED_NAVIGATION]
- [REQ:LINKED_NAVIGATION]
- [REQ:MOUSE_CROSS_WINDOW_SYNC]
- [ARCH:MOUSE_CROSS_WINDOW_SYNC]

## Decision

Implement unified cursor synchronization that respects the Linked navigation toggle for both mouse and keyboard input.

## Implementation Details

### Mouse Click Handling (`app/goful.go`)

The `handleLeftClick` function conditionally calls `SetCursorByNameAll` based on linked state:

```go
// [IMPL:LINKED_CURSOR_SYNC] [IMPL:MOUSE_CROSS_WINDOW_SYNC]
// Sync cursor to same filename in all other windows when linked mode is ON
if g.IsLinkedNav() {
    ws.SetCursorByNameAll(filename)
}
```

- **Before**: Always called `ws.SetCursorByNameAll(filename)` regardless of linked state
- **After**: Only calls when `g.IsLinkedNav()` is true

### Keyboard Movement Wrappers (`app/goful.go`)

New methods wrap cursor movement with conditional sync:

- `MoveCursorLinked(amount int)` - wraps `Dir().MoveCursor()`
- `MoveTopLinked()` - wraps `Dir().MoveTop()`
- `MoveBottomLinked()` - wraps `Dir().MoveBottom()`
- `PageUpLinked()` - wraps `Dir().PageUp()`
- `PageDownLinked()` - wraps `Dir().PageDown()`

Each wrapper follows the same pattern:

```go
// [IMPL:LINKED_CURSOR_SYNC] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
func (g *Goful) MoveCursorLinked(amount int) {
    g.Dir().MoveCursor(amount)
    if g.IsLinkedNav() {
        g.Workspace().SetCursorByNameAll(g.File().Name())
    }
}
```

### Keymap Updates (`main.go`)

The `filerKeymap` function uses the new wrapper methods instead of direct calls:

```go
"C-n":  func() { g.MoveCursorLinked(1) },
"C-p":  func() { g.MoveCursorLinked(-1) },
"down": func() { g.MoveCursorLinked(1) },
"up":   func() { g.MoveCursorLinked(-1) },
"j":    func() { g.MoveCursorLinked(1) },
"k":    func() { g.MoveCursorLinked(-1) },
"C-d":  func() { g.MoveCursorLinked(5) },
"C-u":  func() { g.MoveCursorLinked(-5) },
"C-a":  func() { g.MoveTopLinked() },
"C-e":  func() { g.MoveBottomLinked() },
// etc.
```

## Token Coverage `[PROC:TOKEN_AUDIT]`

- `app/goful.go` includes `[IMPL:LINKED_CURSOR_SYNC] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]`
- `main.go` includes `[IMPL:LINKED_CURSOR_SYNC]` in keymap comments
- Tests reference `[REQ:LINKED_NAVIGATION]` in names/comments

## Validation Evidence (2026-01-18)

- `go test ./app/... -run "REQ_LINKED_NAVIGATION"` (darwin/arm64, Go 1.24.3) - 5 tests pass
- `/opt/homebrew/bin/bash ./scripts/validate_tokens.sh` → `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified 1491 token references across 78 files.`
- Unit tests in `app/linked_cursor_test.go`:
  - `TestMoveCursorLinked_REQ_LINKED_NAVIGATION`
  - `TestMoveCursorLinkedOff_REQ_LINKED_NAVIGATION`
  - `TestMoveTopLinked_REQ_LINKED_NAVIGATION`
  - `TestMoveBottomLinked_REQ_LINKED_NAVIGATION`
  - `TestLinkedCursorSyncMissingFile_REQ_LINKED_NAVIGATION`
- Manual verification: Toggle linked mode with `L` key, verify mouse and keyboard cursor movements sync/don't sync accordingly
