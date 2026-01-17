package cmdline

import (
	"errors"
	"os"
	"path/filepath"
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

// [REQ:DEBT_TRIAGE] [IMPL:HISTORY_ERROR_HANDLING] [ARCH:DEBT_MANAGEMENT]
// Test that LoadHistory treats missing files (first-run) as success.
func TestLoadHistory_FirstRunSuccess_REQ_DEBT_TRIAGE(t *testing.T) {
	// Create a path that definitely does not exist
	nonExistentPath := filepath.Join(t.TempDir(), "nonexistent", "history")
	err := LoadHistory(nonExistentPath)
	if err != nil {
		t.Fatalf("LoadHistory should return nil for missing file (first-run), got: %v", err)
	}
}

// [REQ:DEBT_TRIAGE] [IMPL:HISTORY_ERROR_HANDLING] [ARCH:DEBT_MANAGEMENT]
// Test that LoadHistory returns HistoryError for permission errors.
func TestLoadHistory_PermissionError_REQ_DEBT_TRIAGE(t *testing.T) {
	// Create a file with no read permissions
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "noperm_history")
	if err := os.WriteFile(testFile, []byte("test"), 0000); err != nil {
		t.Skipf("cannot create file without permissions: %v", err)
	}
	t.Cleanup(func() { os.Chmod(testFile, 0644) }) // Restore for cleanup

	err := LoadHistory(testFile)
	if err == nil {
		t.Fatal("LoadHistory should return error for unreadable file")
	}

	var histErr *HistoryError
	if !errors.As(err, &histErr) {
		t.Fatalf("error should be *HistoryError, got %T", err)
	}
	if histErr.Op != "load" {
		t.Errorf("Op should be 'load', got %q", histErr.Op)
	}
	if histErr.IsFirstRun() {
		t.Error("IsFirstRun should be false for permission errors")
	}
}

// [REQ:DEBT_TRIAGE] [IMPL:HISTORY_ERROR_HANDLING] [ARCH:DEBT_MANAGEMENT]
// Test that LoadHistory successfully loads valid history files.
func TestLoadHistory_ValidFile_REQ_DEBT_TRIAGE(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "valid_history")
	if err := os.WriteFile(testFile, []byte("cmd1\ncmd2\ncmd3\n"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err := LoadHistory(testFile)
	if err != nil {
		t.Fatalf("LoadHistory should succeed for valid file, got: %v", err)
	}

	key := filepath.Base(testFile)
	history, ok := historyMap[key]
	if !ok {
		t.Fatalf("historyMap should contain key %q", key)
	}
	if len(history) != 3 {
		t.Errorf("expected 3 entries, got %d", len(history))
	}
	if history[0] != "cmd1" || history[1] != "cmd2" || history[2] != "cmd3" {
		t.Errorf("unexpected history content: %v", history)
	}

	t.Cleanup(func() { delete(historyMap, key) })
}

// [REQ:DEBT_TRIAGE] [IMPL:HISTORY_ERROR_HANDLING] [ARCH:DEBT_MANAGEMENT]
// Test that SaveHistory returns HistoryError for permission errors.
func TestSaveHistory_PermissionError_REQ_DEBT_TRIAGE(t *testing.T) {
	// Create a directory without write permissions
	tmpDir := t.TempDir()
	noWriteDir := filepath.Join(tmpDir, "nowrite")
	if err := os.Mkdir(noWriteDir, 0555); err != nil {
		t.Skipf("cannot create directory without write permissions: %v", err)
	}
	t.Cleanup(func() { os.Chmod(noWriteDir, 0755) }) // Restore for cleanup

	testFile := filepath.Join(noWriteDir, "history")
	err := SaveHistory(testFile)
	if err == nil {
		t.Fatal("SaveHistory should return error for unwritable directory")
	}

	var histErr *HistoryError
	if !errors.As(err, &histErr) {
		t.Fatalf("error should be *HistoryError, got %T", err)
	}
	if histErr.Op != "save" {
		t.Errorf("Op should be 'save', got %q", histErr.Op)
	}
}

// [REQ:DEBT_TRIAGE] [IMPL:HISTORY_ERROR_HANDLING] [ARCH:DEBT_MANAGEMENT]
// Test that SaveHistory successfully saves history.
func TestSaveHistory_ValidSave_REQ_DEBT_TRIAGE(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "subdir", "save_history")

	key := filepath.Base(testFile)
	historyMap[key] = []string{"saved1", "saved2"}
	t.Cleanup(func() { delete(historyMap, key) })

	err := SaveHistory(testFile)
	if err != nil {
		t.Fatalf("SaveHistory should succeed, got: %v", err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}
	expected := "saved1\nsaved2\n"
	if string(content) != expected {
		t.Errorf("expected %q, got %q", expected, string(content))
	}
}

// [REQ:DEBT_TRIAGE] [IMPL:HISTORY_ERROR_HANDLING] [ARCH:DEBT_MANAGEMENT]
// Test HistoryError Error() and Unwrap() methods.
func TestHistoryError_Methods_REQ_DEBT_TRIAGE(t *testing.T) {
	baseErr := errors.New("underlying error")
	histErr := &HistoryError{
		Path: "/path/to/history",
		Op:   "load",
		Err:  baseErr,
	}

	errStr := histErr.Error()
	if errStr != "load history /path/to/history: underlying error" {
		t.Errorf("unexpected Error() output: %s", errStr)
	}

	unwrapped := histErr.Unwrap()
	if unwrapped != baseErr {
		t.Errorf("Unwrap() should return base error")
	}

	// Test IsFirstRun for os.ErrNotExist
	histErr.Err = os.ErrNotExist
	if !histErr.IsFirstRun() {
		t.Error("IsFirstRun should be true for os.ErrNotExist")
	}

	histErr.Err = baseErr
	if histErr.IsFirstRun() {
		t.Error("IsFirstRun should be false for non-ErrNotExist")
	}
}

// [REQ:DEBT_TRIAGE] [IMPL:HISTORY_CACHE_LIMIT] [ARCH:DEBT_MANAGEMENT]
// Test that trimHistory enforces HistoryLimit.
func TestTrimHistory_REQ_DEBT_TRIAGE(t *testing.T) {
	// Test with limit disabled
	oldLimit := HistoryLimit
	HistoryLimit = 0
	t.Cleanup(func() { HistoryLimit = oldLimit })

	history := []string{"a", "b", "c", "d", "e"}
	result := trimHistory(history)
	if len(result) != 5 {
		t.Errorf("limit=0 should not trim, got %d entries", len(result))
	}

	// Test with limit enabled
	HistoryLimit = 3
	result = trimHistory(history)
	if len(result) != 3 {
		t.Errorf("limit=3 should trim to 3 entries, got %d", len(result))
	}
	// Should keep most recent (tail)
	if result[0] != "c" || result[1] != "d" || result[2] != "e" {
		t.Errorf("should keep most recent entries, got %v", result)
	}

	// Test with history smaller than limit
	HistoryLimit = 10
	result = trimHistory(history)
	if len(result) != 5 {
		t.Errorf("history smaller than limit should not trim, got %d", len(result))
	}

	// Test exact limit boundary
	HistoryLimit = 5
	result = trimHistory(history)
	if len(result) != 5 {
		t.Errorf("history at limit should not trim, got %d", len(result))
	}
}

// [REQ:DEBT_TRIAGE] [IMPL:HISTORY_CACHE_LIMIT] [ARCH:DEBT_MANAGEMENT]
// Test that history add respects HistoryLimit.
func TestHistoryAddWithLimit_REQ_DEBT_TRIAGE(t *testing.T) {
	cmd := newTestCmdline(t)
	mode := cmd.mode.String()

	oldLimit := HistoryLimit
	HistoryLimit = 3
	t.Cleanup(func() {
		HistoryLimit = oldLimit
		delete(historyMap, mode)
	})

	// Add more entries than the limit
	for i := 0; i < 5; i++ {
		cmd.SetText("cmd" + string(rune('A'+i)))
		cmd.History.add()
	}

	history := historyMap[mode]
	if len(history) != 3 {
		t.Errorf("history should be capped at 3, got %d", len(history))
	}
	// Should have most recent 3 entries
	if history[0] != "cmdC" || history[1] != "cmdD" || history[2] != "cmdE" {
		t.Errorf("should keep most recent entries, got %v", history)
	}
}

// [REQ:DEBT_TRIAGE] [IMPL:HISTORY_CACHE_LIMIT] [ARCH:DEBT_MANAGEMENT]
// Test that SaveHistory applies trimming before persisting.
func TestSaveHistory_TrimsBeforeSave_REQ_DEBT_TRIAGE(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "trim_history")
	key := filepath.Base(testFile)

	oldLimit := HistoryLimit
	HistoryLimit = 2
	t.Cleanup(func() {
		HistoryLimit = oldLimit
		delete(historyMap, key)
	})

	// Put 5 entries in historyMap
	historyMap[key] = []string{"old1", "old2", "old3", "recent1", "recent2"}

	err := SaveHistory(testFile)
	if err != nil {
		t.Fatalf("SaveHistory failed: %v", err)
	}

	// Check that historyMap was also trimmed in memory
	if len(historyMap[key]) != 2 {
		t.Errorf("historyMap should be trimmed to 2, got %d", len(historyMap[key]))
	}

	// Check persisted file content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}
	expected := "recent1\nrecent2\n"
	if string(content) != expected {
		t.Errorf("expected %q, got %q", expected, string(content))
	}
}
