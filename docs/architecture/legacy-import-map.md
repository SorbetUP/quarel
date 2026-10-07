# Existing Quarel package and runtime map

This is the M01 starting map for the pinned source tree. It records observed
imports and runtime ownership; it does not claim that the current package
graph already follows the target architecture. The Go edges were collected
with `go list -f '{{.ImportPath}} {{.Imports}}' ./internal/...` at source SHA
`d3afef2605838b38d98e8a0b60c0f8e63f1b0065`.

## Runtime and package ownership today

| Area | Current owner and dependencies | Classification / intended direction |
|---|---|---|
| Community | `internal/community` contains HTTP handlers, domain rules, SQLite queries and service wiring. It imports `httpapi`, `realtime`, `secret`, `sqlitedb`, `voice`, `ratelimit`, settings, network, TLS and domain packages. | Legacy service package. Split application/domain/adapters only with real use cases; current direct edges stay on an explicit allowlist until migrated. |
| Identity | `internal/identity` contains HTTP handlers, identity/session/business rules, SQLite queries and service wiring. It imports `httpapi`, `realtime`, `secret`, `sqlitedb`, settings, network, TLS, TURN and token packages. | Legacy service package. Identity is not yet an Authentik adapter or global-user model. |
| HTTP helpers | `internal/httpapi` owns JSON decoding/writing and HTTP error representation. `community`, `identity`, `adminui`, `diskspace` and `ratelimit` import it. | Transport concern currently shared as utility. The domain-to-transport edges are legacy; new application/platform packages must not add them. |
| Realtime | `internal/realtime.Hub` owns WebSocket sessions, groups, subscriptions and ephemeral fanout. `community` and `identity` publish directly to it. | Transient in-process transport only. It is not a durable publisher or event log; a future NATS/outbox contract is deferred to M05. |
| Voice | `internal/voice` manages the embedded LiveKit process and token/signing support; `community` owns room authorization and event handling. `voice` imports `secret`. | Transitional runtime boundary. Keep server-side permissions authoritative; the service split belongs to M08. |
| Backup | `internal/backup` owns archive/restore of service data and is used by `internal/adminui`; tests also exercise SQLite migration state. | Operations adapter for the current standalone SQLite deployment, not a future cell-level backup contract. |
| SQLite | `internal/sqlitedb` owns SQLite driver setup, migration execution and database file configuration. Community and Identity stores use it. | Persistence adapter used by both current services. PostgreSQL and tenant schema work is deferred to M03. |
| Secrets and IDs | `internal/secret` owns CSPRNG-derived IDs, bearer tokens and key-file operations. Community and Identity use it for IDs and cryptographic material; `voice` uses it for API credentials. | Security utility. `NewToken`, setup codes, nonces, signing/encryption keys and other secrets must not be routed through an ordinary ID port. |
| Public packages | `pkg/idtoken`, `pkg/tlsbind` and `pkg/e2ekeys` are imported by server code; they do not depend on `internal` service packages. | Stable reusable protocol/crypto helpers. Keep backend imports out of client packages and do not move these into service core without a concrete reason. |
| Supporting packages | `diskspace` and `ratelimit` import `httpapi`; `adminui` imports `backup` and `httpapi`. Network, TLS, settings and theme packages are consumed by service packages. | Existing shared utilities; their dependency on transport is legacy and included in the boundary allowlist. |

## Container and deployment ownership today

| Deliverable | Current definition | Portability status |
|---|---|---|
| Identity service | `Dockerfile.identity`, `deploy/identity/compose.yaml`, `deploy/identity/compose.proxy.yaml` | OCI image builds and health smoke have passed under Docker. Compose health state and actual Podman operation are not yet verified. The reverse proxy publishes privileged TLS port 443. |
| Community service | `Dockerfile.server`, `deploy/server/compose.yaml` | OCI image builds and health smoke have passed under Docker. The service still embeds LiveKit. Deployment currently uses host networking/UPnP and voice UDP ports; these are known constraints, not portable Compose proof. |
| Documentation MCP | `Dockerfile.docs-mcp` | OCI build and `/health` smoke passed under Docker. It has no deployment Compose service yet. |
| Site and clients | `site/` Astro build; `client/` React/Electron; no OCI site image or versioned Docker development/test stack yet. | Site should be served as static OCI output. Electron packages remain native artifacts but should build/test inside a container. |
| Planned services | Authentik, PostgreSQL, NATS, Centrifugo, Valkey, LiveKit, object storage, search and observability are not part of current runtime. | Add pinned image refs and explicit health/readiness, persistence and secret ownership in their dependent milestones. |

## Directed migration rules

1. Existing `community`/`identity` → `httpapi`, `realtime`, `sqlitedb`, and
   concrete service helpers are recorded legacy edges. New code must not add
   dependencies from future `internal/platform` or `internal/application`
   packages into `internal/transport` or concrete adapters.
2. The current `realtime.Hub` remains an ephemeral `RealtimeSink`; durable
   event delivery is a separate M05 contract and must not be hidden behind the
   same port.
3. `Server.now` is already an injectable clock seam. Preserve it until a real
   shared consumer justifies a common clock contract.
4. Any ordinary ID abstraction is limited to ordinary resource/database IDs.
   It must not generate bearer secrets, setup-authentication codes, nonces,
   cryptographic keys, or tokens.
5. The current `permSnapshot` evaluator is loaded from SQLite but its
   evaluation methods are pure. Keep the SQL snapshot loader at the adapter
   edge; permission decisions consume only the snapshot.
6. `internal/application`, `internal/adapters` and `internal/transport` are
   migration destinations, not permission to move every handler or store in
   M01. Add code when a reviewed consumer/provider exists; do not create empty
   packages just to make the directory tree look complete.

## Provenance of this map

Observed Go imports came from `go list`; Compose/Docker ownership came from
the tracked files listed above. The behavioral and schema baselines remain in
`docs/baseline/`. This map should be revised when an edge is migrated, with
the old allowlist entry removed in the same change.
