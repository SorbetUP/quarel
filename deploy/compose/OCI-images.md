# OCI image scope and pins

## Current deliverables

| Deliverable | OCI status | Runtime base / architecture |
|---|---|---|
| Identity API | `Dockerfile.identity` | Static Go binary in pinned Distroless nonroot; OCI base index digests are in the Dockerfile. |
| Community API / realtime / embedded LiveKit | `Dockerfile.server` | Static Go binary in pinned Distroless nonroot; LiveKit v1.13.7 is pinned by the multiarch manifest digest. |
| Documentation MCP | `Dockerfile.docs-mcp` | Static Go binary in pinned Distroless nonroot. |
| Public site + wiki | `Dockerfile.site` | Astro static output served by the project's Go HTTP server in pinned Distroless nonroot; Node and Go build bases are digest pinned. |
| Development and test toolchain | `Dockerfile.dev` | Go 1.27, Node 24 and LiveKit v1.13.7 bases pinned by digest; SQLite, Python cryptography and lockfile dependencies are installed in the image. Runs as UID 10001. |
| Electron desktop | Native installer/AppImage/Windows artifacts | Built and tested with the dev OCI image; distributed as native artifacts, not as a server container. |

The Node build base used for the site is pinned to the multi-platform index
`node:24-bookworm-slim@sha256:d6aa754f16b3197301076f047b5def2f02ea1dbbc2ca920407d46d7ec7f87b20`.
Product image IDs and the tested Linux architecture are recorded in the CI job
and `docs/baseline/results.md` after each accepted gate.

## Planned infrastructure services

Authentik, PostgreSQL, PgBouncer, NATS JetStream, Centrifugo, Valkey, standalone
LiveKit, S3-compatible storage, search, and observability remain milestone
dependencies; they are not started by the current Compose stack. The milestone
that introduces each service must add its upstream image version, immutable
multiarch digest, architecture, license/provenance record, health/readiness
check, and data/secret mounts before adding it to Compose. This keeps the M01
development deployment small while making OCI packaging a required part of
each later service milestone.

No mutable `latest` image reference is used for a runtime or build base.
