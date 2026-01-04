# Contributing Guide [REQ:CONTRIBUTING_GUIDE] [ARCH:CONTRIBUTION_PROCESS] [IMPL:DOC_CONTRIBUTING]

Thanks for helping improve goful! This guide captures the workflow expectations mandated by Semantic Token-Driven Development (STDD).

## Quickstart [REQ:CONTRIBUTING_GUIDE]

1. Install Go `1.24.x` (we track LTS per `[REQ:GO_TOOLCHAIN_LTS]`).
2. Clone the repo and install dependencies:
   ```bash
   git clone https://github.com/anmitsu/goful.git
   cd goful
   go mod download
   ```
3. Build/Test helpers:
   ```bash
   make lint       # gofmt + go vet
   make test       # go test ./...
   go test ./...   # preferred during development
   ./scripts/validate_tokens.sh
   make release PLATFORM=$(go env GOOS)/$(go env GOARCH)  # optional sanity check for [REQ:RELEASE_BUILD_MATRIX]
   ```

## Development Workflow [REQ:CONTRIBUTING_GUIDE] [PROC:TOKEN_AUDIT]

1. **Plan with STDD docs**: update `stdd/requirements.md`, `architecture-decisions.md`, `implementation-decisions.md`, and `stdd/tasks.md` **before** coding.
2. **Code changes**:
   - Keep modules isolated and validated separately per `[REQ:MODULE_VALIDATION]`.
   - Add semantic tokens to new code/comments/tests (`[IMPL:*] [ARCH:*] [REQ:*]` triplets).
3. **Local verification**:
   ```bash
   go fmt ./...
   go vet ./...
   go test ./...
   ./scripts/validate_tokens.sh
   ```
4. **CI mirrors workflow**: GitHub Actions runs fmt/vet/test, staticcheck, race tests, and token validation; ensure local runs are clean first.
5. **Release hygiene**: When tagging (`git tag vX.Y.Z && git push origin vX.Y.Z`), `.github/workflows/release.yml` rebuilds the same matrix via `make release PLATFORM=os/arch` and publishes the binaries + `.sha256` digests to the GitHub Release (`[REQ:RELEASE_BUILD_MATRIX]`). Confirm the workflow succeeds before announcing the release.

## Semantic Token Discipline [REQ:CONTRIBUTING_GUIDE] [PROC:TOKEN_AUDIT] [PROC:TOKEN_VALIDATION]

- Define `[REQ:*]` tokens in `stdd/requirements.md` before referencing them elsewhere.
- Record architecture decisions (`[ARCH:*]`) and implementation details (`[IMPL:*]`) immediately.
- Reference tokens in:
  - Code comments (`// [IMPL:STATE_PATH_RESOLVER] [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS]`).
  - Test names/comments (`TestResolvePaths_REQ_CONFIGURABLE_STATE_PATHS`).
  - Documentation (this file, `ARCHITECTURE.md`, README).
- Run `./scripts/validate_tokens.sh` after every substantive change and paste the resulting `DIAGNOSTIC: [PROC:TOKEN_VALIDATION] ...` line into your task log or PR description.

## Module Validation & Testing [REQ:MODULE_VALIDATION] [REQ:BEHAVIOR_BASELINE]

- Identify modules and validation criteria upfront (e.g., `KeymapBaselineSuite`, `widget.ListBox`, `cmdline.History`).
- Validate modules independently with unit tests, mocks, or snapshot fixtures before integrating them.
- Keep baseline coverage up to date:
  - `widget`/`filer` unit tests guard rendering and filesystem helpers.
  - `cmdline` tests cover history/completion behavior.
  - `filer/integration_test.go` verifies open/navigate/rename/delete flows.
  - `main_keymap_test.go` (`[TEST:KEYMAP_BASELINE]`) snapshots keybindings so we notice remaps immediately.
- Document validation outcomes inside `stdd/tasks.md` and `implementation-decisions.md`.

## Debug Logging Policy [REQ:CONTRIBUTING_GUIDE]

- Retain `DEBUG:`, `TRACE:`, and `DIAGNOSTIC:` logs that explain architectural or implementation decisions—these are treated as inline documentation.
- Guard logs with env variables or debug flags if necessary, but **do not remove** them unless explicitly directed.
- When adding new diagnostics, include semantic tokens referencing the behavior being explained.

## Review Checklist [REQ:CONTRIBUTING_GUIDE]

Before opening a PR:

- [ ] Requirements/architecture/implementation docs updated with new tokens and decisions.
- [ ] `ARCHITECTURE.md`, `CONTRIBUTING.md`, or README cross-links added/updated when relevant.
- [ ] `make fmt` and `make test` pass; `go test ./...` produces clean output.
- [ ] `./scripts/validate_tokens.sh` succeeded (paste the diagnostic line into the PR/task).
- [ ] Module validation evidence cited (tests or docs) for every touched module.
- [ ] Debug logs retained or added to explain new decision points.

Following this checklist keeps `[REQ:CONTRIBUTING_GUIDE]` satisfied and accelerates reviews.

