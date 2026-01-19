package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fareedst/goful/filer"
)

// newTestGoful creates a minimal Goful instance with a workspace containing two directories.
// [IMPL:LINKED_CURSOR_SYNC] [REQ:LINKED_NAVIGATION] Test helper for linked cursor tests.
func newTestGoful(t *testing.T, tmp1, tmp2 string) *Goful {
	t.Helper()
	ws := filer.NewWorkspace(0, 0, 80, 20, "test")
	dir1 := filer.NewDirectory(0, 0, 40, 20)
	dir1.Chdir(tmp1)
	dir2 := filer.NewDirectory(40, 0, 40, 20)
	dir2.Chdir(tmp2)
	ws.Dirs = []*filer.Directory{dir1, dir2}
	ws.SetFocus(0)

	f := filer.New(0, 0, 80, 24)
	f.Workspaces = []*filer.Workspace{ws}
	f.Current = 0

	return &Goful{
		Filer:     f,
		linkedNav: true, // Linked mode ON by default
	}
}

// TestMoveCursorLinked_REQ_LINKED_NAVIGATION tests cursor movement syncs to other windows when linked mode is ON.
// [IMPL:LINKED_CURSOR_SYNC] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
func TestMoveCursorLinked_REQ_LINKED_NAVIGATION(t *testing.T) {
	tmp1 := t.TempDir()
	tmp2 := t.TempDir()

	// Create common files in both directories
	commonFiles := []string{"alpha.txt", "bravo.txt", "charlie.txt", "delta.txt"}
	for _, name := range commonFiles {
		for _, tmp := range []string{tmp1, tmp2} {
			if err := os.WriteFile(filepath.Join(tmp, name), []byte(name), 0o644); err != nil {
				t.Fatalf("write file: %v", err)
			}
		}
	}

	g := newTestGoful(t, tmp1, tmp2)

	// Verify linked mode is ON
	if !g.IsLinkedNav() {
		t.Fatal("linked mode should be ON")
	}

	// Both directories should start with cursor at position 0
	dir1 := g.Workspace().Dirs[0]
	dir2 := g.Workspace().Dirs[1]

	// Move cursor down in focused window
	g.MoveCursorLinked(1)

	// Verify both windows have cursor on the same file
	if dir1.File().Name() != dir2.File().Name() {
		t.Errorf("cursors not synced: dir1=%s, dir2=%s", dir1.File().Name(), dir2.File().Name())
	}

	// Move cursor down again
	g.MoveCursorLinked(1)

	// Verify both windows still synced
	if dir1.File().Name() != dir2.File().Name() {
		t.Errorf("cursors not synced after second move: dir1=%s, dir2=%s", dir1.File().Name(), dir2.File().Name())
	}
}

// TestMoveCursorLinkedOff_REQ_LINKED_NAVIGATION tests cursor movement does NOT sync when linked mode is OFF.
// [IMPL:LINKED_CURSOR_SYNC] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
func TestMoveCursorLinkedOff_REQ_LINKED_NAVIGATION(t *testing.T) {
	tmp1 := t.TempDir()
	tmp2 := t.TempDir()

	// Create common files in both directories
	commonFiles := []string{"alpha.txt", "bravo.txt", "charlie.txt", "delta.txt"}
	for _, name := range commonFiles {
		for _, tmp := range []string{tmp1, tmp2} {
			if err := os.WriteFile(filepath.Join(tmp, name), []byte(name), 0o644); err != nil {
				t.Fatalf("write file: %v", err)
			}
		}
	}

	g := newTestGoful(t, tmp1, tmp2)

	// Turn OFF linked mode
	g.SetLinkedNav(false)

	// Verify linked mode is OFF
	if g.IsLinkedNav() {
		t.Fatal("linked mode should be OFF")
	}

	// Get initial positions
	dir1 := g.Workspace().Dirs[0]
	dir2 := g.Workspace().Dirs[1]
	initial1 := dir1.Cursor()
	initial2 := dir2.Cursor()

	// Move cursor down in focused window
	g.MoveCursorLinked(1)

	// Dir1 (focused) should have moved, dir2 should NOT have moved
	if dir1.Cursor() == initial1 {
		t.Error("focused window cursor should have moved")
	}
	if dir2.Cursor() != initial2 {
		t.Errorf("unfocused window cursor should NOT have moved: was %d, now %d", initial2, dir2.Cursor())
	}
}

// TestMoveTopLinked_REQ_LINKED_NAVIGATION tests MoveTopLinked syncs when linked mode is ON.
// [IMPL:LINKED_CURSOR_SYNC] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
func TestMoveTopLinked_REQ_LINKED_NAVIGATION(t *testing.T) {
	tmp1 := t.TempDir()
	tmp2 := t.TempDir()

	// Create common files in both directories
	commonFiles := []string{"alpha.txt", "bravo.txt", "charlie.txt"}
	for _, name := range commonFiles {
		for _, tmp := range []string{tmp1, tmp2} {
			if err := os.WriteFile(filepath.Join(tmp, name), []byte(name), 0o644); err != nil {
				t.Fatalf("write file: %v", err)
			}
		}
	}

	g := newTestGoful(t, tmp1, tmp2)
	dir1 := g.Workspace().Dirs[0]
	dir2 := g.Workspace().Dirs[1]

	// Move cursors down first
	g.MoveCursorLinked(2)

	// Now move to top
	g.MoveTopLinked()

	// Both should be at top (first file after ..)
	if dir1.File().Name() != dir2.File().Name() {
		t.Errorf("cursors not synced at top: dir1=%s, dir2=%s", dir1.File().Name(), dir2.File().Name())
	}
}

// TestMoveBottomLinked_REQ_LINKED_NAVIGATION tests MoveBottomLinked syncs when linked mode is ON.
// [IMPL:LINKED_CURSOR_SYNC] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
func TestMoveBottomLinked_REQ_LINKED_NAVIGATION(t *testing.T) {
	tmp1 := t.TempDir()
	tmp2 := t.TempDir()

	// Create common files in both directories
	commonFiles := []string{"alpha.txt", "bravo.txt", "charlie.txt"}
	for _, name := range commonFiles {
		for _, tmp := range []string{tmp1, tmp2} {
			if err := os.WriteFile(filepath.Join(tmp, name), []byte(name), 0o644); err != nil {
				t.Fatalf("write file: %v", err)
			}
		}
	}

	g := newTestGoful(t, tmp1, tmp2)
	dir1 := g.Workspace().Dirs[0]
	dir2 := g.Workspace().Dirs[1]

	// Move to bottom
	g.MoveBottomLinked()

	// Both should be at bottom (last file)
	if dir1.File().Name() != dir2.File().Name() {
		t.Errorf("cursors not synced at bottom: dir1=%s, dir2=%s", dir1.File().Name(), dir2.File().Name())
	}
	if dir1.File().Name() != "charlie.txt" {
		t.Errorf("cursor not at last file: got %s, want charlie.txt", dir1.File().Name())
	}
}

// TestLinkedCursorSyncMissingFile_REQ_LINKED_NAVIGATION tests cursor hides in windows without matching file.
// [IMPL:LINKED_CURSOR_SYNC] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
func TestLinkedCursorSyncMissingFile_REQ_LINKED_NAVIGATION(t *testing.T) {
	tmp1 := t.TempDir()
	tmp2 := t.TempDir()

	// Create different files in each directory
	if err := os.WriteFile(filepath.Join(tmp1, "alpha.txt"), []byte("a"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp1, "unique_to_1.txt"), []byte("1"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp2, "alpha.txt"), []byte("a"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp2, "unique_to_2.txt"), []byte("2"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	g := newTestGoful(t, tmp1, tmp2)
	dir1 := g.Workspace().Dirs[0]
	dir2 := g.Workspace().Dirs[1]

	// Move to unique_to_1.txt in dir1
	dir1.SetCursorByName("unique_to_1.txt")

	// Trigger linked sync
	if g.IsLinkedNav() {
		g.Workspace().SetCursorByNameAll(g.File().Name())
	}

	// Dir2 should have cursor hidden (file doesn't exist there)
	if !dir2.IsCursorHidden() {
		t.Error("dir2 cursor should be hidden when file doesn't exist")
	}
}
