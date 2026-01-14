package filer

import (
	"os"
	"path/filepath"

	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/widget"
)

type layoutType int

const (
	layoutTile layoutType = iota
	layoutTileTop
	layoutTileBottom
	layoutOneline
	layoutOneColumn
	layoutFullscreen
)

// Workspace is a box storing and layouting directories.
type Workspace struct {
	*widget.Window
	Dirs            []*Directory     `json:"directories"`
	Layout          layoutType       `json:"layout"`
	Title           string           `json:"title"`
	Focus           int              `json:"focus"`
	comparisonIndex *ComparisonIndex // [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
	diffSearch      *DiffSearchState // [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
}

// NewWorkspace returns a new workspace of specified sizes.
func NewWorkspace(x, y, width, height int, title string) *Workspace {
	return &Workspace{
		Window:          widget.NewWindow(x, y, width, height),
		Dirs:            []*Directory{},
		Layout:          layoutTile,
		Title:           title,
		Focus:           0,
		comparisonIndex: nil,
	}
}

func (w *Workspace) init4json(x, y, width, height int) {
	w.Window = widget.NewWindow(x, y, width, height)
}

// CreateDir adds the home directory to the head.
func (w *Workspace) CreateDir() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dir := NewDirectory(0, 0, 0, 0)
	dir.Chdir(home)
	w.Dirs = append(w.Dirs, nil)
	copy(w.Dirs[1:], w.Dirs[:len(w.Dirs)-1])
	w.Dirs[0] = dir
	w.SetFocus(0)
	w.allocate()
	// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
	w.RebuildComparisonIndex()
}

// CloseDir closes the focused directory.
func (w *Workspace) CloseDir() {
	if len(w.Dirs) < 2 {
		return
	}
	i := w.Focus
	w.Dirs = append(w.Dirs[:i], w.Dirs[i+1:]...)
	if w.Focus >= len(w.Dirs) {
		w.Focus = len(w.Dirs) - 1
	}
	w.attach()
	w.allocate()
	// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
	w.RebuildComparisonIndex()
}

// ChdirNeighbor changes the focused path a neighbor directory path.
func (w *Workspace) ChdirNeighbor() {
	w.Dir().Chdir(w.NextDir().Path)
}

func (w *Workspace) visible(visible bool) {
	if visible {
		w.ReloadAll()
	} else {
		for _, d := range w.Dirs {
			d.ClearList()
		}
	}
}

// MoveFocus moves the focus with specified amounts.
func (w *Workspace) MoveFocus(amount int) {
	w.Focus += amount
	if len(w.Dirs) <= w.Focus {
		w.Focus = 0
	} else if w.Focus < 0 {
		w.Focus = len(w.Dirs) - 1
	}
	w.attach()
}

// SetFocus sets the focus to a specified position.
func (w *Workspace) SetFocus(x int) {
	w.Focus = x
	if w.Focus < 0 {
		w.Focus = 0
	} else if w.Focus > len(w.Dirs)-1 {
		w.Focus = len(w.Dirs) - 1
	}
	w.attach()
}

func (w *Workspace) attach() {
	err := os.Chdir(w.Dir().Path)
	if err != nil {
		message.Error(err)
		home, _ := os.UserHomeDir()
		w.Dir().Chdir(home)
	}
}

// ReloadAll reloads all directories.
func (w *Workspace) ReloadAll() {
	for _, d := range w.Dirs {
		d.reload()
	}
	err := os.Chdir(w.Dir().Path)
	if err != nil {
		message.Error(err)
		home, _ := os.UserHomeDir()
		w.Dir().Chdir(home)
	}
	// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
	// Rebuild comparison index after all directories are loaded
	w.RebuildComparisonIndex()
}

// RebuildComparisonIndex rebuilds the comparison index from current directory contents.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func (w *Workspace) RebuildComparisonIndex() {
	w.comparisonIndex = BuildComparisonIndex(w.Dirs)
}

// ComparisonIndex returns the current comparison index.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func (w *Workspace) ComparisonIndex() *ComparisonIndex {
	return w.comparisonIndex
}

// GetCompareState returns the comparison state for a file in a specific directory.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func (w *Workspace) GetCompareState(dirIndex int, filename string) *CompareState {
	if w.comparisonIndex == nil {
		return nil
	}
	return w.comparisonIndex.Get(dirIndex, filename)
}

// CalculateDigestForFile calculates and updates digest states for files with the given name.
// Only computes digests for files with equal sizes across directories.
// Returns the number of files that had their digest calculated.
// [IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func (w *Workspace) CalculateDigestForFile(filename string) int {
	if w.comparisonIndex == nil {
		return 0
	}
	return w.comparisonIndex.UpdateDigestStates(filename, w.Dirs)
}

// Dir returns the focused directory.
func (w *Workspace) Dir() *Directory {
	return w.Dirs[w.Focus]
}

// NextDir returns the next directory.
func (w *Workspace) NextDir() *Directory {
	return w.Dirs[w.nextIndex()]
}

// PrevDir returns the previous directory.
func (w *Workspace) PrevDir() *Directory {
	return w.Dirs[w.prevIndex()]
}

// SwapNextDir swaps focus and next directories.
func (w *Workspace) SwapNextDir() {
	next := w.nextIndex()
	w.Dirs[w.Focus], w.Dirs[next] = w.Dirs[next], w.Dirs[w.Focus]
	w.MoveFocus(1)
	w.allocate()
}

// SwapPrevDir swaps focus and previous directories.
func (w *Workspace) SwapPrevDir() {
	prev := w.prevIndex()
	w.Dirs[w.Focus], w.Dirs[prev] = w.Dirs[prev], w.Dirs[w.Focus]
	w.MoveFocus(-1)
	w.allocate()
}

func (w *Workspace) nextIndex() int {
	i := w.Focus + 1
	if i >= len(w.Dirs) {
		return 0
	}
	return i
}

func (w *Workspace) prevIndex() int {
	i := w.Focus - 1
	if i < 0 {
		return len(w.Dirs) - 1
	}
	return i
}

// SetTitle sets the workspace title.
func (w *Workspace) SetTitle(title string) {
	w.Title = title
}

// ChdirAllToSubdirNoRebuild navigates all non-focused directories to a subdirectory with the given name,
// if that subdirectory exists in each directory's current path.
// Returns (navigated, skipped) counts: navigated is the number of non-focused windows that successfully
// navigated, skipped is the number where the subdirectory does not exist.
// Does NOT rebuild the comparison index - caller is responsible for calling RebuildComparisonIndex().
// [IMPL:LINKED_NAVIGATION] [IMPL:LINKED_NAVIGATION_AUTO_DISABLE] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
func (w *Workspace) ChdirAllToSubdirNoRebuild(name string) (navigated, skipped int) {
	for i, d := range w.Dirs {
		if i == w.Focus {
			continue // Skip focused directory; caller handles it
		}
		targetPath := filepath.Join(d.Path, name)
		if info, err := os.Stat(targetPath); err == nil && info.IsDir() {
			d.Chdir(name)
			navigated++
		} else {
			skipped++
		}
	}
	return navigated, skipped
}

// ChdirAllToSubdir navigates all non-focused directories to a subdirectory with the given name,
// if that subdirectory exists in each directory's current path.
// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
func (w *Workspace) ChdirAllToSubdir(name string) {
	for i, d := range w.Dirs {
		if i == w.Focus {
			continue // Skip focused directory; caller handles it
		}
		targetPath := filepath.Join(d.Path, name)
		if info, err := os.Stat(targetPath); err == nil && info.IsDir() {
			d.Chdir(name)
		}
	}
	w.RebuildComparisonIndex()
}

// ChdirAllToParent navigates all directories (including focused) to their respective parent directories.
// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
func (w *Workspace) ChdirAllToParent() {
	for i, d := range w.Dirs {
		if i == w.Focus {
			continue // Skip focused directory; caller handles it
		}
		d.Chdir("..")
	}
	w.RebuildComparisonIndex()
}

// SortAllBy applies the given sort type to all directories in the workspace.
// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
func (w *Workspace) SortAllBy(typ SortType) {
	for _, d := range w.Dirs {
		d.SortBy(typ)
	}
	w.RebuildComparisonIndex()
}

// StartDiffSearch initializes a new difference search session.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (w *Workspace) StartDiffSearch() {
	w.diffSearch = NewDiffSearchState(w.Dirs)
}

// ClearDiffSearch clears the difference search state.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (w *Workspace) ClearDiffSearch() {
	if w.diffSearch != nil {
		w.diffSearch.Clear()
	}
}

// DiffSearchState returns the current difference search state.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (w *Workspace) DiffSearchState() *DiffSearchState {
	return w.diffSearch
}

// IsDiffSearchActive returns whether a difference search is active.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (w *Workspace) IsDiffSearchActive() bool {
	return w.diffSearch.IsActive()
}

// SetCursorByNameAll moves the cursor to the named entry in all directories where it exists.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (w *Workspace) SetCursorByNameAll(name string) {
	for _, d := range w.Dirs {
		d.SetCursorByName(name)
	}
}

// ChdirAll changes all directories to the specified subdirectory name.
// Used during diff search to descend into matching subdirectories.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (w *Workspace) ChdirAll(name string) {
	for _, d := range w.Dirs {
		d.Chdir(name)
	}
	w.RebuildComparisonIndex()
}

// ChdirAllToInitial navigates all directories back to their initial paths.
// Returns true if successful, false if any directory failed.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (w *Workspace) ChdirAllToInitial() bool {
	if w.diffSearch == nil || len(w.diffSearch.InitialDirs) != len(w.Dirs) {
		return false
	}
	for i, d := range w.Dirs {
		d.Chdir(w.diffSearch.InitialDirs[i])
	}
	w.RebuildComparisonIndex()
	return true
}

// LayoutTile allocates to the tile layout.
func (w *Workspace) LayoutTile() {
	w.Layout = layoutTile
	x, y := w.LeftTop()
	k := len(w.Dirs) - 1
	if k < 1 {
		w.Dirs[0].Resize(x, y, w.Width(), w.Height())
		return
	}
	width := w.Width() / 2
	w.Dirs[0].Resize(x, y, width, w.Height())
	height := w.Height() / k
	hodd := w.Height() % k
	wodd := w.Width() % 2
	for i, d := range w.Dirs[1:k] {
		d.Resize(x+width, y+height*i, width+wodd, height)
	}
	w.Dirs[k].Resize(x+width, y+height*(k-1), width+wodd, height+hodd)
}

// LayoutTileTop allocates to the tile top layout.
func (w *Workspace) LayoutTileTop() {
	w.Layout = layoutTileTop
	x, y := w.LeftTop()
	k := len(w.Dirs) - 1
	if k < 1 {
		w.Dirs[0].Resize(x, y, w.Width(), w.Height())
		return
	}
	height := w.Height() / 2
	hodd := w.Height() % 2

	width := w.Width() / k
	wodd := w.Width() % 2

	w.Dirs[0].Resize(x, y, width, height)
	w.Dirs[k].Resize(x, y+height, w.Width(), height+hodd)
	if k < 2 {
		return
	}
	for i, d := range w.Dirs[1 : k-1] {
		d.Resize(x+width*(i+1), y, width, height)
	}
	w.Dirs[k-1].Resize(x+width*(k-1), y, width+wodd, height)
}

// LayoutTileBottom allocates to the tile bottom layout.
func (w *Workspace) LayoutTileBottom() {
	w.Layout = layoutTileBottom
	x, y := w.LeftTop()
	k := len(w.Dirs) - 1
	if k < 1 {
		w.Dirs[0].Resize(x, y, w.Width(), w.Height())
		return
	}
	height := w.Height() / 2
	hodd := w.Height() % 2

	w.Dirs[0].Resize(x, y, w.Width(), height)

	width := w.Width() / k
	for i, d := range w.Dirs[1:k] {
		d.Resize(x+width*i, y+height, width, height+hodd)
	}
	wodd := w.Width() % 2
	w.Dirs[k].Resize(x+width*(k-1), y+height, width+wodd, height+hodd)
}

// LayoutOnerow allocates to the one line layout.
func (w *Workspace) LayoutOnerow() {
	w.Layout = layoutOneline
	x, y := w.LeftTop()
	k := len(w.Dirs)
	width := w.Width() / k
	for i, d := range w.Dirs[:k-1] {
		d.Resize(x+width*i, y, width, w.Height())
	}
	wodd := w.Width() % k
	w.Dirs[k-1].Resize(x+width*(k-1), y, width+wodd, w.Height())
}

// LayoutOnecolumn allocates to the one column layout.
func (w *Workspace) LayoutOnecolumn() {
	w.Layout = layoutOneColumn
	x, y := w.LeftTop()
	k := len(w.Dirs)
	height := w.Height() / k
	for i, d := range w.Dirs[:k-1] {
		d.Resize(x, y+height*i, w.Width(), height)
	}
	hodd := w.Height() % k
	w.Dirs[k-1].Resize(x, y+height*(k-1), w.Width(), height+hodd)
}

// LayoutFullscreen allocates to the full screen layout.
func (w *Workspace) LayoutFullscreen() {
	w.Layout = layoutFullscreen
	for _, d := range w.Dirs {
		x, y := w.LeftTop()
		d.Resize(x, y, w.Width(), w.Height())
	}
}

func (w *Workspace) allocate() {
	switch w.Layout {
	case layoutTile:
		w.LayoutTile()
	case layoutTileTop:
		w.LayoutTileTop()
	case layoutTileBottom:
		w.LayoutTileBottom()
	case layoutOneline:
		w.LayoutOnerow()
	case layoutOneColumn:
		w.LayoutOnecolumn()
	case layoutFullscreen:
		w.LayoutFullscreen()
	}
}

// Resize and layout allocates.
func (w *Workspace) Resize(x, y, width, height int) {
	w.Window.Resize(x, y, width, height)
	w.allocate()
}

// ResizeRelative relative resizes and layout allocates.
func (w *Workspace) ResizeRelative(x, y, width, height int) {
	w.Window.ResizeRelative(x, y, width, height)
	w.allocate()
}

// Draw all directories and hide a cursor if all finders not active.
func (w *Workspace) Draw() {
	if w.Layout == layoutFullscreen {
		w.Dir().draw(true)
	} else {
		w.draw()
	}
	if !w.isShowCursor() {
		widget.HideCursor()
	}
}

func (w *Workspace) isShowCursor() bool {
	for i, d := range w.Dirs {
		if d.finder != nil && i == w.Focus {
			return true
		}
	}
	return false
}

// [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func (w *Workspace) draw() {
	for i, d := range w.Dirs {
		d.drawWithComparisonIndex(i == w.Focus, i, w.comparisonIndex)
	}
}

// WorkspaceNavigator adapts Workspace to the Navigator interface.
// This allows the TreeWalker to operate on a Workspace while keeping
// the navigation logic decoupled from TUI concerns.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
type WorkspaceNavigator struct {
	ws *Workspace
}

// NewWorkspaceNavigator creates a navigator adapter for the given workspace.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func NewWorkspaceNavigator(ws *Workspace) *WorkspaceNavigator {
	return &WorkspaceNavigator{ws: ws}
}

// GetDirs returns the current directories being compared.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (n *WorkspaceNavigator) GetDirs() []*Directory {
	return n.ws.Dirs
}

// ChdirAll changes to the named subdirectory in all directories.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (n *WorkspaceNavigator) ChdirAll(name string) {
	n.ws.ChdirAll(name)
}

// ChdirParentAll changes to parent directory in all directories.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (n *WorkspaceNavigator) ChdirParentAll() {
	for _, d := range n.ws.Dirs {
		d.Chdir("..")
	}
}

// CurrentPath returns the path of the first directory.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (n *WorkspaceNavigator) CurrentPath() string {
	return n.ws.Dir().Path
}

// RebuildComparisonIndex rebuilds the comparison index after directory changes.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (n *WorkspaceNavigator) RebuildComparisonIndex() {
	n.ws.RebuildComparisonIndex()
}
