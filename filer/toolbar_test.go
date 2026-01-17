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
