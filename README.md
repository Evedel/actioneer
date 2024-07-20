# k8s-automaton

That is a pet project to have some fun with golang, & testing, & ci/cd, & tdd.

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

## TODO:
- [ ] split into classes
    - [ ] server
        - [ ] pass args to server
    - should not return on errors inside of "for each alert in notification"
- [ ] run command only on "fired" (and parametrise the key)
- [ ] add logging for skipped alerts
- [ ] prep ci and images
    - [ ] only automaton image
    - [ ] k8s-ready image (automaton + kubectl + yq)
- [ ] toggle to never log commands (in case if there are env vars)
- [ ] add dependabot
- [ ] add templating to k8s e2e
- [ ] add delay to each action
