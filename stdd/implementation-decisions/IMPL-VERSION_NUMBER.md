# [IMPL:VERSION_NUMBER] Version Number Display Implementation

**Cross-References**: [ARCH:VERSION_DISPLAY] [REQ:VERSION_NUMBER]  
**Status**: Active  
**Created**: 2026-01-20  
**Last Updated**: 2026-01-20

---

## Decision

Define version constant (1.0.0), add `--version` flag, override `flag.Usage()` to include version in CLI help, and add version display to help popup.

## Rationale

- Single version constant prevents version drift between different display locations
- `--version` flag enables version checking without TUI startup overhead
- CLI help display provides version information for command-line users
- Help popup display provides version information for TUI users
- Consistent implementation across CLI and TUI contexts improves user experience

## Implementation Approach

### Define Version Constant in `main.go`

- Add `const version = "1.0.0"` near the top of the file
- This serves as the single source of truth for the version number

### Add `--version` Flag

- Add `versionFlag` boolean flag using `flag.Bool("version", false, "Print version and exit")`
- In `main()`, check the flag immediately after `flag.Parse()`
- If `--version` is set, print version to stdout and exit with code 0
- This check must occur before any TUI initialization to avoid unnecessary startup

### Override `flag.Usage()` Function

- Define custom `usage()` function that includes version at the top
- Format: "goful version 1.0.0" followed by usage instructions and flag descriptions
- Assign to `flag.Usage` before `flag.Parse()` is called
- This ensures version appears when users run `goful --help` or `goful -h`

### Add Version to Help Popup

- Modify `help/help.go` to include version in the keystroke catalog
- Add version entry to `keystrokeCatalog` slice in the "=== Application ===" section
- Format: "Version: 1.0.0" or similar
- This makes version visible when users press `?` key in the TUI

## Code Markers

- `main.go` includes `// [IMPL:VERSION_NUMBER] [ARCH:VERSION_DISPLAY] [REQ:VERSION_NUMBER]` comments for:
  - Version constant definition
  - `--version` flag definition
  - Version flag check in `main()`
  - `flag.Usage()` override function
- `help/help.go` includes `// [IMPL:VERSION_NUMBER] [ARCH:VERSION_DISPLAY] [REQ:VERSION_NUMBER]` comment for version entry in catalog

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [x] `main.go` - version constant, `versionFlag`, `usage()` function, version check in `main()`
- [x] `help/help.go` - version entry in `keystrokeCatalog`

Tests that must reference `[REQ:VERSION_NUMBER]`:
- [x] Manual verification for CLI help output
- [x] Manual verification for help popup display
- [x] Manual verification for `--version` flag behavior

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-20 | — | ✅ Pass | Implementation complete, all features verified |

## Related Decisions

- Depends on: —
- See also: [ARCH:VERSION_DISPLAY], [REQ:VERSION_NUMBER]

---

*Created for version number display feature*
