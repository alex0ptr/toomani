#!/usr/bin/env bash

set -euo pipefail

declare -a repositories=(
  "test-space-1/test;test-space-ssh-1;test-space-http-1"
  "test-space-2/test;test-space-ssh-2;test-space-http-2"
)

preCommitAvailable="$(command -v pre-commit || echo "$?")"

# clones all repositories in the list to the current directory
function clone() {
  for repository in "${repositories[@]}"
  do
    path="$(cut -d ';' -f1 <<< "$repository")"
    url="$(cut -d ';' -f2 <<< "$repository")"
    mkdir -p "$path"
    git clone "$url" "$path"
    if [[ "$preCommitAvailable" == "1" && -f ".pre-commit-config.yaml" ]]; then
      eval "cd \"$path\" && pre-commit install"
    fi
  done
}

"$@"