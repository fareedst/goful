package terminalcmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/shlex"
)

const (
	// EnvTerminalCommand selects a custom terminal executable (e.g., "alacritty -e").
	EnvTerminalCommand = "GOFUL_TERMINAL_CMD"
	// EnvDebugTerminal enables debug logs for branch selection.
	EnvDebugTerminal = "GOFUL_DEBUG_TERMINAL"
	// KeepOpenTail preserves the historical pause behaviour after running a command.
	KeepOpenTail = `;read -p "HIT ENTER KEY"`
)

// Options drive how the terminal adapter behaves.
type Options struct {
	GOOS     string
	IsTmux   bool
	Override []string
	Tail     string
	Debug    bool
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
	return Factory{opts: opts}
}

// Command returns the command/args slice to launch the requested shell command.
// [IMPL:TERMINAL_ADAPTER] [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY]
func (f Factory) Command(cmd string) []string {
	payload := cmd + f.opts.Tail

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
		f.log("macOS Terminal (osascript) branch selected")
		return buildAppleScriptCommand(payload)
	}

	f.log("default gnome-terminal branch selected")
	title := linuxTitle(cmd)
	return []string{"gnome-terminal", "--", "bash", "-c", title + payload}
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

func buildAppleScriptCommand(payload string) []string {
	shellLine := fmt.Sprintf("bash -lc %s", strconv.Quote(payload))
	escaped := appleScriptEscape(shellLine)
	run := fmt.Sprintf("tell application \"Terminal\" to do script \"%s\"", escaped)
	return []string{
		"osascript",
		"-e", run,
		"-e", "tell application \"Terminal\" to activate",
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
func Apply(cfg Configurator, factory Factory) {
	cfg.ConfigTerminal(func(cmd string) []string {
		return factory.Command(cmd)
	})
}
