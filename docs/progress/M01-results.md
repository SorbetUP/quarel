# M01 verification record

**Date:** 2026-10-08  
**Source:** `ab9c10070a9982c546f13c1ef6356a99685b9f6d` in `nsb/quarel-m00-baseline`  
**Gate:** independent integrated review accepted; local verification complete. M01 remains open because remote CI for this branch has not run. All current deployment manifests parse with both Docker Compose v2.39.3 and the selected `podman-compose` 1.6.0 provider; actual Podman runtime compatibility remains unverified and is required at M17/M19.

## Implementation

- Added injected `Clock` and ordinary resource `IDGenerator` ports, wired to the existing ID source from both command roots. Setup authentication codes, bearer tokens, nonces and cryptographic material remain on their existing dedicated paths.
- Added the direct permission snapshot matrix and import boundary checker. The checker reads actual `go list -test -json` package metadata, applies rules to `platform`, `application`, `adapters` and `transport`, and checks regular, internal-test and external-test imports. Synthetic fixtures exercise forbidden edges from application, adapters, transport, and test imports.
- Added the package caller map and proposed modular-monolith and cell-architecture ADRs. The ADRs leave event delivery, transactions and outbox protocols for their consumer milestones.
- Added pinned OCI images for Identity, community (with the temporary embedded LiveKit exception), docs MCP and the static site, plus a pinned Linux dev/test image and Compose development stack. Images run as non-root. The stack binds host ports to loopback and carries real healthchecks.
- Product health probes select HTTP or HTTPS from the effective `QUAREL_TLS` environment or data-directory setting. The probe rejects redirects and accepts only 2xx. Compose uses the Identity and community image healthchecks directly; no TLS override is needed for the local `QUAREL_TLS=off` stack.
- The import-boundary policy now checks `internal/transport` and test imports. Platform may import only platform packages; application may import only application/platform packages. Synthetic fixtures also reject platform-to-service edges.
- The slash-command gateway test now waits for `READY` after connecting, and starts its event timeout after setup. This removes a race where the command update could be emitted before the test connection registered.
- Independent integrated review accepted ADR-001/002 and the scope decisions in this record. The reviewer recommended no additional M01 code changes. The adversarial review found no new code defect and confirmed the remaining remote CI and Podman parsing proofs are acceptance gaps, so those remain open rather than being reported as complete.

## Verification evidence

All checks ran from this committed source under Docker/Colima. Go/client/site commands used the pinned Linux dev image; Compose stack and voice checks used dedicated OCI images.

- `docker-compose -f deploy/compose/compose.yaml config --quiet` — passed using standalone Docker Compose v2.39.3.
- Docker Compose v2.39.3 `config --no-env-resolution --no-interpolate --quiet` and podman-compose 1.6.0 `config` both passed for `deploy/identity/compose.yaml`, `deploy/identity/compose.proxy.yaml`, `deploy/server/compose.yaml`, and `deploy/compose/compose.yaml`. Docker Compose was instructed not to resolve env files or interpolate variables, so no `.env` contents were read. Podman parser output retained the quoted `"off"` values as strings. The local dev manifest was corrected to quote these YAML values after parser output showed YAML 1.1 boolean coercion.
- `docker-compose -f deploy/compose/compose.yaml --profile dev build` — all five images built from the committed source using pinned base images.
- `docker-compose -f deploy/compose/compose.yaml up -d --wait` — all four product/site containers reached `healthy`; the identity, community, docs MCP and site health endpoints each returned success. Compose used the image-provided Identity and community probes.
- Started the community OCI image with `QUAREL_TLS=self-signed` outside Compose — its image healthcheck also reached `healthy`, verifying the HTTPS probe path as well as Compose's HTTP path.
- `make e2e-servers` — passed all server scenarios and seven Discord importer tests. A first run exposed missing `curl` in the dev image; `curl` was added and the complete suite then passed.
- `make test`, `make test-race`, `make build`, `make vet`, `go test ./tools/import-boundaries`, and `go run ./tools/import-boundaries` — passed inside the committed Linux dev image. The race run covered all Go packages; the boundary checker inspected 48 internal packages.
- `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` — no vulnerabilities found in called code. The scanner reported one vulnerability in a required module that the code does not appear to call.
- `make client-test` — 54 tests passed; `make client-build`, `make client-interop`, and `make site-build` passed.
- `make e2e-servers` — all server scenarios, including calls through TURN, and seven Discord importer tests passed.
- `make e2e-voice-docker` — passed on Linux/amd64 using Colima 0.9.1 with an arm64 VZ/Rosetta VM. Two Chromium clients exchanged audio/video and screen share, and exercised moderation and LiveKit webhooks without page errors. The test image ID was `sha256:a4eec34e7a450f8aa25a5404e6cebd402eaf65f79f4b3878a6eb9f54bb08e6c0`; embedded LiveKit binary SHA-256 was `defde46e938e431cb23fc92fe793e503034fb5117cfa506796776a818f97afa3`.
- `go test -run TestSlashCommands -count=1 -v ./internal/community` — passed after the gateway test synchronization fix.
- `git diff --check` — passed.

An earlier `make test` attempt failed in `TestSlashCommands` because the test issued an update before the gateway connection had readied; the test synchronization was fixed, and the full suite passed afterward. A bind-mounted client check also failed because the mount replaced the Linux `node_modules` with the host macOS modules; the same checks passed in the image's own Linux workspace.

## Built image identities

These local images were built for Linux/arm64 on Docker/Colima. The listed local IDs are not published registry references.

| Image | Local image ID |
|---|---|
| `localhost/quarel-identity:local` | `sha256:81bc1d21e5354d65b0363b17688a363fcca9620773fea3d2a08ff522fcd4963a` |
| `localhost/quarel-server:local` | `sha256:4ed13d95afad19dd647b471f102aae5fba445dc41c0c38b27cc63dde3e0e063c` |
| `localhost/quarel-docs-mcp:local` | `sha256:7ca53554ce46ee48ba727c2d8c6b7d180f7d000664ea2913a0bbce579a1ca273` |
| `localhost/quarel-site:local` | `sha256:ee4e01b9bf278ed08922f1caba252c896d2ad87db019dcea40f25c902e5a7219` |
| `localhost/quarel-dev:local` | `sha256:cd4554190dd99496fea838cda6e2e698092b98595bfc433131de421ba4125a96` |

## Scope decisions and remaining limits

- No `CommunityContext` was fabricated: M01 has no trustworthy global CommunityID, ActorUserID or trace source. Add and propagate it after the M02/M03 identity contracts and observability provide those values.
- No durable `EventPublisher` was added because the current gateway is transient fanout, with no durable event/outbox consumer in M01. Define the durable contract in the reviewed M05 outbox design before introducing a port.
- No `ObjectStore` was added because the current local file store has atomicity, permissions and cleanup semantics that the future S3-compatible contract must preserve. Define and test that contract at M10, where object storage is introduced.
- No `PermissionEvaluator` interface was added because every current caller evaluates the same in-package SQL snapshot and there is no alternate provider to substitute. The direct matrix tests the pure algorithm. Revisit the interface when a use case crosses the module boundary.
- Docker/Colima and standalone Compose are verified here. Podman 6.1.3 and `podman-compose` 1.6.0 were installed for P8 validation. The selected provider parses all four current deployment manifests without an engine. A dedicated AppleHV VM entered emergency mode because the configured root filesystem UUID was missing. A separate libkrun attempt reached Linux `multi-user.target`, but Podman never received its readiness signal, the API socket remained inaccessible, and no containers could be run. The task-created VMs were removed. Parser-only evidence is recorded above; actual Podman runtime compatibility remains unverified until a working Linux host is available, and deployment/upgrade verification is required at M17/M19.
- The Compose stack is a local development smoke topology: TLS, voice, TURN and UPnP are disabled, and exposed TCP ports bind to loopback. It is not the production network topology. Future service images and their exact digests are pinned when their milestone introduces those services.
- The previously green GitHub Actions run applies to the exact M00 baseline commit, not this branch. New CI definitions are present but require a later remote run for remote evidence. The site build reports its existing missing/empty `i18n` collection and missing 404 content warnings; no site build failure resulted.
