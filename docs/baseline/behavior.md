# Quarel behavior baseline

**Source:** `d3afef2605838b38d98e8a0b60c0f8e63f1b0065` (2026-10-07). This document records observable architecture and contracts before Len(scape) changes. It is a concise map, not a replacement for the endpoint handlers or a generated OpenAPI specification.

## Services and persistence

- `quarel-identity` is the identity, account, friend, device-key, direct-message, inbox, push and TURN-configuration service. It exposes HTTP/JSON under `/v1/`, a real-time gateway at `/v1/gateway`, and the OIDC/JWKS discovery endpoint. Its embedded SQLite database is under `QUAREL_DATA_DIR`.
- `quarel-server` is one community server. It owns the community's channels, membership, roles, permissions, messages, moderation, audit, attachments, voice room state and settings in its own SQLite database. Each instance has an independent `QUAREL_DATA_DIR`; the community HTTP API is under `/v1/`, with its gateway at `/v1/gateway` and local administration UI on a separate listener.
- `quarel-docs-mcp` serves the checked-in site documentation over MCP. It is not on the user data path.
- The desktop/web client is React + TypeScript, with Electron for desktop packaging. Message and DM cryptography uses the committed vodozemac WebAssembly module; the test CLI uses goolm for interoperability checks.
- Voice signaling is proxied by the community server to an embedded LiveKit process when configured. TURN is provided by the Identity service when enabled.

## HTTP and realtime contracts

The public APIs are registered in `internal/identity/server.go` and the community server's `Handler` implementation. Requests use JSON for structured payloads and bearer identity tokens for authenticated HTTP calls. Health checks are `GET /v1/health`. Authentication includes register, email verification, login, logout, password recovery, sessions, passkeys and optional two-factor authentication. Identity also owns `/v1/me`, profile, friends, blocks, device keys, DMs, inbox and backup endpoints. The community service owns join/invite, channels, roles, members, messages, moderation, audit, attachment and community settings endpoints.

Realtime clients open the corresponding `/v1/gateway` WebSocket. Event names and payloads are defined in `pkg/realtime`; community membership and Identity DM event streams are separate. The baseline has no distinct global control plane or cell-routing protocol.

## Data and encryption boundaries

- Identity and each community server have independent SQLite files. There is no shared relational database or tenant ID on community rows at this baseline.
- Server-held community messages are service-readable; client-to-client DMs and uploaded DM files use end-to-end encryption with device keys and Olm/Megolm-compatible crypto.
- Identity tokens identify a Quarel account to community servers. Community authorization is evaluated by each community server using its own membership, roles, permission overwrites and moderation state.
- Attachments and backups are stored by the owning service under its configured data directory. API callers must have service-level authorization before retrieval.

## Schema and API capture pointers

The exact ordered SQLite migrations are captured in [`schema/identity-migrations.sql`](schema/identity-migrations.sql) and [`schema/community-migrations.sql`](schema/community-migrations.sql). Direct HTTP route registrations are captured in [`api-routes.md`](api-routes.md); the route registrations are the source of truth for the baseline API, and there is no checked-in OpenAPI document at this pin. The baseline tests listed in [`results.md`](results.md) capture representative behavior, including message flows, moderation, account changes, encrypted DMs, voice and operational backup/restore.

Before schema or route changes, regenerate a concrete schema/API snapshot and compare it to this commit. Do not treat this summary as a complete schema dump.
