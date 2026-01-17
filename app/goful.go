// Package app is goful application components.
package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/anmitsu/goful/diffstatus"
	"github.com/anmitsu/goful/filer"
	"github.com/anmitsu/goful/help"
	"github.com/anmitsu/goful/info"
	"github.com/anmitsu/goful/menu"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/progress"
	"github.com/anmitsu/goful/widget"
	"github.com/gdamore/tcell/v2"
)

// doubleClickThreshold is the maximum time between clicks to count as a double-click.
// [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]
const doubleClickThreshold = 400 * time.Millisecond

// Goful represents a main application.
type Goful struct {
	*filer.Filer
	shell     func(cmd string) []string
	terminal  func(cmd string) []string
	next      widget.Widget
	event     chan tcell.Event
	interrupt chan int
	callback  chan func()
	task      chan int
	exit      bool
	linkedNav bool // [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION] Linked navigation mode state
	// Double-click state tracking [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]
	lastClickTime time.Time
	lastClickX    int
	lastClickY    int
}

// NewGoful creates a new goful client based recording a previous state.
func NewGoful(path string) *Goful {
	message.Init()
	info.Init()
	progress.Init()
	diffstatus.Init() // [IMPL:DIFF_SEARCH] Initialize diff status display
	width, height := widget.Size()
	goful := &Goful{
		Filer:     filer.NewFromState(path, 0, 0, width, height-2),
		shell:     nil,
		terminal:  nil,
		next:      widget.Nil(),
		event:     make(chan tcell.Event, 1),
		interrupt: make(chan int, 2),
		callback:  make(chan func()),
		task:      make(chan int, 1),
		exit:      false,
		linkedNav: true, // [IMPL:LINKED_NAVIGATION] Enabled by default
	}
	return goful
}

// ConfigShell sets a function that returns a shell name and options.
func (g *Goful) ConfigShell(config func(cmd string) []string) {
	g.shell = config
}

// ConfigTerminal sets a function that returns a terminal name and options.
func (g *Goful) ConfigTerminal(config func(cmd string) []string) {
	g.terminal = config
}

// ConfigFiler sets a keymap function for the filer.
func (g *Goful) ConfigFiler(f func(*Goful) widget.Keymap) {
	g.MergeKeymap(f(g))
}

// Next returns a next widget for drawing and input.
func (g *Goful) Next() widget.Widget { return g.next }

// Disconnect references to a next widget for exiting.
func (g *Goful) Disconnect() { g.next = widget.Nil() }

// Resize all widgets.
func (g *Goful) Resize(x, y, width, height int) {
	offset := 0
	progressActive := !progress.IsFinished()
	diffSearchActive := diffstatus.IsActive()

	if progressActive {
		offset = 2 // progress uses 2 rows (task + gauge)
	}
	// [IMPL:DIFF_SEARCH] Add offset for diff status line when active
	if diffSearchActive {
		offset++
	}
	g.Filer.Resize(x, y, width, height-2-offset)
	g.Next().Resize(x, y, width, height-2-offset)

	// Position status rows from bottom up:
	// - info at height-1 (always)
	// - message at height-2 (always)
	// - diffstatus at height-3 when active (or height-5 if progress is also active)
	// - progress at height-4 and height-3 (gauge) when active
	info.Resize(0, height-1, width, 1)
	message.Resize(0, height-2, width, 1)

	if progressActive && diffSearchActive {
		// Both active: diffstatus above progress
		// [IMPL:DIFF_SEARCH] Position diff status above progress gauge
		diffstatus.Resize(0, height-5, width, 1)
		progress.Resize(0, height-4, width, 1)
	} else if progressActive {
		// Only progress active
		progress.Resize(0, height-4, width, 1)
		diffstatus.Resize(0, height-3, width, 1) // Not drawn, but positioned
	} else if diffSearchActive {
		// Only diff search active
		// [IMPL:DIFF_SEARCH] Position diff status above message line
		diffstatus.Resize(0, height-3, width, 1)
		progress.Resize(0, height-4, width, 1) // Not drawn, but positioned
	} else {
		// Neither active
		diffstatus.Resize(0, height-3, width, 1)
		progress.Resize(0, height-4, width, 1)
	}
}

// Draw all widgets.
func (g *Goful) Draw() {
	// [IMPL:DIFF_SEARCH] Ensure layout is correct when diff search is active
	// This handles cases where navigation might not trigger a resize
	if diffstatus.IsActive() {
		width, height := widget.Size()
		g.Resize(0, 0, width, height)
	}
	g.Filer.Draw()
	g.Next().Draw()
	progress.Draw()
	diffstatus.Draw() // [IMPL:DIFF_SEARCH] Draw diff status line
	message.Draw()
	info.Draw(g.File())
}

// Input to a current widget.
func (g *Goful) Input(key string) {
	if !widget.IsNil(g.Next()) {
		g.Next().Input(key)
	} else {
		g.Filer.Input(key)
	}
}

// Menu runs a menu mode.
func (g *Goful) Menu(name string) {
	m, err := menu.New(name, g)
	if err != nil {
		message.Error(err)
		return
	}
	g.next = m
}

// Help opens the Help popup displaying the keystroke catalog.
// [IMPL:HELP_POPUP] [ARCH:HELP_WIDGET] [REQ:HELP_POPUP]
func (g *Goful) Help() {
	h := help.New(g)
	g.next = h
}

// Run the goful client.
func (g *Goful) Run() {
	message.Info("Welcome to goful")
	g.Workspace().ReloadAll()

	// TODO(goful-maintainers) [IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]:
	// add cancellation so the poller goroutine stops pushing into g.event after the UI exits;
	// the current infinite loop keeps running after Run returns, leaking a goroutine and busy-waiting on the channel.
	go func() {
		for {
			g.event <- widget.PollEvent()
		}
	}()

	for !g.exit {
		g.Draw()
		widget.Show()
		select {
		case ev := <-g.event:
			g.eventHandler(ev)
		case <-g.interrupt:
			<-g.interrupt
		case callback := <-g.callback:
			callback()
		}
	}
}

func (g *Goful) syncCallback(callback func()) {
	g.callback <- callback
}

func (g *Goful) eventHandler(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		key := widget.EventToString(ev)
		g.Input(key)
	case *tcell.EventResize:
		width, height := ev.Size()
		g.Resize(0, 0, width, height)
	case *tcell.EventMouse:
		// [IMPL:MOUSE_FILE_SELECT] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]
		g.mouseHandler(ev)
	}
}

// mouseHandler handles mouse events for file selection, focus switching, and scrolling.
// [IMPL:MOUSE_FILE_SELECT] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]
func (g *Goful) mouseHandler(ev *tcell.EventMouse) {
	x, y := ev.Position()
	buttons := ev.Buttons()

	// Handle modal widgets first - ignore mouse in modals for now
	if !widget.IsNil(g.Next()) {
		return
	}

	// Handle left click for file selection
	// [IMPL:MOUSE_FILE_SELECT] Left-click selects files and switches focus
	if buttons&tcell.Button1 != 0 {
		g.handleLeftClick(x, y)
		return
	}

	// Handle mouse wheel for scrolling
	// [IMPL:MOUSE_FILE_SELECT] Wheel scrolls the file list
	if buttons&tcell.WheelUp != 0 {
		g.handleWheelUp(x, y)
		return
	}
	if buttons&tcell.WheelDown != 0 {
		g.handleWheelDown(x, y)
		return
	}
}

// isDoubleClick checks if this click is a double-click based on timing and position.
// Returns true if the click occurs within the threshold at the same position as the last click.
// Updates the click tracking state regardless of result.
// [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]
func (g *Goful) isDoubleClick(x, y int) bool {
	now := time.Now()
	isDouble := now.Sub(g.lastClickTime) < doubleClickThreshold &&
		g.lastClickX == x && g.lastClickY == y
	g.lastClickTime = now
	g.lastClickX = x
	g.lastClickY = y
	return isDouble
}

// handleLeftClick processes a left mouse click at (x, y).
// Switches focus if clicking in an unfocused window and moves cursor to the clicked file.
// Detects double-clicks and dispatches to appropriate handler.
// [IMPL:MOUSE_FILE_SELECT] [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT] [REQ:MOUSE_DOUBLE_CLICK]
func (g *Goful) handleLeftClick(x, y int) {
	ws := g.Workspace()
	dir, idx := ws.DirectoryAt(x, y)
	if dir == nil {
		return
	}

	// Switch focus if clicking in unfocused window
	if idx != ws.Focus {
		ws.SetFocus(idx)
	}

	// Convert Y to file index and move cursor
	fileIdx := dir.FileIndexAtY(y)
	if fileIdx >= 0 {
		dir.SetCursor(fileIdx)
	}

	// Check for double-click after selection
	// [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]
	if g.isDoubleClick(x, y) && fileIdx >= 0 {
		file := dir.File()
		if file.IsDir() {
			g.handleDoubleClickDir(dir)
		} else {
			g.handleDoubleClickFile(dir)
		}
	}
}

// handleDoubleClickDir navigates into a directory, respecting linked mode.
// When linked mode is ON, navigates all windows to matching subdirectory.
// [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]
func (g *Goful) handleDoubleClickDir(dir *filer.Directory) {
	if g.IsLinkedNav() {
		name := dir.File().Name()
		// Navigate other directories but DON'T rebuild index yet
		navigated, skipped := g.Workspace().ChdirAllToSubdirNoRebuild(name)
		// [IMPL:LINKED_NAVIGATION_AUTO_DISABLE] Auto-disable if any window couldn't navigate
		if skipped > 0 {
			g.SetLinkedNav(false)
			message.Infof("linked navigation disabled: %d window(s) missing '%s'", skipped, name)
		}
		_ = navigated
	}
	dir.EnterDir()
	// Rebuild comparison index AFTER all directories have navigated
	if g.IsLinkedNav() {
		g.Workspace().RebuildComparisonIndex()
	}
}

// handleDoubleClickFile opens a file, and opens same-named files in all windows when linked.
// When linked mode is ON, moves cursor to same-named file in all windows before triggering open.
// [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]
func (g *Goful) handleDoubleClickFile(dir *filer.Directory) {
	filename := dir.File().Name()

	if g.IsLinkedNav() {
		// Move cursor to same-named file in all windows
		for _, d := range g.Workspace().Dirs {
			if d.FindFileByName(filename) != nil {
				d.SetCursorByName(filename)
			}
		}
	}
	// Trigger open action (uses extmap via C-m or o key)
	g.Input("C-m")
}

// handleWheelUp scrolls the directory under the mouse cursor up.
// [IMPL:MOUSE_FILE_SELECT] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]
func (g *Goful) handleWheelUp(x, y int) {
	ws := g.Workspace()
	dir, _ := ws.DirectoryAt(x, y)
	if dir == nil {
		// Default to focused directory if click outside any directory
		dir = ws.Dir()
	}
	dir.MoveCursor(-3)
}

// handleWheelDown scrolls the directory under the mouse cursor down.
// [IMPL:MOUSE_FILE_SELECT] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]
func (g *Goful) handleWheelDown(x, y int) {
	ws := g.Workspace()
	dir, _ := ws.DirectoryAt(x, y)
	if dir == nil {
		// Default to focused directory if click outside any directory
		dir = ws.Dir()
	}
	dir.MoveCursor(3)
}

// ToggleLinkedNav toggles the linked navigation mode and returns the new state.
// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
func (g *Goful) ToggleLinkedNav() bool {
	g.linkedNav = !g.linkedNav
	return g.linkedNav
}

// IsLinkedNav returns true if linked navigation mode is enabled.
// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
func (g *Goful) IsLinkedNav() bool {
	return g.linkedNav
}

// SetLinkedNav sets the linked navigation mode state.
// [IMPL:LINKED_NAVIGATION_AUTO_DISABLE] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
func (g *Goful) SetLinkedNav(enabled bool) {
	g.linkedNav = enabled
}

// DiffSearchStatus returns the current diff search status text for the header.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (g *Goful) DiffSearchStatus() string {
	state := g.Workspace().DiffSearchState()
	if state == nil {
		return ""
	}
	return state.StatusText()
}

// IsDiffSearchActive returns true if a difference search is currently active.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (g *Goful) IsDiffSearchActive() bool {
	return g.Workspace().IsDiffSearchActive()
}

// StartDiffSearch begins a new difference search across all windows.
// Records initial directories and finds the first difference.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (g *Goful) StartDiffSearch() {
	ws := g.Workspace()
	if len(ws.Dirs) < 2 {
		message.Errorf("Difference search requires at least 2 windows")
		return
	}

	// Start a new search session
	ws.StartDiffSearch()
	state := ws.DiffSearchState()
	state.SetSearching(true)
	state.SetCurrentPath(ws.Dir().Path)

	// Resize to allocate space for the diff status line
	width, height := widget.Size()
	g.Resize(0, 0, width, height)

	// Find the first difference
	g.findNextDiff("")
}

// ContinueDiffSearch continues the difference search from the current cursor position.
// Skips the file at cursor and finds the next difference.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (g *Goful) ContinueDiffSearch() {
	ws := g.Workspace()
	if !ws.IsDiffSearchActive() {
		message.Errorf("No active difference search. Use start diff search first.")
		return
	}

	// Set searching state
	state := ws.DiffSearchState()
	state.SetSearching(true)
	state.SetCurrentPath(ws.Dir().Path)

	// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
	// Get the name to skip from lastDiffName if available, otherwise from cursor position.
	// This ensures we continue from the correct position even if the cursor couldn't be set
	// (e.g., when a subdirectory is missing in some windows).
	var startAfter string
	if state.LastDiffName != "" {
		// Remove trailing "/" if present (for subdirectories)
		startAfter = strings.TrimSuffix(state.LastDiffName, "/")
	} else {
		// Fallback to cursor position
		startAfter = g.File().Name()
		if startAfter == ".." {
			startAfter = ""
		}
	}

	// Continue from the next entry
	g.findNextDiff(startAfter)
}

// findNextDiff is the core search loop that finds the next difference.
// Uses the TreeWalker to handle traversal logic, keeping TUI concerns separate.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (g *Goful) findNextDiff(startAfter string) {
	ws := g.Workspace()
	state := ws.DiffSearchState()
	nav := filer.NewWorkspaceNavigator(ws)

	// Start periodic UI refresh goroutine (updates once per second)
	// [IMPL:DIFF_SEARCH] Periodic refresh for progress display
	quit := make(chan struct{})
	defer close(quit)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				g.Draw()
				widget.Show()
			case <-quit:
				return
			}
		}
	}()

	// Create and run the tree walker
	walker := filer.NewTreeWalker(nav, state, startAfter)
	step := walker.Run(func() {
		state.IncrementFilesChecked()
		state.SetCurrentPath(nav.CurrentPath())
	})

	// Handle the result
	switch step.Type {
	case filer.StepFoundDiff:
		// Found a difference - pause search and record result
		// [IMPL:DIFF_SEARCH] Route status to dedicated diffstatus row, not ephemeral message
		state.SetLastDiff(step.Name, step.Reason)
		ws.SetCursorByNameAll(strings.TrimSuffix(step.Name, "/"))
		diffstatus.SetMessage(fmt.Sprintf("Different: %s - %s", step.Name, step.Reason))
	case filer.StepComplete:
		// Search complete - all differences have been found
		// [IMPL:DIFF_SEARCH] Use ephemeral message for completion since diffstatus row disappears
		ws.ClearDiffSearch()
		diffstatus.ClearMessage()
		message.Info("Difference search complete - all differences found")
		// Resize to reclaim space from diff status line
		width, height := widget.Size()
		g.Resize(0, 0, width, height)
	}
}

// SetBorderStyle sets the filer border style.
func (g *Goful) SetBorderStyle(style widget.BorderStyle) {
	filer.SetBorderStyle(style)
	for _, ws := range g.Workspaces {
		for _, d := range ws.Dirs {
			d.SetBorderStyle(style)
		}
	}
}
