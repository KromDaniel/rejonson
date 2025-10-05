#!/bin/bash

packages=("." "v7" "v8" "v9")
script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

for pkg in "${packages[@]}"; do
  cd "$script_dir/$pkg"
  go test . || exit 1
done