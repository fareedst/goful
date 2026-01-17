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

// TestSortAllBy_REQ_LINKED_NAVIGATION tests that linked sort applies
// the same sort type to all directories in the workspace.
// [REQ:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [IMPL:LINKED_NAVIGATION]
func TestSortAllBy_REQ_LINKED_NAVIGATION(t *testing.T) {
	// Create temp directories with files
	tmp1 := t.TempDir()
	tmp2 := t.TempDir()

	// Create test files
	for _, name := range []string{"alpha.txt", "zeta.txt", "beta.txt"} {
		if err := os.WriteFile(filepath.Join(tmp1, name), []byte(name), 0o644); err != nil {
			t.Fatalf("write file tmp1: %v", err)
		}
		if err := os.WriteFile(filepath.Join(tmp2, name), []byte(name), 0o644); err != nil {
			t.Fatalf("write file tmp2: %v", err)
		}
	}

	// Create workspace with 2 directories
	ws := NewWorkspace(0, 0, 80, 40, "test")
	dir1 := newTestDirectory(t, tmp1)
	dir2 := newTestDirectory(t, tmp2)
	ws.Dirs = []*Directory{dir1, dir2}
	ws.Focus = 0

	// Verify initial sort type (default is SortName)
	if dir1.Sort != SortName {
		t.Errorf("dir1 initial sort = %s, want %s", dir1.Sort, SortName)
	}
	if dir2.Sort != SortName {
		t.Errorf("dir2 initial sort = %s, want %s", dir2.Sort, SortName)
	}

	// Apply SortNameRev to all directories
	ws.SortAllBy(SortNameRev)

	// Both directories should now have the same sort type
	if dir1.Sort != SortNameRev {
		t.Errorf("dir1 sort = %s, want %s", dir1.Sort, SortNameRev)
	}
	if dir2.Sort != SortNameRev {
		t.Errorf("dir2 sort = %s, want %s", dir2.Sort, SortNameRev)
	}

	// Apply SortSize to all directories
	ws.SortAllBy(SortSize)

	if dir1.Sort != SortSize {
		t.Errorf("dir1 sort = %s, want %s", dir1.Sort, SortSize)
	}
	if dir2.Sort != SortSize {
		t.Errorf("dir2 sort = %s, want %s", dir2.Sort, SortSize)
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

// TestDirectoryAt_REQ_MOUSE_FILE_SELECT tests workspace hit-testing for mouse events.
// [IMPL:MOUSE_HIT_TEST] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]
func TestDirectoryAt_REQ_MOUSE_FILE_SELECT(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "test.txt"), []byte("data"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// Create workspace with two directories side by side
	ws := NewWorkspace(0, 0, 80, 20, "test")
	dir1 := NewDirectory(0, 0, 40, 20)
	dir1.Chdir(tmp)
	dir2 := NewDirectory(40, 0, 40, 20)
	dir2.Chdir(tmp)
	ws.Dirs = []*Directory{dir1, dir2}
	ws.SetFocus(0)

	// Test clicking in first directory
	gotDir, gotIdx := ws.DirectoryAt(20, 10)
	if gotDir != dir1 || gotIdx != 0 {
		t.Errorf("DirectoryAt(20,10) = (%p, %d), want (%p, 0)", gotDir, gotIdx, dir1)
	}

	// Test clicking in second directory
	gotDir, gotIdx = ws.DirectoryAt(60, 10)
	if gotDir != dir2 || gotIdx != 1 {
		t.Errorf("DirectoryAt(60,10) = (%p, %d), want (%p, 1)", gotDir, gotIdx, dir2)
	}

	// Test clicking outside all directories
	gotDir, gotIdx = ws.DirectoryAt(100, 10)
	if gotDir != nil || gotIdx != -1 {
		t.Errorf("DirectoryAt(100,10) = (%p, %d), want (nil, -1)", gotDir, gotIdx)
	}
}

// TestFileIndexAtY_REQ_MOUSE_FILE_SELECT tests file list hit-testing for mouse events.
// [IMPL:MOUSE_HIT_TEST] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]
func TestFileIndexAtY_REQ_MOUSE_FILE_SELECT(t *testing.T) {
	tmp := t.TempDir()
	// Create 5 files to have a scrollable list
	for i := 0; i < 5; i++ {
		name := filepath.Join(tmp, string(rune('a'+i))+".txt")
		if err := os.WriteFile(name, []byte("data"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}

	// Create directory at position (10, 5) with size 60x15
	dir := NewDirectory(10, 5, 60, 15)
	dir.Chdir(tmp)

	// Directory content starts at y=6 (after header at y=5)
	// File indices should map: y=6 -> 0, y=7 -> 1, etc.
	for _, tc := range []struct {
		name     string
		y        int
		expected int
	}{
		{name: "first_file", y: 6, expected: 0},
		{name: "second_file", y: 7, expected: 1},
		{name: "third_file", y: 8, expected: 2},
		{name: "header_row", y: 5, expected: -1},     // Header row
		{name: "above_window", y: 4, expected: -1},   // Above window
		{name: "footer_row", y: 18, expected: -1},    // Footer row (y=5+15-2=18)
		{name: "below_content", y: 19, expected: -1}, // Below window
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := dir.FileIndexAtY(tc.y)
			if got != tc.expected {
				t.Errorf("FileIndexAtY(%d) = %d, want %d", tc.y, got, tc.expected)
			}
		})
	}
}

// TestDirectoryContains_REQ_MOUSE_FILE_SELECT tests directory boundary detection.
// [IMPL:MOUSE_HIT_TEST] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]
func TestDirectoryContains_REQ_MOUSE_FILE_SELECT(t *testing.T) {
	dir := NewDirectory(10, 5, 40, 15)

	for _, tc := range []struct {
		name     string
		x, y     int
		expected bool
	}{
		{name: "inside", x: 30, y: 10, expected: true},
		{name: "top_left", x: 10, y: 5, expected: true},
		{name: "bottom_right", x: 49, y: 19, expected: true},
		{name: "left_edge_outside", x: 9, y: 10, expected: false},
		{name: "right_edge_outside", x: 50, y: 10, expected: false},
		{name: "above", x: 30, y: 4, expected: false},
		{name: "below", x: 30, y: 20, expected: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := dir.Contains(tc.x, tc.y)
			if got != tc.expected {
				t.Errorf("Directory.Contains(%d,%d) = %v, want %v", tc.x, tc.y, got, tc.expected)
			}
		})
	}
}

// [REQ:DEBT_TRIAGE] [IMPL:EXTMAP_API_SAFETY] [ARCH:DEBT_MANAGEMENT]
// Test that AddExtmap works without prior MergeExtmap call (regression test for nil map panic).
func TestAddExtmap_NilMapSafe_REQ_DEBT_TRIAGE(t *testing.T) {
	// Create a filer with empty extmap (no MergeExtmap called)
	f := &Filer{
		extmap: make(map[string]map[string]func()),
	}

	called := false
	callback := func() { called = true }

	// This should not panic even though f.extmap["testkey"] is nil
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("AddExtmap panicked: %v", r)
		}
	}()

	f.AddExtmap("testkey", ".txt", callback)

	// Verify the callback was added correctly
	if f.extmap["testkey"] == nil {
		t.Fatal("AddExtmap should have created inner map")
	}
	if f.extmap["testkey"][".txt"] == nil {
		t.Fatal("AddExtmap should have stored callback")
	}

	// Call the callback to verify it's the right one
	f.extmap["testkey"][".txt"]()
	if !called {
		t.Fatal("stored callback should be callable")
	}
}

// [REQ:DEBT_TRIAGE] [IMPL:EXTMAP_API_SAFETY] [ARCH:DEBT_MANAGEMENT]
// Test AddExtmap with multiple entries and pre-existing keys.
func TestAddExtmap_MultipleEntries_REQ_DEBT_TRIAGE(t *testing.T) {
	f := &Filer{
		extmap: make(map[string]map[string]func()),
	}

	count := 0
	f.AddExtmap(
		"key1", ".go", func() { count += 1 },
		"key1", ".py", func() { count += 10 },
		"key2", ".rs", func() { count += 100 },
	)

	// Verify all entries were added
	if len(f.extmap["key1"]) != 2 {
		t.Errorf("expected 2 entries for key1, got %d", len(f.extmap["key1"]))
	}
	if len(f.extmap["key2"]) != 1 {
		t.Errorf("expected 1 entry for key2, got %d", len(f.extmap["key2"]))
	}

	// Call all callbacks
	f.extmap["key1"][".go"]()
	f.extmap["key1"][".py"]()
	f.extmap["key2"][".rs"]()
	if count != 111 {
		t.Errorf("expected count=111, got %d", count)
	}
}
