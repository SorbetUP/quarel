# ADR-001: Keep the product as a modular Go monolith

**Status:** Accepted for M01  
**Date:** 2026-10-08

Independent integrated review found no remaining architecture or scope blocker.
The M01 verification record preserves the still-pending remote CI and Podman
provider evidence; this ADR does not claim those checks have run.

## Context

The pinned Quarel source combines HTTP handlers, product rules, SQLite access,
realtime fanout and some runtime wiring inside the `community` and `identity`
packages. The product must retain Quarel's existing API, permissions, E2EE and
media behavior while the foundations change over several milestones. A broad
package move or early service split would increase integration risk without
providing a measured scaling or failure-isolation benefit.

## Decision

Keep the product in one Go module and deployable API processes while creating
directed module seams incrementally. Introduce a port with its actual consumer
and provider; do not add placeholder packages or move all existing handlers
and stores as a preparatory exercise. Keep composition in the command entry
points. Existing service-to-`httpapi`, service-to-`realtime`, and
service-to-SQLite imports remain explicit legacy edges until a relevant
feature milestone migrates them.

New dependencies follow inward direction: platform contracts have no service
or adapter imports; application code may depend on platform/domain contracts;
adapters implement inward-facing contracts; transport invokes application
behavior. The M01 import checker enforces the edges introduced by the new
packages and tests its negative cases. It treats `internal/platform` as a leaf
within the Go module, allows `internal/application` to import only application
and platform packages, and rejects adapters or transport importing legacy
service/core packages. It checks ordinary and test-only imports.

## Consequences

- Existing behavior remains the compatibility oracle; this ADR does not
  authorize endpoint, schema, event, permission-bit or storage changes.
- Services can be tested through narrow ports without requiring every current
  dependency to move at once.
- Runtime processes remain independently deployable containers where the
  current architecture already has separate processes (Identity and
  community); embedded LiveKit remains a documented transition until M08.
- A module is extracted into a service only after measured independent
  scaling, failure isolation, deployment ownership or runtime needs justify
  it, and after its contracts are versioned.
- `EventPublisher`/outbox and `ObjectStore` are not settled by this ADR; they
  remain scheduled for M05 and M10 respectively.
- M01 does not add the master-plan `PermissionEvaluator` interface. Every
  current caller evaluates the same pure methods on `permSnapshot`, which is
  populated from the current SQLite schema, and no alternate evaluator or
  application-layer consumer exists yet. The SQL-free matrix tests this
  algorithm directly. Revisit the interface when permission evaluation moves
  behind a real module boundary; this is an explicit M01 scope deviation.

## Accepted tradeoffs

Some service packages will retain legacy dependency edges during the gradual
migration. The explicit allowlist creates temporary maintenance work, but it
lets new code follow the target direction without making M01 a risky rewrite.

## Revisit when

A concrete module needs independent scaling or failure isolation, an import
boundary migration reaches a real use case, or a future milestone defines a
stable service contract.
