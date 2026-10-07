# Local OCI stack

The Compose file runs the current server products and the documentation site on
loopback ports. It disables UPnP, TURN and embedded voice so a local start never
changes router state or binds public media ports. Identity and community data
live in named volumes. This is a development stack, not the production topology.

From this directory:

```sh
docker compose -f compose.yaml build
docker compose -f compose.yaml up --wait
docker compose -f compose.yaml ps
```

Open the site at `http://127.0.0.1:18000/`. The Identity API, community API,
and docs MCP health endpoints are available on ports 18080, 18090 and 18095.
`docker compose ... ps` reports each service's health. Stop services with `down`;
named data volumes remain. Use `down --volumes` only when intentionally deleting
local Identity and community data.

The dev image is built from the checked out source and includes pinned Go and
Node bases, SQLite, Python cryptography, and client/site lockfile dependencies.
It builds/tests the desktop source without downloading the native Electron
runtime; release/E2E artifacts remain handled by their dedicated build path.
Run checks in that image with:

```sh
docker compose -f compose.yaml --profile dev build dev
docker compose -f compose.yaml --profile dev run --rm dev make test
docker compose -f compose.yaml --profile dev run --rm dev make test-race
docker compose -f compose.yaml --profile dev run --rm dev make client-test
docker compose -f compose.yaml --profile dev run --rm dev make client-build
docker compose -f compose.yaml --profile dev run --rm dev make client-interop
docker compose -f compose.yaml --profile dev run --rm dev make site-build
docker compose -f compose.yaml --profile dev run --rm dev make e2e-servers
```

The dev image contains a source snapshot from build time. Rebuild it after
changing source or lockfiles. It has no Docker or Podman socket and does not run
privileged. The voice browser scenario is separately run by
`make e2e-voice-docker`, which starts and cleans up its own Linux/amd64 test
container from the host.

Production Compose/Podman deployment is gated on a later milestone. This file
does not claim Podman compatibility; the selected Podman version/provider and
real deployment verification are recorded before the M17/M19 release gates.
