package externalcmd

import (
	"os"
	"strings"
	"testing"
)

func TestLoadFallsBackToDefaultsWhenFileMissing_REQ_EXTERNAL_COMMAND_CONFIG(t *testing.T) {
	// [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_LOADER]
	entries, err := Load(Options{
		Path: "/does/not/exist.json",
		GOOS: "linux",
		ReadFile: func(string) ([]byte, error) {
			return nil, os.ErrNotExist
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) == 0 {
		t.Fatalf("expected defaults to be returned")
	}
}

func TestLoadParsesCustomConfig_REQ_EXTERNAL_COMMAND_CONFIG(t *testing.T) {
	// [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_LOADER]
	const raw = `[{"key":"z","label":"zip %m","command":"zip -roD %x.zip %m"}]`
	var calledPath string
	defaultsLen := len(Defaults("linux"))

	entries, err := Load(Options{
		Path: "~/commands.json",
		GOOS: "linux",
		ReadFile: func(path string) ([]byte, error) {
			calledPath = path
			return []byte(raw), nil
		},
	})
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if strings.Contains(calledPath, "~") {
		t.Fatalf("expected tilde in path to expand, got %q", calledPath)
	}
	if len(entries) != defaultsLen+1 {
		t.Fatalf("expected defaults (%d) plus custom entries, got %d", defaultsLen, len(entries))
	}
	custom := entries[0]
	if custom.Key != "z" || custom.Command != "zip -roD %x.zip %m" {
		t.Fatalf("unexpected custom entry at head: %+v", custom)
	}
}

func TestLoadRejectsDuplicates_REQ_EXTERNAL_COMMAND_CONFIG(t *testing.T) {
	// [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_LOADER]
	const raw = `[{"key":"c","label":"copy","command":"cp %m %D"}, {"key":"c","label":"dup","command":"echo dup"}]`
	_, err := Load(Options{
		Path: "/ignore.json",
		GOOS: "linux",
		ReadFile: func(string) ([]byte, error) {
			return []byte(raw), nil
		},
	})
	if err == nil {
		t.Fatalf("expected duplicate error")
	}
	if !strings.Contains(err.Error(), "duplicate shortcut") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadSkipsDisabledAndPlatformMismatchedEntries_REQ_EXTERNAL_COMMAND_CONFIG(t *testing.T) {
	// [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_LOADER]
	const raw = `
{
	"inheritDefaults": false,
	"commands": [
		{"key":"l","label":"linux only","command":"echo linux","platforms":["linux"]},
		{"key":"w","label":"windows only","command":"echo win","platforms":["windows"]},
		{"key":"d","label":"disabled","command":"echo nope","disabled":true},
		{"key":"m","label":"menu","runMenu":"archive"}
	]
}`
	entries, err := Load(Options{
		Path: "/ignore.json",
		GOOS: "linux",
		ReadFile: func(string) ([]byte, error) {
			return []byte(raw), nil
		},
		Debug: true,
		Logf:  func(string, ...interface{}) {},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries after filtering, got %d (%+v)", len(entries), entries)
	}
	if entries[0].Key != "l" || entries[1].RunMenu != "archive" {
		t.Fatalf("unexpected entries: %+v", entries)
	}
}

func TestLoadSupportsWrapperObject_REQ_EXTERNAL_COMMAND_CONFIG(t *testing.T) {
	// [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_LOADER]
	const raw = `{"commands":[{"key":"x","label":"expr","command":"expr 1 + 1"}]}`
	defaultsLen := len(Defaults("linux"))
	entries, err := Load(Options{
		Path: "/ignore.json",
		GOOS: "linux",
		ReadFile: func(string) ([]byte, error) {
			return []byte(raw), nil
		},
	})
	if err != nil {
		t.Fatalf("wrapper parse failed: %v", err)
	}
	if len(entries) != defaultsLen+1 {
		t.Fatalf("expected defaults (%d) plus wrapper entries, got %d", defaultsLen, len(entries))
	}
	if entries[0].Key != "x" {
		t.Fatalf("unexpected entries (custom should be first): %+v", entries)
	}
}

func TestLoadParsesYAMLArray_REQ_EXTERNAL_COMMAND_CONFIG(t *testing.T) {
	// [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_LOADER]
	const raw = `
- key: y
  label: "yaml copy"
  command: "cp %m %D2"
 `
	defaultsLen := len(Defaults("linux"))
	entries, err := Load(Options{
		Path: "/ignore.yaml",
		GOOS: "linux",
		ReadFile: func(string) ([]byte, error) {
			return []byte(raw), nil
		},
	})
	if err != nil {
		t.Fatalf("yaml array parse failed: %v", err)
	}
	if len(entries) != defaultsLen+1 {
		t.Fatalf("expected defaults (%d) plus yaml entry, got %d", defaultsLen, len(entries))
	}
	if entries[0].Key != "y" {
		t.Fatalf("unexpected entries (yaml entry should be first): %+v", entries)
	}
}

func TestLoadParsesYAMLWrapper_REQ_EXTERNAL_COMMAND_CONFIG(t *testing.T) {
	// [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_LOADER]
	const raw = `
inheritDefaults: false
commands:
  - key: a
    label: "archives menu"
    runMenu: "archive"
 `
	entries, err := Load(Options{
		Path: "/ignore.yaml",
		GOOS: "linux",
		ReadFile: func(string) ([]byte, error) {
			return []byte(raw), nil
		},
	})
	if err != nil {
		t.Fatalf("yaml wrapper parse failed: %v", err)
	}
	if len(entries) != 1 || entries[0].RunMenu != "archive" {
		t.Fatalf("unexpected entries: %+v", entries)
	}
}

func TestLoadPrependsDefaultsByDefault_REQ_EXTERNAL_COMMAND_CONFIG(t *testing.T) {
	// [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_APPEND]
	const raw = `[{"key":"o","label":"open","command":"open %f"}]`
	defaultsLen := len(Defaults("windows"))
	entries, err := Load(Options{
		Path: "/windows.json",
		GOOS: "windows",
		ReadFile: func(string) ([]byte, error) {
			return []byte(raw), nil
		},
	})
	if err != nil {
		t.Fatalf("prepend defaults failed: %v", err)
	}
	if len(entries) != defaultsLen+1 {
		t.Fatalf("expected defaults (%d) plus custom entry, got %d", defaultsLen, len(entries))
	}
	if entries[0].Key != "o" {
		t.Fatalf("prepend order incorrect: %+v", entries)
	}
}

func TestLoadCanDisableDefaults_REQ_EXTERNAL_COMMAND_CONFIG(t *testing.T) {
	// [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_APPEND]
	const raw = `{"inheritDefaults":false,"commands":[{"key":"n","label":"noop","command":"echo noop"}]}`
	entries, err := Load(Options{
		Path: "/replace.json",
		GOOS: "linux",
		ReadFile: func(string) ([]byte, error) {
			return []byte(raw), nil
		},
	})
	if err != nil {
		t.Fatalf("disable defaults failed: %v", err)
	}
	if len(entries) != 1 || entries[0].Key != "n" {
		t.Fatalf("expected only custom entry, got %+v", entries)
	}
	for _, entry := range entries {
		if entry.Key == "c" {
			t.Fatalf("found compiled default when inheritDefaults=false: %+v", entries)
		}
	}
}
