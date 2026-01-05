package terminalcmd

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestCommandFactoryTmux_REQ_TERMINAL_PORTABILITY(t *testing.T) {
	// [REQ:TERMINAL_PORTABILITY] tmux selection keeps legacy behaviour.
	factory := NewFactory(Options{
		GOOS:   "linux",
		IsTmux: true,
	})
	got := factory.Command("echo hi")
	want := []string{"tmux", "new-window", "-n", "echo hi", "echo hi" + KeepOpenTail}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tmux command mismatch\nwant: %v\ngot:  %v", want, got)
	}
}

func TestCommandFactoryOverride_REQ_TERMINAL_PORTABILITY(t *testing.T) {
	factory := NewFactory(Options{
		GOOS:     "linux",
		Override: []string{"alacritty", "-e"},
	})
	got := factory.Command("echo hi")
	want := []string{"alacritty", "-e", "bash", "-c", "echo hi" + KeepOpenTail}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("override command mismatch\nwant: %v\ngot:  %v", want, got)
	}
}

func TestCommandFactoryDarwin_REQ_TERMINAL_PORTABILITY(t *testing.T) {
	factory := NewFactory(Options{
		GOOS: "darwin",
	})
	got := factory.CommandWithCwd("echo hi", "/tmp/demo")
	assertAppleScriptCommand(t, got, "Terminal", "bash", "/tmp/demo")
}

func TestCommandFactoryDarwinCustomApp_REQ_TERMINAL_PORTABILITY(t *testing.T) {
	factory := NewFactory(Options{
		GOOS:        "darwin",
		TerminalApp: "iTerm2",
	})
	got := factory.CommandWithCwd("echo hi", "/tmp/demo")
	assertAppleScriptCommand(t, got, "iTerm2", "bash", "/tmp/demo")
}

func TestCommandFactoryDarwinCustomShell_REQ_TERMINAL_PORTABILITY(t *testing.T) {
	factory := NewFactory(Options{
		GOOS:          "darwin",
		TerminalShell: "zsh",
	})
	got := factory.CommandWithCwd("echo hi", "/tmp/demo")
	assertAppleScriptCommand(t, got, "Terminal", "zsh", "/tmp/demo")
}

func TestCommandFactoryOverrideDarwin_REQ_TERMINAL_CWD(t *testing.T) {
	// [REQ:TERMINAL_CWD] overrides still inherit the auto-cd preamble.
	factory := NewFactory(Options{
		GOOS:     "darwin",
		Override: []string{"iTerm2", "-e"},
	})
	got := factory.CommandWithCwd("echo hi", "/tmp/demo")
	want := []string{"iTerm2", "-e", "bash", "-c", `cd "/tmp/demo"; echo hi` + KeepOpenTail}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("darwin override mismatch\nwant: %v\ngot:  %v", want, got)
	}
}

func TestCommandFactoryLinux_REQ_TERMINAL_PORTABILITY(t *testing.T) {
	factory := NewFactory(Options{
		GOOS: "linux",
	})
	got := factory.Command("echo hi")
	want := []string{
		"gnome-terminal",
		"--",
		"bash",
		"-c",
		"echo -n '\\033]0;echo hi\\007';echo hi" + KeepOpenTail,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("linux command mismatch\nwant: %v\ngot:  %v", want, got)
	}
}

func TestParseOverride_REQ_TERMINAL_PORTABILITY(t *testing.T) {
	// [REQ:TERMINAL_PORTABILITY] override accepts quoted strings.
	raw := `alacritty --class "Goful Terminal" -e`
	args, err := ParseOverride(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"alacritty", "--class", "Goful Terminal", "-e"}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("parsed mismatch\nwant: %v\ngot:  %v", want, args)
	}
	_, err = ParseOverride(`"`)
	if err == nil {
		t.Fatalf("expected error for malformed override")
	}
}

type stubConfigurator struct {
	last []string
}

func (s *stubConfigurator) ConfigTerminal(fn func(string) []string) {
	s.last = fn("echo hi")
}

func TestApply_REQ_TERMINAL_PORTABILITY(t *testing.T) {
	cfg := &stubConfigurator{}
	factory := NewFactory(Options{GOOS: "linux"})
	Apply(cfg, factory, func() string { return "/tmp/demo" })
	want := []string{"gnome-terminal", "--", "bash", "-c", "echo -n '\\033]0;echo hi\\007';echo hi" + KeepOpenTail}
	if !reflect.DeepEqual(cfg.last, want) {
		t.Fatalf("apply produced unexpected command\nwant: %v\ngot:  %v", want, cfg.last)
	}
}

func TestApplyDarwinCwd_REQ_TERMINAL_CWD(t *testing.T) {
	cfg := &stubConfigurator{}
	factory := NewFactory(Options{GOOS: "darwin"})
	Apply(cfg, factory, func() string { return "/tmp/demo" })
	assertAppleScriptCommand(t, cfg.last, "Terminal", "bash", "/tmp/demo")
}

func assertAppleScriptCommand(t *testing.T, got []string, app, shell, dir string) {
	t.Helper()
	if len(got) != 5 {
		t.Fatalf("unexpected osascript command shape: %v", got)
	}
	if got[0] != "osascript" {
		t.Fatalf("expected osascript binary, got %v", got)
	}
	if got[1] != "-e" || got[3] != "-e" {
		t.Fatalf("unexpected osascript flags: %v", got)
	}
	activate := fmt.Sprintf("tell application \"%s\" to activate", app)
	if got[4] != activate {
		t.Fatalf("activate clause mismatch\nwant: %s\ngot:  %s", activate, got[4])
	}
	script := got[2]
	prefix := fmt.Sprintf("tell application \"%s\" to do script", app)
	if !strings.Contains(script, prefix) {
		t.Fatalf("missing app script prefix in %q", script)
	}
	if !strings.Contains(script, fmt.Sprintf("%s -c", shell)) {
		t.Fatalf("missing shell %q in %q", shell, script)
	}
	if !strings.Contains(script, "cd ") || !strings.Contains(script, dir) {
		t.Fatalf("missing cd clause with dir %q in %q", dir, script)
	}
	if !strings.Contains(script, "echo hi") {
		t.Fatalf("missing payload in %q", script)
	}
	if !strings.Contains(script, "read -p") || !strings.Contains(script, "HIT ENTER KEY") {
		t.Fatalf("missing pause tail in %q", script)
	}
}
