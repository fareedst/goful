package filer

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func newTestDirectory(t *testing.T, path string) *Directory {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	dir := NewDirectory(0, 0, 80, 20)
	dir.Chdir(path)
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})
	return dir
}

func namesOf(dir *Directory) []string {
	names := make([]string, 0, dir.Len())
	for _, item := range dir.List() {
		names = append(names, item.Name())
	}
	sort.Strings(names)
	return names
}

func TestFlowOpenDirectory_REQ_INTEGRATION_FLOWS(t *testing.T) {
	// [REQ:INTEGRATION_FLOWS] [ARCH:TEST_STRATEGY_INTEGRATION] [IMPL:TEST_INTEGRATION_FLOWS]
	tmp := t.TempDir()
	files := []string{"alpha.txt", "bravo.txt"}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(tmp, name), []byte(name), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}

	dir := newTestDirectory(t, tmp)
	ws := NewWorkspace(0, 0, 80, 20, "test")
	ws.Dirs = []*Directory{dir}
	ws.SetFocus(0)

	got := namesOf(dir)
	for _, want := range files {
		if !contains(got, want) {
			t.Fatalf("open directory missing %s in %v", want, got)
		}
	}
}

func TestFlowNavigateRename_REQ_INTEGRATION_FLOWS(t *testing.T) {
	// [REQ:INTEGRATION_FLOWS] [ARCH:TEST_STRATEGY_INTEGRATION] [IMPL:TEST_INTEGRATION_FLOWS]
	tmp := t.TempDir()
	subdir := filepath.Join(tmp, "nested")
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	orig := filepath.Join(tmp, "file.txt")
	if err := os.WriteFile(orig, []byte("data"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	dir := newTestDirectory(t, tmp)
	dir.SetCursorByName("nested")
	dir.EnterDir()
	if dir.Path != filepath.Clean(subdir) {
		t.Fatalf("EnterDir path = %s, want %s", dir.Path, subdir)
	}

	newName := filepath.Join(tmp, "file-renamed.txt")
	if err := os.Rename(orig, newName); err != nil {
		t.Fatalf("rename: %v", err)
	}
	dir.Chdir(tmp)
	dir.reload()
	if idx := dir.IndexByName("file-renamed.txt"); idx == dir.Lower() {
		t.Fatalf("rename result not visible after reload")
	}
}

func TestFlowDelete_REQ_INTEGRATION_FLOWS(t *testing.T) {
	// [REQ:INTEGRATION_FLOWS] [ARCH:TEST_STRATEGY_INTEGRATION] [IMPL:TEST_INTEGRATION_FLOWS]
	tmp := t.TempDir()
	target := filepath.Join(tmp, "remove.me")
	if err := os.WriteFile(target, []byte("bye"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	dir := newTestDirectory(t, tmp)
	if err := os.Remove(target); err != nil {
		t.Fatalf("remove: %v", err)
	}
	dir.reload()
	if idx := dir.IndexByName("remove.me"); idx != dir.Lower() {
		t.Fatalf("deleted entry still present with index %d", idx)
	}
}

func contains(slice []string, target string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

// TestChdirAllToSubdir_REQ_LINKED_NAVIGATION tests that linked navigation propagates
// subdirectory changes to all windows that have a matching subdirectory.
// [REQ:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [IMPL:LINKED_NAVIGATION]
func TestChdirAllToSubdir_REQ_LINKED_NAVIGATION(t *testing.T) {
	// Create temp directories with shared subdirectory structure
	tmp1 := t.TempDir()
	tmp2 := t.TempDir()
	tmp3 := t.TempDir()

	// Create "shared" subdirectory in tmp1 and tmp2, but not tmp3
	shared1 := filepath.Join(tmp1, "shared")
	shared2 := filepath.Join(tmp2, "shared")
	if err := os.Mkdir(shared1, 0o755); err != nil {
		t.Fatalf("mkdir shared1: %v", err)
	}
	if err := os.Mkdir(shared2, 0o755); err != nil {
		t.Fatalf("mkdir shared2: %v", err)
	}

	// Create workspace with 3 directories
	ws := NewWorkspace(0, 0, 80, 60, "test")
	dir1 := newTestDirectory(t, tmp1)
	dir2 := newTestDirectory(t, tmp2)
	dir3 := newTestDirectory(t, tmp3)
	ws.Dirs = []*Directory{dir1, dir2, dir3}
	ws.Focus = 0

	// Navigate to "shared" subdirectory on all non-focused windows
	ws.ChdirAllToSubdir("shared")

	// dir1 is focused, should not be changed
	if dir1.Path != filepath.Clean(tmp1) {
		t.Errorf("focused dir1 changed unexpectedly to %s", dir1.Path)
	}
	// dir2 should have navigated to shared
	if dir2.Path != filepath.Clean(shared2) {
		t.Errorf("dir2 should navigate to shared, got %s, want %s", dir2.Path, shared2)
	}
	// dir3 should stay unchanged (no shared subdirectory)
	if dir3.Path != filepath.Clean(tmp3) {
		t.Errorf("dir3 changed unexpectedly to %s (should stay %s)", dir3.Path, tmp3)
	}
}

// TestChdirAllToParent_REQ_LINKED_NAVIGATION tests that linked parent navigation
// propagates to all non-focused windows.
// [REQ:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [IMPL:LINKED_NAVIGATION]
func TestChdirAllToParent_REQ_LINKED_NAVIGATION(t *testing.T) {
	// Create temp directories with subdirectories
	tmp1 := t.TempDir()
	tmp2 := t.TempDir()
	sub1 := filepath.Join(tmp1, "child1")
	sub2 := filepath.Join(tmp2, "child2")
	if err := os.Mkdir(sub1, 0o755); err != nil {
		t.Fatalf("mkdir sub1: %v", err)
	}
	if err := os.Mkdir(sub2, 0o755); err != nil {
		t.Fatalf("mkdir sub2: %v", err)
	}

	// Create workspace with 2 directories, starting in subdirectories
	ws := NewWorkspace(0, 0, 80, 40, "test")
	dir1 := newTestDirectory(t, sub1)
	dir2 := newTestDirectory(t, sub2)
	ws.Dirs = []*Directory{dir1, dir2}
	ws.Focus = 0

	// Navigate to parent on all non-focused windows
	ws.ChdirAllToParent()

	// dir1 is focused, should not be changed
	if dir1.Path != filepath.Clean(sub1) {
		t.Errorf("focused dir1 changed unexpectedly to %s", dir1.Path)
	}
	// dir2 should have navigated to parent
	if dir2.Path != filepath.Clean(tmp2) {
		t.Errorf("dir2 should navigate to parent, got %s, want %s", dir2.Path, tmp2)
	}
}

// TestLinkedNavigationSingleWindow_REQ_LINKED_NAVIGATION tests that linked navigation
// is a no-op when there is only one window.
// [REQ:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [IMPL:LINKED_NAVIGATION]
func TestLinkedNavigationSingleWindow_REQ_LINKED_NAVIGATION(t *testing.T) {
	tmp := t.TempDir()
	shared := filepath.Join(tmp, "shared")
	if err := os.Mkdir(shared, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	ws := NewWorkspace(0, 0, 80, 20, "test")
	dir := newTestDirectory(t, tmp)
	ws.Dirs = []*Directory{dir}
	ws.Focus = 0

	// These should be no-ops without panic
	ws.ChdirAllToSubdir("shared")
	ws.ChdirAllToParent()

	// Directory should remain unchanged (focused dir is skipped)
	if dir.Path != filepath.Clean(tmp) {
		t.Errorf("dir changed unexpectedly to %s", dir.Path)
	}
}

func TestExcludedNamesHideEntries_REQ_FILER_EXCLUDE_NAMES(t *testing.T) {
	// [REQ:FILER_EXCLUDE_NAMES] [ARCH:FILER_EXCLUDE_FILTER] [IMPL:FILER_EXCLUDE_RULES]
	t.Cleanup(func() { ConfigureExcludedNames(nil, false) })
	tmp := t.TempDir()
	noisePath := filepath.Join(tmp, ".DS_Store")
	keepPath := filepath.Join(tmp, "keep.txt")
	if err := os.WriteFile(noisePath, []byte("noise"), 0o644); err != nil {
		t.Fatalf("write noise file: %v", err)
	}
	if err := os.WriteFile(keepPath, []byte("keep"), 0o644); err != nil {
		t.Fatalf("write keep file: %v", err)
	}

	ConfigureExcludedNames([]string{".ds_store"}, true)
	dir := newTestDirectory(t, tmp)
	names := namesOf(dir)
	if contains(names, ".DS_Store") {
		t.Fatalf("excluded file was still listed: %v", names)
	}
	if !contains(names, "keep.txt") {
		t.Fatalf("expected keep.txt to remain visible, got %v", names)
	}

	ToggleExcludedNames() // disable filter
	dir.reload()
	names = namesOf(dir)
	if !contains(names, ".DS_Store") {
		t.Fatalf("excluded file should reappear after toggling filter off, got %v", names)
	}
}
