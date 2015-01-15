services
=====

This directory defines services as CoreOS's `fleet` units, which are systemd units with a few extra properties:
https://coreos.com/docs/launching-containers/launching/launching-containers-fleet

Conventions are to have `foo-prod.service` define how to prepare,
start, and stop the service `foo` for stage `prod`, while
`foo-prod-discovery.service` is a "sidekick service", which runs
everywhere `foo-prod` does and registers information on where to find
the service (host, port, version) under `/services/foo-prod` in etcd.
