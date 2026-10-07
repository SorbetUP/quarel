# Production container runtime target

The planned production target is **rootless Podman 6.1.3 on Linux amd64 and
arm64**, with **podman-compose 1.6.0** as its explicit Compose provider. Pin the
provider in the host configuration (`PODMAN_COMPOSE_PROVIDER=podman-compose`)
so an installed `docker-compose` does not silently take precedence. Install the
exact host and provider versions through the production image/package lock
before the M17 deployment gate. The supported Podman release is published at
the upstream [Podman releases page](https://github.com/podman-container-tools/podman/releases);
provider release notes are at
[podman-compose releases](https://github.com/containers/podman-compose/releases).

Podman documents `podman compose` as a wrapper around an external provider and
documents provider selection through `PODMAN_COMPOSE_PROVIDER` or
`containers.conf`. That makes both runtime and provider versions part of the
deployment contract ([official Podman compose manual](https://docs.podman.io/en/stable/markdown/podman-compose.1.html)).

This selection is a target, not a compatibility claim. This development host
has Docker/Colima and no Podman binary. Before M17/M19, install the selected
versions on the actual Linux target and use the same Compose file to verify
start/stop and signal handling, readiness, persistence and restore, upgrades,
secret and volume ownership, public TCP 443, UDP media, source-IP handling,
and voice join/media. Record the exact host kernel, distribution, provider
binary/package hash, image architecture and test result there.

The current local Compose stack explicitly disables UPnP and embedded voice,
binds development HTTP ports to loopback, and does not test public media
networking. The current server Dockerfile still embeds LiveKit until M08.
Rootless binding of 443, media UDP behavior, webhook source-IP preservation,
and the later reverse-proxy/LiveKit split remain deployment design items; a
successful Compose parse does not clear them.
