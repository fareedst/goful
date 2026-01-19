// Package diffstatus displays a persistent status line for difference search.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
package diffstatus

import (
	"github.com/fareedst/goful/look"
	"github.com/fareedst/goful/widget"
	"github.com/mattn/go-runewidth"
)

var status *statusWindow

// statusFn is a callback to get the current status text.
var statusFn func() string

// activeFn is a callback to check if diff search is active.
var activeFn func() bool

// customMessage holds a custom status message that takes priority over statusFn.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
var customMessage string

// statusWindow displays the diff search status.
type statusWindow struct {
	*widget.Window
}

// Init initializes the diff status window at the specified position.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func Init() {
	width, height := widget.Size()
	status = &statusWindow{
		Window: widget.NewWindow(0, height-3, width, 1),
	}
}

// SetStatusFn sets the callback to get the status text.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func SetStatusFn(fn func() string) {
	statusFn = fn
}

// SetActiveFn sets the callback to check if diff search is active.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func SetActiveFn(fn func() bool) {
	activeFn = fn
}

// SetMessage sets a custom status message that takes priority over the statusFn callback.
// Use this for persistent status updates like "Different: X - Y" that should not auto-dismiss.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func SetMessage(text string) {
	customMessage = text
}

// ClearMessage clears the custom status message, reverting to statusFn callback.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func ClearMessage() {
	customMessage = ""
}

// IsActive returns true if diff search is active and status should be shown.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func IsActive() bool {
	if activeFn == nil {
		return false
	}
	return activeFn()
}

// Draw the diff search status line.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func Draw() {
	if !IsActive() || status == nil {
		return
	}
	status.draw()
}

// Clear the diff status display.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func Clear() {
	if status != nil {
		status.Clear()
	}
}

// Resize the diff status window.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func Resize(x, y, width, height int) {
	if status != nil {
		status.Resize(x, y, width, height)
	}
}

func (w *statusWindow) draw() {
	w.Clear()

	// [IMPL:DIFF_SEARCH] Prioritize custom message over statusFn callback
	var text string
	if customMessage != "" {
		text = customMessage
	} else if statusFn != nil {
		text = statusFn()
	}
	if text == "" {
		text = "[DIFF SEARCH ACTIVE]"
	}

	x, y := w.LeftTop()
	text = runewidth.Truncate(text, w.Width(), "~")
	text = runewidth.FillRight(text, w.Width())
	widget.SetCells(x, y, text, look.Default().Reverse(true))
}
