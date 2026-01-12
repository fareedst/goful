// Package look comparison colors for cross-directory file comparison.
// [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
package look

import (
	"sync"

	"github.com/anmitsu/goful/filer/comparecolors"
	"github.com/gdamore/tcell/v2"
)

var (
	comparisonMu sync.RWMutex
	// comparisonEnabled tracks whether comparison coloring is active.
	// [REQ:FILE_COMPARISON_COLORS] Enabled by default for immediate visual feedback.
	comparisonEnabled bool = true
	// comparisonStyles holds the parsed comparison color styles.
	comparisonStyles *comparecolors.ParsedConfig
)

func init() {
	// Initialize with default styles
	comparisonStyles = comparecolors.DefaultConfig().Parse()
}

// ConfigureComparisonColors sets the comparison color styles from a parsed config.
// [IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func ConfigureComparisonColors(cfg *comparecolors.ParsedConfig) {
	comparisonMu.Lock()
	defer comparisonMu.Unlock()
	if cfg != nil {
		comparisonStyles = cfg
	}
}

// ComparisonEnabled reports whether comparison coloring is currently active.
// [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func ComparisonEnabled() bool {
	comparisonMu.RLock()
	defer comparisonMu.RUnlock()
	return comparisonEnabled
}

// SetComparisonEnabled sets whether comparison coloring is active.
// [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func SetComparisonEnabled(enabled bool) {
	comparisonMu.Lock()
	defer comparisonMu.Unlock()
	comparisonEnabled = enabled
}

// ToggleComparisonEnabled toggles comparison coloring on/off and returns the new state.
// [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func ToggleComparisonEnabled() bool {
	comparisonMu.Lock()
	defer comparisonMu.Unlock()
	comparisonEnabled = !comparisonEnabled
	return comparisonEnabled
}

// CompareNamePresent returns the style for file names present in multiple directories.
// [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func CompareNamePresent() tcell.Style {
	comparisonMu.RLock()
	defer comparisonMu.RUnlock()
	return comparisonStyles.NamePresent
}

// CompareSizeEqual returns the style for equal file sizes across directories.
// [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func CompareSizeEqual() tcell.Style {
	comparisonMu.RLock()
	defer comparisonMu.RUnlock()
	return comparisonStyles.SizeEqual
}

// CompareSizeSmallest returns the style for the smallest file size among same-named files.
// [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func CompareSizeSmallest() tcell.Style {
	comparisonMu.RLock()
	defer comparisonMu.RUnlock()
	return comparisonStyles.SizeSmallest
}

// CompareSizeLargest returns the style for the largest file size among same-named files.
// [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func CompareSizeLargest() tcell.Style {
	comparisonMu.RLock()
	defer comparisonMu.RUnlock()
	return comparisonStyles.SizeLargest
}

// CompareTimeEqual returns the style for equal modification times across directories.
// [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func CompareTimeEqual() tcell.Style {
	comparisonMu.RLock()
	defer comparisonMu.RUnlock()
	return comparisonStyles.TimeEqual
}

// CompareTimeEarliest returns the style for the earliest modification time among same-named files.
// [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func CompareTimeEarliest() tcell.Style {
	comparisonMu.RLock()
	defer comparisonMu.RUnlock()
	return comparisonStyles.TimeEarliest
}

// CompareTimeLatest returns the style for the latest modification time among same-named files.
// [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func CompareTimeLatest() tcell.Style {
	comparisonMu.RLock()
	defer comparisonMu.RUnlock()
	return comparisonStyles.TimeLatest
}
