# [IMPL:LINKED_NAVIGATION_AUTO_DISABLE] Linked Navigation Auto-Disable

**Cross-References**: [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]  
**Status**: Active  
**Created**: 2026-01-11  
**Last Updated**: 2026-01-17

---

## Decision

Automatically disable linked navigation when subdirectory navigation partially fails across windows.

## Rationale

- When linked navigation is enabled and the user enters a subdirectory that doesn't exist in all windows, the directory structures have diverged
- Keeping linked mode enabled after divergence breaks the mental model of synchronized navigation - windows are now at different depths in their respective hierarchies
- Automatic disabling with a user message signals the divergence clearly and prevents confusion
- The user can re-enable linked mode manually if they want to continue synchronized navigation from the new state

## Implementation Approach

### Navigation Result Tracking (`filer/workspace.go`)

- Modify `ChdirAllToSubdirNoRebuild(name string)` to return `(navigated, skipped int)` counts
- `navigated`: number of non-focused windows that successfully navigated to the subdirectory
- `skipped`: number of non-focused windows where the subdirectory does not exist

### Direct State Setter (`app/goful.go`)

- Add `func (g *Goful) SetLinkedNav(enabled bool)` to directly set the linked navigation state
- This allows disabling without toggling (clearer intent than double-toggle)

### Auto-Disable Logic (`main.go`)

- In `linkedEnterDir` helper, check the results from `ChdirAllToSubdirNoRebuild`
- If `skipped > 0` (any window couldn't navigate), disable linked mode via `SetLinkedNav(false)`
- Display message: "linked navigation disabled: N window(s) missing 'dirname'"

### Edge Cases

- **Single window**: No other windows to check, linked mode remains enabled
- **All windows succeed**: Linked mode remains enabled
- **All other windows fail**: Linked mode disabled (focused window navigates alone)
- **Partial success**: Linked mode disabled with message

## Code Markers

- `filer/workspace.go`: `// [IMPL:LINKED_NAVIGATION_AUTO_DISABLE] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]`
- `app/goful.go`: `// [IMPL:LINKED_NAVIGATION_AUTO_DISABLE] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]`
- `main.go`: `// [IMPL:LINKED_NAVIGATION_AUTO_DISABLE] [ARCH:LINKED_NAVIGATION] [REQ:LINKED_NAVIGATION]`

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `filer/workspace.go` - ChdirAllToSubdirNoRebuild
- [ ] `app/goful.go` - SetLinkedNav
- [ ] `main.go` - linkedEnterDir helper

Tests that must reference `[REQ:LINKED_NAVIGATION]`:
- [ ] Auto-disable behavior tests

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-11 | — | ✅ Pass | Auto-disable working |

## Related Decisions

- Depends on: [IMPL:LINKED_NAVIGATION]
- See also: [ARCH:LINKED_NAVIGATION], [REQ:LINKED_NAVIGATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
