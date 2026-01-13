// Package app is goful application components.
package app

import (
	"fmt"
	"time"

	"github.com/anmitsu/goful/diffstatus"
	"github.com/anmitsu/goful/filer"
	"github.com/anmitsu/goful/info"
	"github.com/anmitsu/goful/menu"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/progress"
	"github.com/anmitsu/goful/widget"
	"github.com/gdamore/tcell/v2"
)

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
	}
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

	// Get the current file name to skip
	currentFile := g.File().Name()
	if currentFile == ".." {
		currentFile = ""
	}

	// Set searching state
	state := ws.DiffSearchState()
	state.SetSearching(true)
	state.SetCurrentPath(ws.Dir().Path)

	// Continue from the next entry
	g.findNextDiff(currentFile)
}

// findNextDiff is the core search loop that finds the next difference.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (g *Goful) findNextDiff(startAfter string) {
	ws := g.Workspace()
	state := ws.DiffSearchState()

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

	for {
		// Update current path being searched
		state.SetCurrentPath(ws.Dir().Path)
		state.IncrementFilesChecked()

		// First, check files for differences
		result := filer.FindNextDifference(ws.Dirs, startAfter, true)
		if result.Found {
			// Found a file difference - pause search and record result
			// [IMPL:DIFF_SEARCH] Route status to dedicated diffstatus row, not ephemeral message
			state.SetLastDiff(result.Name, result.Reason)
			ws.SetCursorByNameAll(result.Name)
			diffstatus.SetMessage(fmt.Sprintf("Different: %s - %s", result.Name, result.Reason))
			return
		}

		// No file differences, check subdirectories
		subdirResult := filer.FindNextDifference(ws.Dirs, startAfter, false)
		if subdirResult.Found && subdirResult.IsDir {
			// Found a subdir that differs (missing in some window)
			// [IMPL:DIFF_SEARCH] Route status to dedicated diffstatus row, not ephemeral message
			state.SetLastDiff(subdirResult.Name+"/", subdirResult.Reason)
			ws.SetCursorByNameAll(subdirResult.Name)
			diffstatus.SetMessage(fmt.Sprintf("Different: %s/ - %s", subdirResult.Name, subdirResult.Reason))
			return
		}

		// No differences at this level, try to descend into a subdir that exists in all
		// Use FindNextSubdirInAll to respect startAfter and avoid re-searching already-visited subdirs
		subdir, found := filer.FindNextSubdirInAll(ws.Dirs, startAfter)
		if found {
			// Descend into this subdir in all windows
			ws.ChdirAll(subdir)
			startAfter = "" // Start from beginning in new directory
			continue
		}

		// No more subdirs to descend into
		// Check if we're back at initial dirs
		if state.AtInitialDirs(ws.Dirs) {
			// We've completed the search
			// [IMPL:DIFF_SEARCH] Use ephemeral message for completion since diffstatus row disappears
			ws.ClearDiffSearch()
			diffstatus.ClearMessage()
			message.Info("Difference search complete - no differences found")
			// Resize to reclaim space from diff status line
			width, height := widget.Size()
			g.Resize(0, 0, width, height)
			return
		}

		// Go back to parent in all directories and continue
		childDirName := ws.Dir().Base() // Save the child directory name BEFORE going up
		for _, d := range ws.Dirs {
			d.Chdir("..")
		}
		ws.RebuildComparisonIndex()

		// Continue searching from the child directory we just exited
		// so FindNextSubdirInAll can find the next sibling
		startAfter = childDirName

		// If we're back at initial dirs after going up, we're done
		if state.AtInitialDirs(ws.Dirs) {
			// Continue search from where we left off at this level
			continue
		}
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
