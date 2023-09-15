#!/usr/bin/env bash
set -eoux pipefail

docker build --tag k8s-automaton:v0.1.0 .
kind load docker-image k8s-automaton:v0.1.0
kubectl apply -f "./examples/*"
kubectl rollout restart deployment/k8s-automaton -n monitoring
