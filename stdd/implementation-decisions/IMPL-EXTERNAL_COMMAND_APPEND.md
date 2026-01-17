# [IMPL:EXTERNAL_COMMAND_APPEND] External Command Append Toggle

**Cross-References**: [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]  
**Status**: Active  
**Created**: 2026-01-02  
**Last Updated**: 2026-01-17

---

## Decision

Preserve built-in Windows/POSIX menu entries whenever a commands file is present, **prepending** file-defined entries by default and only replacing defaults when the file opts out via `inheritDefaults: false`.

## Rationale

- Operators expect historical shortcuts (`cp`, `mv`, etc.) to remain available unless they explicitly suppress them; this keeps onboarding friction low
- Some environments must ship a clean slate for security reasons, so the same config file needs a deterministic switch to drop defaults entirely
- Encoding the behavior in semantic tokens makes the inheritance contract testable and discoverable beyond requirements prose

## Implementation Approach

- Extend the loader parser to recognize either an array of entries (prepends defaults implicitly) or an object wrapper containing `commands` and `inheritDefaults` (JSON or YAML). Missing flags default to `true` so existing configs pick up prepend semantics automatically
- After sanitizing file entries, merge them with `externalcmd.Defaults` when inheritance is enabled (custom entries first) or return only the sanitized entries when disabled. Emit `DEBUG: [IMPL:EXTERNAL_COMMAND_APPEND]` logs describing whether defaults were included
- Surface the new `[IMPL:EXTERNAL_COMMAND_APPEND]` token in loader code comments, docs, and tests so audits can trace the behavior end-to-end

## Code Markers

- `externalcmd/loader.go` merge logic and debug output include `[IMPL:EXTERNAL_COMMAND_APPEND] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]`
- README + STDD docs document the `inheritDefaults` flag and reference this token

## Token Coverage `[PROC:TOKEN_AUDIT]`

Files/functions that must carry annotations:
- [ ] `externalcmd/loader.go` - merge logic with `[IMPL:EXTERNAL_COMMAND_APPEND]`

Tests that must reference `[REQ:EXTERNAL_COMMAND_CONFIG]`:
- [ ] `TestLoadAppendsDefaultsByDefault_REQ_EXTERNAL_COMMAND_CONFIG` - includes `[IMPL:EXTERNAL_COMMAND_APPEND]` in comments
- [ ] `TestLoadCanDisableDefaults_REQ_EXTERNAL_COMMAND_CONFIG` - includes `[IMPL:EXTERNAL_COMMAND_APPEND]` in comments

## Validation Evidence `[PROC:TOKEN_VALIDATION]`

| Date | Commit | Validation Result | Notes |
|------|--------|-------------------|-------|
| 2026-01-02 | — | ✅ Pass | `go test ./externalcmd` covers prepend vs. replace behaviors |

## Related Decisions

- Depends on: [IMPL:EXTERNAL_COMMAND_LOADER]
- See also: [ARCH:EXTERNAL_COMMAND_REGISTRY], [REQ:EXTERNAL_COMMAND_CONFIG], [REQ:MODULE_VALIDATION]

---

*Migrated from monolithic implementation-decisions.md on 2026-01-17*
