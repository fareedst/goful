#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
xform_script="${script_dir}/xform.sh"
tmpdir="$(mktemp -d)"
trap 'rm -rf "${tmpdir}"' EXIT

tests_run=0
failures=0

run_test() {
  local name="$1"
  shift
  if "$@"; then
    printf 'ok - %s\n' "$name"
  else
    printf 'not ok - %s\n' "$name"
    failures=$((failures + 1))
  fi
  tests_run=$((tests_run + 1))
}

# Validate dry-run output builds the rewritten argv with --to between each target.
# [REQ:CLI_TO_CHAINING]
test_dry_run_inserts_targets_REQ_CLI_TO_CHAINING() {
  local output
  output="$("$xform_script" -n rsync src hostA hostB)"
  [[ "$output" == "rsync src --to hostA --to hostB" ]]
}

# Ensure quoting is preserved for targets containing spaces when previewing.
# [REQ:CLI_TO_CHAINING]
test_dry_run_handles_spaces_REQ_CLI_TO_CHAINING() {
  local output
  output="$("$xform_script" -n scp "My File.txt" "Remote Host" "/tmp/Backup Folder")"
  [[ "$output" == "scp My\\ File.txt --to Remote\\ Host --to /tmp/Backup\\ Folder" ]]
}

# Fewer than three positional arguments should exit with code 64 and print usage text.
# [REQ:CLI_TO_CHAINING]
test_invalid_arg_count_REQ_CLI_TO_CHAINING() {
  local stdout_file="${tmpdir}/stdout"
  local stderr_file="${tmpdir}/stderr"
  local status
  if "$xform_script" scp onlyTwo >"$stdout_file" 2>"$stderr_file"; then
    return 1
  else
    status=$?
  fi
  [[ $status -eq 64 ]] || return 1
  grep -q 'expected at least 3 positional arguments' "$stderr_file"
}

# Custom prefix should replace the default when provided.
# [REQ:CLI_TO_CHAINING]
test_custom_prefix_REQ_CLI_TO_CHAINING() {
  local output
  output="$("$xform_script" -n -p --dest rsync src hostA hostB)"
  [[ "$output" == "rsync src --dest hostA --dest hostB" ]]
}

# Keep count determines how many leading args remain untouched.
# [REQ:CLI_TO_CHAINING]
test_custom_keep_REQ_CLI_TO_CHAINING() {
  local output
  output="$("$xform_script" -n -k 3 deploy key1 key2 hostA hostB)"
  [[ "$output" == "deploy key1 key2 --to hostA --to hostB" ]]
}

# Invalid keep values should fail with an actionable message.
# [REQ:CLI_TO_CHAINING]
test_invalid_keep_value_REQ_CLI_TO_CHAINING() {
  local stderr_file="${tmpdir}/stderr_keep"
  local status
  if "$xform_script" -k 0 mv a b > /dev/null 2>"$stderr_file"; then
    return 1
  else
    status=$?
  fi
  [[ $status -eq 64 ]] || return 1
  grep -q 'must be an integer >= 1' "$stderr_file"
}

main() {
  run_test test_dry_run_inserts_targets_REQ_CLI_TO_CHAINING test_dry_run_inserts_targets_REQ_CLI_TO_CHAINING
  run_test test_dry_run_handles_spaces_REQ_CLI_TO_CHAINING test_dry_run_handles_spaces_REQ_CLI_TO_CHAINING
  run_test test_invalid_arg_count_REQ_CLI_TO_CHAINING test_invalid_arg_count_REQ_CLI_TO_CHAINING
  run_test test_custom_prefix_REQ_CLI_TO_CHAINING test_custom_prefix_REQ_CLI_TO_CHAINING
  run_test test_custom_keep_REQ_CLI_TO_CHAINING test_custom_keep_REQ_CLI_TO_CHAINING
  run_test test_invalid_keep_value_REQ_CLI_TO_CHAINING test_invalid_keep_value_REQ_CLI_TO_CHAINING

  printf '%d tests run, %d failures\n' "$tests_run" "$failures"
  exit "$failures"
}

main "$@"

