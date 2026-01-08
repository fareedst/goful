package filer

import (
	"strings"
	"sync"
)

var (
	excludedNamesMu sync.RWMutex
	excludedNames   map[string]struct{}
	excludeEnabled  bool
)

// ConfigureExcludedNames replaces the current exclude set and activates it when requested.
// [IMPL:FILER_EXCLUDE_RULES] [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]
func ConfigureExcludedNames(names []string, activate bool) int {
	excludedNamesMu.Lock()
	defer excludedNamesMu.Unlock()

	if len(names) == 0 {
		excludedNames = nil
		excludeEnabled = false
		return 0
	}

	set := make(map[string]struct{}, len(names))
	for _, name := range names {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}
		set[strings.ToLower(trimmed)] = struct{}{}
	}

	excludedNames = set
	excludeEnabled = activate && len(set) > 0
	return len(set)
}

// ToggleExcludedNames flips the active state when rules exist.
// Returns (enabled, hasRules, ruleCount).
// [IMPL:FILER_EXCLUDE_RULES] [ARCH:FILER_EXCLUDE_FILTER] [REQ:FILER_EXCLUDE_NAMES]
func ToggleExcludedNames() (bool, bool, int) {
	excludedNamesMu.Lock()
	defer excludedNamesMu.Unlock()

	count := len(excludedNames)
	if count == 0 {
		excludeEnabled = false
		return false, false, 0
	}
	excludeEnabled = !excludeEnabled
	return excludeEnabled, true, count
}

// ExcludedNamesEnabled reports whether the filter is on.
func ExcludedNamesEnabled() bool {
	excludedNamesMu.RLock()
	defer excludedNamesMu.RUnlock()
	return excludeEnabled && len(excludedNames) > 0
}

func shouldExcludeName(name string) bool {
	excludedNamesMu.RLock()
	defer excludedNamesMu.RUnlock()

	if !excludeEnabled || len(excludedNames) == 0 {
		return false
	}
	_, ok := excludedNames[strings.ToLower(name)]
	return ok
}
