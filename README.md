# k8s-automaton

[![Go Coverage](https://github.com/Evedel/automaton/wiki/coverage.svg)](https://raw.githack.com/wiki/Evedel/automaton/coverage.html)

That is a pet project to have some fun with golang, & testing, & ci/cd, & tdd.

The approach that this 'operator' enables is a clear ANTI-pattern and goes against any reasonable operational models.

However, it can also be used in a sane ways. I.e. do some real time aggregation of alerts, files report files, and store them into some external space. 

If there is an alert or problem in the cluster => root cause should be fixed and not the bandaids automated.

## Idea
[automaton/examples/k8s](https://github.com/Evedel/automaton/tree/main/examples/k8s)

Imagine a cluster with:
 - **Prometheus** -- scraping and alerting configured -- firing alerts are send to **Alertmanager**
 - **Alertmanager** -- has **k8s-automaton** set as a reciever
 - **k8s-automaton** -- container with the binary + needed tools, has a configmap with rules in the form of
    ```
    actions:
      - alertname: "Test Alert"
        command: "kubectl delete pod $pod -n $namespace"
    ```
    once the alert is sent **k8s-automaton** way -- `command` is executed

    `command` is just a shell executable, make sure that needed tools present on the docker image and you will be able to trigger execution of anything via alertmanager webhook.

## TODO:
- [ ] parametrise alertname key
- [ ] add specific lable to only run action for alerts with that label
- [ ] prep ci and images
    - [ ] only automaton image
    - [ ] k8s-ready image (automaton + kubectl + yq)
- [ ] toggle to never log commands (in case if there are env vars)
- [ ] add dependabot
- [ ] add templating to k8s e2e
- [ ] add delay to each action
