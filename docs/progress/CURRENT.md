# Len(scape) plan ledger

Master plan: `/Users/sorbet/Downloads/LENSCAPE_QUAREL_MASTER_PLAN.md` (user-provided, 2026-10-07). Progress is tracked here; the master plan remains the design source. Status is never inferred from code presence alone: each gate needs recorded evidence and independent review.

## Cross-cutting runtime constraint (user clarification, 2026-10-07)

The full stack must build and run with Docker now, and remain deployable with Podman in production. Treat OCI images and a documented Compose-compatible deployment as deliverables, not optional packaging. Avoid Docker-engine-only assumptions (Docker socket access, privileged containers, daemon-managed side effects) in services. Pin image digests/tags, run as non-root, use explicit volumes/secrets and health checks. A current temporary exception is `quarel-server`, which embeds LiveKit; M08 is expected to split it into a separate service. Docker verification is mandatory. Podman compatibility is not claimed until the same images and orchestration are actually exercised on a Podman host. Podman 6.1.3 is installed and `podman-compose` 1.6.0 is pinned in a task venv. The selected provider parses all four current deployment manifests. The AppleHV VM entered emergency mode because its root filesystem UUID was missing; an alternative libkrun attempt reached Linux `multi-user.target` but never yielded a usable Podman API/socket. Both task-created VMs were removed. Docker Engine 28.4.0 is available. The Compose CLI plugin is absent, but standalone `docker-compose` v2.39.3 is available: local config validation, image builds, and runtime health checks have been exercised. The new remote CI jobs have not yet run on this branch. Add actual Podman deployment/upgrade verification before production readiness (M17/M19), not only a Compose parse.

| Milestone | Status | Notes |
|---|---|---|
| M00 — freeze, baseline, provenance | COMPLETE | Exact upstream pin and baseline are archived and independently reviewed. Every M00 check passes locally; voice E2E passes in Linux/amd64 Docker under Colima VZ+Rosetta. GitHub Actions run 36354609611 is green on the exact pinned source SHA. The added Docker image/voice jobs have local build/runtime evidence but are not yet included in that remote run. Podman was absent at the time of M00. |
| M01 — architecture seams | IN PROGRESS | P1–P6 and P7 implementation are complete; local Docker/Compose CI-equivalent checks and integrated review pass. P7 remote CI has not run on this branch. P8 selects Podman 6.1.3 + podman-compose 1.6.0; both Docker Compose and podman-compose parse all four current deployment manifests. Actual Podman runtime/deployment remains unverified and is required at M17/M19. See `docs/progress/M01-results.md`. |
| M02–M20 | BLOCKED BY DEPENDENCIES | Not started. |

## Sequencing decisions required by plan review

1. M02 creates Authentik global tables, while the plan places PostgreSQL and migrations in M03. Before implementation, add a reviewed ADR and split the necessary Postgres/migration/UUID foundation out of M03; execute that foundation after M01, then M02, then the remaining M03 tenancy work.
2. Before M05/M06/M09/M12, specify event IDs, idempotency, per-cell outboxes and a separate global outbox in a reviewed contract/ADR.
3. Before M16, replace unsafe post-reopen rollback with fenced writes, drain/watermark and atomic cutover; rollback after reopening requires reverse replication.
4. M20 S0/S1 may be verified in the available environment. Do not claim S2–S4 certification without an approved target environment and recorded hardware/budget evidence.

These are sequencing/acceptance resolutions, not changes to the target architecture. Each must be carried into an ADR or acceptance record before its dependent milestone starts.
