#!/usr/bin/env bash
set -eoux pipefail
project_path="$( cd "$(dirname "$0")/../" >/dev/null 2>&1 ; pwd -P )"
cd $project_path

if [[ "$#" != "1" ]]; then
    echo "Usage: $0 <version>"
    exit 1
fi

$project_path/bin/build-images.sh $1
kind load docker-image k8s-automaton:$1
kubectl apply -f "./examples/k8s/*"
kubectl rollout restart deployment k8s-automaton -n monitoring
kubectl rollout restart deployment alertmanager -n monitoring
