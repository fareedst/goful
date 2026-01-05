package terminalcmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/anmitsu/goful/util"
	"github.com/google/shlex"
)

const (
	// EnvTerminalCommand selects a custom terminal executable (e.g., "alacritty -e").
	EnvTerminalCommand = "GOFUL_TERMINAL_CMD"
	// EnvDebugTerminal enables debug logs for branch selection.
	EnvDebugTerminal = "GOFUL_DEBUG_TERMINAL"
	// EnvTerminalApp chooses which macOS application receives the AppleScript (default Terminal).
	EnvTerminalApp = "GOFUL_TERMINAL_APP"
	// EnvTerminalShell chooses which shell binary runs inside the macOS window (default bash).
	EnvTerminalShell = "GOFUL_TERMINAL_SHELL"
	// KeepOpenTail preserves the historical pause behaviour after running a command.
	KeepOpenTail = `;read -p "HIT ENTER KEY"`
)

const (
	defaultTerminalApp   = "Terminal"
	defaultTerminalShell = "bash"
)

// Options drive how the terminal adapter behaves.
type Options struct {
	GOOS          string
	IsTmux        bool
	Override      []string
	Tail          string
	Debug         bool
	TerminalApp   string // macOS-only application name for AppleScript
	TerminalShell string // macOS-only shell binary inserted into AppleScript
}

// Factory builds terminal command invocations per platform.
// [IMPL:TERMINAL_ADAPTER] [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY]
type Factory struct {
	opts Options
}

// NewFactory returns a factory with sane defaults.
func NewFactory(opts Options) Factory {
	if opts.Tail == "" {
		opts.Tail = KeepOpenTail
	}
	if opts.TerminalApp == "" {
		opts.TerminalApp = defaultTerminalApp
	}
	if opts.TerminalShell == "" {
		opts.TerminalShell = defaultTerminalShell
	}
	return Factory{opts: opts}
}

// Command returns the command/args slice to launch the requested shell command.
// [IMPL:TERMINAL_ADAPTER] [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY]
func (f Factory) Command(cmd string) []string {
	return f.CommandWithCwd(cmd, "")
}

// CommandWithCwd allows callers to supply the focused directory, ensuring macOS sessions cd first.
// [IMPL:TERMINAL_ADAPTER] [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY] [REQ:TERMINAL_CWD]
func (f Factory) CommandWithCwd(cmd string, cwd string) []string {
	payload := f.buildPayload(cmd, cwd)

	switch {
	case len(f.opts.Override) > 0:
		f.log(fmt.Sprintf("override via %v", f.opts.Override))
		args := append([]string{}, f.opts.Override...)
		args = append(args, "bash", "-c", payload)
		return args
	case f.opts.IsTmux:
		f.log("tmux branch selected")
		return []string{"tmux", "new-window", "-n", cmd, payload}
	case strings.EqualFold(f.opts.GOOS, "darwin"):
		f.log(fmt.Sprintf("macOS AppleScript branch selected (app=%s shell=%s)", f.opts.TerminalApp, f.opts.TerminalShell))
		return f.buildAppleScriptCommand(payload)
	}

	f.log("default gnome-terminal branch selected")
	title := linuxTitle(cmd)
	return []string{"gnome-terminal", "--", "bash", "-c", title + payload}
}

func (f Factory) buildPayload(cmd string, cwd string) string {
	payload := cmd + f.opts.Tail
	if cwd != "" && strings.EqualFold(f.opts.GOOS, "darwin") {
		payload = fmt.Sprintf("cd %s; %s", util.Quote(cwd), payload)
	}
	return payload
}

// ParseOverride converts an override string (e.g., EnvTerminalCommand) into args.
func ParseOverride(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	args, err := shlex.Split(raw)
	if err != nil {
		return nil, err
	}
	return args, nil
}

func (f Factory) log(msg string) {
	if !f.opts.Debug {
		return
	}
	_, _ = fmt.Fprintf(os.Stderr, "DEBUG: [IMPL:TERMINAL_ADAPTER] %s\n", msg)
}

func linuxTitle(cmd string) string {
	return "echo -n '\\033]0;" + cmd + "\\007';"
}

func (f Factory) buildAppleScriptCommand(payload string) []string {
	shell := f.opts.TerminalShell
	shellLine := fmt.Sprintf("%s -c %s; exit", shell, strconv.Quote(payload))
	escaped := appleScriptEscape(shellLine)
	app := f.opts.TerminalApp
	run := fmt.Sprintf("tell application \"%s\" to do script \"%s\"", app, escaped)
	activate := fmt.Sprintf("tell application \"%s\" to activate", app)
	return []string{
		"osascript",
		"-e", run,
		"-e", activate,
	}
}

func appleScriptEscape(s string) string {
	replacer := strings.NewReplacer(
		`"`, `\"`,
		`\`, `\\`,
		"\n", `\n`,
	)
	return replacer.Replace(s)
}

// Configurator matches app.Goful's ConfigTerminal method.
type Configurator interface {
	ConfigTerminal(func(cmd string) []string)
}

// Apply wires the factory into the provided configurator.
// [IMPL:TERMINAL_ADAPTER] [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY]
func Apply(cfg Configurator, factory Factory, cwd func() string) {
	cfg.ConfigTerminal(func(cmd string) []string {
		return factory.CommandWithCwd(cmd, cwd())
	})
}
