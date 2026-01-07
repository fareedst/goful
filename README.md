# Goful

[![Go Report Card](https://goreportcard.com/badge/github.com/anmitsu/goful)](https://goreportcard.com/report/github.com/anmitsu/goful)
[![Go Reference](https://pkg.go.dev/badge/github.com/anmitsu/goful.svg)](https://pkg.go.dev/github.com/anmitsu/goful)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/anmitsu/goful/blob/master/LICENSE)

Goful is a CUI file manager written in Go.

* Works on cross-platform such as gnome-terminal and cmd.exe
* Displays multiple windows and workspaces
* A command line to execute using such as bash and tmux
* Provides filtering search, async copy, glob, bulk rename, etc.

![demo](.github/demo.gif)

## Install

### Pre-build binaries

See [releases](https://github.com/anmitsu/goful/releases).

### Go version >= 1.16

    $ go install github.com/anmitsu/goful@latest
    ...
    $ goful

### Go version < 1.16

    $ go get github.com/anmitsu/goful
    ...
    $ goful

## Usage

key                  | function
---------------------|-----------
`C-n` `down` `j`     | Move cursor down
`C-p` `up` `k`       | Move cursor up
`C-a` `home` `^`     | Move cursor top
`C-e` `end` `$`      | Move cursor bottom
`C-f` `C-i` `right` `l`| Move cursor right
`C-b` `left` `h`     | Move cursor left
`C-d`                | More move cursor down
`C-u`                | More move cursor up
`C-v` `pgdn`         | Page down
`M-v` `pgup`         | Page up
`M-n`                | Scroll down
`M-p`                | Scroll up
`C-h` `backspace` `u`| Change to upper directory
`~`                  | Change to home directory
`\`                  | Change to root directory
`w`                  | Change to neighbor directory
`C-o`                | Create directory window
`C-w`                | Close directory window
`M-f`                | Move next workspace
`M-b`                | Move previous workspace
`M-C-o`              | Create workspace
`M-C-w`              | Close workspace
`space`              | Toggle mark
`M-=`                | Invert mark
`C-l`                | Reload
`C-m` `o`            | Open
`i`                  | Open by pager
`s`                  | Sort
`v`                  | View
`b`                  | Bookmark
`e`                  | Editor
`x`                  | Command
`X`                  | External command
`f` `/`              | Find
`;`                  | Shell
`:`                  | Shell suspend
`n`                  | Make file
`K`                  | Make directory
`c`                  | Copy
`m`                  | Move
`r`                  | Rename
`R`                  | Bulk rename by regexp
`D`                  | Remove
`d`                  | Change directory
`g`                  | Glob
`G`                  | Glob recursive
`C-g` `C-[`          | Cancel
`q` `Q`              | Quit

For more see [main.go](main.go)

## Demos

### Copy and Move

Copy (default `c`) and move (default `m`) mark (default `space` and invert
`C-space`) files.

![demo_copy](.github/demo_copy.gif)

First input a copy source file name (or path).  The default source is a file
name on the cursor.  If files are marked, this step is skipped.

Second input a copy destination path and start copy processing.  The default
destination is a neighboring directory path.

During processing draws copying file count, size and name, progress percent,
gauge, bps and estimated time of arrival.

If the source file type is a directory, recursively copy.  Also copy
modification time and permissions.

Rise a override confirm dialog `[y/n/Y/N]` if the name same as source file
exists in the destination.  This dialog means:

* `y` is overwrite only this file
* `n` is not overwrite only this file
* `Y` is overwrite all later file
* `N` is not overwrite all later file

Copy process works asynchronously.  And processed in the order if you run
multiple copies.

Note that copy process can not interrupt.  If you want to interrupt, please quit
the application (default `q` `Q`).

### Bulk Rename

Bulk renaming (default `R`) for mark (default `space` and invert `C-space`)
files.

Rename by the regexp pattern.  Input like the vim substituting style
(regexp/replaced).  Display and confirm matched files and replaced names before
rename.

![demo_bulk](.github/demo_bulk.gif)

### Finder (Filtering search)

The finder (default `f` `/`) filters files in the directory.

Input characters recognizes as the regexp.  Case insensitive when inputs
lowercase only, on the other hand case sensitive when contains uppercase.

Delete characters by `C-h` and `backspace` (default).  Can select input
histories by `M-p` and `M-n` (default).

Other than character inputs (exclude a space) and the finder keymap pass to the
main input.

Hit reset key (default `C-g` `C-[` means `Esc`) to clear filtering.

![demo_finder](.github/demo_finder.gif)

### Glob

Glob is matched by wild card pattern in the current directory (default `g` and
recursive `G`).

Hit reset key (default `C-g` `C-[` means `Esc`) to clear glob patterns.

![demo_glob](.github/demo_glob.gif)

### Layout

Directory windows position are allocated by layouts of tile, tile-top,
tile-bottom, one-row, one-column and fullscreen.

View menu (default `v`), run layout menu and select layout:

![demo_layout](.github/demo_layout.gif)

### Execute Terminal and Shell

Shell mode (default `:` and suspended `;`) runs a terminal and execute shell
such as bash and tmux.  The cmdline completion (file names and commands in
$PATH) is available (default `C-i` that means `tab`).

For example, spawns commands by bash in a gnome-terminal new tab:

![demo_shell](.github/demo_shell.gif)

The terminal immediately doesn't close when command finished because check
outputs.

If goful is running in tmux, it creates a new window and executes the command.

`[REQ:TERMINAL_PORTABILITY]` keeps this workflow intact on macOS and alternative
terminal emulators:

- On macOS, goful now opens Terminal.app via `osascript` so commands run in a new
  tab/window without editing `main.go`. The adapter automatically runs
  `cd "%D";` before your command so relative paths match the focused window, and
  the command stays visible until you hit enter because the historical
  `read -p "HIT ENTER KEY"` prompt remains—after acknowledging it, the adapter
  issues `exit` so the temporary window closes cleanly.
- Want another AppleScript-friendly terminal or shell? Set
  `GOFUL_TERMINAL_APP="iTerm2"` (or similar) to reuse the built-in macOS branch
  with a different application, and set `GOFUL_TERMINAL_SHELL="zsh"` (or any
  other binary) to change the inline shell that runs inside that window—defaults
  remain Terminal.app + bash so existing workflows stay unchanged.
- Set `GOFUL_TERMINAL_CMD="alacritty -e"` (or another emulator) to override the
  launcher on Linux/BSD. The override is inserted before the usual `bash -c`
  invocation, so your configured terminal still receives the expanded command
  and inherits the macOS working-directory preamble when applicable.
- Export `GOFUL_DEBUG_TERMINAL=1` to print `DEBUG: [IMPL:TERMINAL_ADAPTER]`
  messages describing which branch was selected if you need to troubleshoot
  environment detection.

### Expand Macro

macro        | expanded string
-------------|------------------
`%f` `%F`   | File name/path on cursor
`%x` `%X`   | File name/path with extension excluded on cursor
`%m` `%M`   | Marked file names/paths joined by spaces
`%d` `%D`   | Directory name/path on cursor
`%d2` `%D2` | Neighbor directory name/path
`%D@` `%~D@` | Current window stays `%D`; both macros append the other directory paths in display order so `echo %D %D@` lists every window. `%D@` quotes each appended path for shell safety, while `%~D@` deliberately emits the same list without quoting so advanced scripts can opt into raw arguments.
`%d@` `%~d@` | Appends the other directory **names** (`Directory.Base()`) in the same order. `%d@` quotes each name so shell invocations stay safe, while `%~d@` emits raw names when scripts only need lightweight labels.
`%~f` ...   | Expand by non quote
`%&`        | Flag to run command in background

The macro is useful if do not want to specify a file name when run the shell.

Macros starts with `%` are expanded surrounded by quote, and those starts with
`%~` are expanded by non quote (including `%~D@` and `%~d@`, which return raw directory paths or names just like other tilde macros). The `%~` mainly uses to for cmd.exe.

Use `%&` when background execute the shell such as GUI apps launching.

![demo_macro](.github/demo_macro.gif)

<!-- demo size 120x35 -->

## Customize

Most of goful is still customized by editing `main.go`, but `[REQ:EXTERNAL_COMMAND_CONFIG]` now provides a JSON or YAML file for the `external-command` menu so you can ship your own automation without rebuilding the binary.

Examples of customizing:

* Change and add keybindings
* Change terminal and shell
* Change file opener (editor, pager and more)
* Adding bookmarks
* Setting colors and looks

Recommend remain original `main.go` and copy to own `main.go` for example:

Go to source directory

    $ cd $GOPATH/src/github.com/anmitsu/goful

Copy original `main.go` to `my/goful` directory

    $ mkdir -p my/goful
    $ cp main.go my/goful
    $ cd my/goful

Install after edit `my/goful/main.go`

    $ go install

### External Command Config (`external-command` menu)

`[REQ:EXTERNAL_COMMAND_CONFIG]` and `[ARCH:EXTERNAL_COMMAND_REGISTRY]` move the `external-command` menu definitions into a JSON **or** YAML file so you can add/edit/remove shell helpers without touching Go code:

- Flag `-commands /path/to/external_commands.yaml` overrides everything.
- Environment variable `GOFUL_COMMANDS_FILE` applies when the CLI flag is unset.
- Defaults fall back to `~/.goful/external_commands.yaml`, with the historical POSIX/Windows bindings baked in if the file does not exist yet.
- Omit `inheritDefaults` (or set it to `true`) to **prepend** your file-based commands ahead of the compiled defaults; set `inheritDefaults: false` to replace them entirely. `[IMPL:EXTERNAL_COMMAND_APPEND]`
- Set `GOFUL_DEBUG_COMMANDS=1` to log loader diagnostics (`DEBUG: [IMPL:EXTERNAL_COMMAND_LOADER] ...`).

Each entry accepts (shown in JSON, but the same fields work in YAML):

```jsonc
[
  {
    "menu": "external-command", // optional; defaults to external-command
    "key": "c",
    "label": "copy %m to %D2    ",
    "command": "cp -vai %m %D2",
    "offset": -2,               // optional cursor offset for g.Shell
    "platforms": ["linux"],     // optional GOOS filter (case-insensitive)
    "disabled": false
  },
  {
    "key": "A",
    "label": "archives menu     ",
    "runMenu": "archive"        // optional: jump into another goful menu instead of running a shell command
  }
]
```

YAML example:

```yaml
inheritDefaults: true # omit for the default prepend behavior
commands:
  - key: c
    label: "copy %m to %D2    "
    command: "cp -vai %m %D2"
    offset: -2
  - key: A
    label: "archives menu     "
    runMenu: "archive"
```

Need an object wrapper in JSON to drop the defaults entirely:

```jsonc
{
  "inheritDefaults": false,
  "commands": [
    {
      "key": "z",
      "label": "zip %m",
      "command": "zip -roD %x.zip %m"
    }
  ]
}
```

- Use goful macros (`%f`, `%D@`, `%~m`, etc.) inside `command` strings the same way the legacy `main.go` menu did.
- Entries are applied in file order, and duplicate `menu/key` combinations are rejected with a descriptive error surfaced via `message.Errorf`.
- If the file disables every entry, the `external-command` menu still appears and displays a placeholder that explains how to re-enable commands.

### Configuring State & History Paths

`[REQ:CONFIGURABLE_STATE_PATHS]` and `[ARCH:STATE_PATH_SELECTION]` make it possible to redirect the persisted UI state and cmdline history without editing the source:

- Defaults remain `~/.goful/state.json` and `~/.goful/history/shell`.
- Set `GOFUL_STATE_PATH` or `GOFUL_HISTORY_PATH` to override the defaults for a shell/session.
- Pass `-state /tmp/state.json` or `-history /tmp/history` on the command line to override everything else (flags win over environment variables).
- Export `GOFUL_DEBUG_PATHS=1` to log which source produced each path (`DEBUG: [IMPL:STATE_PATH_RESOLVER] ...`) for troubleshooting sandboxes and CI jobs.

### Startup Workspace Directories

`[REQ:WORKSPACE_START_DIRS]` and `[ARCH:WORKSPACE_BOOTSTRAP]` let you pass directories **after** the usual CLI flags so goful opens one filer window per argument (ordered). Examples:

```bash
$ goful ~/src/dotfiles ~/src/goful ~/Downloads
# window 1 -> ~/src/dotfiles, window 2 -> ~/src/goful, window 3 -> ~/Downloads
```

- The arguments are normalized (tilde expansion + absolute paths) but otherwise preserved, so duplicates intentionally open multiple panes pointing at the same directory.
- If an argument does not exist or is not a directory, goful prints `message.Errorf` output explaining the issue and continues processing the remaining entries. When every argument fails, startup falls back to the persisted workspace layout.
- Set `GOFUL_DEBUG_WORKSPACE=1` to emit `DEBUG: [IMPL:WORKSPACE_START_DIRS] ...` lines that document the parsed arguments and window assignments—useful when debugging automation scripts or verifying CI launches.
- Launching without trailing directories keeps the historical behavior (state restoration or default layout), so existing workflows continue to work unchanged.

### Startup Workspace Directories

`[REQ:WORKSPACE_START_DIRS]` and `[ARCH:WORKSPACE_BOOTSTRAP]` let you tell goful which directories to open **after** the usual CLI flags—every remaining positional argument seeds one workspace (tab) in the order provided:

```bash
$ goful ~/src/dotfiles ~/src/goful ~/Downloads
```

- When at least one directory is supplied, goful resizes the workspace list to match (creating or closing tabs as needed) and focuses the first entry. Provide duplicates if you want multiple tabs pointing to the same path; the order is preserved exactly as entered.
- When no positional arguments are present, the historical startup behavior (state.json restoration or default layout) remains unchanged.
- Invalid entries surface via `message.Errorf` before the UI launches so you can fix typos without guessing which path failed.
- Set `GOFUL_DEBUG_WORKSPACE=1` to log `DEBUG: [IMPL:WORKSPACE_START_DIRS] ...` lines that describe how CLI arguments were parsed, which workspaces were created/removed, and which directories were applied—handy when automating project layouts.

## Documentation

- [`ARCHITECTURE.md`](ARCHITECTURE.md) – `[REQ:ARCH_DOCUMENTATION]` package/data-flow overview plus module validation map.
- [`CONTRIBUTING.md`](CONTRIBUTING.md) – `[REQ:CONTRIBUTING_GUIDE]` workflow, semantic token, and module-validation requirements.
- `stdd/` directory – STDD requirements, architecture decisions, implementation decisions, semantic token registry, and task tracker.

## Build & Release

- `make lint` / `make test` keep local runs in sync with CI’s fmt/vet/test gates.
- `make release` (`[REQ:RELEASE_BUILD_MATRIX]` / `[IMPL:MAKE_RELEASE_TARGETS]`) produces CGO-disabled binaries plus `.sha256` digests in `dist/`. Set `PLATFORM=linux/amd64` (or `linux/arm64`, `darwin/arm64`) to build a single target; otherwise all default platforms are generated.
- CI job `release-matrix` (in `.github/workflows/ci.yml`) reuses the same Makefile target so every PR proves artifacts remain reproducible before merge.
- Tag pushes `v*` trigger `.github/workflows/release.yml`, which re-runs the matrix, logs deterministic filenames/checksums, and publishes the binaries + digest files to the GitHub Release via `softprops/action-gh-release`.

## Contributing

Read the [Contributing Guide](CONTRIBUTING.md) before opening a PR so you follow the STDD workflow, semantic token discipline, and module validation rules.
