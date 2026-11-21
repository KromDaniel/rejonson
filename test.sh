#!/bin/bash
# Runs `go vet` and `go test` for every module.
#
# With GOTOOLCHAIN=local (as in CI), modules whose `go` directive is newer than
# the installed toolchain are skipped instead of failing, e.g. v9 requires the
# Go version that go-redis/v9 requires.

set -uo pipefail

packages=("." "v7" "v8" "v9")
script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
current=$(go env GOVERSION | sed 's/^go//')

version_ge() { # $1 >= $2
  [ "$(printf '%s\n%s\n' "$2" "$1" | sort -V | head -n1)" = "$2" ]
}

for pkg in "${packages[@]}"; do
  cd "$script_dir/$pkg" || exit 1
  required=$(sed -n 's/^go \([0-9.]*\)$/\1/p' go.mod)

  if [ "${GOTOOLCHAIN:-}" = "local" ] && ! version_ge "$current" "$required"; then
    echo "skip  $pkg (requires go $required, have go $current)"
    continue
  fi

  echo "==> $pkg"
  go vet ./... || exit 1
  go test -race -count=1 ./... || exit 1
done
