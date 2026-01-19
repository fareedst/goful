package app

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/fareedst/goful/filer"
	"github.com/fareedst/goful/message"
	"github.com/fareedst/goful/util"
	"github.com/fareedst/goful/widget"
)

// Spawn a process by the shell or the terminal.
func (g *Goful) Spawn(cmd string) {
	cmd, background := g.expandMacro(cmd)
	var args []string
	if background {
		args = g.shell(cmd)
	} else {
		args = g.terminal(cmd)
	}
	execCmd := exec.Command(args[0], args[1:]...)
	message.Info(strings.Join(execCmd.Args, " "))
	if err := spawn(execCmd); err != nil {
		message.Error(err)
	}
}

func spawn(cmd *exec.Cmd) error {
	var bufout bytes.Buffer
	cmd.Stdout = &bufout
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		var errWait = cmd.Wait()
		var stderr = &exec.ExitError{}
		switch {
		case errors.As(errWait, &stderr):
			message.Errorf("%q: %s", cmd, stderr.Stderr)
			return
		case errWait != nil:
			message.Errorf("%q: %v", cmd, errWait)
			return
		}
		if bufout.Len() > 0 {
			message.Info(bufout.String())
		}
	}()
	return nil
}

// SpawnSuspend spawns a process and suspends screen.
func (g *Goful) SpawnSuspend(cmd string) {
	cmd, _ = g.expandMacro(cmd)
	args := g.shell(cmd)
	execCmd := exec.Command(args[0], args[1:]...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	widget.Fini()
	defer func(cmd string) {
		widget.Init()
		message.Info(cmd)
	}(strings.Join(execCmd.Args, " "))
	_ = execCmd.Run()

	shell := exec.Command(args[0])
	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr
	_ = shell.Run()
}

const (
	macroPrefix             = '%'
	macroEscape             = '\\' // \ is an escape sequence
	macroNonQuote           = '~'  // %~ is expanded non quote
	macroFile               = 'f'  // %f %~f are expanded a file name on the cursor
	macroFilePath           = 'F'  // %F %~F are expanded a file path on the cursor
	macroFileWithoutExt     = 'x'  // %x %~x are expanded a file name excluded the extension on the cursor
	macroFileWithoutExtPath = 'X'  // %x %~X are expanded a file path excluded the extension on the cursor
	macroMarkfile           = 'm'  // %m %~m are expanded mark file names joined by spaces
	macroMarkfilePath       = 'M'  // %M %~M are expanded mark file paths joined by spaces
	macroDir                = 'd'  // %d %~d are expanded a directory name on the cursor
	macroDirPath            = 'D'  // %D %~D are expanded a directory path on the cursor
	macroNextDir            = '2'  // %d2 %D2 %~d2 %~D2 are expanded the neighbor directory name or path
	macroAllOtherDirs       = '@'  // %D@ %~D@ %d@ %~d@ expand to the remaining window directories in display order
	macroRunBackground      = '&'  // %& is a flag runned in background
)

func (g *Goful) expandMacro(cmd string) (result string, background bool) {
	data := []byte(cmd)
	ret := make([]byte, len(data))
	copy(ret, data)

	background = false
	escape := false
	prefix := false
	nonQuote := false
	offset := 0
	for i, b := range data {
		if escape { // skip the escape sequence
			ret = widget.DeleteBytes(ret, offset-1, 1)
			escape = false
			continue
		}

		if prefix {
			if b == macroNonQuote {
				if nonQuote { // continuous ~ is not expand
					prefix = false
					nonQuote = false
					offset++
				} else {
					nonQuote = true
				}
				continue
			}
			prefix = false
			src := ""
			macrolen := 2
			if nonQuote {
				macrolen++
			}
			switch b {
			case macroFile:
				src = g.File().Name()
				if !nonQuote {
					src = util.Quote(src)
				}
			case macroFilePath:
				src = g.File().Path()
				if !nonQuote {
					src = util.Quote(src)
				}
			case macroFileWithoutExt:
				src = util.RemoveExt(g.File().Name())
				if !nonQuote {
					src = util.Quote(src)
				}
			case macroFileWithoutExtPath:
				src = util.RemoveExt(g.File().Path())
				if !nonQuote {
					src = util.Quote(src)
				}
			case macroMarkfile:
				if !nonQuote {
					src = strings.Join(g.Dir().MarkfileQuotedNames(), " ")
				} else {
					src = strings.Join(g.Dir().MarkfileNames(), " ")
				}
			case macroMarkfilePath:
				if !nonQuote {
					src = strings.Join(g.Dir().MarkfileQuotedPaths(), " ")
				} else {
					src = strings.Join(g.Dir().MarkfilePaths(), " ")
				}
			case macroDir:
				src = g.Dir().Base()
				formattedList := false
				if i != len(data)-1 {
					switch data[i+1] {
					case macroNextDir:
						src = g.Workspace().NextDir().Base()
						macrolen++
					case macroAllOtherDirs:
						// `%d@` mirrors `%D@` but emits only the basename for each window.
						// [IMPL:WINDOW_MACRO_ENUMERATION] [ARCH:WINDOW_MACRO_ENUMERATION] [REQ:WINDOW_MACRO_ENUMERATION]
						src = formatDirListForMacro(
							otherWindowDirNames(g.Workspace()),
							!nonQuote,
						)
						macrolen++
						formattedList = true
					}
				}
				if !formattedList && !nonQuote {
					src = util.Quote(src)
				}
			case macroDirPath:
				src = g.Dir().Path
				formattedList := false
				if i != len(data)-1 {
					switch data[i+1] {
					case macroNextDir:
						src = g.Workspace().NextDir().Path
						macrolen++
					case macroAllOtherDirs:
						// `%D@` stays quoted for safety while `%~D@` explicitly opts into raw output.
						// [IMPL:WINDOW_MACRO_ENUMERATION] [ARCH:WINDOW_MACRO_ENUMERATION] [REQ:WINDOW_MACRO_ENUMERATION]
						src = formatDirListForMacro(
							otherWindowDirPaths(g.Workspace()),
							!nonQuote,
						)
						macrolen++
						formattedList = true
					}
				}
				if !formattedList && !nonQuote {
					src = util.Quote(src)
				}
			case macroRunBackground:
				background = true
			default:
				if nonQuote {
					nonQuote = false
					offset++
				}
				goto other
			}
			ret = widget.DeleteBytes(ret, offset-1, macrolen)
			ret = widget.InsertBytes(ret, []byte(src), offset-1)
			offset += len(src) - macrolen
			offset++
			if nonQuote {
				nonQuote = false
				offset++
			}
			continue
		}
	other:
		switch b {
		case macroPrefix:
			prefix = true
		case macroEscape:
			escape = true
		}
		offset++
	}
	return string(ret), background
}

// otherWindowDirPaths returns every non-focused directory path in deterministic order:
// start at the next window and wrap until all panes are listed.
// [IMPL:WINDOW_MACRO_ENUMERATION] [ARCH:WINDOW_MACRO_ENUMERATION] [REQ:WINDOW_MACRO_ENUMERATION]
func otherWindowDirPaths(ws *filer.Workspace) []string {
	if ws == nil || len(ws.Dirs) <= 1 {
		return nil
	}

	paths := make([]string, 0, len(ws.Dirs)-1)
	for offset := 1; offset < len(ws.Dirs); offset++ {
		idx := (ws.Focus + offset) % len(ws.Dirs)
		paths = append(paths, ws.Dirs[idx].Path)
	}
	return paths
}

// otherWindowDirNames returns every non-focused directory basename in deterministic order.
// [IMPL:WINDOW_MACRO_ENUMERATION] [ARCH:WINDOW_MACRO_ENUMERATION] [REQ:WINDOW_MACRO_ENUMERATION]
func otherWindowDirNames(ws *filer.Workspace) []string {
	if ws == nil || len(ws.Dirs) <= 1 {
		return nil
	}

	names := make([]string, 0, len(ws.Dirs)-1)
	for offset := 1; offset < len(ws.Dirs); offset++ {
		idx := (ws.Focus + offset) % len(ws.Dirs)
		names = append(names, ws.Dirs[idx].Base())
	}
	return names
}

// formatDirListForMacro joins directory entries with spaces, optionally quoting each value.
// [IMPL:WINDOW_MACRO_ENUMERATION] [ARCH:WINDOW_MACRO_ENUMERATION] [REQ:WINDOW_MACRO_ENUMERATION]
func formatDirListForMacro(paths []string, quote bool) string {
	if len(paths) == 0 {
		return ""
	}
	formatted := make([]string, len(paths))
	for i, path := range paths {
		if quote {
			formatted[i] = util.Quote(path)
		} else {
			formatted[i] = path
		}
	}
	return strings.Join(formatted, " ")
}
