package main

import (
	"testing"

	"github.com/anmitsu/goful/app"
	"github.com/anmitsu/goful/cmdline"
	"github.com/anmitsu/goful/filer"
	"github.com/anmitsu/goful/menu"
	"github.com/anmitsu/goful/widget"
)

// [TEST:KEYMAP_BASELINE] helper for asserting canonical bindings exist.
func assertKeysPresent(t *testing.T, label string, km widget.Keymap, keys []string) {
	t.Helper()
	missing := make([]string, 0)
	for _, key := range keys {
		if _, ok := km[key]; !ok {
			missing = append(missing, key)
		}
	}
	t.Logf("DEBUG: [REQ:BEHAVIOR_BASELINE] [ARCH:BASELINE_CAPTURE] [IMPL:BASELINE_SNAPSHOTS] [TEST:KEYMAP_BASELINE] asserting %s keys=%v missing=%v", label, keys, missing)
	if len(missing) > 0 {
		t.Fatalf("%s keymap missing bindings: %v", label, missing)
	}
}

func TestFilerKeymapBaseline_REQ_BEHAVIOR_BASELINE(t *testing.T) {
	// [REQ:BEHAVIOR_BASELINE] [ARCH:BASELINE_CAPTURE] [IMPL:BASELINE_SNAPSHOTS] [TEST:KEYMAP_BASELINE]
	required := []string{
		"j", "k", "h", "l", " ", "C-n", "C-p", "C-d", "C-u", "C-a", "C-e",
		"q", "Q", ":", ";", "f", "/", "n", "K", "c", "m", "r", "R", "D", "d", "g", "G",
	}
	km := filerKeymap((*app.Goful)(nil))
	assertKeysPresent(t, "filer", km, required)
}

func TestCmdlineKeymapBaseline_REQ_BEHAVIOR_BASELINE(t *testing.T) {
	// [REQ:BEHAVIOR_BASELINE] [ARCH:BASELINE_CAPTURE] [IMPL:BASELINE_SNAPSHOTS] [TEST:KEYMAP_BASELINE]
	required := []string{
		"C-a", "C-e", "C-f", "C-b", "M-f", "M-b", "C-h", "backspace",
		"C-d", "delete", "M-d", "M-h", "C-k", "C-i", "C-m", "C-g", "C-[",
		"C-n", "C-p", "down", "up", "C-v", "M-v", "pgdn", "pgup", "M-<", "M->",
	}
	km := cmdlineKeymap(&cmdline.Cmdline{})
	assertKeysPresent(t, "cmdline", km, required)
}

func TestAuxKeymapsBaseline_REQ_BEHAVIOR_BASELINE(t *testing.T) {
	// [REQ:BEHAVIOR_BASELINE] [ARCH:BASELINE_CAPTURE] [IMPL:BASELINE_SNAPSHOTS] [TEST:KEYMAP_BASELINE]
	finderRequired := []string{"C-h", "backspace", "M-p", "M-n", "C-g", "C-["}
	completionRequired := []string{"C-n", "C-p", "C-f", "C-b", "C-i", "C-m", "C-g", "C-["}
	menuRequired := []string{"C-n", "C-p", "down", "up", "C-m", "C-g", "C-[", "C-v", "M-v"}

	assertKeysPresent(t, "finder", finderKeymap((*filer.Finder)(nil)), finderRequired)
	assertKeysPresent(t, "completion", completionKeymap((*cmdline.Completion)(nil)), completionRequired)
	assertKeysPresent(t, "menu", menuKeymap((*menu.Menu)(nil)), menuRequired)
}
