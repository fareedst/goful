# [IMPL:COMPARE_COLOR_CONFIG] Comparison Color Configuration

**Cross-References**: [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]  
**Status**: Active  
**Created**: 2026-01-09  
**Last Updated**: 2026-01-17

---

## Decision

Load comparison color scheme from YAML with sensible defaults for missing/invalid configs.

## Rationale

- Allows users to customize colors to match their terminal themes without code changes
- Provides consistent defaults that work on both light and dark terminals
- Reuses the existing path resolver pattern for flag/env/default precedence

## Implementation Approach

Add `filer/comparecolors/` package with:

- `type Config struct` containing color definitions for each comparison state (NamePresent, SizeEqual, SizeSmallest, SizeLargest, TimeEqual, TimeEarliest, TimeLatest)
- `func Load(path string) (*Config, error)` that reads YAML, validates color names, and returns parsed config
- `func DefaultConfig() *Config` providing sensible defaults when file is missing
- Color names map to `tcell.Color` values via a lookup table supporting named colors ("red", "green", etc.) and hex codes

Extend `configpaths.Resolver` with `-compare-colors` flag, `GOFUL_COMPARE_COLORS` env, and default `~/.goful/compare_colors.yaml`.

`main.go` loads config at startup, passes to `look.ConfigureComparisonColors()`.

## Code Markers

- `filer/comparecolors/config.go` includes `[IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `filer/comparecolors/config.go` - all functions
- [ ] `filer/comparecolors/config_test.go` - tests

Tests that must reference `[REQ:FILE_COMPARISON_COLORS]`:
- [ ] YAML parsing tests
- [ ] Default config tests
- [ ] Invalid input handling tests

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-09 | — | ✅ Pass | `go test ./filer/comparecolors` |

## Related Decisions

- Depends on: —
- See also: [ARCH:FILE_COMPARISON_ENGINE], [REQ:FILE_COMPARISON_COLORS], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
