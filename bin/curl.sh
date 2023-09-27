#!/usr/bin/env bash
set -eoux pipefail

pod_name=$(kubectl get pods -n monitoring | grep test-deployment | awk '{print $1}' | head -1)
body='{"status":"firing","alerts":[{"status":"firing","labels":{"alertname":"High Pod Memory","pod":"'$pod_name'","namespace":"monitoring"}}]}'
curl localhost:8080 -d "$body"
