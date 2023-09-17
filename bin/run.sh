#!/usr/bin/env bash
set -eoux pipefail
project_path="$( cd "$(dirname "$0")/../" >/dev/null 2>&1 ; pwd -P )"

cat <<EOF > /home/vscode/config.yaml
version: v1
actions:
  - alertname: "High Pod Memory"
    command: "kubectl delete pod ~pod -n ~namespace"
EOF

cd $project_path/app; go run cmd/main.go --log-level=debug  --config-path=/home/vscode/config.yaml
