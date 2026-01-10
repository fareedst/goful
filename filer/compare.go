// Package filer comparison index for cross-directory file comparison.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
package filer

import (
	"sync"
	"time"
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

// CompareState holds the comparison state for a single file in a directory.
// [IMPL:FILE_COMPARISON_INDEX] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
type CompareState struct {
	NamePresent bool        // File name appears in multiple directories
	SizeState   SizeCompare // Size comparison state
	TimeState   TimeCompare // Time comparison state
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
