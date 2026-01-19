// Package filer comparison index tests.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
package filer

import (
	"os"
	"testing"
	"time"

	"github.com/fareedst/goful/widget"
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

// TestCalculateFileDigest tests digest calculation for a known file.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:DIGEST_COMPARISON]
func TestCalculateFileDigest_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	// Create a temporary file with known content
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/test.txt"
	content := []byte("Hello, World!")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	digest, err := CalculateFileDigest(tmpFile)
	if err != nil {
		t.Fatalf("CalculateFileDigest failed: %v", err)
	}

	// Verify digest is non-zero
	if digest == 0 {
		t.Error("expected non-zero digest")
	}

	// Verify same content produces same digest
	digest2, err := CalculateFileDigest(tmpFile)
	if err != nil {
		t.Fatalf("second CalculateFileDigest failed: %v", err)
	}
	if digest != digest2 {
		t.Errorf("same file produced different digests: %d != %d", digest, digest2)
	}
}

// TestCalculateFileDigest_DifferentContent tests that different content produces different digests.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:DIGEST_COMPARISON]
func TestCalculateFileDigest_DifferentContent_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	tmpDir := t.TempDir()

	file1 := tmpDir + "/file1.txt"
	file2 := tmpDir + "/file2.txt"

	// Same size, different content
	if err := os.WriteFile(file1, []byte("AAAA"), 0644); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("BBBB"), 0644); err != nil {
		t.Fatalf("failed to create file2: %v", err)
	}

	digest1, err := CalculateFileDigest(file1)
	if err != nil {
		t.Fatalf("CalculateFileDigest file1 failed: %v", err)
	}

	digest2, err := CalculateFileDigest(file2)
	if err != nil {
		t.Fatalf("CalculateFileDigest file2 failed: %v", err)
	}

	if digest1 == digest2 {
		t.Errorf("different content should produce different digests: both got %d", digest1)
	}
}

// TestCalculateFileDigest_NonExistent tests error handling for non-existent files.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:DIGEST_COMPARISON]
func TestCalculateFileDigest_NonExistent_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	_, err := CalculateFileDigest("/nonexistent/path/to/file.txt")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

// TestDigestCompare_States tests that DigestCompare constants are distinct.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:DIGEST_COMPARISON]
func TestDigestCompare_States_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	states := []DigestCompare{DigestUnknown, DigestEqual, DigestDifferent, DigestNA}
	seen := make(map[DigestCompare]bool)
	for _, s := range states {
		if seen[s] {
			t.Errorf("duplicate DigestCompare value: %v", s)
		}
		seen[s] = true
	}
}

// TestCompareState_DigestField tests that CompareState includes DigestState field.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:DIGEST_COMPARISON]
func TestCompareState_DigestField_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	state := &CompareState{
		NamePresent: true,
		SizeState:   SizeEqual,
		TimeState:   TimeEqual,
		DigestState: DigestUnknown,
	}

	if state.DigestState != DigestUnknown {
		t.Errorf("expected DigestUnknown, got %v", state.DigestState)
	}

	state.DigestState = DigestEqual
	if state.DigestState != DigestEqual {
		t.Errorf("expected DigestEqual, got %v", state.DigestState)
	}
}

// TestUpdateDigestStates_NilIndex tests that nil index returns 0.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:DIGEST_COMPARISON]
func TestUpdateDigestStates_NilIndex_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	var idx *ComparisonIndex
	count := idx.UpdateDigestStates("file.txt", nil)
	if count != 0 {
		t.Errorf("expected 0 from nil index, got %d", count)
	}
}

// TestUpdateDigestStates_NoMatchingFiles tests that missing files return 0.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:DIGEST_COMPARISON]
func TestUpdateDigestStates_NoMatchingFiles_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	now := time.Now()
	dir1 := mockDirectory(mockFileStat("file1.txt", 100, now))
	dir2 := mockDirectory(mockFileStat("file2.txt", 100, now))

	dirs := []*Directory{dir1, dir2}
	idx := BuildComparisonIndex(dirs)

	count := idx.UpdateDigestStates("nonexistent.txt", dirs)
	if count != 0 {
		t.Errorf("expected 0 for non-existent file, got %d", count)
	}
}

// TestUpdateDigestStates_DifferentSizes tests that files with unique sizes get DigestNA.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:DIGEST_COMPARISON]
func TestUpdateDigestStates_DifferentSizes_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	now := time.Now()
	dir1 := mockDirectory(mockFileStat("file.txt", 100, now))
	dir2 := mockDirectory(mockFileStat("file.txt", 200, now))

	dirs := []*Directory{dir1, dir2}
	idx := BuildComparisonIndex(dirs)

	count := idx.UpdateDigestStates("file.txt", dirs)
	// Each file has a unique size, so digest calculation should be skipped for both
	if count != 0 {
		t.Errorf("expected 0 for files with unique sizes, got %d", count)
	}

	// Both states should be DigestNA (each has unique size)
	state1 := idx.Get(0, "file.txt")
	state2 := idx.Get(1, "file.txt")

	if state1.DigestState != DigestNA {
		t.Errorf("dir1 expected DigestNA, got %v", state1.DigestState)
	}
	if state2.DigestState != DigestNA {
		t.Errorf("dir2 expected DigestNA, got %v", state2.DigestState)
	}
}

// TestUpdateDigestStates_MixedSizes tests that files are grouped by size for digest comparison.
// Files with matching sizes get compared, files with unique sizes get DigestNA.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:DIGEST_COMPARISON]
func TestUpdateDigestStates_MixedSizes_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	now := time.Now()
	tmpDir := t.TempDir()

	// Create actual files so digest calculation works
	file1 := tmpDir + "/dir1/file.txt"
	file2 := tmpDir + "/dir2/file.txt"
	file3 := tmpDir + "/dir3/file.txt"

	os.MkdirAll(tmpDir+"/dir1", 0755)
	os.MkdirAll(tmpDir+"/dir2", 0755)
	os.MkdirAll(tmpDir+"/dir3", 0755)

	// dir1 and dir2 have same content (same size), dir3 has different size
	os.WriteFile(file1, []byte("AAAA"), 0644)
	os.WriteFile(file2, []byte("AAAA"), 0644)
	os.WriteFile(file3, []byte("BBBBBBBB"), 0644) // 8 bytes vs 4 bytes

	// Create mock file stats with correct sizes
	fs1 := &FileStat{
		FileInfo: &mockFileInfo{name: "file.txt", size: 4, modTime: now},
		stat:     &mockFileInfo{name: "file.txt", size: 4, modTime: now},
		path:     file1,
		name:     "file.txt",
	}
	fs2 := &FileStat{
		FileInfo: &mockFileInfo{name: "file.txt", size: 4, modTime: now},
		stat:     &mockFileInfo{name: "file.txt", size: 4, modTime: now},
		path:     file2,
		name:     "file.txt",
	}
	fs3 := &FileStat{
		FileInfo: &mockFileInfo{name: "file.txt", size: 8, modTime: now},
		stat:     &mockFileInfo{name: "file.txt", size: 8, modTime: now},
		path:     file3,
		name:     "file.txt",
	}

	dir1 := mockDirectory(fs1)
	dir2 := mockDirectory(fs2)
	dir3 := mockDirectory(fs3)

	dirs := []*Directory{dir1, dir2, dir3}
	idx := BuildComparisonIndex(dirs)

	count := idx.UpdateDigestStates("file.txt", dirs)
	// dir1 and dir2 have same size and should be compared (2 files)
	// dir3 has unique size and should get DigestNA
	if count != 2 {
		t.Errorf("expected 2 files processed, got %d", count)
	}

	state1 := idx.Get(0, "file.txt")
	state2 := idx.Get(1, "file.txt")
	state3 := idx.Get(2, "file.txt")

	// dir1 and dir2 should have DigestEqual (same content)
	if state1.DigestState != DigestEqual {
		t.Errorf("dir1 expected DigestEqual, got %v", state1.DigestState)
	}
	if state2.DigestState != DigestEqual {
		t.Errorf("dir2 expected DigestEqual, got %v", state2.DigestState)
	}
	// dir3 should have DigestNA (unique size)
	if state3.DigestState != DigestNA {
		t.Errorf("dir3 expected DigestNA, got %v", state3.DigestState)
	}
}

// TestUpdateDigestStates_MixedSizesDifferentContent tests digest comparison with same size but different content.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:DIGEST_COMPARISON]
func TestUpdateDigestStates_MixedSizesDifferentContent_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	now := time.Now()
	tmpDir := t.TempDir()

	file1 := tmpDir + "/dir1/file.txt"
	file2 := tmpDir + "/dir2/file.txt"
	file3 := tmpDir + "/dir3/file.txt"

	os.MkdirAll(tmpDir+"/dir1", 0755)
	os.MkdirAll(tmpDir+"/dir2", 0755)
	os.MkdirAll(tmpDir+"/dir3", 0755)

	// dir1 and dir2 have same size but different content
	os.WriteFile(file1, []byte("AAAA"), 0644)
	os.WriteFile(file2, []byte("BBBB"), 0644) // Same size (4 bytes), different content
	os.WriteFile(file3, []byte("CCCCCCCC"), 0644)

	fs1 := &FileStat{
		FileInfo: &mockFileInfo{name: "file.txt", size: 4, modTime: now},
		stat:     &mockFileInfo{name: "file.txt", size: 4, modTime: now},
		path:     file1,
		name:     "file.txt",
	}
	fs2 := &FileStat{
		FileInfo: &mockFileInfo{name: "file.txt", size: 4, modTime: now},
		stat:     &mockFileInfo{name: "file.txt", size: 4, modTime: now},
		path:     file2,
		name:     "file.txt",
	}
	fs3 := &FileStat{
		FileInfo: &mockFileInfo{name: "file.txt", size: 8, modTime: now},
		stat:     &mockFileInfo{name: "file.txt", size: 8, modTime: now},
		path:     file3,
		name:     "file.txt",
	}

	dir1 := mockDirectory(fs1)
	dir2 := mockDirectory(fs2)
	dir3 := mockDirectory(fs3)

	dirs := []*Directory{dir1, dir2, dir3}
	idx := BuildComparisonIndex(dirs)

	count := idx.UpdateDigestStates("file.txt", dirs)
	if count != 2 {
		t.Errorf("expected 2 files processed, got %d", count)
	}

	state1 := idx.Get(0, "file.txt")
	state2 := idx.Get(1, "file.txt")
	state3 := idx.Get(2, "file.txt")

	// dir1 and dir2 should have DigestDifferent (same size, different content)
	if state1.DigestState != DigestDifferent {
		t.Errorf("dir1 expected DigestDifferent, got %v", state1.DigestState)
	}
	if state2.DigestState != DigestDifferent {
		t.Errorf("dir2 expected DigestDifferent, got %v", state2.DigestState)
	}
	// dir3 should have DigestNA (unique size)
	if state3.DigestState != DigestNA {
		t.Errorf("dir3 expected DigestNA, got %v", state3.DigestState)
	}
}

// TestSharedFilenames_NilIndex_REQ_TOOLBAR_COMPARE_BUTTON tests SharedFilenames with nil index.
// [REQ:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_COMPARE_BUTTON]
func TestSharedFilenames_NilIndex_REQ_TOOLBAR_COMPARE_BUTTON(t *testing.T) {
	var idx *ComparisonIndex = nil
	names := idx.SharedFilenames()
	if names != nil {
		t.Errorf("expected nil for nil index, got %v", names)
	}
}

// TestSharedFilenames_EmptyIndex_REQ_TOOLBAR_COMPARE_BUTTON tests SharedFilenames with empty index.
// [REQ:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_COMPARE_BUTTON]
func TestSharedFilenames_EmptyIndex_REQ_TOOLBAR_COMPARE_BUTTON(t *testing.T) {
	idx := NewComparisonIndex()
	names := idx.SharedFilenames()
	if len(names) != 0 {
		t.Errorf("expected empty slice for empty index, got %v", names)
	}
}

// TestSharedFilenames_WithFiles_REQ_TOOLBAR_COMPARE_BUTTON tests SharedFilenames returns shared filenames.
// [REQ:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_COMPARE_BUTTON]
func TestSharedFilenames_WithFiles_REQ_TOOLBAR_COMPARE_BUTTON(t *testing.T) {
	now := time.Now()

	// Create directories with shared and unique files
	dir1 := mockDirectory(
		mockFileStat("shared1.txt", 100, now),
		mockFileStat("shared2.txt", 200, now),
		mockFileStat("unique1.txt", 300, now),
	)
	dir2 := mockDirectory(
		mockFileStat("shared1.txt", 100, now),
		mockFileStat("shared2.txt", 200, now),
		mockFileStat("unique2.txt", 400, now),
	)

	dirs := []*Directory{dir1, dir2}
	idx := BuildComparisonIndex(dirs)

	names := idx.SharedFilenames()

	// Should have 2 shared filenames (shared1.txt and shared2.txt)
	if len(names) != 2 {
		t.Errorf("expected 2 shared filenames, got %d: %v", len(names), names)
	}

	// Check that both shared files are in the result
	nameSet := make(map[string]bool)
	for _, name := range names {
		nameSet[name] = true
	}

	if !nameSet["shared1.txt"] {
		t.Error("expected shared1.txt in shared filenames")
	}
	if !nameSet["shared2.txt"] {
		t.Error("expected shared2.txt in shared filenames")
	}
	if nameSet["unique1.txt"] {
		t.Error("unique1.txt should not be in shared filenames")
	}
	if nameSet["unique2.txt"] {
		t.Error("unique2.txt should not be in shared filenames")
	}
}
