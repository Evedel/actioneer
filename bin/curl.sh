#!/usr/bin/env bash
set -eoux pipefail

body='{"status":"firing","alerts":[{"status":"firing","labels":{"alertname":"High Pod Memory","pod":"k8s-automaton-5887dbf84b-2gr64","namespace":"monitoring"}}]}'
curl localhost:8080 -d "$body"
