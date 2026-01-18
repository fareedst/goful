package filer

import "testing"

// resetToolbarButtonsForTest clears the toolbar button bounds map.
func resetToolbarButtonsForTest() {
	toolbarButtons = make(map[string]toolbarBounds)
}

// TestToolbarButtonAt_REQ_TOOLBAR_PARENT_BUTTON tests hit-testing for toolbar buttons.
// [REQ:TOOLBAR_PARENT_BUTTON] [REQ:TOOLBAR_BUTTON_STYLING] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_PARENT_BUTTON] [IMPL:TOOLBAR_BUTTON_STYLING]
func TestToolbarButtonAt_REQ_TOOLBAR_PARENT_BUTTON(t *testing.T) {
	t.Cleanup(resetToolbarButtonsForTest)

	// Setup: Register a parent button at x=0, y=0 (single char "^")
	toolbarButtons["parent"] = toolbarBounds{x1: 0, y: 0, x2: 0}

	tests := []struct {
		name     string
		x, y     int
		expected string
	}{
		{"click at x=0, y=0 (on button)", 0, 0, "parent"},
		{"click at x=1, y=0 (just outside)", 1, 0, ""},
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
// [REQ:TOOLBAR_PARENT_BUTTON] [REQ:TOOLBAR_BUTTON_STYLING] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_PARENT_BUTTON] [IMPL:TOOLBAR_BUTTON_STYLING]
func TestToolbarButtonAtMultipleButtons_REQ_TOOLBAR_PARENT_BUTTON(t *testing.T) {
	t.Cleanup(resetToolbarButtonsForTest)

	// Setup: Register multiple buttons on the same row (single char buttons with space between)
	toolbarButtons["parent"] = toolbarBounds{x1: 0, y: 0, x2: 0} // "^"
	toolbarButtons["reload"] = toolbarBounds{x1: 2, y: 0, x2: 2} // hypothetical

	tests := []struct {
		name     string
		x, y     int
		expected string
	}{
		{"click on parent button", 0, 0, "parent"},
		{"click on reload button", 2, 0, "reload"},
		{"click between buttons", 1, 0, ""},
		{"click after all buttons", 3, 0, ""},
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
// [REQ:TOOLBAR_LINKED_TOGGLE] [REQ:TOOLBAR_BUTTON_STYLING] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_LINKED_TOGGLE] [IMPL:TOOLBAR_BUTTON_STYLING]
func TestToolbarLinkedButtonHit_REQ_TOOLBAR_LINKED_TOGGLE(t *testing.T) {
	t.Cleanup(resetToolbarButtonsForTest)

	// Setup: Register parent and linked buttons as they would appear in the header
	// Layout: ^ L ... (single char buttons with space between)
	toolbarButtons["parent"] = toolbarBounds{x1: 0, y: 0, x2: 0} // "^"
	toolbarButtons["linked"] = toolbarBounds{x1: 2, y: 0, x2: 2} // "L"

	tests := []struct {
		name     string
		x, y     int
		expected string
	}{
		{"click on parent button", 0, 0, "parent"},
		{"click on linked button", 2, 0, "linked"},
		{"click between buttons (space)", 1, 0, ""},
		{"click after linked button", 3, 0, ""},
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
// [REQ:TOOLBAR_COMPARE_BUTTON] [REQ:TOOLBAR_BUTTON_STYLING] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_COMPARE_BUTTON] [IMPL:TOOLBAR_BUTTON_STYLING]
func TestToolbarCompareButtonHit_REQ_TOOLBAR_COMPARE_BUTTON(t *testing.T) {
	t.Cleanup(resetToolbarButtonsForTest)

	// Setup: Register parent, linked, and compare buttons as they would appear in the header
	// Layout: ^ L = ... (single char buttons with space between)
	toolbarButtons["parent"] = toolbarBounds{x1: 0, y: 0, x2: 0}  // "^"
	toolbarButtons["linked"] = toolbarBounds{x1: 2, y: 0, x2: 2}  // "L"
	toolbarButtons["compare"] = toolbarBounds{x1: 4, y: 0, x2: 4} // "="

	tests := []struct {
		name     string
		x, y     int
		expected string
	}{
		{"click on parent button", 0, 0, "parent"},
		{"click on linked button", 2, 0, "linked"},
		{"click on compare button", 4, 0, "compare"},
		{"click between linked and compare", 3, 0, ""},
		{"click after compare button", 5, 0, ""},
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
// [REQ:TOOLBAR_SYNC_BUTTONS] [REQ:TOOLBAR_BUTTON_STYLING] [ARCH:TOOLBAR_LAYOUT] [IMPL:TOOLBAR_SYNC_COPY] [IMPL:TOOLBAR_SYNC_DELETE] [IMPL:TOOLBAR_SYNC_RENAME] [IMPL:TOOLBAR_IGNORE_FAILURES] [IMPL:TOOLBAR_BUTTON_STYLING]
func TestToolbarSyncButtonsHit_REQ_TOOLBAR_SYNC_BUTTONS(t *testing.T) {
	t.Cleanup(resetToolbarButtonsForTest)

	// Setup: Register all toolbar buttons as they would appear in the header
	// Layout: ^ L = C D R ! ... (single char buttons with space between)
	toolbarButtons["parent"] = toolbarBounds{x1: 0, y: 0, x2: 0}           // "^"
	toolbarButtons["linked"] = toolbarBounds{x1: 2, y: 0, x2: 2}           // "L"
	toolbarButtons["compare"] = toolbarBounds{x1: 4, y: 0, x2: 4}          // "="
	toolbarButtons["synccopy"] = toolbarBounds{x1: 6, y: 0, x2: 6}         // "C"
	toolbarButtons["syncdelete"] = toolbarBounds{x1: 8, y: 0, x2: 8}       // "D"
	toolbarButtons["syncrename"] = toolbarBounds{x1: 10, y: 0, x2: 10}     // "R"
	toolbarButtons["ignorefailures"] = toolbarBounds{x1: 12, y: 0, x2: 12} // "!"

	tests := []struct {
		name     string
		x, y     int
		expected string
	}{
		{"click on synccopy button", 6, 0, "synccopy"},
		{"click on syncdelete button", 8, 0, "syncdelete"},
		{"click on syncrename button", 10, 0, "syncrename"},
		{"click on ignorefailures button", 12, 0, "ignorefailures"},
		{"click between compare and synccopy", 5, 0, ""},
		{"click after ignorefailures button", 13, 0, ""},
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

// resetWorkspaceTabsForTest clears the workspace tab bounds map.
func resetWorkspaceTabsForTest() {
	workspaceTabs = make(map[int]workspaceTabBounds)
}

// TestWorkspaceTabAt_REQ_CLICKABLE_WORKSPACE_TABS tests hit-testing for workspace tabs.
// [REQ:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [IMPL:CLICKABLE_WORKSPACE_TABS]
func TestWorkspaceTabAt_REQ_CLICKABLE_WORKSPACE_TABS(t *testing.T) {
	t.Cleanup(resetWorkspaceTabsForTest)

	// Setup: Register workspace tabs as they would appear in the header
	// Layout: «1» «2» «3» (guillemet + number + guillemet = 3 chars each, with space between)
	// Assuming they start at x=50 on row 0
	workspaceTabs[0] = workspaceTabBounds{x1: 50, y: 0, x2: 52} // «1»
	workspaceTabs[1] = workspaceTabBounds{x1: 54, y: 0, x2: 56} // «2»
	workspaceTabs[2] = workspaceTabBounds{x1: 58, y: 0, x2: 60} // «3»

	tests := []struct {
		name     string
		x, y     int
		expected int
	}{
		{"click on workspace 1", 50, 0, 0},
		{"click on workspace 1 middle", 51, 0, 0},
		{"click on workspace 1 end", 52, 0, 0},
		{"click on workspace 2", 54, 0, 1},
		{"click on workspace 3", 58, 0, 2},
		{"click between tabs (space)", 53, 0, -1},
		{"click before all tabs", 49, 0, -1},
		{"click after all tabs", 61, 0, -1},
		{"click on wrong row", 50, 1, -1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := WorkspaceTabAt(tc.x, tc.y)
			if result != tc.expected {
				t.Errorf("WorkspaceTabAt(%d, %d) = %d, want %d", tc.x, tc.y, result, tc.expected)
			}
		})
	}
}

// TestInvokeWorkspaceTab_REQ_CLICKABLE_WORKSPACE_TABS tests that tab click invokes callback.
// [REQ:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [IMPL:CLICKABLE_WORKSPACE_TABS]
func TestInvokeWorkspaceTab_REQ_CLICKABLE_WORKSPACE_TABS(t *testing.T) {
	// Setup: Track callback invocation
	var clickedIndex int = -1
	originalFn := workspaceTabClickFn
	t.Cleanup(func() {
		workspaceTabClickFn = originalFn
	})

	SetWorkspaceTabClickFn(func(index int) {
		clickedIndex = index
	})

	// Test invoking tab 0
	handled := InvokeWorkspaceTab(0)
	if !handled {
		t.Errorf("InvokeWorkspaceTab(0) = false, want true")
	}
	if clickedIndex != 0 {
		t.Errorf("Callback received index %d, want 0", clickedIndex)
	}

	// Test invoking tab 2
	handled = InvokeWorkspaceTab(2)
	if !handled {
		t.Errorf("InvokeWorkspaceTab(2) = false, want true")
	}
	if clickedIndex != 2 {
		t.Errorf("Callback received index %d, want 2", clickedIndex)
	}
}

// TestInvokeWorkspaceTabWithNilCallback_REQ_CLICKABLE_WORKSPACE_TABS tests behavior when no callback is set.
// [REQ:CLICKABLE_WORKSPACE_TABS] [ARCH:CLICKABLE_WORKSPACE_TABS] [IMPL:CLICKABLE_WORKSPACE_TABS]
func TestInvokeWorkspaceTabWithNilCallback_REQ_CLICKABLE_WORKSPACE_TABS(t *testing.T) {
	// Setup: Clear the callback
	originalFn := workspaceTabClickFn
	t.Cleanup(func() {
		workspaceTabClickFn = originalFn
	})
	workspaceTabClickFn = nil

	// Should return false when no callback is set
	handled := InvokeWorkspaceTab(0)
	if handled {
		t.Errorf("InvokeWorkspaceTab(0) with nil callback = true, want false")
	}
}
