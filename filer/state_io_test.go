package filer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateStateFileCreatesParents_REQ_CONFIGURABLE_STATE_PATHS(t *testing.T) {
	// [REQ:CONFIGURABLE_STATE_PATHS] [ARCH:STATE_PATH_SELECTION] [IMPL:STATE_PATH_RESOLVER]
	tmp := t.TempDir()
	target := filepath.Join(tmp, "nested", "dir", "state.json")

	file, err := createStateFile(target)
	if err != nil {
		t.Fatalf("expected directory creation to succeed, got error %v", err)
	}
	file.Close()

	if _, err := os.Stat(filepath.Dir(target)); err != nil {
		t.Fatalf("parent directories should exist, stat error %v", err)
	}
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("state file should exist at %s, got %v", target, err)
	}
}
