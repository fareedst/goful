package app

import (
	"testing"
	"time"
)

// TestIsDoubleClick_REQ_MOUSE_DOUBLE_CLICK tests the double-click detection logic.
// [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]
func TestIsDoubleClick_REQ_MOUSE_DOUBLE_CLICK(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupClickTime time.Time
		setupX, setupY int
		clickX, clickY int
		wantDouble     bool
	}{
		{
			name:           "first click is never double",
			setupClickTime: time.Time{}, // zero time
			setupX:         0,
			setupY:         0,
			clickX:         10,
			clickY:         20,
			wantDouble:     false,
		},
		{
			name:           "same position within threshold is double",
			setupClickTime: time.Now().Add(-100 * time.Millisecond),
			setupX:         10,
			setupY:         20,
			clickX:         10,
			clickY:         20,
			wantDouble:     true,
		},
		{
			name:           "same position outside threshold is not double",
			setupClickTime: time.Now().Add(-500 * time.Millisecond),
			setupX:         10,
			setupY:         20,
			clickX:         10,
			clickY:         20,
			wantDouble:     false,
		},
		{
			name:           "different X position is not double",
			setupClickTime: time.Now().Add(-100 * time.Millisecond),
			setupX:         10,
			setupY:         20,
			clickX:         11,
			clickY:         20,
			wantDouble:     false,
		},
		{
			name:           "different Y position is not double",
			setupClickTime: time.Now().Add(-100 * time.Millisecond),
			setupX:         10,
			setupY:         20,
			clickX:         10,
			clickY:         21,
			wantDouble:     false,
		},
		{
			name:           "both positions different is not double",
			setupClickTime: time.Now().Add(-100 * time.Millisecond),
			setupX:         10,
			setupY:         20,
			clickX:         15,
			clickY:         25,
			wantDouble:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &Goful{
				lastClickTime: tc.setupClickTime,
				lastClickX:    tc.setupX,
				lastClickY:    tc.setupY,
			}

			got := g.isDoubleClick(tc.clickX, tc.clickY)
			if got != tc.wantDouble {
				t.Errorf("isDoubleClick(%d, %d) = %v, want %v", tc.clickX, tc.clickY, got, tc.wantDouble)
			}

			// Verify state was updated
			if g.lastClickX != tc.clickX || g.lastClickY != tc.clickY {
				t.Errorf("click position not updated: got (%d, %d), want (%d, %d)",
					g.lastClickX, g.lastClickY, tc.clickX, tc.clickY)
			}
		})
	}
}

// TestDoubleClickThreshold_REQ_MOUSE_DOUBLE_CLICK verifies the threshold constant is reasonable.
// [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]
func TestDoubleClickThreshold_REQ_MOUSE_DOUBLE_CLICK(t *testing.T) {
	t.Parallel()

	// Threshold should be between 200ms and 600ms for usability
	if doubleClickThreshold < 200*time.Millisecond {
		t.Errorf("doubleClickThreshold too short: %v", doubleClickThreshold)
	}
	if doubleClickThreshold > 600*time.Millisecond {
		t.Errorf("doubleClickThreshold too long: %v", doubleClickThreshold)
	}
}

// TestIsDoubleClickUpdatesState_REQ_MOUSE_DOUBLE_CLICK verifies click state is always updated.
// [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]
func TestIsDoubleClickUpdatesState_REQ_MOUSE_DOUBLE_CLICK(t *testing.T) {
	t.Parallel()

	g := &Goful{}

	// First click
	beforeClick := time.Now()
	g.isDoubleClick(100, 200)
	afterClick := time.Now()

	if g.lastClickX != 100 || g.lastClickY != 200 {
		t.Errorf("position not updated: got (%d, %d), want (100, 200)", g.lastClickX, g.lastClickY)
	}
	if g.lastClickTime.Before(beforeClick) || g.lastClickTime.After(afterClick) {
		t.Errorf("time not updated correctly: %v not in [%v, %v]", g.lastClickTime, beforeClick, afterClick)
	}

	// Second click at different position
	g.isDoubleClick(300, 400)
	if g.lastClickX != 300 || g.lastClickY != 400 {
		t.Errorf("position not updated on second click: got (%d, %d), want (300, 400)", g.lastClickX, g.lastClickY)
	}
}

// TestDoubleClickSequence_REQ_MOUSE_DOUBLE_CLICK tests a realistic click sequence.
// [IMPL:MOUSE_DOUBLE_CLICK] [ARCH:MOUSE_DOUBLE_CLICK] [REQ:MOUSE_DOUBLE_CLICK]
func TestDoubleClickSequence_REQ_MOUSE_DOUBLE_CLICK(t *testing.T) {
	t.Parallel()

	g := &Goful{}

	// First click - not a double-click
	if g.isDoubleClick(50, 100) {
		t.Error("first click should not be double-click")
	}

	// Second click at same position - should be double-click
	if !g.isDoubleClick(50, 100) {
		t.Error("second click at same position should be double-click")
	}

	// Third click at same position - also triggers as double-click
	// (continuous rapid clicks are all treated as double-clicks)
	if !g.isDoubleClick(50, 100) {
		t.Error("third click at same position within threshold should still be double-click")
	}

	// Click at different position - not a double-click
	if g.isDoubleClick(60, 110) {
		t.Error("click at different position should not be double-click")
	}

	// Click back at original position - not a double-click (different from previous)
	if g.isDoubleClick(50, 100) {
		t.Error("click at original position after moving should not be double-click")
	}

	// Second click at original - should be double-click again
	if !g.isDoubleClick(50, 100) {
		t.Error("second click at same position should be double-click")
	}
}
