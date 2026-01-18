package filer

import "testing"

// resetToolbarButtonsForTest clears the toolbar button bounds map.
func resetToolbarButtonsForTest() {
	toolbarButtons = make(map[string]toolbarBounds)
}

// TestToolbarButtonAt_REQ_TOOLBAR_PARENT_BUTTON tests hit-testing for toolbar buttons.
// [REQ:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_PARENT_BUTTON]
func TestToolbarButtonAt_REQ_TOOLBAR_PARENT_BUTTON(t *testing.T) {
	t.Cleanup(resetToolbarButtonsForTest)

	// Setup: Register a parent button at x=0 to x=2, y=0 (representing "[^]")
	toolbarButtons["parent"] = toolbarBounds{x1: 0, y: 0, x2: 2}

	tests := []struct {
		name     string
		x, y     int
		expected string
	}{
		{"click at x=0, y=0 (left edge)", 0, 0, "parent"},
		{"click at x=1, y=0 (middle)", 1, 0, "parent"},
		{"click at x=2, y=0 (right edge)", 2, 0, "parent"},
		{"click at x=3, y=0 (just outside)", 3, 0, ""},
		{"click at x=0, y=1 (wrong row)", 0, 1, ""},
		{"click at x=-1, y=0 (before button)", -1, 0, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ToolbarButtonAt(tc.x, tc.y)
			if result != tc.expected {
				t.Errorf("ToolbarButtonAt(%d, %d) = %q, want %q", tc.x, tc.y, result, tc.expected)
			}
		})
	}
}

// TestToolbarButtonAtMultipleButtons_REQ_TOOLBAR_PARENT_BUTTON tests hit-testing with multiple buttons.
// [REQ:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_PARENT_BUTTON]
func TestToolbarButtonAtMultipleButtons_REQ_TOOLBAR_PARENT_BUTTON(t *testing.T) {
	t.Cleanup(resetToolbarButtonsForTest)

	// Setup: Register multiple buttons on the same row
	toolbarButtons["parent"] = toolbarBounds{x1: 0, y: 0, x2: 2} // "[^]"
	toolbarButtons["reload"] = toolbarBounds{x1: 4, y: 0, x2: 6} // "[R]" hypothetical

	tests := []struct {
		name     string
		x, y     int
		expected string
	}{
		{"click on parent button", 1, 0, "parent"},
		{"click on reload button", 5, 0, "reload"},
		{"click between buttons", 3, 0, ""},
		{"click after all buttons", 7, 0, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ToolbarButtonAt(tc.x, tc.y)
			if result != tc.expected {
				t.Errorf("ToolbarButtonAt(%d, %d) = %q, want %q", tc.x, tc.y, result, tc.expected)
			}
		})
	}
}

// TestInvokeToolbarButton_REQ_TOOLBAR_PARENT_BUTTON tests that button invocation triggers callbacks.
// [REQ:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_PARENT_BUTTON]
func TestInvokeToolbarButton_REQ_TOOLBAR_PARENT_BUTTON(t *testing.T) {
	// Setup: Track callback invocation
	callbackCalled := false
	originalFn := toolbarParentNavFn
	t.Cleanup(func() {
		toolbarParentNavFn = originalFn
	})

	SetToolbarParentNavFn(func() {
		callbackCalled = true
	})

	// Test invoking the parent button
	handled := InvokeToolbarButton("parent")
	if !handled {
		t.Errorf("InvokeToolbarButton(\"parent\") = false, want true")
	}
	if !callbackCalled {
		t.Errorf("Parent button callback was not invoked")
	}

	// Test invoking unknown button
	handled = InvokeToolbarButton("unknown")
	if handled {
		t.Errorf("InvokeToolbarButton(\"unknown\") = true, want false")
	}
}

// TestInvokeToolbarButtonWithNilCallback_REQ_TOOLBAR_PARENT_BUTTON tests behavior when no callback is set.
// [REQ:TOOLBAR_PARENT_BUTTON] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_PARENT_BUTTON]
func TestInvokeToolbarButtonWithNilCallback_REQ_TOOLBAR_PARENT_BUTTON(t *testing.T) {
	// Setup: Clear the callback
	originalFn := toolbarParentNavFn
	t.Cleanup(func() {
		toolbarParentNavFn = originalFn
	})
	toolbarParentNavFn = nil

	// Should return false when no callback is set
	handled := InvokeToolbarButton("parent")
	if handled {
		t.Errorf("InvokeToolbarButton(\"parent\") with nil callback = true, want false")
	}
}

// TestToolbarLinkedButtonHit_REQ_TOOLBAR_LINKED_TOGGLE tests hit-testing for the linked button.
// [REQ:TOOLBAR_LINKED_TOGGLE] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_LINKED_TOGGLE]
func TestToolbarLinkedButtonHit_REQ_TOOLBAR_LINKED_TOGGLE(t *testing.T) {
	t.Cleanup(resetToolbarButtonsForTest)

	// Setup: Register parent and linked buttons as they would appear in the header
	// Layout: [^] [L] ...
	toolbarButtons["parent"] = toolbarBounds{x1: 0, y: 0, x2: 2} // "[^]"
	toolbarButtons["linked"] = toolbarBounds{x1: 4, y: 0, x2: 6} // "[L]"

	tests := []struct {
		name     string
		x, y     int
		expected string
	}{
		{"click on parent button", 1, 0, "parent"},
		{"click on linked button left edge", 4, 0, "linked"},
		{"click on linked button middle", 5, 0, "linked"},
		{"click on linked button right edge", 6, 0, "linked"},
		{"click between buttons (space)", 3, 0, ""},
		{"click after linked button", 7, 0, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ToolbarButtonAt(tc.x, tc.y)
			if result != tc.expected {
				t.Errorf("ToolbarButtonAt(%d, %d) = %q, want %q", tc.x, tc.y, result, tc.expected)
			}
		})
	}
}

// TestInvokeToolbarLinkedButton_REQ_TOOLBAR_LINKED_TOGGLE tests that linked button invocation triggers callback.
// [REQ:TOOLBAR_LINKED_TOGGLE] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_LINKED_TOGGLE]
func TestInvokeToolbarLinkedButton_REQ_TOOLBAR_LINKED_TOGGLE(t *testing.T) {
	// Setup: Track callback invocation
	callbackCalled := false
	originalFn := toolbarLinkedToggleFn
	t.Cleanup(func() {
		toolbarLinkedToggleFn = originalFn
	})

	SetToolbarLinkedToggleFn(func() {
		callbackCalled = true
	})

	// Test invoking the linked button
	handled := InvokeToolbarButton("linked")
	if !handled {
		t.Errorf("InvokeToolbarButton(\"linked\") = false, want true")
	}
	if !callbackCalled {
		t.Errorf("Linked button callback was not invoked")
	}
}

// TestInvokeToolbarLinkedButtonWithNilCallback_REQ_TOOLBAR_LINKED_TOGGLE tests behavior when no callback is set.
// [REQ:TOOLBAR_LINKED_TOGGLE] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_LINKED_TOGGLE]
func TestInvokeToolbarLinkedButtonWithNilCallback_REQ_TOOLBAR_LINKED_TOGGLE(t *testing.T) {
	// Setup: Clear the callback
	originalFn := toolbarLinkedToggleFn
	t.Cleanup(func() {
		toolbarLinkedToggleFn = originalFn
	})
	toolbarLinkedToggleFn = nil

	// Should return false when no callback is set
	handled := InvokeToolbarButton("linked")
	if handled {
		t.Errorf("InvokeToolbarButton(\"linked\") with nil callback = true, want false")
	}
}

// TestToolbarCompareButtonHit_REQ_TOOLBAR_COMPARE_BUTTON tests hit-testing for the compare button.
// [REQ:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_COMPARE_BUTTON]
func TestToolbarCompareButtonHit_REQ_TOOLBAR_COMPARE_BUTTON(t *testing.T) {
	t.Cleanup(resetToolbarButtonsForTest)

	// Setup: Register parent, linked, and compare buttons as they would appear in the header
	// Layout: [^] [L] [=] ...
	toolbarButtons["parent"] = toolbarBounds{x1: 0, y: 0, x2: 2}   // "[^]"
	toolbarButtons["linked"] = toolbarBounds{x1: 4, y: 0, x2: 6}   // "[L]"
	toolbarButtons["compare"] = toolbarBounds{x1: 8, y: 0, x2: 10} // "[=]"

	tests := []struct {
		name     string
		x, y     int
		expected string
	}{
		{"click on parent button", 1, 0, "parent"},
		{"click on linked button", 5, 0, "linked"},
		{"click on compare button left edge", 8, 0, "compare"},
		{"click on compare button middle", 9, 0, "compare"},
		{"click on compare button right edge", 10, 0, "compare"},
		{"click between linked and compare", 7, 0, ""},
		{"click after compare button", 11, 0, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ToolbarButtonAt(tc.x, tc.y)
			if result != tc.expected {
				t.Errorf("ToolbarButtonAt(%d, %d) = %q, want %q", tc.x, tc.y, result, tc.expected)
			}
		})
	}
}

// TestInvokeToolbarCompareButton_REQ_TOOLBAR_COMPARE_BUTTON tests that compare button invocation triggers callback.
// [REQ:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_COMPARE_BUTTON]
func TestInvokeToolbarCompareButton_REQ_TOOLBAR_COMPARE_BUTTON(t *testing.T) {
	// Setup: Track callback invocation
	callbackCalled := false
	originalFn := toolbarCompareDigestFn
	t.Cleanup(func() {
		toolbarCompareDigestFn = originalFn
	})

	SetToolbarCompareDigestFn(func() {
		callbackCalled = true
	})

	// Test invoking the compare button
	handled := InvokeToolbarButton("compare")
	if !handled {
		t.Errorf("InvokeToolbarButton(\"compare\") = false, want true")
	}
	if !callbackCalled {
		t.Errorf("Compare button callback was not invoked")
	}
}

// TestInvokeToolbarCompareButtonWithNilCallback_REQ_TOOLBAR_COMPARE_BUTTON tests behavior when no callback is set.
// [REQ:TOOLBAR_COMPARE_BUTTON] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_COMPARE_BUTTON]
func TestInvokeToolbarCompareButtonWithNilCallback_REQ_TOOLBAR_COMPARE_BUTTON(t *testing.T) {
	// Setup: Clear the callback
	originalFn := toolbarCompareDigestFn
	t.Cleanup(func() {
		toolbarCompareDigestFn = originalFn
	})
	toolbarCompareDigestFn = nil

	// Should return false when no callback is set
	handled := InvokeToolbarButton("compare")
	if handled {
		t.Errorf("InvokeToolbarButton(\"compare\") with nil callback = true, want false")
	}
}

// TestToolbarSyncButtonsHit_REQ_TOOLBAR_SYNC_BUTTONS tests hit-testing for all sync buttons.
// [REQ:TOOLBAR_SYNC_BUTTONS] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_SYNC_COPY] [IMPL:TOOLBAR_SYNC_DELETE] [IMPL:TOOLBAR_SYNC_RENAME] [IMPL:TOOLBAR_IGNORE_FAILURES]
func TestToolbarSyncButtonsHit_REQ_TOOLBAR_SYNC_BUTTONS(t *testing.T) {
	t.Cleanup(resetToolbarButtonsForTest)

	// Setup: Register all toolbar buttons as they would appear in the header
	// Layout: [^] [L] [=] [C] [D] [R] [!] ...
	toolbarButtons["parent"] = toolbarBounds{x1: 0, y: 0, x2: 2}           // "[^]"
	toolbarButtons["linked"] = toolbarBounds{x1: 4, y: 0, x2: 6}           // "[L]"
	toolbarButtons["compare"] = toolbarBounds{x1: 8, y: 0, x2: 10}         // "[=]"
	toolbarButtons["synccopy"] = toolbarBounds{x1: 12, y: 0, x2: 14}       // "[C]"
	toolbarButtons["syncdelete"] = toolbarBounds{x1: 16, y: 0, x2: 18}     // "[D]"
	toolbarButtons["syncrename"] = toolbarBounds{x1: 20, y: 0, x2: 22}     // "[R]"
	toolbarButtons["ignorefailures"] = toolbarBounds{x1: 24, y: 0, x2: 26} // "[!]"

	tests := []struct {
		name     string
		x, y     int
		expected string
	}{
		{"click on synccopy button", 13, 0, "synccopy"},
		{"click on syncdelete button", 17, 0, "syncdelete"},
		{"click on syncrename button", 21, 0, "syncrename"},
		{"click on ignorefailures button", 25, 0, "ignorefailures"},
		{"click between compare and synccopy", 11, 0, ""},
		{"click after ignorefailures button", 27, 0, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ToolbarButtonAt(tc.x, tc.y)
			if result != tc.expected {
				t.Errorf("ToolbarButtonAt(%d, %d) = %q, want %q", tc.x, tc.y, result, tc.expected)
			}
		})
	}
}

// TestInvokeToolbarSyncCopyButton_REQ_TOOLBAR_SYNC_BUTTONS tests that sync copy button invocation triggers callback.
// [REQ:TOOLBAR_SYNC_BUTTONS] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_SYNC_COPY]
func TestInvokeToolbarSyncCopyButton_REQ_TOOLBAR_SYNC_BUTTONS(t *testing.T) {
	callbackCalled := false
	originalFn := toolbarSyncCopyFn
	t.Cleanup(func() {
		toolbarSyncCopyFn = originalFn
	})

	SetToolbarSyncCopyFn(func() {
		callbackCalled = true
	})

	handled := InvokeToolbarButton("synccopy")
	if !handled {
		t.Errorf("InvokeToolbarButton(\"synccopy\") = false, want true")
	}
	if !callbackCalled {
		t.Errorf("Sync copy button callback was not invoked")
	}
}

// TestInvokeToolbarSyncDeleteButton_REQ_TOOLBAR_SYNC_BUTTONS tests that sync delete button invocation triggers callback.
// [REQ:TOOLBAR_SYNC_BUTTONS] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_SYNC_DELETE]
func TestInvokeToolbarSyncDeleteButton_REQ_TOOLBAR_SYNC_BUTTONS(t *testing.T) {
	callbackCalled := false
	originalFn := toolbarSyncDeleteFn
	t.Cleanup(func() {
		toolbarSyncDeleteFn = originalFn
	})

	SetToolbarSyncDeleteFn(func() {
		callbackCalled = true
	})

	handled := InvokeToolbarButton("syncdelete")
	if !handled {
		t.Errorf("InvokeToolbarButton(\"syncdelete\") = false, want true")
	}
	if !callbackCalled {
		t.Errorf("Sync delete button callback was not invoked")
	}
}

// TestInvokeToolbarSyncRenameButton_REQ_TOOLBAR_SYNC_BUTTONS tests that sync rename button invocation triggers callback.
// [REQ:TOOLBAR_SYNC_BUTTONS] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_SYNC_RENAME]
func TestInvokeToolbarSyncRenameButton_REQ_TOOLBAR_SYNC_BUTTONS(t *testing.T) {
	callbackCalled := false
	originalFn := toolbarSyncRenameFn
	t.Cleanup(func() {
		toolbarSyncRenameFn = originalFn
	})

	SetToolbarSyncRenameFn(func() {
		callbackCalled = true
	})

	handled := InvokeToolbarButton("syncrename")
	if !handled {
		t.Errorf("InvokeToolbarButton(\"syncrename\") = false, want true")
	}
	if !callbackCalled {
		t.Errorf("Sync rename button callback was not invoked")
	}
}

// TestInvokeToolbarIgnoreFailuresButton_REQ_TOOLBAR_SYNC_BUTTONS tests that ignore-failures button invocation triggers callback.
// [REQ:TOOLBAR_SYNC_BUTTONS] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_IGNORE_FAILURES]
func TestInvokeToolbarIgnoreFailuresButton_REQ_TOOLBAR_SYNC_BUTTONS(t *testing.T) {
	callbackCalled := false
	originalFn := toolbarIgnoreFailuresFn
	t.Cleanup(func() {
		toolbarIgnoreFailuresFn = originalFn
	})

	SetToolbarIgnoreFailuresFn(func() {
		callbackCalled = true
	})

	handled := InvokeToolbarButton("ignorefailures")
	if !handled {
		t.Errorf("InvokeToolbarButton(\"ignorefailures\") = false, want true")
	}
	if !callbackCalled {
		t.Errorf("Ignore failures button callback was not invoked")
	}
}

// TestInvokeToolbarSyncButtonsWithNilCallback_REQ_TOOLBAR_SYNC_BUTTONS tests behavior when no callbacks are set.
// [REQ:TOOLBAR_SYNC_BUTTONS] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_SYNC_COPY] [IMPL:TOOLBAR_SYNC_DELETE] [IMPL:TOOLBAR_SYNC_RENAME] [IMPL:TOOLBAR_IGNORE_FAILURES]
func TestInvokeToolbarSyncButtonsWithNilCallback_REQ_TOOLBAR_SYNC_BUTTONS(t *testing.T) {
	// Save and restore all callbacks
	origCopy := toolbarSyncCopyFn
	origDelete := toolbarSyncDeleteFn
	origRename := toolbarSyncRenameFn
	origIgnore := toolbarIgnoreFailuresFn
	t.Cleanup(func() {
		toolbarSyncCopyFn = origCopy
		toolbarSyncDeleteFn = origDelete
		toolbarSyncRenameFn = origRename
		toolbarIgnoreFailuresFn = origIgnore
	})

	// Clear all callbacks
	toolbarSyncCopyFn = nil
	toolbarSyncDeleteFn = nil
	toolbarSyncRenameFn = nil
	toolbarIgnoreFailuresFn = nil

	// All should return false when no callback is set
	if InvokeToolbarButton("synccopy") {
		t.Errorf("InvokeToolbarButton(\"synccopy\") with nil callback = true, want false")
	}
	if InvokeToolbarButton("syncdelete") {
		t.Errorf("InvokeToolbarButton(\"syncdelete\") with nil callback = true, want false")
	}
	if InvokeToolbarButton("syncrename") {
		t.Errorf("InvokeToolbarButton(\"syncrename\") with nil callback = true, want false")
	}
	if InvokeToolbarButton("ignorefailures") {
		t.Errorf("InvokeToolbarButton(\"ignorefailures\") with nil callback = true, want false")
	}
}

// TestIgnoreFailuresIndicator_REQ_TOOLBAR_SYNC_BUTTONS tests the indicator function for button styling.
// [REQ:TOOLBAR_SYNC_BUTTONS] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_IGNORE_FAILURES]
func TestIgnoreFailuresIndicator_REQ_TOOLBAR_SYNC_BUTTONS(t *testing.T) {
	originalIndicator := syncIgnoreFailuresIndicator
	t.Cleanup(func() {
		syncIgnoreFailuresIndicator = originalIndicator
	})

	// Test with nil indicator
	syncIgnoreFailuresIndicator = nil
	// Should not panic when indicator is nil
	_ = syncIgnoreFailuresIndicator

	// Test with false indicator
	SetSyncIgnoreFailuresIndicator(func() bool { return false })
	if syncIgnoreFailuresIndicator() {
		t.Errorf("syncIgnoreFailuresIndicator() = true, want false")
	}

	// Test with true indicator
	SetSyncIgnoreFailuresIndicator(func() bool { return true })
	if !syncIgnoreFailuresIndicator() {
		t.Errorf("syncIgnoreFailuresIndicator() = false, want true")
	}
}
