// Package help provides a Help popup widget that displays the keystroke catalog.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]
package help

import (
	"strings"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/widget"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
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
	"",
	"=== Multi-Pane Operations ===",
	"C                    Copy All (to all panes)",
	"M                    Move All (to all panes)",
	"",
	"=== Sync Operations ===",
	"S                    Sync mode prefix",
	"  S then c           Copy same-named files (new name)",
	"  S then d           Delete same-named files",
	"  S then r           Rename same-named files",
	"  S then !           Toggle ignore-failures mode",
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

// keyColumnWidth is the fixed width for the key binding column.
// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]
const keyColumnWidth = 21

// Help is a popup widget displaying the keystroke catalog.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
type Help struct {
	*widget.ListBox
	filer widget.Widget
}

// New creates a new Help popup based on filer widget sizes.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]
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

	// Populate the list with styled help entries
	// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]
	for _, entry := range keystrokeCatalog {
		h.AppendList(newHelpEntry(entry))
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

// Draw the help popup with styled borders.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]
func (h *Help) Draw() {
	if h.Upper() < 1 {
		return
	}
	h.AdjustCursor()
	h.AdjustOffset()
	h.clearWithBackground()
	h.drawColoredBorder()
	h.drawColoredHeader()
	h.drawScrollbar()
	h.drawContents()
}

// clearWithBackground clears the window with the help description background.
// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]
func (h *Help) clearWithBackground() {
	style := look.HelpDesc()
	x, y := h.LeftTop()
	xend, yend := h.RightBottom()
	for row := y; row <= yend; row++ {
		for col := x; col <= xend; col++ {
			widget.SetCells(col, row, " ", style)
		}
	}
}

// drawColoredBorder draws the border using HelpBorder style.
// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]
func (h *Help) drawColoredBorder() {
	style := look.HelpBorder()
	x, y := h.LeftTop()
	xend, yend := h.RightBottom()

	// Horizontal lines (top and bottom)
	for col := x; col <= xend; col++ {
		widget.SetCells(col, y, string(tcell.RuneHLine), style)
		widget.SetCells(col, yend, string(tcell.RuneHLine), style)
	}

	// Vertical lines (left and right)
	for row := y + 1; row < yend; row++ {
		widget.SetCells(x, row, string(tcell.RuneVLine), style)
		widget.SetCells(xend, row, string(tcell.RuneVLine), style)
	}

	// Corners
	widget.SetCells(x, y, string(tcell.RuneULCorner), style)
	widget.SetCells(xend, y, string(tcell.RuneURCorner), style)
	widget.SetCells(x, yend, string(tcell.RuneLLCorner), style)
	widget.SetCells(xend, yend, string(tcell.RuneLRCorner), style)
}

// drawColoredHeader draws the title with HelpHeader style.
// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]
func (h *Help) drawColoredHeader() {
	title := h.Title()
	x, y := h.LeftTop()
	style := look.HelpHeader()
	widget.SetCells(x+1, y, " "+title+" ", style)
}

// drawScrollbar draws the scrollbar indicator.
// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]
func (h *Help) drawScrollbar() {
	if h.Upper() <= (h.Height()-2)*h.Column() {
		return
	}
	height := h.Height() - 2
	rowCol := (h.Height() - 2) * h.Column()
	offset := int(float64(h.Offset()) / float64(h.Upper()-rowCol) * float64(height))
	if offset > height-1 {
		offset = height - 1
	}

	x, y := h.RightTop()
	y++
	style := look.HelpBorder()
	for i := 0; i < height; i++ {
		if i == offset {
			widget.SetCells(x, y+i, "█", style)
		} else {
			widget.SetCells(x, y+i, "│", style)
		}
	}
}

// drawContents draws the list entries.
// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]
func (h *Help) drawContents() {
	width, height := h.Width()-2, h.Height()-2
	shift := 2 // AllBorder has 2 char offset
	colwidth := width/h.Column() - shift + 1
	row, col := 1, 0

	for i := h.Offset(); i < h.Upper(); i++ {
		if col >= h.Column() {
			col = 0
			row++
			if row > height {
				break
			}
		}
		x, y := h.LeftTop()
		x += col*colwidth + shift
		y += row
		focus := i == h.Cursor()
		h.List()[i].Draw(x, y, colwidth, focus)
		col++
	}
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

// helpEntry is a custom content drawer for help entries that applies styling.
// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]
type helpEntry struct {
	text     string
	isHeader bool
}

// newHelpEntry creates a new help entry with type detection.
// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]
func newHelpEntry(text string) *helpEntry {
	return &helpEntry{
		text:     text,
		isHeader: strings.HasPrefix(text, "==="),
	}
}

// Name returns the entry text for ListBox compatibility.
// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]
func (e *helpEntry) Name() string { return e.text }

// Draw renders the help entry with appropriate styling.
// [IMPL:HELP_STYLING] [ARCH:HELP_STYLING] [REQ:HELP_POPUP_STYLING]
func (e *helpEntry) Draw(x, y, width int, focus bool) {
	if e.text == "" {
		// Blank line - just fill with background
		style := look.HelpDesc()
		if focus {
			style = style.Reverse(true)
		}
		widget.SetCells(x, y, strings.Repeat(" ", width), style)
		return
	}

	if e.isHeader {
		// Section header - use HelpHeader style
		style := look.HelpHeader()
		if focus {
			style = style.Reverse(true)
		}
		text := runewidth.Truncate(e.text, width, "~")
		text = runewidth.FillRight(text, width)
		widget.SetCells(x, y, text, style)
		return
	}

	// Key binding entry - split into key (left) and description (right)
	keyPart := e.text
	descPart := ""
	if len(e.text) > keyColumnWidth {
		keyPart = e.text[:keyColumnWidth]
		descPart = e.text[keyColumnWidth:]
	}

	// Draw key part with HelpKey style
	keyStyle := look.HelpKey()
	if focus {
		keyStyle = keyStyle.Reverse(true)
	}
	keyText := runewidth.Truncate(keyPart, keyColumnWidth, "")
	keyText = runewidth.FillRight(keyText, keyColumnWidth)
	pos := widget.SetCells(x, y, keyText, keyStyle)

	// Draw description part with HelpDesc style
	descStyle := look.HelpDesc()
	if focus {
		descStyle = descStyle.Reverse(true)
	}
	remaining := width - keyColumnWidth
	if remaining > 0 {
		descText := runewidth.Truncate(descPart, remaining, "~")
		descText = runewidth.FillRight(descText, remaining)
		widget.SetCells(pos, y, descText, descStyle)
	}
}
