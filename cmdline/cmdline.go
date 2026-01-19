// Package cmdline is a command line widget.
package cmdline

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/fareedst/goful/look"
	"github.com/fareedst/goful/util"
	"github.com/fareedst/goful/widget"
	"github.com/mattn/go-runewidth"
)

// Mode describes a cmdline mode.
type Mode interface {
	String() string // uses as a history map name
	Prompt() string // is displayed at the beginning
	Draw(*Cmdline)
	Run(*Cmdline)
}

// Cmdline is one line text box with a specified mode.
type Cmdline struct {
	*widget.TextBox
	filer      widget.Widget
	completion widget.Widget
	mode       Mode
	History    *History
}

var keymap func(*Cmdline) widget.Keymap

// Config sets the cmdline keymap function.
func Config(config func(*Cmdline) widget.Keymap) {
	keymap = config
}

// New creates a new cmdline with a specified mode and a history list box.
// These widget size based on the filer widget.
func New(m Mode, filer widget.Widget) *Cmdline {
	x, y := filer.LeftTop()
	width, height := filer.Width(), filer.Height()
	filer.ResizeRelative(0, 0, 0, -1)
	c := &Cmdline{
		TextBox:    widget.NewTextBox(x, y+height-1, width, 1),
		filer:      filer,
		completion: widget.Nil(),
		mode:       m,
		History:    &History{},
	}

	y = (y + height) * 2 / 3
	height -= y + 1
	c.History = NewHistory(x, y, width, height, c)
	c.Edithook = func() {
		c.History.update()
		c.History.MoveTop()
	}
	c.History.update()
	return c
}

// Next returns a completion widget.
func (c *Cmdline) Next() widget.Widget { return c.completion }

// Disconnect a completion widget.
func (c *Cmdline) Disconnect() { c.completion = widget.Nil() }

// StartCompletion starts a completion based on the cmdline text.
func (c *Cmdline) StartCompletion() {
	x, y := c.History.LeftTop()
	width, height := c.History.Width(), c.History.Height()
	comp := NewCompletion(x, y, width, height, c)
	if comp.IsEmpty() {
		return
	} else if len(comp.List()) == 1 {
		comp.InsertCompletion()
		return
	}
	c.completion = comp
}

// Resize the cmdline and the history list box.
func (c *Cmdline) Resize(x, y, width, height int) {
	c.TextBox.Resize(x, y+height-1, width, 1)
	y = (y + height) * 2 / 3
	height -= y + 1
	c.History.Resize(x, y, width, height)
	c.Next().Resize(x, y, width, height)
	c.filer.ResizeRelative(0, 0, 0, -1)
}

// ResizeRelative resizes relative to cmdline current sizes.
func (c *Cmdline) ResizeRelative(x, y, width, height int) {
	c.TextBox.ResizeRelative(x, y, width, height)
	c.History.ResizeRelative(x, y, width, height)
	c.Next().ResizeRelative(x, y, width, height)
}

// DrawLine draws the cmdline.
func (c *Cmdline) DrawLine() {
	c.Clear()
	x, y := c.LeftTop()
	x = widget.SetCells(x, y, c.mode.Prompt(), look.Prompt())
	w := c.Width() - runewidth.StringWidth(c.mode.Prompt()) - 2
	s := c.String()
	s = runewidth.Truncate(s, w, "")
	if c.Cursor() >= w {
		s = c.TextBeforeCursor()
		s = widget.TruncLeft(s, w, "~")
		x = widget.SetCells(x, y, s, look.Cmdline())
		widget.ShowCursor(x, y)
	} else {
		widget.SetCells(x, y, s, look.Cmdline())
		widget.ShowCursor(x+c.Cursor(), y)
	}
}

// Draw the cmdline and the completion or the histry list box
func (c *Cmdline) Draw() {
	c.mode.Draw(c)
	if !widget.IsNil(c.Next()) {
		c.Next().Draw()
	} else {
		c.History.Draw()
	}
}

// Input to the text box or widget keymaps.
func (c *Cmdline) Input(key string) {
	if !widget.IsNil(c.completion) {
		c.completion.Input(key)
	} else if cb, ok := keymap(c)[key]; ok {
		cb()
	} else {
		if utf8.RuneCountInString(key) == 1 {
			r, _ := utf8.DecodeRuneInString(key)
			c.InsertChar(r)
		}
	}
}

// Exit the cmdline and add a cmdline text to the history.
func (c *Cmdline) Exit() {
	c.History.add()
	widget.HideCursor()
	c.filer.ResizeRelative(0, 0, 0, 1)
	c.filer.Disconnect()
}

// Run the cmdline mode and add a cmdline text to the history.
func (c *Cmdline) Run() {
	c.History.add()
	c.mode.Run(c)
}

// [IMPL:HISTORY_CACHE_LIMIT] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]
// HistoryLimit controls maximum entries per mode to prevent unbounded memory growth.
// Default 1000 entries per mode; set to 0 for unlimited (not recommended for long sessions).
var HistoryLimit = 1000

var historyMap = map[string][]string{}

// HistoryError represents errors during history loading/saving operations.
// [IMPL:HISTORY_ERROR_HANDLING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]
type HistoryError struct {
	Path string
	Op   string // "load" or "save"
	Err  error
}

func (e *HistoryError) Error() string {
	return e.Op + " history " + e.Path + ": " + e.Err.Error()
}

func (e *HistoryError) Unwrap() error { return e.Err }

// IsFirstRun returns true if the error indicates a missing history file (first run).
// [IMPL:HISTORY_ERROR_HANDLING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]
func (e *HistoryError) IsFirstRun() bool {
	return errors.Is(e.Err, os.ErrNotExist)
}

// LoadHistory loads from a path and append to history maps of a key as the file name.
// Returns nil on success or when the file does not exist (first-run behavior).
// Returns a *HistoryError for actual IO failures (permissions, corrupt files, etc.)
// that should be surfaced to the user.
// [IMPL:HISTORY_ERROR_HANDLING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]
func LoadHistory(path string) error {
	path = util.ExpandPath(path)
	file, err := os.Open(path)
	if err != nil {
		// [IMPL:HISTORY_ERROR_HANDLING] First-run: missing file is not an error.
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		// Actual IO failure (permissions, etc.) — return structured error for caller to handle.
		return &HistoryError{Path: path, Op: "load", Err: err}
	}
	defer file.Close()

	history := make([]string, 0)
	rd := bufio.NewReader(file)
	for {
		line, _, err := rd.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &HistoryError{Path: path, Op: "load", Err: err}
		}
		history = append(history, string(line))
	}
	key := filepath.Base(path)
	historyMap[key] = history
	return nil
}

// SaveHistory saves the history to a path.
// Returns a *HistoryError for IO failures (permissions, disk full, etc.)
// that should be surfaced to the user.
// [IMPL:HISTORY_ERROR_HANDLING] [IMPL:HISTORY_CACHE_LIMIT] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]
func SaveHistory(path string) error {
	path = util.ExpandPath(path)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return &HistoryError{Path: path, Op: "save", Err: err}
	}

	file, err := os.Create(path)
	if err != nil {
		return &HistoryError{Path: path, Op: "save", Err: err}
	}
	defer file.Close()

	key := filepath.Base(path)
	if history, ok := historyMap[key]; ok {
		// [IMPL:HISTORY_CACHE_LIMIT] Trim before saving to compact persisted files
		history = trimHistory(history)
		historyMap[key] = history

		writer := bufio.NewWriter(file)
		for _, h := range history {
			if _, err := writer.WriteString(h + "\n"); err != nil {
				return &HistoryError{Path: path, Op: "save", Err: err}
			}
		}
		if err := writer.Flush(); err != nil {
			return &HistoryError{Path: path, Op: "save", Err: err}
		}
	}
	return nil
}

// History is the cmdline mode history.
type History struct {
	*widget.ListBox
	cmdline *Cmdline
	input   string // in the input
}

// NewHistory creates a new history list box.
func NewHistory(x, y, width, height int, cmdline *Cmdline) *History {
	lb := widget.NewListBox(x, y, width, height, "History")
	lb.SetLower(-1)
	lb.SetCursor(lb.Lower())
	return &History{lb, cmdline, ""}
}

func (h *History) update() {
	text := h.cmdline.String()
	name := h.cmdline.mode.String()

	h.input = text

	if history, ok := historyMap[name]; ok {
		h.ClearList()
		for i := len(history) - 1; i > -1; i-- {
			hist := history[i]
			if strings.Contains(hist, text) {
				h.AppendHighlightString(hist, text)
			}
		}
	}
}

func (h *History) add() {
	text := h.cmdline.String()
	mode := h.cmdline.mode.String()

	if text == "" || text[0] == ' ' {
		return
	}

	if history, ok := historyMap[mode]; ok {
		for i, hist := range history {
			if hist == text {
				history = append(history[:i], history[i+1:]...)
				break
			}
		}
		historyMap[mode] = trimHistory(append(history, text))
	} else {
		history := []string{}
		historyMap[mode] = trimHistory(append(history, text))
	}
}

// trimHistory enforces HistoryLimit by dropping oldest entries.
// [IMPL:HISTORY_CACHE_LIMIT] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]
func trimHistory(history []string) []string {
	if HistoryLimit <= 0 || len(history) <= HistoryLimit {
		return history
	}
	// Keep the most recent entries (tail of the slice)
	excess := len(history) - HistoryLimit
	return history[excess:]
}

// Delete a history content on the cursor.
func (h *History) Delete() {
	if h.Cursor() < 0 || h.Upper() < 1 {
		return
	}
	mode := h.cmdline.mode.String()
	if history, ok := historyMap[mode]; ok {
		for i := 0; i < len(history); i++ {
			if history[i] == h.CurrentContent().Name() {
				history = append(history[:i], history[i+1:]...)
				historyMap[mode] = history

				h.cmdline.SetText(h.input)
				h.update()
				h.AdjustCursor()
				break
			}
		}
	}
}

// MoveCursor moves the history list box cursor and sets a text to the cmdline.
func (h *History) MoveCursor(amount int) {
	h.ListBox.MoveCursor(amount)
	if h.Cursor() == h.Lower() {
		h.cmdline.SetText(h.input)
	} else {
		h.cmdline.SetText(h.CurrentContent().Name())
	}
}

// CursorDown downs the history list box cursor and sets a text to the cmdline.
func (h *History) CursorDown() {
	h.ListBox.CursorDown()
	if h.Cursor() == h.Lower() {
		h.cmdline.SetText(h.input)
	} else {
		h.cmdline.SetText(h.CurrentContent().Name())
	}
}

// CursorUp ups the history list box cursor and sets a text to the cmdline.
func (h *History) CursorUp() {
	h.ListBox.CursorUp()
	if h.Cursor() == h.Lower() {
		h.cmdline.SetText(h.input)
	} else {
		h.cmdline.SetText(h.CurrentContent().Name())
	}
}
