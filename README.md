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
`E`                  | Toggle filename excludes
`` ` ``              | Toggle comparison colors
`C`                  | Copy All (to all visible panes)
`M`                  | Move All (to all visible panes)
`S`                  | Sync command mode (synchronized ops)
`=`                  | Calculate file digest
`L` or `M-l`             | Toggle linked navigation mode
`[`                  | Start difference search
`]`                  | Continue difference search
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
`?`                  | Help (keystroke catalog with color styling and mouse scroll)
`q` `Q`              | Quit

For more see [main.go](main.go)

### Mouse Support `[REQ:MOUSE_FILE_SELECT]` `[REQ:TOOLBAR_PARENT_BUTTON]` `[REQ:TOOLBAR_LINKED_TOGGLE]` `[REQ:TOOLBAR_COMPARE_BUTTON]`

Goful provides mouse support for common navigation and file operations, making it accessible to users who prefer graphical interaction alongside keyboard shortcuts.

**Toolbar Buttons**

The header row displays clickable toolbar buttons at the left edge:

| Button | Function | Behavior |
|--------|----------|----------|
| `[^]` | Parent directory | Navigate to parent directory. Respects linked mode: when ON, all windows navigate to their respective parents; when OFF, only the focused window navigates. |
| `[L]` | Linked mode toggle | Toggle linked navigation mode. The button style indicates current state: **highlighted** (reverse) when ON, normal when OFF. Click to toggle and see confirmation message. |
| `[=]` | Compare all digests | Calculate xxHash64 digests for **all** files that appear in multiple panes. Equivalent to pressing `=` on every shared filename. Displays a summary message with the count of files processed. |

**File Selection**

- **Left-click on a file**: Selects the file and moves the cursor. If clicking in an unfocused pane, switches focus to that pane first.
- **Cross-window cursor sync**: When linked mode is ON, clicking a file syncs all other panes to the same filename (if it exists). When linked mode is OFF, only the clicked pane is affected. Panes that don't contain a matching file have their cursor highlight erased to avoid confusion.
- **Double-click on a directory**: Navigates into the directory. When linked mode is ON, all panes that contain a matching subdirectory navigate into it.
- **Double-click on a file**: Opens the file using the system default application. When linked mode is ON, opens the same-named file from all panes where it exists.

**Scrolling**

- **Mouse wheel up/down**: Scrolls the file list in the pane under the cursor (3 lines per scroll).

### Help Popup `[REQ:HELP_POPUP_ENHANCEMENTS]`

Press `?` to open a styled help popup displaying the complete keystroke catalog. The help popup features:

**Color Styling**
- **Border**: Distinct color matching the current theme
- **Header**: Bold title with theme-aware highlighting
- **Keybindings**: Color-coded key combinations for easy scanning
- **Descriptions**: Styled command descriptions

The colors automatically adapt to the active theme (default, midnight, black, or white).

**Navigation**
- **Keyboard**: `j`/`k` or arrow keys move the cursor, `pgup`/`pgdn` for page scrolling, `M-n`/`M-p` for smooth scroll
- **Mouse wheel**: Scroll up/down through the help content
- **Exit**: Press `?`, `q`, `C-g`, or `Esc` to close

### Linked navigation mode `[REQ:LINKED_NAVIGATION]`

Linked navigation synchronizes directory navigation and cursor position across all open panes. When enabled, navigating into a subdirectory, pressing backspace (parent directory), or moving the cursor in the focused window causes all other windows to sync accordingly.

**Toggle**: Press `L` (uppercase L), `M-l` (Alt+l), or click the `[L]` toolbar button to enable or disable linked mode. The `[L]` button in the header is always visible: **highlighted** when linked mode is ON, normal when OFF.

> **macOS note**: The Option key often produces special characters instead of acting as Meta/Alt, so use uppercase `L` (Shift+l) which works reliably across all platforms.

**Cursor synchronization**: When you move the cursor (via keyboard `j`/`k`/arrows/page keys or mouse click), all other panes move their cursors to the same filename if it exists. Panes without a matching file have their cursor highlight erased. When linked mode is OFF, cursor movements only affect the focused pane.

**Subdirectory navigation**: When you enter a subdirectory named `foo`, all other panes that also contain a subdirectory named `foo` will navigate into it. Panes without a matching subdirectory stay on their current path.

**Parent navigation**: When you press backspace (or `C-h` or `u`), all panes navigate to their respective parent directories.

**Sort synchronization**: When you change the sort order (via the `s` menu), all panes adopt the same sort order.

**Use cases**:
- Comparing parallel directory structures (e.g., syncing `src/` and `backup/src/`)
- Navigating release versions side-by-side (`v1.0/`, `v2.0/`)
- Keeping workspace panes aligned when exploring mirrored folder hierarchies
- Quickly finding matching files across panes by moving cursor in any single pane

The mode is **on by default** and does not persist across restarts. Press `L` or click the `[L]` toolbar button to toggle it off if you prefer independent navigation.

### Difference search `[REQ:DIFF_SEARCH]`

Difference search finds files and directories that differ across your workspace panes. It's designed for comparing similar directory structures—like backup copies, release versions, or synced folders—by highlighting entries that are missing or have different sizes.

**Start search**: Press `[` to begin a new difference search. Goful records the current directories in all panes and starts scanning for differences alphabetically.

**Continue search**: Press `]` to advance to the next difference. The search resumes from the current cursor position and finds the next entry that differs.

**What counts as a difference**:
- A file or directory that exists in some panes but is missing from others
- A file that exists in all panes but has different sizes

**How the search works**:
1. Files are compared first, in alphabetical order
2. When all files at the current level match, the search descends into subdirectories that exist in all panes
3. Subdirectories that exist in only some panes are flagged as differences
4. The search continues depth-first until returning to the starting directories

**Status display**: While a search is active, the header shows a `[DIFF: ...]` indicator with either the current search progress or the last found difference.

**Navigating during search**: You can freely navigate within panes between `]` presses. When you continue the search, it resumes from your current position and correctly traverses back through the directory tree to find the next difference.

**Ending the search**: The search ends automatically when all directories have been checked. You can also start a new search with `[` at any time, which replaces the current search state.

**Use cases**:
- Finding files that failed to sync between backup and source directories
- Locating additions or deletions between release versions
- Identifying size mismatches that indicate incomplete copies or corruption

**Tip**: Combine with **linked navigation mode** (`L`) to keep panes synchronized as you explore differences.

### Batch Diff Report (N-Way Comparison) `[REQ:BATCH_DIFF_REPORT]`

Goful offers a **unique n-way directory comparison** capability that goes beyond traditional two-directory diff tools. While most comparison utilities are limited to pairwise comparisons, goful can simultaneously compare 2, 3, 4, or more directory trees in a single operation—revealing differences that would require multiple manual comparisons with conventional tools.

**Use cases for n-way comparison**:
- **Multi-environment validation**: Compare staging, production, and development deployments simultaneously
- **Parallel backup verification**: Validate that local, remote, and archive copies are identical
- **Multi-version analysis**: Track changes across several release versions at once
- **Distributed system checks**: Verify configuration consistency across multiple servers

**CLI command**: Run a non-interactive batch comparison with YAML output:

```bash
# Compare two directories
goful --diff-report dir1 dir2

# Compare three or more directories (n-way)
goful --diff-report prod/ staging/ dev/

# Suppress progress output for scripting
goful --diff-report --quiet backup1/ backup2/ backup3/
```

**Output format**: The command produces a structured YAML report to stdout:

```yaml
directories:
  - /path/to/prod
  - /path/to/staging
  - /path/to/dev
totalFilesChecked: 1542
totalDirectoriesTraversed: 87
durationSeconds: 2.34
differences:
  - name: config.json
    path: app/config.json
    reason: "size mismatch: 1024 vs 1048 vs 1024"
    isDir: false
  - name: cache/
    path: cache/
    reason: "missing in window 2"
    isDir: true
```

**Exit codes** for scripting:
| Code | Meaning |
|------|---------|
| 0 | No differences found |
| 1 | Error (invalid arguments, directory not found) |
| 2 | Differences found |

**Progress reporting**: By default, progress updates are printed to stderr every 2 seconds. Use `--quiet` to suppress these for cleaner pipeline integration.

**Example pipeline**:
```bash
# Check backups and alert on differences
if ! goful --diff-report --quiet /backup/daily /backup/weekly /live; then
    goful --diff-report /backup/daily /backup/weekly /live > /var/log/diff-report.yaml
    mail -s "Backup mismatch detected" admin@example.com < /var/log/diff-report.yaml
fi
```

### Filename exclude list `[REQ:FILER_EXCLUDE_NAMES]`

- Create a newline-delimited file containing the basenames you want to hide (for example `.DS_Store`, `Thumbs.db`). Blank lines and lines starting with `#` are ignored.
- The list is loaded from `~/.goful/excludes` by default. Override the path with `-exclude-names /path/to/file` or `GOFUL_EXCLUDES_FILE=/path/to/file` to share custom lists across machines.
- Matching is case-insensitive and applies to every directory pane, finder result, and macro expansion while the filter is enabled.
- Press `E` (or open the View menu and press `n`) to toggle the filter at runtime. Goful instantly reloads every pane so you can temporarily reveal hidden files, inspect them, and re-enable the filter without restarting the UI.

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

### File comparison colors (`` ` `` toggle)

`[REQ:FILE_COMPARISON_COLORS]` (implemented via `[ARCH:FILE_COMPARISON_ENGINE]` + `[IMPL:COMPARE_COLOR_CONFIG]`) adds a per-column palette that highlights files with the same name across multiple windows. The feature is **on by default** so you immediately see color-coded file listings; press `` ` `` (backtick) at any time (or use `View → toggle comparison colors`) to disable it if you prefer the standard theme, and tap `` ` `` again to re-enable. When enabled, goful rebuilds its comparison index and redraws each pane so you can immediately spot duplicates, newer backups, or size mismatches.

**Default palette**

- Names present in more than one pane render in bold yellow so duplicates stand out even before you inspect metadata.
- Size column colors: cyan when sizes match, red for the smallest copy, green for the largest copy.
- Time column colors: cyan for matching timestamps, red for the earliest copy, green for the latest copy.

### Digest comparison (`=` key)

`[IMPL:DIGEST_COMPARISON]` extends the comparison feature with on-demand content verification. When you need to confirm whether files with identical names **and** sizes contain the same data, press `=` on any file to calculate xxHash64 digests.

**How it works**

1. Press `=` on a file (or use `View → calculate file digest`).
2. Goful finds all files with the **same name** (case-sensitive) across loaded directory panes.
3. Files are grouped by their **actual size**:
   - Files that share the same size with at least one other file participate in digest calculation.
   - Files with a unique size (no other same-named file has that size) are skipped and continue displaying their normal size comparison colors (smallest/largest/middle).
4. Within each size group (where at least two files share the same size), goful computes xxHash64 digests and compares them:
   - **Equal digests** (identical content): the size field displays with an **underline** attribute.
   - **Different digests** (different content despite same size): the size field displays with a **strikethrough** attribute.
5. Files with unique sizes or unique names show no digest indicator—their standard comparison colors remain.

**Example**: If you have `backup.zip` in three directories with sizes 100MB, 100MB, and 150MB, pressing `=` will compare the two 100MB copies (showing underline if identical, strikethrough if different) while the 150MB copy keeps its normal "largest" color with no digest indicator.

**Use cases**

- Verify backup integrity by confirming copies match the original.
- Detect silent corruption where file sizes match but content differs.
- Quickly identify duplicates that can be safely removed.

**Notes**

- Digest calculation reads the entire file, so very large files may take noticeable time.
- Results persist until the comparison index is rebuilt (e.g., toggling comparison colors off/on or reloading directories).
- Terminal support for underline and strikethrough varies; most modern terminals (iTerm2, gnome-terminal, Windows Terminal) render both attributes correctly.

**Configure the palette**

The palette loads from a YAML file resolved in the usual precedence order: the `-compare-colors /path/to/colors.yaml` flag wins, next `GOFUL_COMPARE_COLORS=/path/to/colors.yaml`, then the default `~/.goful/compare_colors.yaml`. Edit or create the file, restart goful so the config is reloaded during startup, and press `` ` `` to see the new scheme:

```yaml
# ~/.goful/compare_colors.yaml
name:
  present: yellow
size:
  equal: cyan
  smallest: "#ff5555"   # accepts named colors or #RRGGBB hex
  largest: green
time:
  equal: cyan
  earliest: orange
  latest: "#32cd32"
```

All entries accept tcell color names (`red`, `cyan`, `magenta`, etc.) or hex values. Leave a field blank to fall back to the default. Combine this file with the runtime toggle to tailor comparison cues to whichever terminal theme you use.

### Multi-target copy/move (`C` / `M` keys)

`[REQ:NSYNC_MULTI_TARGET]` (implemented via `[ARCH:NSYNC_INTEGRATION]`) provides parallel copy/move to all visible workspace panes simultaneously using the nsync SDK.

**How it works**

1. Mark files with `Space` (or use the cursor file if nothing is marked).
2. Press `C` (Copy All) or `M` (Move All).
3. Files are copied/moved in parallel to **all other visible directory panes**.
4. Progress is displayed via goful's standard progress bar.
5. When only one pane is visible, these commands fall back to the regular single-target `c`/`m` behavior.

**Use cases**

- Deploy files to multiple directories at once (e.g., syncing assets across project folders).
- Create backups to multiple destinations simultaneously.
- Distribute configuration files to several locations.

The nsync SDK handles parallel execution, content verification, and move semantics (source deletion only after successful sync to all destinations).

> **Note**: `C`/`M` copy/move files **from the focused pane to other panes**. For synchronized operations on **same-named files across all panes**, see Sync commands below.

### Sync operations (`S` prefix) `[REQ:SYNC_COMMANDS]`

Sync commands execute copy, delete, or rename operations **simultaneously across all workspace panes** on files with the **same name** as the cursor file. This is designed for managing synchronized directory structures where you need to perform identical operations on matching files in every pane.

**How it works**

1. Position the cursor on the target file in any pane.
2. Press `S` to enter sync mode. The prompt shows: `Sync [c]opy [d]elete [r]ename [!]ignore:`
3. Optionally press `!` to toggle "ignore failures" mode (continues through all panes even on errors).
4. Press an operation key:
   - `c` – **Copy**: prompts for a new filename (default: current name); you must enter a different name. The file is copied to the new name in each pane's directory.
   - `d` – **Delete**: confirms with `y/n`; deletes the same-named file in each pane.
   - `r` – **Rename**: prompts for a new name; renames the same-named file in each pane.
5. The operation executes sequentially starting from the focused pane, then proceeding through other panes in order.

**Example**: You have `config.yaml` in three panes (`~/project-a/`, `~/project-b/`, `~/project-c/`). Press `S`, then `r`, enter `config.yaml.bak`, and all three files are renamed to `config.yaml.bak` in one action.

**Failure handling**

- By default, the operation **aborts on first failure** and reports which pane failed.
- Press `!` before the operation key to enable **ignore failures** mode, which continues through all panes and reports all failures at the end.
- If a file with the target name doesn't exist in a pane, that pane is **skipped** (not treated as a failure).

**Comparison with `C`/`M` (multi-target copy/move)**

| Feature | `C`/`M` (Multi-target) | `S` (Sync) |
|---------|------------------------|------------|
| **Source** | Marked files or cursor file in focused pane | Same-named file in **every** pane |
| **Destination** | All other visible panes | Same directory (copy to new name) or in-place (rename/delete) |
| **Use case** | Distribute files from one pane to others | Synchronized operations on matching files across panes |
| **Execution** | Parallel (nsync SDK) | Sequential |

**Use cases**

- Rename a file that exists in multiple mirrored directories with a single command.
- Delete a generated file across all project folders at once.
- Create a backup copy of a config file in every workspace pane simultaneously.

**Tip**: Combine with **linked navigation mode** (`L`) to keep panes synchronized as you navigate, then use `S` commands to batch-operate on matching files.

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
