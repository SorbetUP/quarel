# Third-party notices and M00 inventory scope

Quarel is licensed under Apache-2.0; see [`LICENSE`](LICENSE) and the upstream [`NOTICE`](NOTICE). This file supplements the upstream notice for the M00 Len(scape) baseline. The upstream `NOTICE` is retained unchanged as provenance; its final sentence overstates what manifests and lockfiles prove and must not be used as an exhaustive license inventory.

## Reviewed first-party distributions and known third-party components

| Distribution/component | Known third-party material | Evidence and limits |
|---|---|---|
| `quarel-server` | Go modules used by community server; embedded LiveKit server in Docker; vendored LiveKit JS voice test | See Go SBOM and `NOTICE`. Dockerfile includes Quarel and LiveKit license/notice files under `/licenses/`; LiveKit files are pinned to v1.13.7 in [`licenses/livekit/`](licenses/livekit/). |
| `quarel-identity` | Go modules used by Identity, including WebAuthn, Pion TURN and SQLite | See Go SBOM and `NOTICE`; a per-binary SBOM remains required. |
| `quarelctl` | Go modules including Pion WebRTC and goolm (MPL-2.0) | `make build` produces this CLI, but whether it is a supported distributed product is not established by this baseline. Treat as a distribution for notices until that is decided. |
| `quarel-docs-mcp` | Go modules including the MCP Go SDK | See Go SBOM; per-binary SBOM remains required. |
| Go module `github.com/tinylib/msgp@v1.6.4` | MIT License; `msgp/json.go` contains portions attributed to The Go Authors under a BSD-style license | CycloneDX Go license detection omitted this component's license. The pinned module's `LICENSE` and `msgp/json.go` headers were reviewed manually; retain both attributions in distributions that include this package. |
| Desktop client | Electron, Chromium, Node, React, LiveKit client, crypto libraries, fonts and transitive npm packages | Lockfile SBOM includes dev dependencies. Built packages and their `LICENSES.chromium.html` must be inspected separately; see upstream `NOTICE`. |
| Site | Astro/Starlight, Pagefind, fonts and transitive npm packages | Lockfile SBOM is not a built-site artifact scan; inspect generated `site/dist` before release. |
| Rust/WASM | vodozemac and locked Cargo dependencies | `client-wasm-cargo-metadata.json` records Cargo metadata, including reported crate license fields; review exceptions and generate CycloneDX before rebuilding/releasing WASM. |
| Windows installers / CI tools | NSIS, `go-winres`, Playwright/Chromium, build tools | Build/test tools are not automatically runtime components; retain their own inventory and scan final installers. |

## License inventory process

The SPDX identifiers above are sourced from the upstream `NOTICE` and dependency metadata, not a legal compatibility opinion. Lockfiles pin versions but do not prove complete, correct license coverage. `docs/baseline/sbom/README.md` records generated inventories and limits. Review license files for all shipped components and scan each final binary, container, desktop package, installer and website output before a release. Preserve notices provided by Electron/Chromium and LiveKit. The component that embeds a dependency is responsible for ensuring its distribution carries applicable notices.

No code from the inspiration-only projects listed in the master plan has been copied as part of M00.
