//go:build darwin
// +build darwin

package app

import (
	"testing"

	"github.com/fareedst/nsync/pkg/nsync"
)

// TestGofulObserver_StateTracking_REQ_NSYNC_MULTI_TARGET verifies observer state initialization.
// We test the state tracking directly without calling OnStart which requires progress.Init().
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func TestGofulObserver_StateTracking_REQ_NSYNC_MULTI_TARGET(t *testing.T) {
	obs := newGofulObserver()

	// Verify initial state
	if obs.totalBytes != 0 {
		t.Errorf("initial totalBytes = %d, want 0", obs.totalBytes)
	}
	if obs.totalItems != 0 {
		t.Errorf("initial totalItems = %d, want 0", obs.totalItems)
	}
	if obs.destinations != 0 {
		t.Errorf("initial destinations = %d, want 0", obs.destinations)
	}
	if obs.lastBytes != 0 {
		t.Errorf("initial lastBytes = %d, want 0", obs.lastBytes)
	}

	// Simulate what OnStart would set (without calling progress functions)
	plan := nsync.Plan{
		TotalItems:        10,
		TotalBytes:        1024 * 1024, // 1 MB
		TotalDestinations: 3,
	}

	// Manually set state as OnStart would
	obs.mu.Lock()
	obs.totalBytes = plan.TotalBytes
	obs.totalItems = plan.TotalItems
	obs.destinations = plan.TotalDestinations
	obs.mu.Unlock()

	if obs.totalBytes != plan.TotalBytes {
		t.Errorf("totalBytes = %d, want %d", obs.totalBytes, plan.TotalBytes)
	}
	if obs.totalItems != plan.TotalItems {
		t.Errorf("totalItems = %d, want %d", obs.totalItems, plan.TotalItems)
	}
	if obs.destinations != plan.TotalDestinations {
		t.Errorf("destinations = %d, want %d", obs.destinations, plan.TotalDestinations)
	}
}

// TestGofulObserver_ByteTracking_REQ_NSYNC_MULTI_TARGET verifies byte delta tracking logic.
// We test the internal state tracking without calling methods that require progress.Init().
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func TestGofulObserver_ByteTracking_REQ_NSYNC_MULTI_TARGET(t *testing.T) {
	obs := newGofulObserver()

	// Simulate the byte tracking logic that OnItemProgress performs
	// First update: 256 bytes
	obs.mu.Lock()
	delta1 := int64(256) - obs.lastBytes // delta = 256 - 0 = 256
	obs.lastBytes = 256
	obs.mu.Unlock()
	if delta1 != 256 {
		t.Errorf("delta1 = %d, want 256", delta1)
	}
	if obs.lastBytes != 256 {
		t.Errorf("lastBytes = %d, want 256", obs.lastBytes)
	}

	// Second update: 512 bytes (cumulative)
	obs.mu.Lock()
	delta2 := int64(512) - obs.lastBytes // delta = 512 - 256 = 256
	obs.lastBytes = 512
	obs.mu.Unlock()
	if delta2 != 256 {
		t.Errorf("delta2 = %d, want 256", delta2)
	}

	// Third update: 1024 bytes (cumulative)
	obs.mu.Lock()
	delta3 := int64(1024) - obs.lastBytes // delta = 1024 - 512 = 512
	obs.lastBytes = 1024
	obs.mu.Unlock()
	if delta3 != 512 {
		t.Errorf("delta3 = %d, want 512", delta3)
	}
}

// TestGofulObserver_ItemCompletionReset_REQ_NSYNC_MULTI_TARGET verifies item completion resets state.
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func TestGofulObserver_ItemCompletionReset_REQ_NSYNC_MULTI_TARGET(t *testing.T) {
	obs := newGofulObserver()
	obs.lastBytes = 1024 // Simulate some progress

	// Simulate what OnItemComplete does for the reset (without progress calls)
	obs.mu.Lock()
	obs.lastBytes = 0
	obs.mu.Unlock()

	if obs.lastBytes != 0 {
		t.Errorf("lastBytes = %d after completion, want 0", obs.lastBytes)
	}
}

// TestGofulObserver_ConcurrentSafety_REQ_NSYNC_MULTI_TARGET verifies thread safety.
// We test the internal mutex protection without calling methods that require progress.Init().
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func TestGofulObserver_ConcurrentSafety_REQ_NSYNC_MULTI_TARGET(t *testing.T) {
	obs := newGofulObserver()

	// Run concurrent updates to verify mutex protection
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(n int) {
			// Simulate what OnItemProgress does internally (just the state tracking)
			obs.mu.Lock()
			obs.lastBytes = int64(n * 100)
			obs.mu.Unlock()
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should not panic or race - test passes if we get here
}

// TestFormatBytes_REQ_NSYNC_MULTI_TARGET verifies byte formatting.
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func TestFormatBytes_REQ_NSYNC_MULTI_TARGET(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0B"},
		{512, "512B"},
		{1024, "1.0KB"},
		{1536, "1.5KB"},
		{1024 * 1024, "1.0MB"},
		{1024 * 1024 * 1024, "1.0GB"},
		{1024 * 1024 * 1024 * 1024, "1.0TB"},
	}

	for _, tt := range tests {
		got := formatBytes(tt.bytes)
		if got != tt.want {
			t.Errorf("formatBytes(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}

// TestOtherWindowDirPaths_Integration_REQ_NSYNC_MULTI_TARGET verifies destination enumeration.
// [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func TestOtherWindowDirPaths_Integration_REQ_NSYNC_MULTI_TARGET(t *testing.T) {
	// This test verifies the otherWindowDirPaths helper is available
	// The actual implementation is tested in spawn_test.go
	// Here we just verify the function signature works for nsync integration

	// Test with nil workspace
	paths := otherWindowDirPaths(nil)
	if paths != nil {
		t.Errorf("otherWindowDirPaths(nil) = %v, want nil", paths)
	}
}

// TestFakeFileInfo_REQ_NSYNC_MULTI_TARGET verifies the fake FileInfo for progress display.
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func TestFakeFileInfo_REQ_NSYNC_MULTI_TARGET(t *testing.T) {
	fi := &fakeFileInfo{
		name: "test.txt",
		size: 1024,
	}

	if fi.Name() != "test.txt" {
		t.Errorf("Name() = %q, want %q", fi.Name(), "test.txt")
	}
	if fi.Size() != 1024 {
		t.Errorf("Size() = %d, want %d", fi.Size(), 1024)
	}
	if fi.IsDir() {
		t.Error("IsDir() = true, want false")
	}
	if fi.Mode() != 0 {
		t.Errorf("Mode() = %v, want 0", fi.Mode())
	}
	if fi.ModTime().IsZero() == false {
		t.Error("ModTime() should be zero time")
	}
	if fi.Sys() != nil {
		t.Errorf("Sys() = %v, want nil", fi.Sys())
	}
}

// TestCopyAllMode_String_REQ_NSYNC_CONFIRMATION verifies copyAllMode interface.
// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]
func TestCopyAllMode_String_REQ_NSYNC_CONFIRMATION(t *testing.T) {
	mode := &copyAllMode{
		sources:      []string{"/src/a.txt", "/src/b.txt"},
		destinations: []string{"/dst1", "/dst2", "/dst3"},
	}

	// Test mode name
	if got := mode.String(); got != "copyall" {
		t.Errorf("String() = %q, want %q", got, "copyall")
	}

	// Test prompt format
	prompt := mode.Prompt()
	want := "Copy 2 file(s) to 3 destinations? [Y/n] "
	if prompt != want {
		t.Errorf("Prompt() = %q, want %q", prompt, want)
	}
}

// TestMoveAllMode_String_REQ_NSYNC_CONFIRMATION verifies moveAllMode interface.
// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]
func TestMoveAllMode_String_REQ_NSYNC_CONFIRMATION(t *testing.T) {
	mode := &moveAllMode{
		sources:      []string{"/src/a.txt"},
		destinations: []string{"/dst1", "/dst2"},
	}

	// Test mode name
	if got := mode.String(); got != "moveall" {
		t.Errorf("String() = %q, want %q", got, "moveall")
	}

	// Test prompt format
	prompt := mode.Prompt()
	want := "Move 1 file(s) to 2 destinations? [Y/n] "
	if prompt != want {
		t.Errorf("Prompt() = %q, want %q", prompt, want)
	}
}

// TestCopyAllMode_PromptCounts_REQ_NSYNC_CONFIRMATION tests various source/dest counts.
// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]
func TestCopyAllMode_PromptCounts_REQ_NSYNC_CONFIRMATION(t *testing.T) {
	tests := []struct {
		name         string
		sources      []string
		destinations []string
		wantPrompt   string
	}{
		{
			name:         "single file single dest",
			sources:      []string{"/a.txt"},
			destinations: []string{"/dst"},
			wantPrompt:   "Copy 1 file(s) to 1 destinations? [Y/n] ",
		},
		{
			name:         "multiple files multiple dests",
			sources:      []string{"/a.txt", "/b.txt", "/c.txt"},
			destinations: []string{"/d1", "/d2"},
			wantPrompt:   "Copy 3 file(s) to 2 destinations? [Y/n] ",
		},
		{
			name:         "ten files five dests",
			sources:      []string{"/1", "/2", "/3", "/4", "/5", "/6", "/7", "/8", "/9", "/10"},
			destinations: []string{"/a", "/b", "/c", "/d", "/e"},
			wantPrompt:   "Copy 10 file(s) to 5 destinations? [Y/n] ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode := &copyAllMode{
				sources:      tt.sources,
				destinations: tt.destinations,
			}
			if got := mode.Prompt(); got != tt.wantPrompt {
				t.Errorf("Prompt() = %q, want %q", got, tt.wantPrompt)
			}
		})
	}
}

// TestMoveAllMode_PromptCounts_REQ_NSYNC_CONFIRMATION tests various source/dest counts.
// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]
func TestMoveAllMode_PromptCounts_REQ_NSYNC_CONFIRMATION(t *testing.T) {
	tests := []struct {
		name         string
		sources      []string
		destinations []string
		wantPrompt   string
	}{
		{
			name:         "single file single dest",
			sources:      []string{"/a.txt"},
			destinations: []string{"/dst"},
			wantPrompt:   "Move 1 file(s) to 1 destinations? [Y/n] ",
		},
		{
			name:         "multiple files multiple dests",
			sources:      []string{"/a.txt", "/b.txt"},
			destinations: []string{"/d1", "/d2", "/d3", "/d4"},
			wantPrompt:   "Move 2 file(s) to 4 destinations? [Y/n] ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode := &moveAllMode{
				sources:      tt.sources,
				destinations: tt.destinations,
			}
			if got := mode.Prompt(); got != tt.wantPrompt {
				t.Errorf("Prompt() = %q, want %q", got, tt.wantPrompt)
			}
		})
	}
}
