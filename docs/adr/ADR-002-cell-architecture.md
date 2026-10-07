# ADR-002: Isolate tenant placement into cells

**Status:** Accepted for M01  
**Date:** 2026-10-08

Independent integrated review found no remaining architecture or scope blocker.
The M01 verification record preserves the still-pending remote CI and Podman
provider evidence; this ADR does not claim those checks have run.

## Context

The target product must support small shared communities and larger dedicated
deployments without maintaining different business logic for each tier.
Quarel currently runs as a standalone community service with a local SQLite
database and a key-derived server ID. It does not yet have a global
`CommunityID`, tenant context, cell router or shared PostgreSQL schema. Those
primitives are deferred to M02/M03.

## Decision

Treat a cell as the unit of capacity, failure containment and operational
placement. A community has one primary cell at a time. Shared and dedicated
cells run the same product and schema; only placement and resource allocation
differ. A routing layer resolves an authoritative community placement before
tenant data is read or mutated.

The intended boundaries are:

- **Global plane:** global users, sessions, devices, friend/DM relationships,
  notifications and authoritative community placement.
- **Cell:** community-owned metadata, members, roles, channels, messages,
  moderation, voice authorization and tenant-scoped derived state.
- **Shared platform services:** realtime transport, event workers, media,
  object storage, search and observability, with tenant namespacing and
  explicit failure/ownership boundaries.

This ADR records only boundaries, failure domains and placement. It does not
define database transactions, cross-plane consistency, outbox schemas,
idempotency keys, placement migration or event delivery. Those contracts must
be reviewed in the M03/M05 ADRs before implementation.

## Consequences

- No current Quarel `ServerID`, local member row ID or database path is treated
  as the future global `CommunityID`.
- M01 creates no empty cell router or `CommunityContext`; trustworthy global
  IDs and request tracing must arrive before propagation is implemented.
- Cell health, capacity, backup/restore, failover and placement changes need
  explicit operational evidence before release gates.
- A dedicated cell uses the same domain rules and API contracts as a shared
  cell, avoiding a special product fork.

## Accepted tradeoffs

The current self-hosted server model remains in place while global identity
and tenant storage contracts are built. This delays cell placement code but
avoids inventing IDs or authorization context from incompatible current keys.

## Revisit when

M02/M03 establish global community IDs, tenant context and cell persistence,
or when measured load/failure requirements change the placement model.
