# Baseline SBOMs

These files describe locked dependency graphs at the M00 source pin. They are inventory inputs, not a legal conclusion or a complete attestation of built artifacts.

- `client-npm.cdx.json` and `site-npm.cdx.json`: npm 11.19.0 CycloneDX 1.5 generated from each package lock with `npm sbom --package-lock-only`.
- `client-wasm-cargo-metadata.json`: Cargo metadata for the locked Rust/WASM graph, generated with `cargo metadata --locked`. It is not CycloneDX; crate license fields must be reviewed and a standard CycloneDX SBOM should be generated before shipping a WASM rebuild.
- `quarel-go.cdx.json`: Go module graph CycloneDX generated with `cyclonedx-gomod v1.12.0 mod -json -licenses` on `darwin/arm64`. It is module-wide and not a per-binary runtime inventory. The generator could not resolve the untagged main module's version and reports no detected license for `github.com/tinylib/msgp@v1.6.4`. Manual review of the pinned module source confirmed MIT plus Go Authors BSD-style portions; see `THIRD_PARTY_NOTICES.md`. The tool's omission remains visible in the BOM and should be corrected/enriched in the next SBOM pipeline.

CycloneDX Go `app` generation was attempted separately for `quarel-server`, `quarel-identity`, `quarelctl`, and `quarel-docs-mcp`; v1.12.0 failed with `failed to determine version of main module: git: reference not found` because this baseline has no tag. The generated npm SBOMs include locked dev dependencies, so Electron and tooling are visible; this still does not inspect the actual Electron/Chromium package or built site. The Windows installers, containers, LiveKit binary, Chromium runtime notices and release artifacts require artifact-level scans before distribution. The LiveKit v1.13.7 `LICENSE` and `NOTICE` included with the community Docker image were obtained from its [versioned upstream source](https://github.com/livekit/livekit/tree/v1.13.7).

Do not describe `go.mod`, `go.sum`, or npm lockfiles as exhaustive license inventories. `NOTICE` is inherited from upstream and contains that claim; `THIRD_PARTY_NOTICES.md` records the M00 reviewed scope and known limitations without silently rewriting upstream history.
