// Package filer comparison index for cross-directory file comparison.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
package filer

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/cespare/xxhash/v2"
)

// SizeCompare represents the comparison state for file sizes.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
type SizeCompare int

const (
	SizeUnknown  SizeCompare = iota
	SizeEqual                // All files have the same size
	SizeSmallest             // This file has the smallest size
	SizeLargest              // This file has the largest size
	SizeMiddle               // This file is neither smallest nor largest
)

// TimeCompare represents the comparison state for modification times.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
type TimeCompare int

const (
	TimeUnknown  TimeCompare = iota
	TimeEqual                // All files have the same modification time
	TimeEarliest             // This file has the earliest modification time
	TimeLatest               // This file has the latest modification time
	TimeMiddle               // This file is neither earliest nor latest
)

// DigestCompare represents the comparison state for file content digests.
// [IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
type DigestCompare int

const (
	DigestUnknown   DigestCompare = iota // Digest not yet calculated
	DigestEqual                          // Same digest (content identical)
	DigestDifferent                      // Different digest despite equal size
	DigestNA                             // Not applicable (sizes differ)
)

// CompareState holds the comparison state for a single file in a directory.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
type CompareState struct {
	NamePresent bool          // File name appears in multiple directories
	SizeState   SizeCompare   // Size comparison state
	TimeState   TimeCompare   // Time comparison state
	DigestState DigestCompare // Digest comparison state (on-demand calculation)
}

// fileEntry represents a file in a specific directory for comparison purposes.
type fileEntry struct {
	dirIndex int
	size     int64
	modTime  time.Time
}

// ComparisonIndex holds cached comparison states for files across directories.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
type ComparisonIndex struct {
	mu    sync.RWMutex
	cache map[string]map[int]*CompareState // map[filename][dirIndex]*CompareState
}

// NewComparisonIndex creates a new empty comparison index.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func NewComparisonIndex() *ComparisonIndex {
	return &ComparisonIndex{
		cache: make(map[string]map[int]*CompareState),
	}
}

// Get retrieves the comparison state for a file in a specific directory.
// Returns nil if the file is not in the index (i.e., unique to one directory).
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func (idx *ComparisonIndex) Get(dirIndex int, filename string) *CompareState {
	if idx == nil {
		return nil
	}
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	if dirs, ok := idx.cache[filename]; ok {
		if state, ok := dirs[dirIndex]; ok {
			return state
		}
	}
	return nil
}

// BuildComparisonIndex builds a comparison index from the given directories.
// Files appearing in multiple directories are analyzed for name/size/time comparison.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func BuildComparisonIndex(dirs []*Directory) *ComparisonIndex {
	if len(dirs) < 2 {
		return nil
	}

	// Collect all files by name across directories
	// map[filename][]fileEntry
	filesByName := make(map[string][]fileEntry)

	for dirIdx, dir := range dirs {
		if dir == nil {
			continue
		}
		for _, item := range dir.List() {
			fs, ok := item.(*FileStat)
			if !ok || fs == nil {
				continue
			}
			name := fs.Name()
			if name == ".." {
				continue
			}
			// Case-sensitive name comparison per requirement
			filesByName[name] = append(filesByName[name], fileEntry{
				dirIndex: dirIdx,
				size:     fs.Size(),
				modTime:  fs.ModTime(),
			})
		}
	}

	// Build comparison index for files appearing in multiple directories
	idx := NewComparisonIndex()
	for filename, entries := range filesByName {
		if len(entries) < 2 {
			// File only in one directory, no comparison needed
			continue
		}

		// Compute size and time comparison states
		sizeStates, timeStates := computeComparisonStates(entries)

		// Store in cache
		idx.cache[filename] = make(map[int]*CompareState)
		for i, entry := range entries {
			idx.cache[filename][entry.dirIndex] = &CompareState{
				NamePresent: true,
				SizeState:   sizeStates[i],
				TimeState:   timeStates[i],
			}
		}
	}

	return idx
}

// computeComparisonStates computes size and time comparison states for a set of file entries.
func computeComparisonStates(entries []fileEntry) ([]SizeCompare, []TimeCompare) {
	n := len(entries)
	sizeStates := make([]SizeCompare, n)
	timeStates := make([]TimeCompare, n)

	if n == 0 {
		return sizeStates, timeStates
	}

	// Find min/max size
	minSize, maxSize := entries[0].size, entries[0].size
	allSizesEqual := true
	for i := 1; i < n; i++ {
		if entries[i].size < minSize {
			minSize = entries[i].size
		}
		if entries[i].size > maxSize {
			maxSize = entries[i].size
		}
		if entries[i].size != entries[0].size {
			allSizesEqual = false
		}
	}

	// Find min/max time (truncated to second precision per requirement)
	minTime := entries[0].modTime.Truncate(time.Second)
	maxTime := minTime
	allTimesEqual := true
	for i := 1; i < n; i++ {
		t := entries[i].modTime.Truncate(time.Second)
		if t.Before(minTime) {
			minTime = t
		}
		if t.After(maxTime) {
			maxTime = t
		}
		if !t.Equal(entries[0].modTime.Truncate(time.Second)) {
			allTimesEqual = false
		}
	}

	// Assign states
	for i, entry := range entries {
		// Size state
		if allSizesEqual {
			sizeStates[i] = SizeEqual
		} else if entry.size == minSize {
			sizeStates[i] = SizeSmallest
		} else if entry.size == maxSize {
			sizeStates[i] = SizeLargest
		} else {
			sizeStates[i] = SizeMiddle
		}

		// Time state (truncate to second precision)
		t := entry.modTime.Truncate(time.Second)
		if allTimesEqual {
			timeStates[i] = TimeEqual
		} else if t.Equal(minTime) {
			timeStates[i] = TimeEarliest
		} else if t.Equal(maxTime) {
			timeStates[i] = TimeLatest
		} else {
			timeStates[i] = TimeMiddle
		}
	}

	return sizeStates, timeStates
}

// Clear removes all entries from the comparison index.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func (idx *ComparisonIndex) Clear() {
	if idx == nil {
		return
	}
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.cache = make(map[string]map[int]*CompareState)
}

// CalculateFileDigest computes the xxHash64 digest of a file using streaming.
// [IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func CalculateFileDigest(path string) (uint64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	h := xxhash.New()
	if _, err := io.Copy(h, f); err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}

// UpdateDigestStates calculates and updates digest states for a specific filename.
// Groups files by actual size and compares digests within each size group.
// Files with unique sizes (no other file shares the same size) get DigestNA.
// Returns the number of files processed.
// [IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func (idx *ComparisonIndex) UpdateDigestStates(filename string, dirs []*Directory) int {
	if idx == nil {
		return 0
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()

	dirStates, ok := idx.cache[filename]
	if !ok {
		return 0
	}

	// Collect file info including actual size
	type fileInfo struct {
		dirIndex int
		path     string
		size     int64
		state    *CompareState
	}

	var allFiles []fileInfo
	for dirIdx, state := range dirStates {
		if dirIdx >= len(dirs) || dirs[dirIdx] == nil {
			continue
		}
		dir := dirs[dirIdx]
		for _, item := range dir.List() {
			fs, ok := item.(*FileStat)
			if !ok || fs == nil {
				continue
			}
			if fs.Name() == filename {
				allFiles = append(allFiles, fileInfo{
					dirIndex: dirIdx,
					path:     fs.Path(),
					size:     fs.Size(),
					state:    state,
				})
				break
			}
		}
	}

	// Group files by size
	sizeGroups := make(map[int64][]fileInfo)
	for _, fi := range allFiles {
		sizeGroups[fi.size] = append(sizeGroups[fi.size], fi)
	}

	// Process each size group
	totalProcessed := 0
	for _, group := range sizeGroups {
		if len(group) < 2 {
			// Only one file with this size, mark as DigestNA
			for _, fi := range group {
				fi.state.DigestState = DigestNA
			}
			continue
		}

		// Calculate digests for all files in this size group
		digests := make(map[int]uint64)
		for _, fi := range group {
			digest, err := CalculateFileDigest(fi.path)
			if err != nil {
				fi.state.DigestState = DigestUnknown
				continue
			}
			digests[fi.dirIndex] = digest
		}

		// Check if all digests in this group are equal
		var firstDigest uint64
		allEqual := true
		first := true
		for _, digest := range digests {
			if first {
				firstDigest = digest
				first = false
			} else if digest != firstDigest {
				allEqual = false
				break
			}
		}

		// Update states based on comparison
		for _, fi := range group {
			if _, ok := digests[fi.dirIndex]; !ok {
				// Digest calculation failed
				continue
			}
			if allEqual {
				fi.state.DigestState = DigestEqual
			} else {
				fi.state.DigestState = DigestDifferent
			}
			totalProcessed++
		}
	}

	return totalProcessed
}
