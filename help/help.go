// Package help provides a Help popup widget that displays the keystroke catalog.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
package help

import (
	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/widget"
)

// keystrokeCatalog contains all key bindings displayed in the help popup.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
var keystrokeCatalog = []string{
	"=== Navigation ===",
	"C-n, down, j         Move cursor down",
	"C-p, up, k           Move cursor up",
	"C-a, home, ^         Move cursor top",
	"C-e, end, $          Move cursor bottom",
	"C-f, C-i, right, l   Move cursor right (next pane)",
	"C-b, left, h         Move cursor left (prev pane)",
	"C-d                  More move cursor down",
	"C-u                  More move cursor up",
	"C-v, pgdn            Page down",
	"M-v, pgup            Page up",
	"M-n                  Scroll down",
	"M-p                  Scroll up",
	"",
	"=== Directory ===",
	"C-h, backspace, u    Change to parent directory",
	"~                    Change to home directory",
	"\\                    Change to root directory",
	"w                    Change to neighbor directory",
	"d                    Change directory (prompt)",
	"C-o                  Create directory window",
	"C-w                  Close directory window",
	"C-l                  Reload all",
	"",
	"=== Workspace ===",
	"M-f                  Move next workspace",
	"M-b                  Move previous workspace",
	"M-C-o                Create workspace",
	"M-C-w                Close workspace",
	"M-W                  Change workspace title",
	"",
	"=== Selection ===",
	"space                Toggle mark",
	"M-=                  Invert mark",
	"",
	"=== File Operations ===",
	"C-m, o               Open file/directory",
	"i                    Open by pager",
	"n                    Make file",
	"K                    Make directory",
	"c                    Copy",
	"m                    Move",
	"r                    Rename",
	"R                    Bulk rename by regexp",
	"D                    Remove",
	"C                    Copy All (to all panes)",
	"M                    Move All (to all panes)",
	"",
	"=== Search & Filter ===",
	"f, /                 Find (filter)",
	"g                    Glob",
	"G                    Glob recursive",
	"C-g, C-[             Cancel/Reset",
	"",
	"=== View & Compare ===",
	"s                    Sort menu",
	"v                    View menu",
	"E                    Toggle filename excludes",
	"`                    Toggle comparison colors",
	"=                    Calculate file digest",
	"L, M-l               Toggle linked navigation",
	"[                    Start difference search",
	"]                    Continue difference search",
	"",
	"=== Menus & Commands ===",
	"b                    Bookmark menu",
	"e                    Editor menu",
	"x                    Command menu",
	"X                    External command menu",
	";                    Shell",
	":                    Shell suspend",
	"",
	"=== Application ===",
	"?                    Help (this popup)",
	"q, Q                 Quit",
	"",
	"Press ?, q, C-g, or Esc to close",
}

// Help is a popup widget displaying the keystroke catalog.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
type Help struct {
	*widget.ListBox
	filer widget.Widget
}

// New creates a new Help popup based on filer widget sizes.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
func New(filer widget.Widget) *Help {
	screenWidth, screenHeight := widget.Size()

	// Size the popup to ~80% of screen, centered
	width := screenWidth * 80 / 100
	height := screenHeight * 80 / 100
	if width < 40 {
		width = screenWidth - 4
	}
	if height < 10 {
		height = screenHeight - 4
	}
	// Ensure we don't exceed content size + borders
	maxHeight := len(keystrokeCatalog) + 2
	if height > maxHeight {
		height = maxHeight
	}

	// Center the popup
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2

	h := &Help{
		ListBox: widget.NewListBox(x, y, width, height, "Help - Keystrokes"),
		filer:   filer,
	}

	// Populate the list with keystroke entries
	for _, entry := range keystrokeCatalog {
		h.AppendString(entry)
	}

	h.SetBorderStyle(widget.AllBorder)

	return h
}

// Resize the help window.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
func (h *Help) Resize(x, y, width, height int) {
	screenWidth, screenHeight := widget.Size()

	// Recalculate dimensions
	w := screenWidth * 80 / 100
	ht := screenHeight * 80 / 100
	if w < 40 {
		w = screenWidth - 4
	}
	if ht < 10 {
		ht = screenHeight - 4
	}
	maxHeight := len(keystrokeCatalog) + 2
	if ht > maxHeight {
		ht = maxHeight
	}

	// Center
	newX := (screenWidth - w) / 2
	newY := (screenHeight - ht) / 2

	h.ListBox.Resize(newX, newY, w, ht)
}

// Draw the help popup with a distinct background.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
func (h *Help) Draw() {
	h.ListBox.Draw()
}

// Input handles keyboard input for the help popup.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
func (h *Help) Input(key string) {
	switch key {
	// Exit keys - toggle with ?, standard exit with q/C-g/Esc
	case "?", "q", "Q", "C-g", "C-[":
		h.Exit()
	// Navigation
	case "C-n", "down", "j":
		h.MoveCursor(1)
	case "C-p", "up", "k":
		h.MoveCursor(-1)
	case "C-v", "pgdn":
		h.PageDown()
	case "M-v", "pgup":
		h.PageUp()
	case "C-a", "home", "^":
		h.MoveTop()
	case "C-e", "end", "$":
		h.MoveBottom()
	case "M-n":
		h.Scroll(1)
	case "M-p":
		h.Scroll(-1)
	}
}

// Exit closes the help popup and returns to the filer.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
func (h *Help) Exit() {
	h.filer.Disconnect()
}

// Next implements widget.Widget.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
func (h *Help) Next() widget.Widget {
	return widget.Nil()
}

// Disconnect implements widget.Widget.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
func (h *Help) Disconnect() {}

// Custom content drawer for help entries that applies styling.
type helpContent struct {
	text string
}

func (c *helpContent) Name() string { return c.text }

func (c *helpContent) Draw(x, y, width int, focus bool) {
	style := look.Default()
	if focus {
		style = style.Reverse(true)
	}

	// Pad or truncate to fit width
	text := c.text
	if len(text) > width {
		text = text[:width-1] + "~"
	}
	for len(text) < width {
		text += " "
	}

	widget.SetCells(x, y, text, style)
}
