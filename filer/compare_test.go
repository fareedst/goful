// Package filer comparison index tests.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
package filer

import (
	"os"
	"testing"
	"time"

	"github.com/anmitsu/goful/widget"
)

// mockFileInfo implements os.FileInfo for testing.
type mockFileInfo struct {
	name    string
	size    int64
	modTime time.Time
	isDir   bool
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return m.size }
func (m *mockFileInfo) Mode() os.FileMode  { return 0644 }
func (m *mockFileInfo) ModTime() time.Time { return m.modTime }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }

// mockFileStat creates a mock FileStat for testing.
func mockFileStat(name string, size int64, modTime time.Time) *FileStat {
	info := &mockFileInfo{
		name:    name,
		size:    size,
		modTime: modTime,
		isDir:   false,
	}
	return &FileStat{
		FileInfo: info,
		stat:     info,
		path:     name,
		name:     name,
		display:  name,
		marked:   false,
	}
}

// mockDirectory creates a Directory with the given files for testing.
func mockDirectory(files ...*FileStat) *Directory {
	d := NewDirectory(0, 0, 100, 50)
	drawers := make([]widget.Drawer, len(files))
	for i, f := range files {
		drawers[i] = f
	}
	d.SetList(drawers)
	return d
}

// TestBuildComparisonIndex_SingleDirectory tests that single directories return nil index.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:FILE_COMPARISON_INDEX]
func TestBuildComparisonIndex_SingleDirectory_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	now := time.Now()
	dirs := []*Directory{mockDirectory(mockFileStat("file.txt", 100, now))}
	idx := BuildComparisonIndex(dirs)
	if idx != nil {
		t.Error("expected nil index for single directory")
	}
}

// TestBuildComparisonIndex_NoDirs tests that empty directories return nil index.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:FILE_COMPARISON_INDEX]
func TestBuildComparisonIndex_NoDirs_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	dirs := []*Directory{}
	idx := BuildComparisonIndex(dirs)
	if idx != nil {
		t.Error("expected nil index for empty directories")
	}
}

// TestBuildComparisonIndex_UniqueFiles tests that unique files are not indexed.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:FILE_COMPARISON_INDEX]
func TestBuildComparisonIndex_UniqueFiles_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	now := time.Now()

	// Create two directories with different files
	dir1 := mockDirectory(mockFileStat("file1.txt", 100, now))
	dir2 := mockDirectory(mockFileStat("file2.txt", 200, now))

	dirs := []*Directory{dir1, dir2}
	idx := BuildComparisonIndex(dirs)

	// Files are unique, so no comparison states should exist
	if idx.Get(0, "file1.txt") != nil {
		t.Error("unique file1.txt should not be in index")
	}
	if idx.Get(1, "file2.txt") != nil {
		t.Error("unique file2.txt should not be in index")
	}
}

// TestBuildComparisonIndex_CommonFiles tests that common files are indexed with correct states.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:FILE_COMPARISON_INDEX]
func TestBuildComparisonIndex_CommonFiles_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-1 * time.Hour)

	// Create two directories with the same file, different sizes and times
	dir1 := mockDirectory(mockFileStat("common.txt", 100, earlier))
	dir2 := mockDirectory(mockFileStat("common.txt", 200, now))

	dirs := []*Directory{dir1, dir2}
	idx := BuildComparisonIndex(dirs)

	// Check dir1's entry
	state1 := idx.Get(0, "common.txt")
	if state1 == nil {
		t.Fatal("expected comparison state for dir1 common.txt")
	}
	if !state1.NamePresent {
		t.Error("expected NamePresent=true for common file")
	}
	if state1.SizeState != SizeSmallest {
		t.Errorf("expected SizeSmallest for dir1, got %v", state1.SizeState)
	}
	if state1.TimeState != TimeEarliest {
		t.Errorf("expected TimeEarliest for dir1, got %v", state1.TimeState)
	}

	// Check dir2's entry
	state2 := idx.Get(1, "common.txt")
	if state2 == nil {
		t.Fatal("expected comparison state for dir2 common.txt")
	}
	if !state2.NamePresent {
		t.Error("expected NamePresent=true for common file")
	}
	if state2.SizeState != SizeLargest {
		t.Errorf("expected SizeLargest for dir2, got %v", state2.SizeState)
	}
	if state2.TimeState != TimeLatest {
		t.Errorf("expected TimeLatest for dir2, got %v", state2.TimeState)
	}
}

// TestBuildComparisonIndex_EqualSizesAndTimes tests equal size and time detection.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:FILE_COMPARISON_INDEX]
func TestBuildComparisonIndex_EqualSizesAndTimes_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	now := time.Now().Truncate(time.Second) // Truncate to second precision

	// Create two directories with identical files
	dir1 := mockDirectory(mockFileStat("same.txt", 100, now))
	dir2 := mockDirectory(mockFileStat("same.txt", 100, now))

	dirs := []*Directory{dir1, dir2}
	idx := BuildComparisonIndex(dirs)

	// Both should show equal states
	state1 := idx.Get(0, "same.txt")
	if state1 == nil {
		t.Fatal("expected comparison state for dir1")
	}
	if state1.SizeState != SizeEqual {
		t.Errorf("expected SizeEqual, got %v", state1.SizeState)
	}
	if state1.TimeState != TimeEqual {
		t.Errorf("expected TimeEqual, got %v", state1.TimeState)
	}

	state2 := idx.Get(1, "same.txt")
	if state2 == nil {
		t.Fatal("expected comparison state for dir2")
	}
	if state2.SizeState != SizeEqual {
		t.Errorf("expected SizeEqual, got %v", state2.SizeState)
	}
	if state2.TimeState != TimeEqual {
		t.Errorf("expected TimeEqual, got %v", state2.TimeState)
	}
}

// TestBuildComparisonIndex_ThreeDirectories tests comparison across three directories.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:FILE_COMPARISON_INDEX]
func TestBuildComparisonIndex_ThreeDirectories_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-2 * time.Hour)
	middle := now.Add(-1 * time.Hour)

	dir1 := mockDirectory(mockFileStat("file.txt", 50, earlier)) // smallest, earliest
	dir2 := mockDirectory(mockFileStat("file.txt", 100, middle)) // middle
	dir3 := mockDirectory(mockFileStat("file.txt", 200, now))    // largest, latest

	dirs := []*Directory{dir1, dir2, dir3}
	idx := BuildComparisonIndex(dirs)

	// Check dir1 - smallest and earliest
	state1 := idx.Get(0, "file.txt")
	if state1.SizeState != SizeSmallest {
		t.Errorf("dir1 expected SizeSmallest, got %v", state1.SizeState)
	}
	if state1.TimeState != TimeEarliest {
		t.Errorf("dir1 expected TimeEarliest, got %v", state1.TimeState)
	}

	// Check dir2 - middle
	state2 := idx.Get(1, "file.txt")
	if state2.SizeState != SizeMiddle {
		t.Errorf("dir2 expected SizeMiddle, got %v", state2.SizeState)
	}
	if state2.TimeState != TimeMiddle {
		t.Errorf("dir2 expected TimeMiddle, got %v", state2.TimeState)
	}

	// Check dir3 - largest and latest
	state3 := idx.Get(2, "file.txt")
	if state3.SizeState != SizeLargest {
		t.Errorf("dir3 expected SizeLargest, got %v", state3.SizeState)
	}
	if state3.TimeState != TimeLatest {
		t.Errorf("dir3 expected TimeLatest, got %v", state3.TimeState)
	}
}

// TestBuildComparisonIndex_CaseSensitive tests that file name comparison is case-sensitive.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:FILE_COMPARISON_INDEX]
func TestBuildComparisonIndex_CaseSensitive_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	now := time.Now()

	dir1 := mockDirectory(mockFileStat("File.txt", 100, now))
	dir2 := mockDirectory(mockFileStat("file.txt", 100, now))

	dirs := []*Directory{dir1, dir2}
	idx := BuildComparisonIndex(dirs)

	// Different case = different files, should not be in index
	if idx.Get(0, "File.txt") != nil {
		t.Error("File.txt should not be indexed (no match)")
	}
	if idx.Get(1, "file.txt") != nil {
		t.Error("file.txt should not be indexed (no match)")
	}
}

// TestComparisonIndex_Get_NilIndex tests that Get on nil index returns nil.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:FILE_COMPARISON_INDEX]
func TestComparisonIndex_Get_NilIndex_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	var idx *ComparisonIndex
	if idx.Get(0, "file.txt") != nil {
		t.Error("expected nil from nil index")
	}
}

// TestComparisonIndex_Clear tests that Clear removes all entries.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:FILE_COMPARISON_INDEX]
func TestComparisonIndex_Clear_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	now := time.Now()

	dir1 := mockDirectory(mockFileStat("file.txt", 100, now))
	dir2 := mockDirectory(mockFileStat("file.txt", 200, now))

	dirs := []*Directory{dir1, dir2}
	idx := BuildComparisonIndex(dirs)

	if idx.Get(0, "file.txt") == nil {
		t.Fatal("expected entry before clear")
	}

	idx.Clear()

	if idx.Get(0, "file.txt") != nil {
		t.Error("expected nil after clear")
	}
}

// TestTimeComparisonPrecision tests that time comparison uses second precision.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:FILE_COMPARISON_INDEX]
func TestTimeComparisonPrecision_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	// Same second, different nanoseconds
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	time1 := base.Add(100 * time.Millisecond)
	time2 := base.Add(500 * time.Millisecond)

	dir1 := mockDirectory(mockFileStat("file.txt", 100, time1))
	dir2 := mockDirectory(mockFileStat("file.txt", 100, time2))

	dirs := []*Directory{dir1, dir2}
	idx := BuildComparisonIndex(dirs)

	state1 := idx.Get(0, "file.txt")
	state2 := idx.Get(1, "file.txt")

	// Both should be TimeEqual since they're in the same second
	if state1.TimeState != TimeEqual {
		t.Errorf("dir1 expected TimeEqual (same second), got %v", state1.TimeState)
	}
	if state2.TimeState != TimeEqual {
		t.Errorf("dir2 expected TimeEqual (same second), got %v", state2.TimeState)
	}
}
