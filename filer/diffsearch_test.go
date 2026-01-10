// Package filer difference search tests.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
package filer

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewDiffSearchState_REQ_DIFF_SEARCH tests state initialization.
func TestNewDiffSearchState_REQ_DIFF_SEARCH(t *testing.T) {
	// Create temp directories for testing
	tmpDir := t.TempDir()
	dir1Path := filepath.Join(tmpDir, "dir1")
	dir2Path := filepath.Join(tmpDir, "dir2")
	os.MkdirAll(dir1Path, 0755)
	os.MkdirAll(dir2Path, 0755)

	// Create Directory objects
	dir1 := NewDirectory(0, 0, 10, 10)
	dir1.Chdir(dir1Path)
	dir2 := NewDirectory(0, 0, 10, 10)
	dir2.Chdir(dir2Path)

	dirs := []*Directory{dir1, dir2}

	// Test state creation
	state := NewDiffSearchState(dirs)
	if state == nil {
		t.Fatal("NewDiffSearchState returned nil")
	}

	if !state.Active {
		t.Error("Expected state to be active after creation")
	}

	if len(state.InitialDirs) != 2 {
		t.Errorf("Expected 2 initial dirs, got %d", len(state.InitialDirs))
	}

	if state.InitialDirs[0] != dir1Path {
		t.Errorf("Expected initial dir 0 to be %s, got %s", dir1Path, state.InitialDirs[0])
	}

	if state.InitialDirs[1] != dir2Path {
		t.Errorf("Expected initial dir 1 to be %s, got %s", dir2Path, state.InitialDirs[1])
	}
}

// TestDiffSearchStateIsActive_REQ_DIFF_SEARCH tests the IsActive method.
func TestDiffSearchStateIsActive_REQ_DIFF_SEARCH(t *testing.T) {
	// Nil state
	var nilState *DiffSearchState
	if nilState.IsActive() {
		t.Error("Nil state should not be active")
	}

	// Active state
	state := &DiffSearchState{Active: true}
	if !state.IsActive() {
		t.Error("State with Active=true should be active")
	}

	// Inactive state
	state.Active = false
	if state.IsActive() {
		t.Error("State with Active=false should not be active")
	}
}

// TestDiffSearchStateClear_REQ_DIFF_SEARCH tests the Clear method.
func TestDiffSearchStateClear_REQ_DIFF_SEARCH(t *testing.T) {
	state := &DiffSearchState{Active: true}
	state.Clear()

	if state.Active {
		t.Error("State should not be active after Clear")
	}

	// Clear on nil should not panic
	var nilState *DiffSearchState
	nilState.Clear() // Should not panic
}

// TestAtInitialDirs_REQ_DIFF_SEARCH tests checking if dirs are at initial positions.
func TestAtInitialDirs_REQ_DIFF_SEARCH(t *testing.T) {
	tmpDir := t.TempDir()
	dir1Path := filepath.Join(tmpDir, "dir1")
	dir2Path := filepath.Join(tmpDir, "dir2")
	subPath := filepath.Join(dir1Path, "sub")
	os.MkdirAll(dir1Path, 0755)
	os.MkdirAll(dir2Path, 0755)
	os.MkdirAll(subPath, 0755)

	dir1 := NewDirectory(0, 0, 10, 10)
	dir1.Chdir(dir1Path)
	dir2 := NewDirectory(0, 0, 10, 10)
	dir2.Chdir(dir2Path)

	dirs := []*Directory{dir1, dir2}
	state := NewDiffSearchState(dirs)

	// Should be at initial dirs
	if !state.AtInitialDirs(dirs) {
		t.Error("Should be at initial dirs")
	}

	// Change one dir
	dir1.Chdir(subPath)
	if state.AtInitialDirs(dirs) {
		t.Error("Should not be at initial dirs after chdir")
	}

	// Change back
	dir1.Chdir(dir1Path)
	if !state.AtInitialDirs(dirs) {
		t.Error("Should be at initial dirs after chdir back")
	}
}

// TestCollectAllNames_REQ_DIFF_SEARCH tests name collection from directories.
func TestCollectAllNames_REQ_DIFF_SEARCH(t *testing.T) {
	tmpDir := t.TempDir()
	dir1Path := filepath.Join(tmpDir, "dir1")
	dir2Path := filepath.Join(tmpDir, "dir2")
	os.MkdirAll(dir1Path, 0755)
	os.MkdirAll(dir2Path, 0755)

	// Create files in dir1
	os.WriteFile(filepath.Join(dir1Path, "alpha.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(dir1Path, "beta.txt"), []byte("b"), 0644)
	os.WriteFile(filepath.Join(dir1Path, "common.txt"), []byte("c"), 0644)

	// Create files in dir2
	os.WriteFile(filepath.Join(dir2Path, "common.txt"), []byte("c"), 0644)
	os.WriteFile(filepath.Join(dir2Path, "gamma.txt"), []byte("g"), 0644)

	dir1 := NewDirectory(0, 0, 10, 10)
	dir1.Chdir(dir1Path)
	dir2 := NewDirectory(0, 0, 10, 10)
	dir2.Chdir(dir2Path)

	dirs := []*Directory{dir1, dir2}
	names := CollectAllNames(dirs)

	// Should have 4 unique names, sorted alphabetically
	expected := []string{"alpha.txt", "beta.txt", "common.txt", "gamma.txt"}
	if len(names) != len(expected) {
		t.Errorf("Expected %d names, got %d: %v", len(expected), len(names), names)
	}

	for i, name := range expected {
		if i >= len(names) || names[i] != name {
			t.Errorf("Expected names[%d] = %s, got %v", i, name, names)
			break
		}
	}
}

// TestCheckDifferenceMissing_REQ_DIFF_SEARCH tests detection of missing files.
func TestCheckDifferenceMissing_REQ_DIFF_SEARCH(t *testing.T) {
	tmpDir := t.TempDir()
	dir1Path := filepath.Join(tmpDir, "dir1")
	dir2Path := filepath.Join(tmpDir, "dir2")
	os.MkdirAll(dir1Path, 0755)
	os.MkdirAll(dir2Path, 0755)

	// Create file only in dir1
	os.WriteFile(filepath.Join(dir1Path, "only_in_dir1.txt"), []byte("a"), 0644)

	// Create file in both
	os.WriteFile(filepath.Join(dir1Path, "common.txt"), []byte("c"), 0644)
	os.WriteFile(filepath.Join(dir2Path, "common.txt"), []byte("c"), 0644)

	dir1 := NewDirectory(0, 0, 10, 10)
	dir1.Chdir(dir1Path)
	dir2 := NewDirectory(0, 0, 10, 10)
	dir2.Chdir(dir2Path)

	dirs := []*Directory{dir1, dir2}

	// Check file only in dir1
	isDiff, reason, _ := CheckDifference("only_in_dir1.txt", dirs)
	if !isDiff {
		t.Error("Expected only_in_dir1.txt to be different (missing)")
	}
	if reason == "" {
		t.Error("Expected reason for difference")
	}

	// Check common file
	isDiff, reason, _ = CheckDifference("common.txt", dirs)
	if isDiff {
		t.Errorf("Expected common.txt to be same, got reason: %s", reason)
	}
}

// TestCheckDifferenceSizeMismatch_REQ_DIFF_SEARCH tests detection of size mismatches.
func TestCheckDifferenceSizeMismatch_REQ_DIFF_SEARCH(t *testing.T) {
	tmpDir := t.TempDir()
	dir1Path := filepath.Join(tmpDir, "dir1")
	dir2Path := filepath.Join(tmpDir, "dir2")
	os.MkdirAll(dir1Path, 0755)
	os.MkdirAll(dir2Path, 0755)

	// Create files with different sizes
	os.WriteFile(filepath.Join(dir1Path, "different.txt"), []byte("short"), 0644)
	os.WriteFile(filepath.Join(dir2Path, "different.txt"), []byte("much longer content"), 0644)

	// Create files with same size
	os.WriteFile(filepath.Join(dir1Path, "same.txt"), []byte("equal"), 0644)
	os.WriteFile(filepath.Join(dir2Path, "same.txt"), []byte("equal"), 0644)

	dir1 := NewDirectory(0, 0, 10, 10)
	dir1.Chdir(dir1Path)
	dir2 := NewDirectory(0, 0, 10, 10)
	dir2.Chdir(dir2Path)

	dirs := []*Directory{dir1, dir2}

	// Check file with different sizes
	isDiff, reason, _ := CheckDifference("different.txt", dirs)
	if !isDiff {
		t.Error("Expected different.txt to be different (size mismatch)")
	}
	if reason != "size mismatch" {
		t.Errorf("Expected 'size mismatch' reason, got: %s", reason)
	}

	// Check file with same sizes
	isDiff, reason, _ = CheckDifference("same.txt", dirs)
	if isDiff {
		t.Errorf("Expected same.txt to be same, got reason: %s", reason)
	}
}

// TestFindNextDifference_REQ_DIFF_SEARCH tests finding differences in order.
func TestFindNextDifference_REQ_DIFF_SEARCH(t *testing.T) {
	tmpDir := t.TempDir()
	dir1Path := filepath.Join(tmpDir, "dir1")
	dir2Path := filepath.Join(tmpDir, "dir2")
	os.MkdirAll(dir1Path, 0755)
	os.MkdirAll(dir2Path, 0755)

	// Create files - note alphabetic order: alpha < beta < gamma
	os.WriteFile(filepath.Join(dir1Path, "alpha.txt"), []byte("a"), 0644)    // same
	os.WriteFile(filepath.Join(dir2Path, "alpha.txt"), []byte("a"), 0644)    // same
	os.WriteFile(filepath.Join(dir1Path, "beta.txt"), []byte("diff1"), 0644) // different size
	os.WriteFile(filepath.Join(dir2Path, "beta.txt"), []byte("diff22"), 0644)
	os.WriteFile(filepath.Join(dir1Path, "gamma.txt"), []byte("g"), 0644) // only in dir1

	dir1 := NewDirectory(0, 0, 10, 10)
	dir1.Chdir(dir1Path)
	dir2 := NewDirectory(0, 0, 10, 10)
	dir2.Chdir(dir2Path)

	dirs := []*Directory{dir1, dir2}

	// Find first difference (should be beta.txt - size mismatch)
	result := FindNextDifference(dirs, "", true)
	if !result.Found {
		t.Fatal("Expected to find a difference")
	}
	if result.Name != "beta.txt" {
		t.Errorf("Expected first difference to be beta.txt, got %s", result.Name)
	}

	// Find next difference after beta (should be gamma.txt - missing)
	result = FindNextDifference(dirs, "beta.txt", true)
	if !result.Found {
		t.Fatal("Expected to find second difference")
	}
	if result.Name != "gamma.txt" {
		t.Errorf("Expected second difference to be gamma.txt, got %s", result.Name)
	}

	// Find next difference after gamma (should be none)
	result = FindNextDifference(dirs, "gamma.txt", true)
	if result.Found {
		t.Errorf("Expected no more differences, got %s", result.Name)
	}
}

// TestFindNextDifferenceNoDiffs_REQ_DIFF_SEARCH tests when there are no differences.
func TestFindNextDifferenceNoDiffs_REQ_DIFF_SEARCH(t *testing.T) {
	tmpDir := t.TempDir()
	dir1Path := filepath.Join(tmpDir, "dir1")
	dir2Path := filepath.Join(tmpDir, "dir2")
	os.MkdirAll(dir1Path, 0755)
	os.MkdirAll(dir2Path, 0755)

	// Create identical files in both directories
	os.WriteFile(filepath.Join(dir1Path, "file1.txt"), []byte("same"), 0644)
	os.WriteFile(filepath.Join(dir2Path, "file1.txt"), []byte("same"), 0644)
	os.WriteFile(filepath.Join(dir1Path, "file2.txt"), []byte("also same"), 0644)
	os.WriteFile(filepath.Join(dir2Path, "file2.txt"), []byte("also same"), 0644)

	dir1 := NewDirectory(0, 0, 10, 10)
	dir1.Chdir(dir1Path)
	dir2 := NewDirectory(0, 0, 10, 10)
	dir2.Chdir(dir2Path)

	dirs := []*Directory{dir1, dir2}

	result := FindNextDifference(dirs, "", true)
	if result.Found {
		t.Errorf("Expected no differences, got %s: %s", result.Name, result.Reason)
	}
}

// TestCollectSubdirNames_REQ_DIFF_SEARCH tests subdirectory collection.
func TestCollectSubdirNames_REQ_DIFF_SEARCH(t *testing.T) {
	tmpDir := t.TempDir()
	dir1Path := filepath.Join(tmpDir, "dir1")
	os.MkdirAll(dir1Path, 0755)

	// Create subdirectories
	os.MkdirAll(filepath.Join(dir1Path, "subA"), 0755)
	os.MkdirAll(filepath.Join(dir1Path, "subB"), 0755)
	// Create a file (should not be included)
	os.WriteFile(filepath.Join(dir1Path, "file.txt"), []byte("f"), 0644)

	dir1 := NewDirectory(0, 0, 10, 10)
	dir1.Chdir(dir1Path)

	dirs := []*Directory{dir1}
	subdirs := CollectSubdirNames(dirs)

	if len(subdirs) != 2 {
		t.Errorf("Expected 2 subdirs, got %d: %v", len(subdirs), subdirs)
	}

	expected := []string{"subA", "subB"}
	for i, name := range expected {
		if i >= len(subdirs) || subdirs[i] != name {
			t.Errorf("Expected subdirs[%d] = %s, got %v", i, name, subdirs)
			break
		}
	}
}

// TestFirstSubdirInAll_REQ_DIFF_SEARCH tests finding first common subdirectory.
func TestFirstSubdirInAll_REQ_DIFF_SEARCH(t *testing.T) {
	tmpDir := t.TempDir()
	dir1Path := filepath.Join(tmpDir, "dir1")
	dir2Path := filepath.Join(tmpDir, "dir2")
	os.MkdirAll(dir1Path, 0755)
	os.MkdirAll(dir2Path, 0755)

	// Create subdirs - alpha only in dir1, beta in both, gamma only in dir2
	os.MkdirAll(filepath.Join(dir1Path, "alpha"), 0755)
	os.MkdirAll(filepath.Join(dir1Path, "beta"), 0755)
	os.MkdirAll(filepath.Join(dir2Path, "beta"), 0755)
	os.MkdirAll(filepath.Join(dir2Path, "gamma"), 0755)

	dir1 := NewDirectory(0, 0, 10, 10)
	dir1.Chdir(dir1Path)
	dir2 := NewDirectory(0, 0, 10, 10)
	dir2.Chdir(dir2Path)

	dirs := []*Directory{dir1, dir2}

	name, found := FirstSubdirInAll(dirs)
	if !found {
		t.Fatal("Expected to find a common subdirectory")
	}
	if name != "beta" {
		t.Errorf("Expected first common subdir to be 'beta', got '%s'", name)
	}
}

// TestFirstSubdirInAllNone_REQ_DIFF_SEARCH tests when no common subdirectory exists.
func TestFirstSubdirInAllNone_REQ_DIFF_SEARCH(t *testing.T) {
	tmpDir := t.TempDir()
	dir1Path := filepath.Join(tmpDir, "dir1")
	dir2Path := filepath.Join(tmpDir, "dir2")
	os.MkdirAll(dir1Path, 0755)
	os.MkdirAll(dir2Path, 0755)

	// Create different subdirs in each
	os.MkdirAll(filepath.Join(dir1Path, "onlyDir1"), 0755)
	os.MkdirAll(filepath.Join(dir2Path, "onlyDir2"), 0755)

	dir1 := NewDirectory(0, 0, 10, 10)
	dir1.Chdir(dir1Path)
	dir2 := NewDirectory(0, 0, 10, 10)
	dir2.Chdir(dir2Path)

	dirs := []*Directory{dir1, dir2}

	_, found := FirstSubdirInAll(dirs)
	if found {
		t.Error("Expected no common subdirectory")
	}
}

// TestFindNextSubdirInAll_REQ_DIFF_SEARCH tests finding next common subdirectory after a position.
func TestFindNextSubdirInAll_REQ_DIFF_SEARCH(t *testing.T) {
	tmpDir := t.TempDir()
	dir1Path := filepath.Join(tmpDir, "dir1")
	dir2Path := filepath.Join(tmpDir, "dir2")
	os.MkdirAll(dir1Path, 0755)
	os.MkdirAll(dir2Path, 0755)

	// Create subdirs in both: alpha, beta, gamma all in both
	// delta only in dir1, epsilon only in dir2
	os.MkdirAll(filepath.Join(dir1Path, "alpha"), 0755)
	os.MkdirAll(filepath.Join(dir2Path, "alpha"), 0755)
	os.MkdirAll(filepath.Join(dir1Path, "beta"), 0755)
	os.MkdirAll(filepath.Join(dir2Path, "beta"), 0755)
	os.MkdirAll(filepath.Join(dir1Path, "delta"), 0755)   // only in dir1
	os.MkdirAll(filepath.Join(dir2Path, "epsilon"), 0755) // only in dir2
	os.MkdirAll(filepath.Join(dir1Path, "gamma"), 0755)
	os.MkdirAll(filepath.Join(dir2Path, "gamma"), 0755)

	dir1 := NewDirectory(0, 0, 10, 10)
	dir1.Chdir(dir1Path)
	dir2 := NewDirectory(0, 0, 10, 10)
	dir2.Chdir(dir2Path)

	dirs := []*Directory{dir1, dir2}

	// Test: empty startAfter should behave like FirstSubdirInAll
	name, found := FindNextSubdirInAll(dirs, "")
	if !found {
		t.Fatal("Expected to find a common subdirectory with empty startAfter")
	}
	if name != "alpha" {
		t.Errorf("Expected first common subdir to be 'alpha', got '%s'", name)
	}

	// Test: find next after alpha (should be beta)
	name, found = FindNextSubdirInAll(dirs, "alpha")
	if !found {
		t.Fatal("Expected to find next common subdirectory after alpha")
	}
	if name != "beta" {
		t.Errorf("Expected next common subdir after alpha to be 'beta', got '%s'", name)
	}

	// Test: find next after beta (should be gamma, skipping delta/epsilon which aren't in both)
	name, found = FindNextSubdirInAll(dirs, "beta")
	if !found {
		t.Fatal("Expected to find next common subdirectory after beta")
	}
	if name != "gamma" {
		t.Errorf("Expected next common subdir after beta to be 'gamma', got '%s'", name)
	}

	// Test: find next after gamma (should be none)
	_, found = FindNextSubdirInAll(dirs, "gamma")
	if found {
		t.Error("Expected no common subdirectory after gamma")
	}
}

// TestFindNextSubdirInAllSkipsNonCommon_REQ_DIFF_SEARCH tests that non-common subdirs are skipped.
func TestFindNextSubdirInAllSkipsNonCommon_REQ_DIFF_SEARCH(t *testing.T) {
	tmpDir := t.TempDir()
	dir1Path := filepath.Join(tmpDir, "dir1")
	dir2Path := filepath.Join(tmpDir, "dir2")
	os.MkdirAll(dir1Path, 0755)
	os.MkdirAll(dir2Path, 0755)

	// Create: aaa in both, bbb only in dir1, ccc only in dir2, ddd in both
	os.MkdirAll(filepath.Join(dir1Path, "aaa"), 0755)
	os.MkdirAll(filepath.Join(dir2Path, "aaa"), 0755)
	os.MkdirAll(filepath.Join(dir1Path, "bbb"), 0755) // only in dir1
	os.MkdirAll(filepath.Join(dir2Path, "ccc"), 0755) // only in dir2
	os.MkdirAll(filepath.Join(dir1Path, "ddd"), 0755)
	os.MkdirAll(filepath.Join(dir2Path, "ddd"), 0755)

	dir1 := NewDirectory(0, 0, 10, 10)
	dir1.Chdir(dir1Path)
	dir2 := NewDirectory(0, 0, 10, 10)
	dir2.Chdir(dir2Path)

	dirs := []*Directory{dir1, dir2}

	// After aaa, next common should be ddd (skipping bbb and ccc)
	name, found := FindNextSubdirInAll(dirs, "aaa")
	if !found {
		t.Fatal("Expected to find next common subdirectory after aaa")
	}
	if name != "ddd" {
		t.Errorf("Expected to skip bbb/ccc and find 'ddd', got '%s'", name)
	}
}

// TestFindNextSubdirInAllNoCommon_REQ_DIFF_SEARCH tests when no common subdirectory exists after position.
func TestFindNextSubdirInAllNoCommon_REQ_DIFF_SEARCH(t *testing.T) {
	tmpDir := t.TempDir()
	dir1Path := filepath.Join(tmpDir, "dir1")
	dir2Path := filepath.Join(tmpDir, "dir2")
	os.MkdirAll(dir1Path, 0755)
	os.MkdirAll(dir2Path, 0755)

	// Create: only aaa in both, rest are non-common
	os.MkdirAll(filepath.Join(dir1Path, "aaa"), 0755)
	os.MkdirAll(filepath.Join(dir2Path, "aaa"), 0755)
	os.MkdirAll(filepath.Join(dir1Path, "bbb"), 0755) // only in dir1
	os.MkdirAll(filepath.Join(dir2Path, "ccc"), 0755) // only in dir2

	dir1 := NewDirectory(0, 0, 10, 10)
	dir1.Chdir(dir1Path)
	dir2 := NewDirectory(0, 0, 10, 10)
	dir2.Chdir(dir2Path)

	dirs := []*Directory{dir1, dir2}

	// After aaa, there should be no more common subdirectories
	_, found := FindNextSubdirInAll(dirs, "aaa")
	if found {
		t.Error("Expected no common subdirectory after aaa")
	}
}

// TestAlphabeticSorting_REQ_DIFF_SEARCH tests that names are sorted alphabetically.
func TestAlphabeticSorting_REQ_DIFF_SEARCH(t *testing.T) {
	tmpDir := t.TempDir()
	dirPath := filepath.Join(tmpDir, "dir")
	os.MkdirAll(dirPath, 0755)

	// Create files out of alphabetic order
	os.WriteFile(filepath.Join(dirPath, "charlie.txt"), []byte("c"), 0644)
	os.WriteFile(filepath.Join(dirPath, "alpha.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(dirPath, "delta.txt"), []byte("d"), 0644)
	os.WriteFile(filepath.Join(dirPath, "bravo.txt"), []byte("b"), 0644)

	dir := NewDirectory(0, 0, 10, 10)
	dir.Chdir(dirPath)

	dirs := []*Directory{dir}
	names := CollectAllNames(dirs)

	// Expected order: alphabetically sorted
	expected := []string{"alpha.txt", "bravo.txt", "charlie.txt", "delta.txt"}
	if len(names) != len(expected) {
		t.Errorf("Expected %d names, got %d: %v", len(expected), len(names), names)
		return
	}

	for i, name := range expected {
		if names[i] != name {
			t.Errorf("Expected names[%d] = %s, got %s (full list: %v)", i, name, names[i], names)
		}
	}
}
