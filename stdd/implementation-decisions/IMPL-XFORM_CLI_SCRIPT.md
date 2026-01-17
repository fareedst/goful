# [IMPL:XFORM_CLI_SCRIPT] Xform CLI Script

**Cross-References**: [ARCH:XFORM_CLI_PIPELINE] [REQ:CLI_TO_CHAINING]  
**Status**: Active  
**Created**: 2026-01-06  
**Last Updated**: 2026-01-17

---

## Decision

Implement `scripts/xform.sh` as a portable Bash helper that can be executed directly or sourced, exposing a `xform` function that inserts `--to` before every destination argument while offering a dry-run preview.

## Rationale

- Provides a single, well-tested place to handle argument parsing and quoting for workflows that need the `--to <target>` pattern repeated for multiple paths
- Keeps the helper compatible with macOS `/bin/bash` 3.2 by avoiding Bash 4+ features (no associative arrays or `local -n`) so contributors can run it without Homebrew Bash

## Implementation Approach

- `scripts/xform.sh` starts with `set -euo pipefail` and defines:
  - `xform::usage` — prints help text and exits with status 64 when invoked incorrectly
  - `xform::parse` — consumes `-p/--prefix`, `-k/--keep`, `-n/--dry-run`, `-h/--help`, and `--` flags, ensures at least `keep + 1` positional arguments remain, and exports globals (`XFORM_PREFIX`, `XFORM_KEEP`, `XFORM_DRY_RUN`, `XFORM_ARGS`) for downstream logic. Defaults: prefix `--to`, keep `2`
  - `xform::run` — builds the transformed argv array by preserving the first `keep` positional arguments and interleaving `<prefix>` between the remaining entries. When `dry-run` is enabled it prints the `%q`-formatted command; otherwise it executes the new argv and propagates the exit code
  - `xform` — public wrapper that calls `xform::parse` followed by `xform::run`
- The script checks `[[ "${BASH_SOURCE[0]}" == "$0" ]]` to decide whether to execute immediately or just define the function for callers who `source` it
- `scripts/xform_test.sh` sources the helper and runs two module-validation suites:
  - Parser tests feed different flag combinations (custom prefix/keep along with error paths) and assert correct exit codes/messages for insufficient arguments
  - Builder tests run `xform -n ...` and check the printed, quoted command to ensure interleaving/logging semantics work for arguments with spaces and custom prefixes/keep windows

## Code Markers

- `scripts/xform.sh` contains `# [IMPL:XFORM_CLI_SCRIPT] [ARCH:XFORM_CLI_PIPELINE] [REQ:CLI_TO_CHAINING]` comments near the parser and builder functions
- `scripts/xform_test.sh` comments reference `[REQ:CLI_TO_CHAINING]` when asserting dry-run output and error handling

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `scripts/xform.sh` - parser and builder functions
- [ ] `scripts/xform_test.sh` - test functions
- [ ] `scripts/xform.bats` - Bats specs

Tests that must reference `[REQ:CLI_TO_CHAINING]`:
- [ ] `test_dry_run_inserts_targets_REQ_CLI_TO_CHAINING`
- [ ] Bats: `dry-run uses default prefix [REQ:CLI_TO_CHAINING]`
- [ ] Bats: `invalid keep value fails with guidance [REQ:CLI_TO_CHAINING]`

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-06 | — | ✅ Pass | `bash scripts/xform_test.sh` (macOS 15.1 arm64) |
| 2026-01-06 | — | ✅ Pass | `./scripts/validate_tokens.sh` → verified 269 token references across 55 files |

## Related Decisions

- Depends on: —
- See also: [ARCH:XFORM_CLI_PIPELINE], [REQ:CLI_TO_CHAINING], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
