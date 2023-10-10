# k8s-automaton

TODO:
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
