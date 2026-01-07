# Architecture Overview [REQ:ARCH_DOCUMENTATION] [ARCH:DOCS_STRUCTURE] [IMPL:DOC_ARCH_GUIDE]

Goful is a terminal file manager built around Semantic Token-Driven Development (STDD). This guide records how the CLI entry point, UI widgets, file operations, and validation suites cooperate so future changes remain traceable.

## Runtime Flow [REQ:ARCH_DOCUMENTATION]

```
main.go
 └─ flag.Parse(), configpaths.Resolver.Resolve()  [REQ:CONFIGURABLE_STATE_PATHS] [ARCH:STATE_PATH_SELECTION]
     └─ config(goful, isTMUX?)
         ├─ look/menu/message/info/progress initialization
         ├─ keymap + menu wiring (filer/cmdline/finder/completion/menu) [REQ:BEHAVIOR_BASELINE]
         ├─ app.ParseStartupDirs()/SeedStartupWorkspaces() apply positional CLI directories to filer windows [REQ:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP]
         └─ goful.Run()
             ├─ widget.Init()/PollEvent() loop (tcell)
             ├─ app.Goful.Draw()/Input()/Resize orchestrate widgets
             ├─ filer.Workspace manages panes + directory state
             ├─ menu/menu.go shows pop-up menus on demand
             └─ cmdline package handles command-mode text input/history
```

## Module Directory [REQ:ARCH_DOCUMENTATION]

| Package | Responsibilities | Key Tokens |
|---------|------------------|------------|
| `main` | bootstrap, CLI flags, path resolution, keymap wiring | `[REQ:CONFIGURABLE_STATE_PATHS]` `[ARCH:STATE_PATH_SELECTION]` `[REQ:BEHAVIOR_BASELINE]` |
| `configpaths` | Pure resolver expanding overrides + defaults for persisted state/history | `[IMPL:STATE_PATH_RESOLVER]` `[ARCH:STATE_PATH_SELECTION]` |
| `app` | Application core (`Goful`, menus, async file ops, workspace orchestration) | `[REQ:MODULE_VALIDATION]` (module boundaries enforced through tests) |
| `filer` | Workspace, directories, finder modal, file manipulation flows | `[REQ:INTEGRATION_FLOWS]` `[IMPL:TEST_INTEGRATION_FLOWS]` |
| `widget` | Rendering primitives and input dispatch (keymaps, text boxes, listboxes, gauges) | `[REQ:UI_PRIMITIVE_TESTS]` |
| `cmdline` | Command-line mode textbox, history management, completion widget | `[REQ:CMD_HANDLER_TESTS]` |
| `menu` | Menu widget plus keymap injection for dynamic menus | `[REQ:BEHAVIOR_BASELINE]` |
| `message`, `progress`, `info`, `look` | Status lines, progress bars, info panel, theming | `[ARCH:DOCS_STRUCTURE]` linkage |
| `util`, `configpaths`, `info` | Misc helpers (humanized sizes, OS detection, path expansion) | `[REQ:CONFIGURABLE_STATE_PATHS]` |

## Event Loop & Modes [REQ:ARCH_DOCUMENTATION] [REQ:MODULE_VALIDATION]

- `widget.PollEvent()` produces tcell events that `app.Goful` fans out to either the current modal (`Next()`) or the primary `filer.Filer`.
- `Goful.Run()` multiplexes three asynchronous channels: UI events, interrupt queue (for async file jobs), and callback queue used by `asyncFilectrl`.
- Finder, completion, menus, and cmdline each expose `widget.Keymap` factories so they can be configured centrally in `config()`. This keeps the bindings declarative, enabling the baseline tests outlined in `[ARCH:BASELINE_CAPTURE]`.
- Resizing cascades through `Goful.Resize`, calling `Widget.Resize` on active components plus `progress`, `message`, and `info` footers.

## Persistence & Configuration [REQ:CONFIGURABLE_STATE_PATHS] [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:STATE_PATH_SELECTION] [ARCH:EXTERNAL_COMMAND_REGISTRY]

- `configpaths.Resolver` enforces precedence for state/history/commands: CLI flag (`-state`, `-history`, `-commands`) → environment (`GOFUL_STATE_PATH`, `GOFUL_HISTORY_PATH`, `GOFUL_COMMANDS_FILE`) → defaults (`~/.goful/...`).  
- `main.emitPathDebug` logs provenance for all three paths when `GOFUL_DEBUG_PATHS=1` so operators can confirm overrides without editing code.
- `externalcmd.Load` consumes the resolved commands path, parses either JSON or YAML, and falls back to baked-in defaults while logging `[IMPL:EXTERNAL_COMMAND_LOADER]` diagnostics when configs are missing or filtered. `[IMPL:EXTERNAL_COMMAND_APPEND]` ensures file-defined commands append to the compiled defaults unless `inheritDefaults: false` is supplied, so operators explicitly control whether historical shortcuts persist.
- `filer.SaveState` + `cmdline.{Load,Save}History` receive the resolved paths and are invoked before exit, ensuring persistence remains in sync with overrides, while `registerExternalCommands` wires loader output into the runtime menu.

## Menus, Keymaps, and Associations [REQ:BEHAVIOR_BASELINE] [ARCH:BASELINE_CAPTURE]

- Top-level menus (`sort`, `view`, `layout`, `stat`, `look`, `command`, `external-command`, `archive`, `bookmark`, `editor`, `image`, `media`) are declared in `main.config` with semantic tokens in surrounding comments. Each menu registers keystrokes plus callbacks that mutate `filer.Workspace`, fire shell commands, or open sub-menus.
- Default keymaps:
  - `filerKeymap` handles workspace management, navigation (`hjkl`, `C-n/C-p`), marking, finder toggles, and file operations.
  - `cmdlineKeymap`, `finderKeymap`, `completionKeymap`, and `menuKeymap` each describe chord sets for editing, history navigation, and exit semantics.
  - File-extension associations map keystrokes (`C-m` / `o`) + extension-specific behavior to actions (e.g., `tar`, `unrar`, open image/media viewers).
- Keeping bindings centralized allows the pure `KeymapBaselineSuite` to assert canonical chords are still registered even if handler implementations evolve.

## Validation & Testing Surfaces [REQ:MODULE_VALIDATION]

| Module | Validation Artifacts | Tokens |
|--------|---------------------|--------|
| `configpaths.Resolver` | `configpaths/resolver_test.go` validates precedence/expansion | `[REQ:CONFIGURABLE_STATE_PATHS]` `[IMPL:STATE_PATH_RESOLVER]` |
| `widget` primitives | `widget/listbox_test.go`, `widget/gauge.go` tests | `[REQ:UI_PRIMITIVE_TESTS]` `[IMPL:TEST_WIDGETS]` |
| `cmdline` history/completion | `cmdline/history_test.go`, completion tests (future) | `[REQ:CMD_HANDLER_TESTS]` |
| Integration flows | `filer/integration_test.go` ensures open/navigate/rename/delete flows | `[REQ:INTEGRATION_FLOWS]` |
| Keymap baselines | `main_keymap_test.go` (new) | `[REQ:BEHAVIOR_BASELINE]` `[ARCH:BASELINE_CAPTURE]` `[IMPL:BASELINE_SNAPSHOTS]` `[TEST:KEYMAP_BASELINE]` |
| CI gates | `.github/workflows/ci.yml` runs fmt/vet/test/staticcheck/race + token validation | `[REQ:CI_PIPELINE_CORE]` `[REQ:STATIC_ANALYSIS]` `[REQ:RACE_TESTING]` |

Each suite is designed to validate its module independently before integration, fulfilling `[REQ:MODULE_VALIDATION]`.

## Change Impact Guidance [REQ:ARCH_DOCUMENTATION]

- **UI/Keymap changes**: Update `ARCHITECTURE.md`, `CONTRIBUTING.md`, relevant requirements, and extend `KeymapBaselineSuite`.
- **Persistence changes**: Update `configpaths`, document new precedence rules here, and add tests covering overrides.
- **Menu/menu additions**: Document the new menu flow under Menus section and add baseline assertions for new entry keys.
- **Feature toggles / debug policy**: Keep `CONTRIBUTING.md` and this file in sync so developers understand expectations about diagnostic logging.

Maintaining this document alongside code/tests ensures future contributors can trace `[REQ:*] → [ARCH:*] → [IMPL:*] → tests` without spelunking through source files first.

