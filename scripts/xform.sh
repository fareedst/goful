#!/usr/bin/env bash
# v2026-01-06

# Transform a command like:
#   xform rsync source target1 target2
# into:
#   rsync source --to target1 --to target2
# Works as a standalone CLI or a sourced function on macOS /bin/bash 3.2+.
# [IMPL:XFORM_CLI_SCRIPT] [ARCH:XFORM_CLI_PIPELINE] [REQ:CLI_TO_CHAINING]

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
  set -euo pipefail
fi

XFORM_DRY_RUN=0
XFORM_PREFIX="--to"
XFORM_KEEP=2
XFORM_ARGS=()
XFORM_USAGE=$'Usage: xform [-p|--prefix STR] [-k|--keep N] [-n|--dry-run] command arg1 arg2 [targets ...]\n\nRewrites the command so every argument after the first N entries is\npreceded by the provided prefix (defaults shown below).\nExamples:\n  xform mv source dest1 dest2        # mv source --to dest1 --to dest2\n  xform -p --dest mv source dest1    # mv source --dest dest1\n  xform -k 3 cmd fixed1 fixed2 fixed3 t1 t2\n\nDefaults:\n  prefix: "--to"\n  keep:   2\n\nOptions:\n  -p, --prefix STR   String to insert before each target (default: --to)\n  -k, --keep N       Number of leading args to preserve (default: 2, must be >=1)\n  -n, --dry-run      Print the rewritten command with %q quoting, do not run it\n  -h, --help         Show this help text\n  --                 Stop parsing options; treat the rest as positional args\n\nRequires at least keep+1 positional arguments (command + keep untouched args + targets).'

xform::usage() {
  printf '%s\n' "$XFORM_USAGE"
}

xform::print_command() {
  local first=1
  local arg
  for arg in "$@"; do
    if (( first )); then
      printf '%q' "$arg"
      first=0
    else
      printf ' %q' "$arg"
    fi
  done
  printf '\n'
}

# Parse CLI flags without relying on Bash 4+ features.
# [IMPL:XFORM_CLI_SCRIPT] [ARCH:XFORM_CLI_PIPELINE] [REQ:CLI_TO_CHAINING]
xform::parse() {
  XFORM_DRY_RUN=0
  XFORM_PREFIX="--to"
  XFORM_KEEP=2
  while (($#)); do
    case "$1" in
      -p|--prefix)
        if (( $# < 2 )); then
          printf 'xform: --prefix requires a value\n' >&2
          xform::usage >&2
          return 64
        fi
        XFORM_PREFIX="$2"
        shift 2
        ;;
      -k|--keep)
        if (( $# < 2 )); then
          printf 'xform: --keep requires a numeric value\n' >&2
          xform::usage >&2
          return 64
        fi
        if [[ ! "$2" =~ ^[0-9]+$ ]] || (( 10#$2 < 1 )); then
          printf 'xform: --keep must be an integer >= 1 (got %s)\n' "$2" >&2
          xform::usage >&2
          return 64
        fi
        XFORM_KEEP=$((10#$2))
        shift 2
        ;;
      -n|--dry-run)
        XFORM_DRY_RUN=1
        shift
        ;;
      -h|--help)
        xform::usage
        return 65
        ;;
      --)
        shift
        break
        ;;
      -*)
        printf 'xform: unknown option "%s"\n' "$1" >&2
        xform::usage >&2
        return 64
        ;;
      *)
        break
        ;;
    esac
  done

  if (( $# <= XFORM_KEEP )); then
    printf 'xform: expected at least %d positional arguments (got %d) with keep=%d\n' "$((XFORM_KEEP + 1))" "$#" "$XFORM_KEEP" >&2
    xform::usage >&2
    return 64
  fi

  XFORM_ARGS=( "$@" )
  return 0
}

# Build the transformed argv and either print (%q) or execute it.
# [IMPL:XFORM_CLI_SCRIPT] [ARCH:XFORM_CLI_PIPELINE] [REQ:CLI_TO_CHAINING]
xform::interleave_and_run() {
  local -a argv
  argv=( "${XFORM_ARGS[@]:0:XFORM_KEEP}" )
  local i
  for (( i = XFORM_KEEP; i < ${#XFORM_ARGS[@]}; i++ )); do
    argv+=( "${XFORM_PREFIX}" "${XFORM_ARGS[i]}" )
  done

  if (( XFORM_DRY_RUN )); then
    xform::print_command "${argv[@]}"
    return 0
  fi

  "${argv[@]}"
}

xform() {
  xform::parse "$@"
  local parse_status=$?
  if (( parse_status != 0 )); then
    if (( parse_status == 65 )); then
      return 0
    fi
    return "$parse_status"
  fi
  xform::interleave_and_run
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
  xform "$@"
fi