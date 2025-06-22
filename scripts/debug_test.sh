#!/usr/bin/env bash

set -euo pipefail

input="$1"

# Expect input in the form "model::TestTimelineOffsetCount"
if [[ "$input" != *::* ]]; then
    echo "Expected input in the form 'package::TestName', got: $input"
    exit 1
fi

pkg="${input%%::*}"
test_name="${input##*::}"

case "$pkg" in
    model) path="./internal/model" ;;
    controller) path="./internal/controller" ;;
    api) path="./internal/api" ;;
    *) echo "Unknown package: $pkg" && exit 1 ;;
esac

# Run delve with the extracted path and test name
dlv test "$path" -- -test.run "^$test_name"
