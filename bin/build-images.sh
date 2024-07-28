#!/usr/bin/env bash
set -eoux pipefail
project_path="$( cd "$(dirname "$0")/../" >/dev/null 2>&1 ; pwd -P )"
cd $project_path

if [[ "$#" != "1" ]]; then
    echo "Usage: $0 <version>"
    exit 1
fi
tag=$1

docker build --target=k8s --tag actioneer:$tag .
