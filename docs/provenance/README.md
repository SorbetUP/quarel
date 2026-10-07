# Provenance records

`UPSTREAM.md` records the exact Quarel source pin. For each later code port or copied implementation, add a record containing source project, repository, commit, source files, license and modifications as required by the master plan. Inspiration-only work must remain clean-room and be documented in an ADR without copying source code.

M00 has no ported implementation. Test portability edits are local test harness fixes and do not copy external implementation code. Baseline dependency inventories are in `docs/baseline/sbom/`; `THIRD_PARTY_NOTICES.md` describes scope and unresolved items.
