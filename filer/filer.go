// Package filer draws directories and files and handles inputs.
package filer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/util"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
)

// Filer is a file manager with workspaces to layout directorires to list files.
type Filer struct {
	*widget.Window
	keymap     widget.Keymap
	extmap     widget.Extmap
	Workspaces []*Workspace `json:"workspaces"`
	Current    int          `json:"current"`
}

// linkedNavIndicator is a callback that returns whether linked navigation is enabled.
// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
var linkedNavIndicator func() bool

// SetLinkedNavIndicator sets the callback used to check if linked navigation is enabled.
// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
func SetLinkedNavIndicator(fn func() bool) {
	linkedNavIndicator = fn
}

// diffSearchStatusFn is a callback that returns the diff search status text.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
var diffSearchStatusFn func() string

// SetDiffSearchStatusFn sets the callback used to get the diff search status text.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func SetDiffSearchStatusFn(fn func() string) {
	diffSearchStatusFn = fn
}

// toolbarButtonBounds stores the screen bounds of toolbar buttons for hit-testing.
// Key is the button identifier (e.g., "parent"), value is the bounds.
// [IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]
type toolbarBounds struct {
	x1, y, x2 int
}

var toolbarButtons = make(map[string]toolbarBounds)

// toolbarParentNavFn is a callback invoked when the parent toolbar button is clicked.
// [IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]
var toolbarParentNavFn func()

// SetToolbarParentNavFn sets the callback for the parent toolbar button.
// [IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]
func SetToolbarParentNavFn(fn func()) {
	toolbarParentNavFn = fn
}

// toolbarLinkedToggleFn is a callback invoked when the linked toggle button is clicked.
// [IMPL:TOOLBAR_LINKED_TOGGLE] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_LINKED_TOGGLE]
var toolbarLinkedToggleFn func()

// SetToolbarLinkedToggleFn sets the callback for the linked toggle button.
// [IMPL:TOOLBAR_LINKED_TOGGLE] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_LINKED_TOGGLE]
func SetToolbarLinkedToggleFn(fn func()) {
	toolbarLinkedToggleFn = fn
}

// toolbarCompareDigestFn is a callback invoked when the compare digest button is clicked.
// [IMPL:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_COMPARE_BUTTON]
var toolbarCompareDigestFn func()

// SetToolbarCompareDigestFn sets the callback for the compare digest button.
// [IMPL:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_COMPARE_BUTTON]
func SetToolbarCompareDigestFn(fn func()) {
	toolbarCompareDigestFn = fn
}

// toolbarSyncCopyFn is a callback invoked when the sync copy button is clicked.
// [IMPL:TOOLBAR_SYNC_COPY] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
var toolbarSyncCopyFn func()

// SetToolbarSyncCopyFn sets the callback for the sync copy button.
// [IMPL:TOOLBAR_SYNC_COPY] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
func SetToolbarSyncCopyFn(fn func()) {
	toolbarSyncCopyFn = fn
}

// toolbarSyncDeleteFn is a callback invoked when the sync delete button is clicked.
// [IMPL:TOOLBAR_SYNC_DELETE] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
var toolbarSyncDeleteFn func()

// SetToolbarSyncDeleteFn sets the callback for the sync delete button.
// [IMPL:TOOLBAR_SYNC_DELETE] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
func SetToolbarSyncDeleteFn(fn func()) {
	toolbarSyncDeleteFn = fn
}

// toolbarSyncRenameFn is a callback invoked when the sync rename button is clicked.
// [IMPL:TOOLBAR_SYNC_RENAME] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
var toolbarSyncRenameFn func()

// SetToolbarSyncRenameFn sets the callback for the sync rename button.
// [IMPL:TOOLBAR_SYNC_RENAME] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
func SetToolbarSyncRenameFn(fn func()) {
	toolbarSyncRenameFn = fn
}

// toolbarIgnoreFailuresFn is a callback invoked when the ignore-failures toggle button is clicked.
// [IMPL:TOOLBAR_IGNORE_FAILURES] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
var toolbarIgnoreFailuresFn func()

// SetToolbarIgnoreFailuresFn sets the callback for the ignore-failures toggle button.
// [IMPL:TOOLBAR_IGNORE_FAILURES] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
func SetToolbarIgnoreFailuresFn(fn func()) {
	toolbarIgnoreFailuresFn = fn
}

// syncIgnoreFailuresIndicator returns whether ignore-failures mode is enabled (for button styling).
// [IMPL:TOOLBAR_IGNORE_FAILURES] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
var syncIgnoreFailuresIndicator func() bool

// SetSyncIgnoreFailuresIndicator sets the callback used to check if ignore-failures mode is enabled.
// [IMPL:TOOLBAR_IGNORE_FAILURES] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
func SetSyncIgnoreFailuresIndicator(fn func() bool) {
	syncIgnoreFailuresIndicator = fn
}

// ToolbarButtonAt returns the toolbar button identifier at coordinates (x, y).
// Returns empty string if no button is at that position.
// [IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]
func ToolbarButtonAt(x, y int) string {
	for name, bounds := range toolbarButtons {
		if y == bounds.y && x >= bounds.x1 && x <= bounds.x2 {
			return name
		}
	}
	return ""
}

// InvokeToolbarButton invokes the action for the named toolbar button.
// Returns true if the button was handled, false if unknown.
// [IMPL:TOOLBAR_PARENT_BUTTON] [IMPL:TOOLBAR_LINKED_TOGGLE] [IMPL:TOOLBAR_COMPARE_BUTTON] [IMPL:TOOLBAR_SYNC_COPY] [IMPL:TOOLBAR_SYNC_DELETE] [IMPL:TOOLBAR_SYNC_RENAME] [IMPL:TOOLBAR_IGNORE_FAILURES] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON] [REQ:TOOLBAR_LINKED_TOGGLE] [REQ:TOOLBAR_COMPARE_BUTTON] [REQ:TOOLBAR_SYNC_BUTTONS]
func InvokeToolbarButton(name string) bool {
	switch name {
	case "parent":
		if toolbarParentNavFn != nil {
			toolbarParentNavFn()
			return true
		}
	case "linked":
		// [IMPL:TOOLBAR_LINKED_TOGGLE] Toggle linked navigation mode
		if toolbarLinkedToggleFn != nil {
			toolbarLinkedToggleFn()
			return true
		}
	case "compare":
		// [IMPL:TOOLBAR_COMPARE_BUTTON] Calculate digests for all shared files
		if toolbarCompareDigestFn != nil {
			toolbarCompareDigestFn()
			return true
		}
	case "synccopy":
		// [IMPL:TOOLBAR_SYNC_COPY] Trigger sync or single-window copy based on Linked mode
		if toolbarSyncCopyFn != nil {
			toolbarSyncCopyFn()
			return true
		}
	case "syncdelete":
		// [IMPL:TOOLBAR_SYNC_DELETE] Trigger sync or single-window delete based on Linked mode
		if toolbarSyncDeleteFn != nil {
			toolbarSyncDeleteFn()
			return true
		}
	case "syncrename":
		// [IMPL:TOOLBAR_SYNC_RENAME] Trigger sync or single-window rename based on Linked mode
		if toolbarSyncRenameFn != nil {
			toolbarSyncRenameFn()
			return true
		}
	case "ignorefailures":
		// [IMPL:TOOLBAR_IGNORE_FAILURES] Toggle ignore-failures mode for sync operations
		if toolbarIgnoreFailuresFn != nil {
			toolbarIgnoreFailuresFn()
			return true
		}
	}
	return false
}

// New creates a new filer based on specified size and coordinates.
// Creates five workspaces and default path is home directory.
func New(x, y, width, height int) *Filer {
	home, err := os.UserHomeDir()
	if err != nil {
		message.Error(err)
		home = "/"
	}

	workspaces := make([]*Workspace, 3)
	for i := 0; i < 3; i++ {
		title := fmt.Sprintf("%d", i+1)
		ws := NewWorkspace(x, y+1, width, height-1, title)
		ws.Dirs = make([]*Directory, 2)
		for j := 0; j < 2; j++ {
			ws.Dirs[j] = NewDirectory(0, 0, 0, 0)
			ws.Dirs[j].Path = home
			ws.Dirs[j].SetTitle(util.AbbrPath(home))
		}
		ws.allocate()
		workspaces[i] = ws
	}
	return &Filer{
		Window:     widget.NewWindow(x, y, width, height),
		keymap:     widget.Keymap{},
		extmap:     widget.Extmap{},
		Workspaces: workspaces,
		Current:    0,
	}
}

// NewFromState creates a new filer form the state json file.
func NewFromState(path string, x, y, width, height int) *Filer {
	file, err := os.Open(util.ExpandPath(path))
	if err != nil {
		return New(x, y, width, height)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return New(x, y, width, height)
	}

	filer := &Filer{
		Window:     widget.NewWindow(x, y, width, height),
		keymap:     widget.Keymap{},
		extmap:     widget.Extmap{},
		Workspaces: []*Workspace{},
		Current:    0,
	}
	if err := json.Unmarshal(data, filer); err != nil {
		return New(x, y, width, height)
	}
	if len(filer.Workspaces) < 1 {
		return New(x, y, width, height)
	}
	for _, ws := range filer.Workspaces {
		if len(ws.Dirs) < 1 {
			return New(x, y, width, height)
		}
		ws.init4json(x, y+1, width, height-1)
		for _, dir := range ws.Dirs {
			dir.init4json()
		}
		ws.allocate()
	}
	return filer
}

// SaveState saves the filer state to the file.
func (f *Filer) SaveState(path string) error {
	jsondata, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	file, err := createStateFile(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(jsondata); err != nil {
		return err
	}
	return nil
}

func createStateFile(path string) (*os.File, error) {
	statePath := util.ExpandPath(path)
	stateDir := filepath.Dir(statePath)
	// [IMPL:STATE_PATH_RESOLVER] [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS]
	// Ensure overrides work even when parent directories do not yet exist.
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return nil, err
	}
	return os.Create(statePath)
}

// CreateWorkspace creates and adds a workspace to the end.
func (f *Filer) CreateWorkspace() {
	title := fmt.Sprintf("%d", len(f.Workspaces)+1)
	x, y := f.LeftTop()
	width, height := f.Width(), f.Height()
	ws := NewWorkspace(x, y+1, width, height-1, title)
	ws.CreateDir()
	ws.CreateDir()
	f.Workspaces = append(f.Workspaces, ws)
}

// CloseWorkspace closes a workspace on the current.
func (f *Filer) CloseWorkspace() {
	if len(f.Workspaces) < 2 {
		return
	}
	i := f.Current
	f.Workspaces[i].visible(false)
	f.Workspaces[i] = nil
	f.Workspaces = append(f.Workspaces[:i], f.Workspaces[i+1:]...)
	if f.Current > len(f.Workspaces)-1 {
		f.Current = len(f.Workspaces) - 1
	}
	f.Workspace().visible(true)
}

// MoveWorkspace moves to the other workspace.
func (f *Filer) MoveWorkspace(amount int) {
	f.Workspace().visible(false)
	f.Current += amount
	if f.Current >= len(f.Workspaces) {
		f.Current = 0
	} else if f.Current < 0 {
		f.Current = len(f.Workspaces) - 1
	}
	f.Workspace().visible(true)
}

// Workspace returns the current workspace.
func (f *Filer) Workspace() *Workspace {
	return f.Workspaces[f.Current]
}

// Dir returns the focused directory on the current workspace.
func (f *Filer) Dir() *Directory {
	return f.Workspace().Dir()
}

// File returns the cursor file in the focused directory on the current workspace.
func (f *Filer) File() *FileStat {
	return f.Dir().File()
}

// AddKeymap adds to the filer keymap.
func (f *Filer) AddKeymap(keys ...interface{}) {
	if len(keys)%2 != 0 {
		panic("items must be a multiple of 2")
	}

	for i := 0; i < len(keys); i += 2 {
		key := keys[i].(string)
		callback := keys[i+1].(func())
		f.keymap[key] = callback
	}
}

// MergeKeymap merges to the filer keymap.
func (f *Filer) MergeKeymap(m widget.Keymap) {
	for key, callback := range m {
		f.keymap[key] = callback
	}
}

// AddExtmap adds to the filer extmap.
// [IMPL:EXTMAP_API_SAFETY] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]
// Safe for third-party integrations: allocates inner map if missing.
func (f *Filer) AddExtmap(a ...interface{}) {
	if len(a)%3 != 0 {
		panic("items must be a multiple of 3")
	}

	for i := 0; i < len(a); i += 3 {
		key := a[i].(string)
		ext := a[i+1].(string)
		callback := a[i+2].(func())
		// [IMPL:EXTMAP_API_SAFETY] Allocate inner map if missing to prevent nil map panic.
		if _, found := f.extmap[key]; !found {
			f.extmap[key] = map[string]func(){}
		}
		f.extmap[key][ext] = callback
	}
}

// MergeExtmap merges to the filer extmap.
func (f *Filer) MergeExtmap(m widget.Extmap) {
	for key, submap := range m {
		if _, found := f.extmap[key]; !found {
			f.extmap[key] = map[string]func(){}
		}
		for ext, callback := range submap {
			f.extmap[key][ext] = callback
		}
	}
}

// Input for key events.
func (f *Filer) Input(key string) {
	if finder := f.Dir().finder; finder != nil {
		if callback, ok := finderKeymap(finder)[key]; ok {
			callback()
			return
		} else if utf8.RuneCountInString(key) == 1 && key != " " {
			r, _ := utf8.DecodeRuneInString(key)
			finder.InsertChar(r)
			return
		}
	}

	if ext, ok := f.extmap[key]; ok {
		if callback, ok := ext[".dir"]; ok && (f.File().IsDir() || f.File().stat.IsDir()) {
			callback()
		} else if callback, ok := ext[".exec"]; ok && f.File().IsExec() {
			callback()
		} else if callback, ok := ext[f.File().Ext()]; ok {
			callback()
		} else if callback, ok := f.keymap[key]; ok {
			callback()
		}
	} else if callback, ok := f.keymap[key]; ok {
		callback()
	}
}

func (f *Filer) drawHeader() {
	x, y := f.LeftTop()

	// [IMPL:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_PARENT_BUTTON]
	// Draw parent directory button at the start of the header
	parentBtn := "[^]"
	parentX1 := x
	x = widget.SetCells(x, y, parentBtn, look.Default().Reverse(true))
	parentX2 := x - 1
	toolbarButtons["parent"] = toolbarBounds{x1: parentX1, y: y, x2: parentX2}
	x = widget.SetCells(x, y, " ", look.Default())

	// [IMPL:TOOLBAR_LINKED_TOGGLE] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_LINKED_TOGGLE]
	// Draw linked mode toggle button - reverse style when ON, normal when OFF
	linkedBtn := "[L]"
	linkedX1 := x
	linkedStyle := look.Default()
	if linkedNavIndicator != nil && linkedNavIndicator() {
		linkedStyle = linkedStyle.Reverse(true)
	}
	x = widget.SetCells(x, y, linkedBtn, linkedStyle)
	linkedX2 := x - 1
	toolbarButtons["linked"] = toolbarBounds{x1: linkedX1, y: y, x2: linkedX2}
	x = widget.SetCells(x, y, " ", look.Default())

	// [IMPL:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_COMPARE_BUTTON]
	// Draw compare digest button - normal style (action button, not state toggle)
	compareBtn := "[=]"
	compareX1 := x
	x = widget.SetCells(x, y, compareBtn, look.Default())
	compareX2 := x - 1
	toolbarButtons["compare"] = toolbarBounds{x1: compareX1, y: y, x2: compareX2}
	x = widget.SetCells(x, y, " ", look.Default())

	// [IMPL:TOOLBAR_SYNC_COPY] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
	// Draw sync copy button - normal style (action button)
	syncCopyBtn := "[C]"
	syncCopyX1 := x
	x = widget.SetCells(x, y, syncCopyBtn, look.Default())
	syncCopyX2 := x - 1
	toolbarButtons["synccopy"] = toolbarBounds{x1: syncCopyX1, y: y, x2: syncCopyX2}
	x = widget.SetCells(x, y, " ", look.Default())

	// [IMPL:TOOLBAR_SYNC_DELETE] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
	// Draw sync delete button - normal style (action button)
	syncDeleteBtn := "[D]"
	syncDeleteX1 := x
	x = widget.SetCells(x, y, syncDeleteBtn, look.Default())
	syncDeleteX2 := x - 1
	toolbarButtons["syncdelete"] = toolbarBounds{x1: syncDeleteX1, y: y, x2: syncDeleteX2}
	x = widget.SetCells(x, y, " ", look.Default())

	// [IMPL:TOOLBAR_SYNC_RENAME] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
	// Draw sync rename button - normal style (action button)
	syncRenameBtn := "[R]"
	syncRenameX1 := x
	x = widget.SetCells(x, y, syncRenameBtn, look.Default())
	syncRenameX2 := x - 1
	toolbarButtons["syncrename"] = toolbarBounds{x1: syncRenameX1, y: y, x2: syncRenameX2}
	x = widget.SetCells(x, y, " ", look.Default())

	// [IMPL:TOOLBAR_IGNORE_FAILURES] [ARCH:TOOLBAR_LAYOUT] [REQ:TOOLBAR_SYNC_BUTTONS]
	// Draw ignore-failures toggle button - reverse style when ON, normal when OFF
	ignoreBtn := "[!]"
	ignoreX1 := x
	ignoreStyle := look.Default()
	if syncIgnoreFailuresIndicator != nil && syncIgnoreFailuresIndicator() {
		ignoreStyle = ignoreStyle.Reverse(true)
	}
	x = widget.SetCells(x, y, ignoreBtn, ignoreStyle)
	ignoreX2 := x - 1
	toolbarButtons["ignorefailures"] = toolbarBounds{x1: ignoreX1, y: y, x2: ignoreX2}
	x = widget.SetCells(x, y, " ", look.Default())

	for i, ws := range f.Workspaces {
		s := fmt.Sprintf(" %s ", ws.Title)
		if f.Current != i {
			x = widget.SetCells(x, y, s, look.Default())
		} else {
			x = widget.SetCells(x, y, s, look.Default().Reverse(true))
		}
	}
	x = widget.SetCells(x, y, " | ", look.Default())

	// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
	// Show diff search status when active
	if diffSearchStatusFn != nil {
		if status := diffSearchStatusFn(); status != "" {
			x = widget.SetCells(x, y, status, look.Default().Reverse(true))
			x = widget.SetCells(x, y, " ", look.Default())
		}
	}

	ws := f.Workspace()
	width := (f.Width() - x) / len(ws.Dirs)
	for i := 0; i < len(ws.Dirs); i++ {
		style := look.Default()
		if ws.Focus == i {
			style = style.Reverse(true)
		}
		s := fmt.Sprintf("[%d] ", i+1)
		x = widget.SetCells(x, y, s, style)
		w := width - len(s)
		s = util.ShortenPath(ws.Dirs[i].Title(), w)
		s = runewidth.Truncate(s, w, "~")
		s = runewidth.FillRight(s, w)
		x = widget.SetCells(x, y, s, style)
	}
}

// Draw the current workspace.
func (f *Filer) Draw() {
	f.Clear()
	f.drawHeader()
	f.Workspace().Draw()
}

// Resize all workspaces.
func (f *Filer) Resize(x, y, width, height int) {
	f.Window.Resize(x, y, width, height)
	for _, ws := range f.Workspaces {
		ws.Resize(x, y+1, width, height-1)
	}
}

// ResizeRelative resize relative to current sizes.
func (f *Filer) ResizeRelative(x, y, width, height int) {
	f.Window.ResizeRelative(x, y, width, height)
	for _, ws := range f.Workspaces {
		ws.ResizeRelative(x, y, width, height)
	}
}
