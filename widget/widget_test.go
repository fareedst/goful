package widget

import (
	"bytes"
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestInsertBytes(t *testing.T) {
	for _, d := range []struct {
		s      []byte
		data   []byte
		offset int
		result []byte
	}{
		{
			[]byte("Hello world"),
			[]byte("my "),
			6,
			[]byte("Hello my world"),
		},
		{
			[]byte("こんにちは世界"),
			[]byte("私の"),
			15,
			[]byte("こんにちは私の世界"),
		},
		{
			[]byte("こんにちは△□の世界"),
			[]byte("○✕"),
			15,
			[]byte("こんにちは○✕△□の世界"),
		},
	} {
		s := InsertBytes(d.s, d.data, d.offset)
		if !bytes.Equal(s, d.result) {
			t.Errorf("InsertBytes(%q, %q, %q)=%q, want %q", d.s, d.data, d.offset, s, d.result)
		}
	}
}

func TestDeleteBytes(t *testing.T) {
	for _, d := range []struct {
		s      []byte
		offset int
		length int
		result []byte
	}{
		{
			[]byte("Hello my world"),
			6,
			3,
			[]byte("Hello world"),
		},
		{
			[]byte("こんにちは私の世界"),
			15,
			6,
			[]byte("こんにちは世界"),
		},
		{
			[]byte("こんにちは○✕△□の世界"),
			15,
			15,
			[]byte("こんにちは世界"),
		},
	} {
		s := DeleteBytes(d.s, d.offset, d.length)
		if !bytes.Equal(s, d.result) {
			t.Errorf("DeleteBytes(%q, %q, %q)=%q, want %q", d.s, d.offset, d.length, s, d.result)
		}
	}
}

func TestEventToStringReturnKey_REQ_QUIT_DIALOG_DEFAULT(t *testing.T) {
	t.Helper()
	// [IMPL:QUIT_DIALOG_ENTER] [ARCH:QUIT_DIALOG_KEYS] [REQ:QUIT_DIALOG_DEFAULT]
	for _, tc := range []struct {
		name string
		key  tcell.Key
	}{
		{name: "enter", key: tcell.KeyEnter},
		{name: "ctrl-m", key: tcell.KeyCtrlM},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ev := tcell.NewEventKey(tc.key, 0, tcell.ModNone)
			if got := EventToString(ev); got != "C-m" {
				t.Fatalf("EventToString(%s)=%q, want %q", tc.name, got, "C-m")
			}
		})
	}
}

func TestEventToStringBackspace_REQ_BACKSPACE_BEHAVIOR(t *testing.T) {
	t.Helper()
	// [IMPL:BACKSPACE_TRANSLATION] [ARCH:BACKSPACE_TRANSLATION] [REQ:BACKSPACE_BEHAVIOR]
	for _, tc := range []struct {
		name string
		key  tcell.Key
	}{
		{name: "backspace", key: tcell.KeyBackspace},
		{name: "backspace2", key: tcell.KeyBackspace2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ev := tcell.NewEventKey(tc.key, 0, tcell.ModNone)
			if got := EventToString(ev); got != "backspace" {
				t.Fatalf("EventToString(%s)=%q, want %q", tc.name, got, "backspace")
			}
		})
	}
}

func TestWindowContains_REQ_MOUSE_FILE_SELECT(t *testing.T) {
	// [IMPL:MOUSE_HIT_TEST] [ARCH:MOUSE_EVENT_ROUTING] [REQ:MOUSE_FILE_SELECT]
	// Test hit-testing for mouse coordinate detection
	for _, tc := range []struct {
		name     string
		x, y     int
		w, h     int
		testX    int
		testY    int
		expected bool
	}{
		{name: "inside_center", x: 10, y: 5, w: 20, h: 10, testX: 15, testY: 8, expected: true},
		{name: "top_left_corner", x: 10, y: 5, w: 20, h: 10, testX: 10, testY: 5, expected: true},
		{name: "bottom_right_corner", x: 10, y: 5, w: 20, h: 10, testX: 29, testY: 14, expected: true},
		{name: "left_of_window", x: 10, y: 5, w: 20, h: 10, testX: 9, testY: 8, expected: false},
		{name: "right_of_window", x: 10, y: 5, w: 20, h: 10, testX: 30, testY: 8, expected: false},
		{name: "above_window", x: 10, y: 5, w: 20, h: 10, testX: 15, testY: 4, expected: false},
		{name: "below_window", x: 10, y: 5, w: 20, h: 10, testX: 15, testY: 15, expected: false},
		{name: "origin_window", x: 0, y: 0, w: 80, h: 24, testX: 0, testY: 0, expected: true},
		{name: "origin_outside", x: 0, y: 0, w: 80, h: 24, testX: 80, testY: 24, expected: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := NewWindow(tc.x, tc.y, tc.w, tc.h)
			got := w.Contains(tc.testX, tc.testY)
			if got != tc.expected {
				t.Errorf("Window(%d,%d,%d,%d).Contains(%d,%d)=%v, want %v",
					tc.x, tc.y, tc.w, tc.h, tc.testX, tc.testY, got, tc.expected)
			}
		})
	}
}
