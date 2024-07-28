# actioneer

[![Go Coverage](https://github.com/Evedel/actioneer/wiki/coverage.svg)](https://raw.githack.com/wiki/Evedel/actioneer/coverage.html)

That is a pet project to have some fun with golang, & testing, & ci/cd, & tdd.

The approach that this 'operator' enables is a clear ANTI-pattern and goes against any reasonable operational models.

If there is an alert or problem in the cluster => the root cause should be fixed and not the bandaids automated.

However, it can also be used in a sane way. I.e. do some real-time aggregation of alerts, report combination, wire notification into another external store not easily reachable. 

## Idea
[actioneer/examples/k8s](https://github.com/Evedel/actioneer/tree/main/examples/k8s)

Imagine a cluster with:
 - **Prometheus** -- scraping and alerting configured -- firing alerts are send to **Alertmanager**
    ```
    - name: test alert
      rules:
      - alert: Test Alert
        expr: max(container_memory_usage_bytes{pod=~"test-deployment.*",namespace!~"kube-system"}) by (pod,namespace) > 0
        for: 10m
    ```
 - **Alertmanager** -- has **actioneer** set as a reciever
    ```
    receivers:
      - name: 'actioneer'
        webhook_configs:
          - url: 'http://actioneer:8080/'
    ```
 - **actioneer** -- container with the binary + needed tools, has a configmap with rules in the form of
    ```
    actions:
      - name: "Restart Flake Pod"
        alertname: "Alert That Is Triggered Without Reason And Fixed By Pod Restart"
        command: "kubectl delete pod $pod -n $namespace"
    ```
    once the alert is sent **actioneer** way -- `command` is executed

    `command` is just a shell executable, make sure that needed tools present on the docker image and you will be able to trigger execution of anything via alertmanager webhook.
    