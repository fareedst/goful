#!/usr/bin/env bash

# [IMPL:TOKEN_VALIDATION_SCRIPT] [ARCH:TOKEN_VALIDATION_AUTOMATION] [REQ:STDD_SETUP]
# Validates that every semantic token used in tracked source files exists in
# `stdd/semantic-tokens.md`. Extend the default globs or pass explicit paths to
# scan additional files (e.g., docs) once placeholders are removed.

set -euo pipefail

if ! command -v git >/dev/null 2>&1; then
  echo "ERROR: git is required for [PROC:TOKEN_VALIDATION]" >&2
  exit 1
fi

if ! command -v rg >/dev/null 2>&1; then
  echo "ERROR: ripgrep (rg) is required for [PROC:TOKEN_VALIDATION]" >&2
  exit 1
fi

repo_root=$(git rev-parse --show-toplevel 2>/dev/null || true)
if [[ -z "${repo_root}" ]]; then
  echo "ERROR: Unable to determine repository root via git rev-parse" >&2
  exit 1
fi

cd "${repo_root}"

tokens_file="stdd/semantic-tokens.md"
if [[ ! -f "${tokens_file}" ]]; then
  echo "ERROR: Expected token registry ${tokens_file} was not found" >&2
  exit 1
fi

match_pattern='\[[A-Z]+:[A-Z0-9_]+\]'

# Build registry from semantic-tokens.md
mapfile -t registry_tokens < <(rg -o "${match_pattern}" "${tokens_file}" | sort -u)
if [[ ${#registry_tokens[@]} -eq 0 ]]; then
  echo "ERROR: No semantic tokens registered in ${tokens_file}" >&2
  exit 1
fi

declare -A registry_lookup=()
for token in "${registry_tokens[@]}"; do
  registry_lookup["${token}"]=1
done

declare -a search_targets=()
if [[ $# -gt 0 ]]; then
  search_targets=("$@")
else
  # Default to tracked source files where placeholder tokens are not expected.
  mapfile -d '' search_targets < <(git ls-files -z -- '*.go' '*.mod' '*.sum' '*.sh' '*.yml' '*.yaml' 'Makefile')
fi

if [[ ${#search_targets[@]} -eq 0 ]]; then
  echo "DEBUG: No files matched the default search globs; skipping token validation." >&2
  exit 0
fi

check_refs=0
declare -A missing_map=()

while IFS=: read -r file line token; do
  [[ -z "${file}" || -z "${token}" ]] && continue
  ((++check_refs))
  if [[ -z "${registry_lookup["${token}"]:-}" ]]; then
    missing_map["${token}"]+="${file}:${line}"$'\n'
  fi
done < <(rg --color=never --no-heading --with-filename -n -o "${match_pattern}" "${search_targets[@]}" 2>/dev/null || true)

if [[ ${#missing_map[@]} -gt 0 ]]; then
  echo "ERROR: [PROC:TOKEN_VALIDATION] detected tokens missing from ${tokens_file}:" >&2
  for token in "${!missing_map[@]}"; do
    echo "  ${token}" >&2
    while IFS= read -r loc; do
      [[ -z "${loc}" ]] && continue
      echo "    - ${loc}" >&2
    done <<< "${missing_map[${token}]}"
  done
  exit 1
fi

echo "DIAGNOSTIC: [PROC:TOKEN_VALIDATION] verified ${check_refs} token references across ${#search_targets[@]} files."
exit 0

