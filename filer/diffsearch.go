// Package filer difference search for cross-directory comparison.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
package filer

import (
	"fmt"
	"sort"
)

// DiffSearchState holds the state for a difference search session.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
type DiffSearchState struct {
	InitialDirs    []string // Initial directory paths when search started
	Active         bool     // Whether a search is in progress
	LastDiffName   string   // Name of last found difference (for status display)
	LastDiffReason string   // Reason for last difference (for status display)
	CurrentPath    string   // Current directory being searched (for progress display)
	FilesChecked   int      // Count of files checked (for progress display)
	Searching      bool     // Whether actively searching (vs paused at a difference)
}

// NewDiffSearchState creates a new search state from the current directories.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func NewDiffSearchState(dirs []*Directory) *DiffSearchState {
	initialDirs := make([]string, len(dirs))
	for i, d := range dirs {
		initialDirs[i] = d.Path
	}
	return &DiffSearchState{
		InitialDirs: initialDirs,
		Active:      true,
	}
}

// IsActive returns whether the search is active.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (s *DiffSearchState) IsActive() bool {
	if s == nil {
		return false
	}
	return s.Active
}

// Clear deactivates the search state and resets status fields.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (s *DiffSearchState) Clear() {
	if s == nil {
		return
	}
	s.Active = false
	s.Searching = false
	s.LastDiffName = ""
	s.LastDiffReason = ""
	s.CurrentPath = ""
	s.FilesChecked = 0
}

// SetSearching marks the search as actively running.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (s *DiffSearchState) SetSearching(searching bool) {
	if s == nil {
		return
	}
	s.Searching = searching
}

// SetCurrentPath updates the current path being searched.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (s *DiffSearchState) SetCurrentPath(path string) {
	if s == nil {
		return
	}
	s.CurrentPath = path
}

// IncrementFilesChecked increments the files checked counter.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (s *DiffSearchState) IncrementFilesChecked() {
	if s == nil {
		return
	}
	s.FilesChecked++
}

// SetLastDiff records the last found difference.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (s *DiffSearchState) SetLastDiff(name, reason string) {
	if s == nil {
		return
	}
	s.LastDiffName = name
	s.LastDiffReason = reason
	s.Searching = false // Paused at a difference
}

// StatusText returns a formatted status string for display.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (s *DiffSearchState) StatusText() string {
	if s == nil || !s.Active {
		return ""
	}
	if s.Searching {
		if s.CurrentPath != "" {
			return fmt.Sprintf("[DIFF: searching %s (%d files)]", s.CurrentPath, s.FilesChecked)
		}
		return fmt.Sprintf("[DIFF: searching (%d files)]", s.FilesChecked)
	}
	if s.LastDiffName != "" {
		return fmt.Sprintf("[DIFF: %s - %s]", s.LastDiffName, s.LastDiffReason)
	}
	return "[DIFF SEARCH]"
}

// AtInitialDirs returns true if all directories are at their initial positions.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func (s *DiffSearchState) AtInitialDirs(dirs []*Directory) bool {
	if s == nil || len(s.InitialDirs) != len(dirs) {
		return false
	}
	for i, d := range dirs {
		if d.Path != s.InitialDirs[i] {
			return false
		}
	}
	return true
}

// DiffResult represents the result of a difference search.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
type DiffResult struct {
	Name   string // Name of the different entry
	Reason string // Why it's different
	Found  bool   // Whether a difference was found
	IsDir  bool   // Whether the entry is a directory
}

// CollectAllNames returns the union of all file/directory names across directories,
// sorted alphabetically (case-sensitive). Excludes ".." entries.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func CollectAllNames(dirs []*Directory) []string {
	nameSet := make(map[string]struct{})

	for _, dir := range dirs {
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
			nameSet[name] = struct{}{}
		}
	}

	names := make([]string, 0, len(nameSet))
	for name := range nameSet {
		names = append(names, name)
	}
	sort.Strings(names) // Case-sensitive alphabetic sort
	return names
}

// CollectFileNames returns only regular file names (not directories),
// sorted alphabetically (case-sensitive).
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func CollectFileNames(dirs []*Directory) []string {
	nameSet := make(map[string]struct{})

	for _, dir := range dirs {
		if dir == nil {
			continue
		}
		for _, item := range dir.List() {
			fs, ok := item.(*FileStat)
			if !ok || fs == nil {
				continue
			}
			name := fs.Name()
			if name == ".." || fs.IsDir() {
				continue
			}
			nameSet[name] = struct{}{}
		}
	}

	names := make([]string, 0, len(nameSet))
	for name := range nameSet {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// CollectSubdirNames returns only directory names,
// sorted alphabetically (case-sensitive).
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func CollectSubdirNames(dirs []*Directory) []string {
	nameSet := make(map[string]struct{})

	for _, dir := range dirs {
		if dir == nil {
			continue
		}
		for _, item := range dir.List() {
			fs, ok := item.(*FileStat)
			if !ok || fs == nil {
				continue
			}
			name := fs.Name()
			if name == ".." || !fs.IsDir() {
				continue
			}
			nameSet[name] = struct{}{}
		}
	}

	names := make([]string, 0, len(nameSet))
	for name := range nameSet {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// findEntryInDir finds a FileStat by name in a directory.
// Returns nil if not found.
func findEntryInDir(dir *Directory, name string) *FileStat {
	if dir == nil {
		return nil
	}
	for _, item := range dir.List() {
		fs, ok := item.(*FileStat)
		if !ok || fs == nil {
			continue
		}
		if fs.Name() == name {
			return fs
		}
	}
	return nil
}

// CheckDifference checks if a named entry is different across directories.
// Returns (isDifferent, reason).
// An entry is different if:
//   - It's missing from any directory, OR
//   - It has different sizes across directories (for files)
//
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func CheckDifference(name string, dirs []*Directory) (isDiff bool, reason string, isDir bool) {
	var entries []*FileStat
	var presentIn []int
	var missingIn []int

	for i, dir := range dirs {
		entry := findEntryInDir(dir, name)
		if entry != nil {
			entries = append(entries, entry)
			presentIn = append(presentIn, i+1) // 1-indexed for display
		} else {
			missingIn = append(missingIn, i+1)
		}
	}

	// Check if missing from any window
	if len(missingIn) > 0 {
		if len(missingIn) == 1 {
			reason = fmt.Sprintf("missing in window %d", missingIn[0])
		} else {
			reason = fmt.Sprintf("missing in windows %v", missingIn)
		}
		isDir = len(entries) > 0 && entries[0].IsDir()
		return true, reason, isDir
	}

	// All present, check if it's a directory
	if entries[0].IsDir() {
		// Directories are compared by presence only, not size
		return false, "", true
	}

	// Compare sizes for files
	firstSize := entries[0].Size()
	for i := 1; i < len(entries); i++ {
		if entries[i].Size() != firstSize {
			return true, "size mismatch", false
		}
	}

	return false, "", false
}

// FindNextDifference searches for the next different entry starting after startAfter.
// If startAfter is empty, starts from the beginning.
// Uses alphabetical comparison so startAfter can be a name not in the current list.
// Returns the result of the search.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func FindNextDifference(dirs []*Directory, startAfter string, filesOnly bool) DiffResult {
	var names []string
	if filesOnly {
		names = CollectFileNames(dirs)
	} else {
		names = CollectAllNames(dirs)
	}

	for _, name := range names {
		// Skip entries that come before or equal to startAfter alphabetically
		if startAfter != "" && name <= startAfter {
			continue
		}

		isDiff, reason, isDir := CheckDifference(name, dirs)
		if isDiff {
			return DiffResult{
				Name:   name,
				Reason: reason,
				Found:  true,
				IsDir:  isDir,
			}
		}
	}

	return DiffResult{Found: false}
}

// FindNextSubdir finds the next subdirectory after startAfter.
// Returns the name, whether it exists in all directories, and whether any was found.
// Uses alphabetical comparison so startAfter can be a filename not in the subdirs list.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func FindNextSubdir(dirs []*Directory, startAfter string) (name string, existsInAll bool, found bool) {
	subdirs := CollectSubdirNames(dirs)

	for _, subdir := range subdirs {
		// Skip entries that come before or equal to startAfter alphabetically
		if startAfter != "" && subdir <= startAfter {
			continue
		}

		// Check if this subdir exists in all directories
		allHave := true
		for _, dir := range dirs {
			entry := findEntryInDir(dir, subdir)
			if entry == nil || !entry.IsDir() {
				allHave = false
				break
			}
		}

		return subdir, allHave, true
	}

	return "", false, false
}

// FindNextSubdirDifference searches for the next subdirectory that differs across directories.
// This is separate from FindNextDifference to ensure subdirectories are checked independently
// of files, per the requirement "files first, then dirs".
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func FindNextSubdirDifference(dirs []*Directory, startAfter string) DiffResult {
	subdirs := CollectSubdirNames(dirs)

	for _, subdir := range subdirs {
		// Skip entries that come before or equal to startAfter alphabetically
		if startAfter != "" && subdir <= startAfter {
			continue
		}

		isDiff, reason, isDir := CheckDifference(subdir, dirs)
		if isDiff {
			return DiffResult{
				Name:   subdir,
				Reason: reason,
				Found:  true,
				IsDir:  isDir,
			}
		}
	}

	return DiffResult{Found: false}
}

// FirstSubdirInAll returns the first subdirectory that exists in all directories.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func FirstSubdirInAll(dirs []*Directory) (name string, found bool) {
	subdirs := CollectSubdirNames(dirs)

	for _, subdir := range subdirs {
		allHave := true
		for _, dir := range dirs {
			entry := findEntryInDir(dir, subdir)
			if entry == nil || !entry.IsDir() {
				allHave = false
				break
			}
		}
		if allHave {
			return subdir, true
		}
	}

	return "", false
}

// FindNextSubdirInAll returns the next subdirectory after startAfter that exists in all directories.
// If startAfter is empty, behaves like FirstSubdirInAll.
// This is used during diff search traversal to respect the current search position.
// Uses alphabetical comparison so startAfter can be a filename not in the subdirs list.
// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
func FindNextSubdirInAll(dirs []*Directory, startAfter string) (name string, found bool) {
	subdirs := CollectSubdirNames(dirs)

	for _, subdir := range subdirs {
		// Skip entries that come before or equal to startAfter alphabetically
		// This handles the case where startAfter is a filename not in the subdirs list
		if startAfter != "" && subdir <= startAfter {
			continue
		}

		// Check if this subdir exists in all directories
		allHave := true
		for _, dir := range dirs {
			entry := findEntryInDir(dir, subdir)
			if entry == nil || !entry.IsDir() {
				allHave = false
				break
			}
		}
		if allHave {
			return subdir, true
		}
	}

	return "", false
}
