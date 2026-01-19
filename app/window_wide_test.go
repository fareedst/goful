// Package app provides tests for sync command operations.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/anmitsu/goful/filer"
)

// TestFindFileByName_REQ_SYNC_COMMANDS tests the FindFileByName helper.
func TestFindFileByName_REQ_SYNC_COMMANDS(t *testing.T) {
	// Save original directory at the start (before any other tests might change it)
	origDir, err := os.Getwd()
	if err != nil {
		// Fallback to system temp directory if current directory is deleted
		// (can happen during parallel test execution on Ubuntu)
		origDir = os.TempDir()
	}

	// Create a temp directory with test files
	tempDir := t.TempDir()

	// Create test files
	testFiles := []string{"file1.txt", "file2.txt", "subdir"}
	for _, name := range testFiles[:2] {
		f, err := os.Create(filepath.Join(tempDir, name))
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		f.Close()
	}
	if err := os.Mkdir(filepath.Join(tempDir, "subdir"), 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	// Change to temp directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}
	defer func() {
		// Try to restore, but don't fail if the directory was already cleaned up
		_ = os.Chdir(origDir)
	}()

	// Create a directory widget
	dir := filer.NewDirectory(0, 0, 80, 24)
	dir.Chdir(tempDir)

	t.Run("find existing file", func(t *testing.T) {
		file := dir.FindFileByName("file1.txt")
		if file == nil {
			t.Error("expected to find file1.txt, got nil")
		} else if file.Name() != "file1.txt" {
			t.Errorf("expected name 'file1.txt', got '%s'", file.Name())
		}
	})

	t.Run("find existing directory", func(t *testing.T) {
		file := dir.FindFileByName("subdir")
		if file == nil {
			t.Error("expected to find subdir, got nil")
		} else if file.Name() != "subdir" {
			t.Errorf("expected name 'subdir', got '%s'", file.Name())
		}
	})

	t.Run("not found returns nil", func(t *testing.T) {
		file := dir.FindFileByName("nonexistent.txt")
		if file != nil {
			t.Errorf("expected nil for nonexistent file, got %v", file)
		}
	})

	t.Run("exact match only", func(t *testing.T) {
		// Should not match partial names
		file := dir.FindFileByName("file1")
		if file != nil {
			t.Errorf("expected nil for partial match, got %v", file)
		}
	})
}

// TestSyncResult_REQ_SYNC_COMMANDS tests the result structure.
func TestSyncResult_REQ_SYNC_COMMANDS(t *testing.T) {
	t.Run("empty result", func(t *testing.T) {
		result := SyncResult{}
		if result.Succeeded != 0 {
			t.Errorf("expected 0 succeeded, got %d", result.Succeeded)
		}
		if result.Skipped != 0 {
			t.Errorf("expected 0 skipped, got %d", result.Skipped)
		}
		if len(result.Failures) != 0 {
			t.Errorf("expected 0 failures, got %d", len(result.Failures))
		}
	})

	t.Run("record failures", func(t *testing.T) {
		result := SyncResult{
			Succeeded: 2,
			Skipped:   1,
			Failures:  []SyncFailure{{PaneIndex: 0, Error: os.ErrNotExist}},
		}
		if result.Succeeded != 2 {
			t.Errorf("expected 2 succeeded, got %d", result.Succeeded)
		}
		if result.Skipped != 1 {
			t.Errorf("expected 1 skipped, got %d", result.Skipped)
		}
		if len(result.Failures) != 1 {
			t.Errorf("expected 1 failure, got %d", len(result.Failures))
		}
		if result.Failures[0].PaneIndex != 0 {
			t.Errorf("expected failure in pane 0, got %d", result.Failures[0].PaneIndex)
		}
	})
}

// TestCopyFileSimple_REQ_SYNC_COMMANDS tests the simple file copy helper.
func TestCopyFileSimple_REQ_SYNC_COMMANDS(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file with content
	srcPath := filepath.Join(tempDir, "source.txt")
	content := []byte("test content for copy")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	// Copy to destination
	dstPath := filepath.Join(tempDir, "destination.txt")
	if err := copyFileSimple(srcPath, dstPath); err != nil {
		t.Fatalf("copyFileSimple failed: %v", err)
	}

	// Verify destination exists and has correct content
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}
	if string(dstContent) != string(content) {
		t.Errorf("expected content %q, got %q", string(content), string(dstContent))
	}
}

// TestCopyDirRecursive_REQ_SYNC_COMMANDS tests the recursive directory copy helper.
func TestCopyDirRecursive_REQ_SYNC_COMMANDS(t *testing.T) {
	tempDir := t.TempDir()

	// Create source directory structure
	srcDir := filepath.Join(tempDir, "source")
	if err := os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755); err != nil {
		t.Fatalf("failed to create source directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("file1"), 0644); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "subdir", "file2.txt"), []byte("file2"), 0644); err != nil {
		t.Fatalf("failed to create file2: %v", err)
	}

	// Copy to destination
	dstDir := filepath.Join(tempDir, "destination")
	if err := copyDirRecursive(srcDir, dstDir); err != nil {
		t.Fatalf("copyDirRecursive failed: %v", err)
	}

	// Verify structure
	if _, err := os.Stat(filepath.Join(dstDir, "file1.txt")); err != nil {
		t.Errorf("expected file1.txt in destination: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dstDir, "subdir", "file2.txt")); err != nil {
		t.Errorf("expected subdir/file2.txt in destination: %v", err)
	}
}

// TestSyncMode_REQ_SYNC_COMMANDS tests the mode string methods.
func TestSyncMode_REQ_SYNC_COMMANDS(t *testing.T) {
	// Test mode struct directly (without full Goful context)
	t.Run("mode string", func(t *testing.T) {
		mode := &syncMode{ignoreFailures: false}
		if mode.String() != "sync" {
			t.Errorf("expected mode string 'sync', got '%s'", mode.String())
		}
	})

	t.Run("prompt normal mode", func(t *testing.T) {
		mode := &syncMode{ignoreFailures: false}
		prompt := mode.Prompt()
		if prompt == "" {
			t.Error("expected non-empty prompt")
		}
		// Prompt should contain operation keys
		if !containsSubstring(prompt, "c") || !containsSubstring(prompt, "d") || !containsSubstring(prompt, "r") {
			t.Errorf("prompt should contain c/d/r keys: %s", prompt)
		}
	})

	t.Run("prompt ignore mode", func(t *testing.T) {
		mode := &syncMode{ignoreFailures: true}
		prompt := mode.Prompt()
		// Should indicate ignore mode
		if !containsSubstring(prompt, "ignore") && !containsSubstring(prompt, "!") {
			t.Errorf("prompt should indicate ignore mode: %s", prompt)
		}
	})
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s[1:], substr) || (len(s) >= len(substr) && s[:len(substr)] == substr))
}
