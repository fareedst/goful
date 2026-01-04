package cmdline

import (
	"testing"

	"github.com/anmitsu/goful/widget"
)

type fakeMode struct {
	name   string
	prompt string
}

func (m fakeMode) String() string { return m.name }
func (m fakeMode) Prompt() string { return m.prompt }
func (fakeMode) Draw(*Cmdline)    {}
func (fakeMode) Run(*Cmdline)     {}

func newTestCmdline(t *testing.T) *Cmdline {
	t.Helper()
	mode := fakeMode{name: "test-mode", prompt: ":"}
	cmd := &Cmdline{
		TextBox: widget.NewTextBox(0, 0, 40, 1),
		mode:    mode,
	}
	cmd.History = NewHistory(0, 0, 40, 4, cmd)
	return cmd
}

func TestHistoryAddDedup_REQ_CMD_HANDLER_TESTS(t *testing.T) {
	// [REQ:CMD_HANDLER_TESTS] [ARCH:TEST_STRATEGY_CMD] [IMPL:TEST_CMDLINE]
	cmd := newTestCmdline(t)
	mode := cmd.mode.String()
	historyMap[mode] = []string{"first", "second"}
	t.Cleanup(func() {
		delete(historyMap, mode)
	})

	cmd.SetText("second")
	cmd.History.add()
	if got := historyMap[mode]; len(got) != 2 || got[1] != "second" {
		t.Fatalf("duplicate should move to tail, got %v", got)
	}

	cmd.SetText("")
	cmd.History.add()
	if got := historyMap[mode]; len(got) != 2 {
		t.Fatalf("empty input should not change history, got %v", got)
	}

	cmd.SetText(" leading")
	cmd.History.add()
	if got := historyMap[mode]; len(got) != 2 {
		t.Fatalf("inputs starting with space ignored, got %v", got)
	}

	cmd.SetText("third")
	cmd.History.add()
	if got := historyMap[mode]; len(got) != 3 || got[2] != "third" {
		t.Fatalf("new entry should append, got %v", got)
	}
}

func TestHistoryCursorMovement_REQ_CMD_HANDLER_TESTS(t *testing.T) {
	// [REQ:CMD_HANDLER_TESTS] [ARCH:TEST_STRATEGY_CMD] [IMPL:TEST_CMDLINE]
	cmd := newTestCmdline(t)
	cmd.SetText("seed")
	cmd.History.input = cmd.String()
	cmd.History.AppendString("alpha", "beta")
	cmd.History.SetCursor(cmd.History.Lower())

	cmd.History.MoveCursor(1)
	if got := cmd.String(); got != "alpha" {
		t.Fatalf("MoveCursor onto history should copy entry, got %q", got)
	}

	cmd.History.CursorDown()
	if got := cmd.String(); got != "beta" {
		t.Fatalf("CursorDown should move to next entry, got %q", got)
	}

	cmd.History.MoveCursor(cmd.History.Lower() - cmd.History.Cursor())
	if got := cmd.String(); got != cmd.History.input {
		t.Fatalf("Returning to lower should restore input, got %q", got)
	}
}
