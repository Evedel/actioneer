#!/usr/bin/env bash
set -eoux pipefail

cat <<EOF > /home/vscode/config.yaml
version: v1
actions:
  - alertname: "High Pod Memory"
    command: "kubectl delete pod ~pod -n ~namespace"
EOF

go run main.go --log-level=debug  --config-path=/home/vscode/config.yaml
