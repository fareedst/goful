package filer

import "testing"

func resetExcludesForTest() {
	ConfigureExcludedNames(nil, false)
}

func TestConfigureExcludedNamesNormalizes_REQ_FILER_EXCLUDE_NAMES(t *testing.T) {
	// [REQ:FILER_EXCLUDE_NAMES] [ARCH:FILER_EXCLUDE_FILTER] [IMPL:FILER_EXCLUDE_RULES]
	t.Cleanup(resetExcludesForTest)
	count := ConfigureExcludedNames([]string{"  Noise.tmp  ", "README.MD"}, true)
	if count != 2 {
		t.Fatalf("ConfigureExcludedNames count=%d, want 2", count)
	}
	if !ExcludedNamesEnabled() {
		t.Fatalf("ExcludedNamesEnabled should be true when activate=true and names exist")
	}
	if !shouldExcludeName("noise.tmp") || !shouldExcludeName("readme.md") {
		t.Fatalf("shouldExcludeName failed to match normalized entries")
	}
	if shouldExcludeName("other.txt") {
		t.Fatalf("unexpected exclude match for other.txt")
	}
}

func TestToggleExcludedNames_REQ_FILER_EXCLUDE_NAMES(t *testing.T) {
	// [REQ:FILER_EXCLUDE_NAMES] [ARCH:FILER_EXCLUDE_FILTER] [IMPL:FILER_EXCLUDE_RULES]
	t.Cleanup(resetExcludesForTest)

	if enabled, hasRules, _ := ToggleExcludedNames(); hasRules || enabled {
		t.Fatalf("toggle with no rules should report hasRules=false, enabled=false")
	}

	ConfigureExcludedNames([]string{"temp.tmp"}, true)
	enabled, hasRules, _ := ToggleExcludedNames()
	if !hasRules || enabled {
		t.Fatalf("first toggle should disable the filter (enabled=%v, hasRules=%v)", enabled, hasRules)
	}
	enabled, hasRules, _ = ToggleExcludedNames()
	if !hasRules || !enabled {
		t.Fatalf("second toggle should re-enable the filter (enabled=%v, hasRules=%v)", enabled, hasRules)
	}
}
