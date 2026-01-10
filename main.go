package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/anmitsu/goful/app"
	"github.com/anmitsu/goful/cmdline"
	"github.com/anmitsu/goful/configpaths"
	"github.com/anmitsu/goful/diffstatus"
	"github.com/anmitsu/goful/externalcmd"
	"github.com/anmitsu/goful/filer"
	"github.com/anmitsu/goful/filer/comparecolors"
	"github.com/anmitsu/goful/internal/externalmenu"
	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/menu"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/terminalcmd"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
)

const debugWorkspaceEnv = "GOFUL_DEBUG_WORKSPACE"

var (
	stateFlag = flag.String(
		"state",
		"",
		"Override path to state.json (default "+configpaths.DefaultStatePath+" or "+configpaths.EnvStateKey+")",
	)
	historyFlag = flag.String(
		"history",
		"",
		"Override path to cmdline history (default "+configpaths.DefaultHistoryPath+" or "+configpaths.EnvHistoryKey+")",
	)
	commandsFlag = flag.String(
		"commands",
		"",
		"Override path to external-command config (default "+configpaths.DefaultCommandsPath+" or "+configpaths.EnvCommandsKey+")",
	)
	excludeNamesFlag = flag.String(
		"exclude-names",
		"",
		"Override path to filename exclude list (default "+configpaths.DefaultExcludesPath+" or "+configpaths.EnvExcludesKey+")",
	)
	// [IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
	compareColorsFlag = flag.String(
		"compare-colors",
		"",
		"Override path to comparison colors config (default "+configpaths.DefaultCompareColorsPath+" or "+configpaths.EnvCompareColorsKey+")",
	)
)

func main() {
	flag.Parse()

	pathsResolver := configpaths.Resolver{}
	runtimePaths := pathsResolver.Resolve(*stateFlag, *historyFlag, *commandsFlag, *excludeNamesFlag, *compareColorsFlag)
	emitPathDebug(runtimePaths)
	loadExcludedNames(runtimePaths.Excludes)
	// [IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
	loadCompareColors(runtimePaths.CompareColors)

	is_tmux := false
	widget.Init()
	defer widget.Fini()

	if runtime.GOOS == "darwin" {
		is_tmux = strings.Contains(os.Getenv("TERM_PROGRAM"), "tmux")
	} else {
		is_tmux = strings.Contains(os.Getenv("TERM"), "screen")
	}
	// Change a terminal title.
	if is_tmux {
		os.Stdout.WriteString("\033kgoful\033") // for tmux
	} else {
		os.Stdout.WriteString("\033]0;goful\007") // for otherwise
	}

	goful := app.NewGoful(runtimePaths.State)
	config(goful, is_tmux, runtimePaths)
	// TODO(goful-maintainers) [IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]:
	// plumb LoadHistory errors (distinguish first-run missing files vs actual IO failures) so we can alert users instead of
	// silently discarding history and ingest errors.
	_ = cmdline.LoadHistory(runtimePaths.History)

	startupDirs, startupWarnings := app.ParseStartupDirs(flag.Args())
	for _, warn := range startupWarnings {
		message.Errorf("[REQ:WORKSPACE_START_DIRS] %s", warn)
	}
	// [IMPL:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]
	app.SeedStartupWorkspaces(goful, startupDirs, os.Getenv(debugWorkspaceEnv) != "")

	goful.Run()

	_ = goful.SaveState(runtimePaths.State)
	// TODO(goful-maintainers) [IMPL:DEBT_TRACKING] [ARCH:DEBT_MANAGEMENT] [REQ:DEBT_TRIAGE]:
	// surface SaveHistory failures (permissions/full disk) via message logging instead of quietly ignoring write errors.
	_ = cmdline.SaveHistory(runtimePaths.History)
}

func config(g *app.Goful, is_tmux bool, paths configpaths.Paths) {
	look.Set("default") // default, midnight, black, white

	// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
	// Wire linked navigation indicator to filer header
	filer.SetLinkedNavIndicator(g.IsLinkedNav)

	// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
	// Wire diff search status indicator to filer header
	filer.SetDiffSearchStatusFn(g.DiffSearchStatus)

	// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
	// Wire diff search status to dedicated status line
	diffstatus.SetStatusFn(g.DiffSearchStatus)
	diffstatus.SetActiveFn(g.IsDiffSearchActive)

	if runewidth.EastAsianWidth {
		// Because layout collapsing for ambiguous runes if LANG=ja_JP.
		widget.SetBorder('|', '-', '+', '+', '+', '+')
	} else {
		// Look good if environment variable RUNEWIDTH_EASTASIAN=0 and
		// ambiguous char setting is half-width for gnome-terminal.
		widget.SetBorder('│', '─', '┌', '┐', '└', '┘') // 0x2502, 0x2500, 0x250c, 0x2510, 0x2514, 0x2518
	}
	g.SetBorderStyle(widget.AllBorder) // AllBorder, ULBorder, NoBorder

	message.SetInfoLog("~/.goful/log/info.log")   // "" is not logging
	message.SetErrorLog("~/.goful/log/error.log") // "" is not logging
	message.Sec(5)                                // display second for a message

	// Setup widget keymaps.
	g.ConfigFiler(filerKeymap)
	filer.ConfigFinder(finderKeymap)
	cmdline.Config(cmdlineKeymap)
	cmdline.ConfigCompletion(completionKeymap)
	menu.Config(menuKeymap)

	filer.SetStatView(true, false, true)  // size, permission and time
	filer.SetTimeFormat("06-01-02 15:04") // ex: "Jan _2 15:04"

	toggleExcludedNames := func() {
		enabled, hasRules, count := filer.ToggleExcludedNames()
		if !hasRules {
			message.Infof("[REQ:FILER_EXCLUDE_NAMES] exclude list inactive (path %s)", paths.Excludes)
			return
		}
		state := "disabled"
		if enabled {
			state = "enabled"
		}
		// [REQ:FILER_EXCLUDE_NAMES] [IMPL:FILER_EXCLUDE_LOADER] Runtime toggle feedback + reload.
		message.Infof("[REQ:FILER_EXCLUDE_NAMES] filename excludes %s (%d entries)", state, count)
		g.Workspace().ReloadAll()
	}

	// [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
	toggleComparisonColors := func() {
		enabled := look.ToggleComparisonEnabled()
		state := "disabled"
		if enabled {
			state = "enabled"
		}
		message.Infof("[REQ:FILE_COMPARISON_COLORS] comparison colors %s", state)
		// Rebuild comparison index if enabling
		if enabled {
			g.Workspace().RebuildComparisonIndex()
		}
		g.Workspace().ReloadAll()
	}

	// [IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
	calculateDigest := func() {
		filename := g.File().Name()
		if filename == ".." {
			message.Info("[REQ:FILE_COMPARISON_COLORS] cannot calculate digest for parent directory")
			return
		}
		count := g.Workspace().CalculateDigestForFile(filename)
		if count > 0 {
			message.Infof("[REQ:FILE_COMPARISON_COLORS] calculated digest for %d files named %q", count, filename)
		} else {
			message.Infof("[REQ:FILE_COMPARISON_COLORS] no matching files with equal size for %q", filename)
		}
	}

	// Setup open command for C-m (when the enter key is pressed)
	// The macro %f means expanded to a file name, for more see (spawn.go)
	opener := "xdg-open %f %&"
	switch runtime.GOOS {
	case "windows":
		opener = "explorer %~f %&"
	case "darwin":
		opener = "open %f %&"
	}
	g.MergeKeymap(widget.Keymap{
		"C-m": func() { g.Spawn(opener) },
		"o":   func() { g.Spawn(opener) },
	})

	// Setup pager by $PAGER
	pager := os.Getenv("PAGER")
	if pager == "" {
		if runtime.GOOS == "windows" {
			pager = "more"
		} else {
			pager = "less"
		}
	}
	if runtime.GOOS == "windows" {
		pager += " %~f"
	} else {
		pager += " %f"
	}
	g.AddKeymap("i", func() { g.Spawn(pager) })

	// Setup a shell and a terminal to execute external commands.
	// The shell is called when execute on background by the macro %&.
	// The terminal is called when the other.
	if runtime.GOOS == "windows" {
		g.ConfigShell(func(cmd string) []string {
			return []string{"cmd", "/c", cmd}
		})
		g.ConfigTerminal(func(cmd string) []string {
			return []string{"cmd", "/c", "start", "cmd", "/c", cmd + "& pause"}
		})
	} else {
		g.ConfigShell(func(cmd string) []string {
			return []string{"bash", "-c", cmd}
		})
		overrideArgs, overrideErr := terminalcmd.ParseOverride(os.Getenv(terminalcmd.EnvTerminalCommand))
		if overrideErr != nil {
			message.Errorf("[REQ:TERMINAL_PORTABILITY] failed to parse %s: %v", terminalcmd.EnvTerminalCommand, overrideErr)
		}
		macApp := os.Getenv(terminalcmd.EnvTerminalApp)
		macShell := os.Getenv(terminalcmd.EnvTerminalShell)
		factory := terminalcmd.NewFactory(terminalcmd.Options{
			GOOS:          runtime.GOOS,
			IsTmux:        is_tmux,
			Override:      overrideArgs,
			Tail:          terminalcmd.KeepOpenTail,
			Debug:         os.Getenv(terminalcmd.EnvDebugTerminal) != "",
			TerminalApp:   macApp,
			TerminalShell: macShell,
		})
		// [IMPL:TERMINAL_ADAPTER] [ARCH:TERMINAL_LAUNCHER] [REQ:TERMINAL_PORTABILITY] [REQ:TERMINAL_CWD]
		terminalcmd.Apply(g, factory, func() string {
			return g.Dir().Path
		})
	}

	// Setup menus and add to keymap.
	// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
	// Sort menu with linked mode support - when linked, sort applies to all windows
	menu.Add("sort",
		"n", "sort name          ", func() {
			if g.IsLinkedNav() {
				g.Workspace().SortAllBy(filer.SortName)
			} else {
				g.Dir().SortName()
			}
		},
		"N", "sort name decending", func() {
			if g.IsLinkedNav() {
				g.Workspace().SortAllBy(filer.SortNameRev)
			} else {
				g.Dir().SortNameDec()
			}
		},
		"s", "sort size          ", func() {
			if g.IsLinkedNav() {
				g.Workspace().SortAllBy(filer.SortSize)
			} else {
				g.Dir().SortSize()
			}
		},
		"S", "sort size decending", func() {
			if g.IsLinkedNav() {
				g.Workspace().SortAllBy(filer.SortSizeRev)
			} else {
				g.Dir().SortSizeDec()
			}
		},
		"t", "sort time          ", func() {
			if g.IsLinkedNav() {
				g.Workspace().SortAllBy(filer.SortMtime)
			} else {
				g.Dir().SortMtime()
			}
		},
		"T", "sort time decending", func() {
			if g.IsLinkedNav() {
				g.Workspace().SortAllBy(filer.SortMtimeRev)
			} else {
				g.Dir().SortMtimeDec()
			}
		},
		"e", "sort ext           ", func() {
			if g.IsLinkedNav() {
				g.Workspace().SortAllBy(filer.SortExt)
			} else {
				g.Dir().SortExt()
			}
		},
		"E", "sort ext decending ", func() {
			if g.IsLinkedNav() {
				g.Workspace().SortAllBy(filer.SortExtRev)
			} else {
				g.Dir().SortExtDec()
			}
		},
		".", "toggle priority    ", func() { filer.TogglePriority(); g.Workspace().ReloadAll() },
	)
	g.AddKeymap("s", func() { g.Menu("sort") })

	menu.Add("view",
		"s", "stat menu    ", func() { g.Menu("stat") },
		"l", "layout menu  ", func() { g.Menu("layout") },
		"L", "look menu    ", func() { g.Menu("look") },
		"n", "toggle filename excludes", func() { toggleExcludedNames() },
		".", "toggle show hidden files", func() { filer.ToggleShowHiddens(); g.Workspace().ReloadAll() },
		"c", "toggle comparison colors", func() { toggleComparisonColors() }, // [REQ:FILE_COMPARISON_COLORS]
		"=", "calculate file digest   ", func() { calculateDigest() }, // [REQ:FILE_COMPARISON_COLORS] [IMPL:DIGEST_COMPARISON]
		"[", "start diff search       ", func() { g.StartDiffSearch() }, // [REQ:DIFF_SEARCH] [IMPL:DIFF_SEARCH]
		"]", "continue diff search    ", func() { g.ContinueDiffSearch() }, // [REQ:DIFF_SEARCH] [IMPL:DIFF_SEARCH]
	)
	g.AddKeymap("v", func() { g.Menu("view") })
	g.AddKeymap("E", toggleExcludedNames)
	g.AddKeymap("C", toggleComparisonColors) // [REQ:FILE_COMPARISON_COLORS]
	g.AddKeymap("=", calculateDigest)        // [REQ:FILE_COMPARISON_COLORS] [IMPL:DIGEST_COMPARISON]

	menu.Add("layout",
		"t", "tile       ", func() { g.Workspace().LayoutTile() },
		"T", "tile-top   ", func() { g.Workspace().LayoutTileTop() },
		"b", "tile-bottom", func() { g.Workspace().LayoutTileBottom() },
		"r", "one-row    ", func() { g.Workspace().LayoutOnerow() },
		"c", "one-column ", func() { g.Workspace().LayoutOnecolumn() },
		"f", "fullscreen ", func() { g.Workspace().LayoutFullscreen() },
	)

	menu.Add("stat",
		"s", "toggle size  ", func() { filer.ToggleSizeView() },
		"p", "toggle perm  ", func() { filer.TogglePermView() },
		"t", "toggle time  ", func() { filer.ToggleTimeView() },
		"1", "all stat     ", func() { filer.SetStatView(true, true, true) },
		"0", "no stat      ", func() { filer.SetStatView(false, false, false) },
	)

	menu.Add("look",
		"d", "default      ", func() { look.Set("default") },
		"n", "midnight     ", func() { look.Set("midnight") },
		"b", "black        ", func() { look.Set("black") },
		"w", "white        ", func() { look.Set("white") },
		"a", "all border   ", func() { g.SetBorderStyle(widget.AllBorder) },
		"u", "ul border    ", func() { g.SetBorderStyle(widget.ULBorder) },
		"0", "no border    ", func() { g.SetBorderStyle(widget.NoBorder) },
	)

	menu.Add("command",
		"c", "copy         ", func() { g.Copy() },
		"m", "move         ", func() { g.Move() },
		"D", "delete       ", func() { g.Remove() },
		"k", "mkdir        ", func() { g.Mkdir() },
		"n", "newfile      ", func() { g.Touch() },
		"M", "chmod        ", func() { g.Chmod() },
		"r", "rename       ", func() { g.Rename() },
		"R", "bulk rename  ", func() { g.BulkRename() },
		"d", "chdir        ", func() { g.Chdir() },
		"g", "glob         ", func() { g.Glob() },
		"G", "globdir      ", func() { g.Globdir() },
	)
	g.AddKeymap("x", func() { g.Menu("command") })

	// [IMPL:EXTERNAL_COMMAND_LOADER] [IMPL:EXTERNAL_COMMAND_BINDER] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]
	commandEntries, loadErr := externalcmd.Load(externalcmd.Options{
		Path:  paths.Commands,
		GOOS:  runtime.GOOS,
		Debug: os.Getenv(externalcmd.EnvDebugCommands) != "",
		Logf: func(format string, args ...interface{}) {
			message.Infof("DEBUG: [IMPL:EXTERNAL_COMMAND_LOADER] "+format, args...)
		},
	})
	if loadErr != nil {
		message.Errorf("[REQ:EXTERNAL_COMMAND_CONFIG] failed to load %s: %v", paths.Commands, loadErr)
	}
	externalmenu.Register(g, commandEntries)
	g.AddKeymap("X", func() { g.Menu(externalcmd.MenuName) })

	menu.Add("archive",
		"z", "zip     ", func() { g.Shell(`zip -roD %x.zip %m`, -7) },
		"t", "tar     ", func() { g.Shell(`tar cvf %x.tar %m`, -7) },
		"g", "tar.gz  ", func() { g.Shell(`tar cvfz %x.tgz %m`, -7) },
		"b", "tar.bz2 ", func() { g.Shell(`tar cvfj %x.bz2 %m`, -7) },
		"x", "tar.xz  ", func() { g.Shell(`tar cvfJ %x.txz %m`, -7) },
		"r", "rar     ", func() { g.Shell(`rar u %x.rar %m`, -7) },

		"Z", "extract zip for %m", func() { g.Shell(`for i in %m; do unzip "$i" -d ./; done`, -6) },
		"T", "extract tar for %m", func() { g.Shell(`for i in %m; do tar xvf "$i" -C ./; done`, -6) },
		"G", "extract tgz for %m", func() { g.Shell(`for i in %m; do tar xvfz "$i" -C ./; done`, -6) },
		"B", "extract bz2 for %m", func() { g.Shell(`for i in %m; do tar xvfj "$i" -C ./; done`, -6) },
		"X", "extract txz for %m", func() { g.Shell(`for i in %m; do tar xvfJ "$i" -C ./; done`, -6) },
		"R", "extract rar for %m", func() { g.Shell(`for i in %m; do unrar x "$i" -C ./; done`, -6) },

		"1", "find . *.zip extract", func() { g.Shell(`find . -name "*.zip" -type f -prune -print0 | xargs -n1 -0 unzip -d ./`) },
		"2", "find . *.tar extract", func() { g.Shell(`find . -name "*.tar" -type f -prune -print0 | xargs -n1 -0 tar xvf -C ./`) },
		"3", "find . *.tgz extract", func() { g.Shell(`find . -name "*.tgz" -type f -prune -print0 | xargs -n1 -0 tar xvfz -C ./`) },
		"4", "find . *.bz2 extract", func() { g.Shell(`find . -name "*.bz2" -type f -prune -print0 | xargs -n1 -0 tar xvfj -C ./`) },
		"5", "find . *.txz extract", func() { g.Shell(`find . -name "*.txz" -type f -prune -print0 | xargs -n1 -0 tar xvfJ -C ./`) },
		"6", "find . *.rar extract", func() { g.Shell(`find . -name "*.rar" -type f -prune -print0 | xargs -n1 -0 unrar x -C ./`) },
	)

	menu.Add("bookmark",
		"t", "~/Desktop  ", func() { g.Dir().Chdir("~/Desktop") },
		"c", "~/Documents", func() { g.Dir().Chdir("~/Documents") },
		"d", "~/Downloads", func() { g.Dir().Chdir("~/Downloads") },
		"m", "~/Music    ", func() { g.Dir().Chdir("~/Music") },
		"p", "~/Pictures ", func() { g.Dir().Chdir("~/Pictures") },
		"v", "~/Videos   ", func() { g.Dir().Chdir("~/Videos") },
	)
	if runtime.GOOS == "windows" {
		menu.Add("bookmark",
			"C", "C:/", func() { g.Dir().Chdir("C:/") },
			"D", "D:/", func() { g.Dir().Chdir("D:/") },
			"E", "E:/", func() { g.Dir().Chdir("E:/") },
		)
	} else {
		menu.Add("bookmark",
			"e", "/etc   ", func() { g.Dir().Chdir("/etc") },
			"u", "/usr   ", func() { g.Dir().Chdir("/usr") },
			"x", "/media ", func() { g.Dir().Chdir("/media") },
		)
	}
	g.AddKeymap("b", func() { g.Menu("bookmark") })

	menu.Add("editor",
		"c", "vscode        ", func() { g.Spawn("code %f %&") },
		"e", "emacs client  ", func() { g.Spawn("emacsclient -n %f %&") },
		"v", "vim           ", func() { g.Spawn("vim %f") },
	)
	g.AddKeymap("e", func() { g.Menu("editor") })

	menu.Add("image",
		"x", "default    ", func() { g.Spawn(opener) },
		"e", "eog        ", func() { g.Spawn("eog %f %&") },
		"g", "gimp       ", func() { g.Spawn("gimp %m %&") },
	)

	menu.Add("media",
		"x", "default ", func() { g.Spawn(opener) },
		"m", "mpv     ", func() { g.Spawn("mpv %f") },
		"v", "vlc     ", func() { g.Spawn("vlc %f %&") },
	)

	// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
	// Helper for linked directory entry
	linkedEnterDir := func() {
		if g.IsLinkedNav() {
			name := g.File().Name()
			g.Workspace().ChdirAllToSubdir(name)
		}
		g.Dir().EnterDir()
	}

	var associate widget.Keymap
	if runtime.GOOS == "windows" {
		associate = widget.Keymap{
			".dir": linkedEnterDir, // [IMPL:LINKED_NAVIGATION]
			".go":  func() { g.Shell("go run %~f") },
			".py":  func() { g.Shell("python %~f") },
			".rb":  func() { g.Shell("ruby %~f") },
			".js":  func() { g.Shell("node %~f") },
		}
	} else {
		associate = widget.Keymap{
			".dir":  linkedEnterDir, // [IMPL:LINKED_NAVIGATION]
			".exec": func() { g.Shell(" ./" + g.File().Name()) },

			".zip": func() { g.Shell("unzip %f -d %D") },
			".tar": func() { g.Shell("tar xvf %f -C %D") },
			".gz":  func() { g.Shell("tar xvfz %f -C %D") },
			".tgz": func() { g.Shell("tar xvfz %f -C %D") },
			".bz2": func() { g.Shell("tar xvfj %f -C %D") },
			".xz":  func() { g.Shell("tar xvfJ %f -C %D") },
			".txz": func() { g.Shell("tar xvfJ %f -C %D") },
			".rar": func() { g.Shell("unrar x %f -C %D") },

			".go": func() { g.Shell("go run %f") },
			".py": func() { g.Shell("python %f") },
			".rb": func() { g.Shell("ruby %f") },
			".js": func() { g.Shell("node %f") },

			".jpg":  func() { g.Menu("image") },
			".jpeg": func() { g.Menu("image") },
			".gif":  func() { g.Menu("image") },
			".png":  func() { g.Menu("image") },
			".bmp":  func() { g.Menu("image") },

			".avi":  func() { g.Menu("media") },
			".mp4":  func() { g.Menu("media") },
			".mkv":  func() { g.Menu("media") },
			".wmv":  func() { g.Menu("media") },
			".flv":  func() { g.Menu("media") },
			".mp3":  func() { g.Menu("media") },
			".flac": func() { g.Menu("media") },
			".tta":  func() { g.Menu("media") },
		}
	}

	g.MergeExtmap(widget.Extmap{
		"C-m": associate,
		"o":   associate,
	})
}

func loadExcludedNames(path string) {
	if path == "" {
		filer.ConfigureExcludedNames(nil, false)
		return
	}
	names, err := parseExcludeFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			filer.ConfigureExcludedNames(nil, false)
			return
		}
		fmt.Fprintf(os.Stderr, "WARN: [REQ:FILER_EXCLUDE_NAMES] failed to read exclude list %s: %v\n", path, err)
		filer.ConfigureExcludedNames(nil, false)
		return
	}
	count := filer.ConfigureExcludedNames(names, true)
	if count == 0 {
		filer.ConfigureExcludedNames(nil, false)
		return
	}
	if os.Getenv("GOFUL_DEBUG_PATHS") != "" {
		fmt.Fprintf(os.Stderr, "DEBUG: [REQ:FILER_EXCLUDE_NAMES] loaded %d excluded names from %s\n", count, path)
	}
}

func parseExcludeFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseExcludeLines(file)
}

func parseExcludeLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	names := make([]string, 0, 16)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		names = append(names, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return names, nil
}

func emitPathDebug(paths configpaths.Paths) {
	if os.Getenv("GOFUL_DEBUG_PATHS") == "" {
		return
	}
	fmt.Fprintf(
		os.Stderr,
		"DEBUG: [IMPL:STATE_PATH_RESOLVER] [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS] [REQ:EXTERNAL_COMMAND_CONFIG] [REQ:FILER_EXCLUDE_NAMES] [REQ:FILE_COMPARISON_COLORS] state=%s (%s) history=%s (%s) commands=%s (%s) excludes=%s (%s) compare_colors=%s (%s)\n",
		paths.State,
		paths.StateSource,
		paths.History,
		paths.HistorySource,
		paths.Commands,
		paths.CommandsSource,
		paths.Excludes,
		paths.ExcludesSource,
		paths.CompareColors,
		paths.CompareColorsSource,
	)
}

// loadCompareColors loads the comparison color configuration from the specified path.
// [IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func loadCompareColors(path string) {
	cfg, err := comparecolors.Load(path)
	if err != nil {
		if os.Getenv("GOFUL_DEBUG_PATHS") != "" {
			fmt.Fprintf(os.Stderr, "DEBUG: [IMPL:COMPARE_COLOR_CONFIG] failed to load compare colors from %s: %v\n", path, err)
		}
	}
	parsed := cfg.Parse()
	look.ConfigureComparisonColors(parsed)
	if os.Getenv("GOFUL_DEBUG_PATHS") != "" {
		fmt.Fprintf(os.Stderr, "DEBUG: [IMPL:COMPARE_COLOR_CONFIG] loaded comparison colors from %s\n", path)
	}
}

// Widget keymap functions.

func filerKeymap(g *app.Goful) widget.Keymap {
	// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
	// Helper for linked parent navigation
	linkedParentNav := func() {
		if g.IsLinkedNav() {
			g.Workspace().ChdirAllToParent()
		}
		g.Dir().Chdir("..")
	}

	return widget.Keymap{
		"M-C-o":     func() { g.CreateWorkspace() },
		"M-C-w":     func() { g.CloseWorkspace() },
		"M-f":       func() { g.MoveWorkspace(1) },
		"M-b":       func() { g.MoveWorkspace(-1) },
		"C-o":       func() { g.Workspace().CreateDir() },
		"C-w":       func() { g.Workspace().CloseDir() },
		"C-l":       func() { g.Workspace().ReloadAll() },
		"C-f":       func() { g.Workspace().MoveFocus(1) },
		"C-b":       func() { g.Workspace().MoveFocus(-1) },
		"right":     func() { g.Workspace().MoveFocus(1) },
		"left":      func() { g.Workspace().MoveFocus(-1) },
		"C-i":       func() { g.Workspace().MoveFocus(1) },
		"l":         func() { g.Workspace().MoveFocus(1) },
		"h":         func() { g.Workspace().MoveFocus(-1) },
		"F":         func() { g.Workspace().SwapNextDir() },
		"B":         func() { g.Workspace().SwapPrevDir() },
		"w":         func() { g.Workspace().ChdirNeighbor() },
		"C-h":       linkedParentNav, // [IMPL:LINKED_NAVIGATION]
		"backspace": linkedParentNav, // [IMPL:LINKED_NAVIGATION]
		"u":         linkedParentNav, // [IMPL:LINKED_NAVIGATION]
		// [IMPL:LINKED_NAVIGATION] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]
		// Toggle linked navigation mode with Alt+l or L (uppercase for macOS compatibility)
		"M-l": func() {
			enabled := g.ToggleLinkedNav()
			state := "disabled"
			if enabled {
				state = "enabled"
			}
			message.Infof("[REQ:LINKED_NAVIGATION] linked navigation %s", state)
		},
		// [IMPL:LINKED_NAVIGATION] macOS-friendly alternative (Option key often produces special chars)
		"L": func() {
			enabled := g.ToggleLinkedNav()
			state := "disabled"
			if enabled {
				state = "enabled"
			}
			message.Infof("[REQ:LINKED_NAVIGATION] linked navigation %s", state)
		},
		"~":    func() { g.Dir().Chdir("~") },
		"\\":   func() { g.Dir().Chdir("/") },
		"C-n":  func() { g.Dir().MoveCursor(1) },
		"C-p":  func() { g.Dir().MoveCursor(-1) },
		"down": func() { g.Dir().MoveCursor(1) },
		"up":   func() { g.Dir().MoveCursor(-1) },
		"j":    func() { g.Dir().MoveCursor(1) },
		"k":    func() { g.Dir().MoveCursor(-1) },
		"C-d":  func() { g.Dir().MoveCursor(5) },
		"C-u":  func() { g.Dir().MoveCursor(-5) },
		"C-a":  func() { g.Dir().MoveTop() },
		"C-e":  func() { g.Dir().MoveBottom() },
		"home": func() { g.Dir().MoveTop() },
		"end":  func() { g.Dir().MoveBottom() },
		"^":    func() { g.Dir().MoveTop() },
		"$":    func() { g.Dir().MoveBottom() },
		"M-n":  func() { g.Dir().Scroll(1) },
		"M-p":  func() { g.Dir().Scroll(-1) },
		"C-v":  func() { g.Dir().PageDown() },
		"M-v":  func() { g.Dir().PageUp() },
		"pgdn": func() { g.Dir().PageDown() },
		"pgup": func() { g.Dir().PageUp() },
		" ":    func() { g.Dir().ToggleMark() },
		"M-=":  func() { g.Dir().InvertMark() },
		"C-g":  func() { g.Dir().Reset() },
		"C-[":  func() { g.Dir().Reset() }, // C-[ means ESC
		"f":    func() { g.Dir().Finder() },
		"/":    func() { g.Dir().Finder() },
		"q":    func() { g.Quit() },
		"Q":    func() { g.Quit() },
		";":    func() { g.Shell("") },
		":":    func() { g.ShellSuspend("") },
		"M-W":  func() { g.ChangeWorkspaceTitle() },
		"n":    func() { g.Touch() },
		"K":    func() { g.Mkdir() },
		"c":    func() { g.Copy() },
		"m":    func() { g.Move() },
		"r":    func() { g.Rename() },
		"R":    func() { g.BulkRename() },
		"D":    func() { g.Remove() },
		"d":    func() { g.Chdir() },
		"g":    func() { g.Glob() },
		"G":    func() { g.Globdir() },
		// [IMPL:DIFF_SEARCH] [ARCH:DIFF_SEARCH] [REQ:DIFF_SEARCH]
		// Difference search commands
		"[": func() { g.StartDiffSearch() },    // Start difference search
		"]": func() { g.ContinueDiffSearch() }, // Continue to next difference
	}
}

func finderKeymap(w *filer.Finder) widget.Keymap {
	return widget.Keymap{
		"C-h":       func() { w.DeleteBackwardChar() },
		"backspace": func() { w.DeleteBackwardChar() },
		"M-p":       func() { w.MoveHistory(1) },
		"M-n":       func() { w.MoveHistory(-1) },
		"C-g":       func() { w.Exit() },
		"C-[":       func() { w.Exit() },
	}
}

func cmdlineKeymap(w *cmdline.Cmdline) widget.Keymap {
	return widget.Keymap{
		"C-a":       func() { w.MoveTop() },
		"C-e":       func() { w.MoveBottom() },
		"C-f":       func() { w.ForwardChar() },
		"C-b":       func() { w.BackwardChar() },
		"right":     func() { w.ForwardChar() },
		"left":      func() { w.BackwardChar() },
		"M-f":       func() { w.ForwardWord() },
		"M-b":       func() { w.BackwardWord() },
		"C-d":       func() { w.DeleteChar() },
		"delete":    func() { w.DeleteChar() },
		"C-h":       func() { w.DeleteBackwardChar() },
		"backspace": func() { w.DeleteBackwardChar() },
		"M-d":       func() { w.DeleteForwardWord() },
		"M-h":       func() { w.DeleteBackwardWord() },
		"C-k":       func() { w.KillLine() },
		"C-i":       func() { w.StartCompletion() },
		"C-m":       func() { w.Run() },
		"C-g":       func() { w.Exit() },
		"C-[":       func() { w.Exit() },
		"C-n":       func() { w.History.CursorDown() },
		"C-p":       func() { w.History.CursorUp() },
		"down":      func() { w.History.CursorDown() },
		"up":        func() { w.History.CursorUp() },
		"C-v":       func() { w.History.PageDown() },
		"M-v":       func() { w.History.PageUp() },
		"pgdn":      func() { w.History.PageDown() },
		"pgup":      func() { w.History.PageUp() },
		"M-<":       func() { w.History.MoveTop() },
		"M->":       func() { w.History.MoveBottom() },
		"home":      func() { w.History.MoveTop() },
		"end":       func() { w.History.MoveBottom() },
		"M-n":       func() { w.History.Scroll(1) },
		"M-p":       func() { w.History.Scroll(-1) },
		"C-x":       func() { w.History.Delete() },
	}
}

func completionKeymap(w *cmdline.Completion) widget.Keymap {
	return widget.Keymap{
		"C-n":   func() { w.CursorDown() },
		"C-p":   func() { w.CursorUp() },
		"down":  func() { w.CursorDown() },
		"up":    func() { w.CursorUp() },
		"C-f":   func() { w.CursorToRight() },
		"C-b":   func() { w.CursorToLeft() },
		"right": func() { w.CursorToRight() },
		"left":  func() { w.CursorToLeft() },
		"C-i":   func() { w.CursorToRight() },
		"C-v":   func() { w.PageDown() },
		"M-v":   func() { w.PageUp() },
		"pgdn":  func() { w.PageDown() },
		"pgup":  func() { w.PageUp() },
		"M-<":   func() { w.MoveTop() },
		"M->":   func() { w.MoveBottom() },
		"home":  func() { w.MoveTop() },
		"end":   func() { w.MoveBottom() },
		"M-n":   func() { w.Scroll(1) },
		"M-p":   func() { w.Scroll(-1) },
		"C-m":   func() { w.InsertCompletion() },
		"C-g":   func() { w.Exit() },
		"C-[":   func() { w.Exit() },
	}
}

func menuKeymap(w *menu.Menu) widget.Keymap {
	return widget.Keymap{
		"C-n":  func() { w.MoveCursor(1) },
		"C-p":  func() { w.MoveCursor(-1) },
		"down": func() { w.MoveCursor(1) },
		"up":   func() { w.MoveCursor(-1) },
		"C-v":  func() { w.PageDown() },
		"M-v":  func() { w.PageUp() },
		"M->":  func() { w.MoveBottom() },
		"M-<":  func() { w.MoveTop() },
		"C-m":  func() { w.Exec() },
		"C-g":  func() { w.Exit() },
		"C-[":  func() { w.Exit() },
	}
}
