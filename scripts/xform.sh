# Transform N>=2 args from:
#   xform a1 a2 a3
# into:
#   func2 a1 --to a2 --to a3
#
# Notes:
# - Requires Bash (works on macOS /bin/bash 3.2).
# - Preserves arguments exactly (spaces/globs/etc.) via quoting and arrays.
# xform () { 
set -x

  if (( $# < 2 )); then
    printf 'usage: %s <a1> <a2> [a3 ...]\n' "${FUNCNAME[0]}" >&2
    return 64
  fi

  declare -a argv
  argv=( "$1" )
  shift

  while (( $# )); do
    argv+=( --to "$1" )
    shift
  done

  # func2 -- "${argv[@]}"
  /Users/fareed/Documents/dev/go/goful/scripts/log.sh "${argv[@]}"
# }
# xform "$@"