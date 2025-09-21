#!/bin/sh

set -e

cfg="$PWD/.golangci.yaml"
modules=$(find . -name "go.mod" -exec dirname {} \;)

for mod in $modules; do
  echo "Linting $mod"
  
  (
  cd "$mod"
  golangci-lint run --fix --config "$cfg"
  )
done
