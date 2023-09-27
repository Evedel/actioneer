#!/usr/bin/env bash
set -eoux pipefail
project_path="$( cd "$(dirname "$0")/../" >/dev/null 2>&1 ; pwd -P )"

go=""

if [[ -x "$(command -v docker)" ]]; then
    echo "container: no"
    go="docker compose run --rm go"
else
    echo "container: yes"
    cd $project_path/app
    go="go"
fi

set +e
$go test -v -cover -coverprofile=coverage.out ./...
set -e

$go tool cover -html=coverage.out -o coverage.html
