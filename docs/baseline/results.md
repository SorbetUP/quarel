# Baseline verification record

**Source commit:** `d3afef2605838b38d98e8a0b60c0f8e63f1b0065`  
**Run date:** 2026-10-07  
**Host:** macOS arm64; Go 1.27.1; Node.js/npm from the configured workspace runtime.  
**Worktree:** isolated managed worktree `nsb/quarel-m00-baseline`.

## Docker image verification

Docker Engine 28.4.0 on Colima (Linux/arm64) built all three pinned upstream images. Runtime smoke tests were performed as non-root containers and queried from the host:

| Image | Build result | Image ID | Runtime observation |
|---|---|---|---|
| `Dockerfile.identity` | PASS (`linux/arm64`) | `sha256:71c20f943a27bd786ebd2b03ebb31490b9dd36f51cfb9f0b1ece3f9a5e243c25` | `GET /v1/health` returned `{"status":"ok"}` |
| `Dockerfile.server` | PASS (`linux/arm64`) | `sha256:f04ccc1a90eac1acd03ff8ad8d27b552725b75d47222ae7adadf181185f03c7d` | Embedded LiveKit process started; TLS `GET /v1/health` returned `{"status":"ok"}` |
| `Dockerfile.docs-mcp` | PASS (`linux/arm64`) | `sha256:c8989dc3193f58df4c2b091ef17ccd78bb780f3956440507bcc051a12050b535` | `GET /health` returned `ok` |

The first Identity build exposed that `.dockerignore` excluded `THIRD_PARTY_NOTICES.md`; the ignore rule now explicitly includes that file, and all images rebuilt successfully. These are local image IDs, not published registry digests. Podman is not installed on this host, so Podman runtime compatibility remains unverified.

The server image filesystem was inspected with `docker cp`: `/licenses/quarel/THIRD_PARTY_NOTICES.md` (3002 bytes) and `/licenses/livekit/NOTICE` (553 bytes) were present. The exact pinned LiveKit license-file provenance and checksums are recorded in `docs/provenance/livekit-license-v1.13.7.md`.

| Check | Result | Evidence / limitation |
|---|---|---|
| `make test` | PASS | `go test -timeout 120s ./...` |
| `make test-race` | PASS | `CGO_ENABLED=1 go test -race -timeout 300s ./...` |
| `make vet` | PASS | `go vet ./...` |
| `make client-test` | PASS | TypeScript check and Vitest: 12 files, 54 tests |
| `make client-interop` | PASS | vodozemac ↔ goolm interoperability assertions |
| `make e2e-messages` | PASS | Real local Identity/community services and `quarelctl` |
| `make e2e-moderation` | PASS | Same |
| `make e2e-accounts` | PASS | Same |
| `make e2e-dm` | PASS | Same |
| `make e2e-dm-groups` | PASS | Same |
| `make e2e-voice` | PASS | Docker `linux/amd64` image built with Buildx and run under a separate Colima 0.9.1 VZ+Rosetta profile (Docker Engine 28.4.0). Image ID `sha256:202cbda40e6a037f3afcd5c56009e981ce251a3329b6d088c7015be0bb365390`; embedded LiveKit `sha256:defde46e938e431cb23fc92fe793e503034fb5117cfa506796776a818f97afa3`. Two headless Chromium clients exchanged synthetic audio/video, shared a screen, and passed moderation/webhook checks with no page errors. This test uses one container network namespace and does not verify host port mapping, NAT, LAN reachability, or real-device audio quality. |
| `make e2e-security` | PASS | Real local services and `quarelctl` |
| `make e2e-ops` | PASS | Backup/restore scenario |
| GitHub Actions | PASS | CI run [36354609611](https://github.com/anlekg/quarel/actions/runs/36354609611) on the exact pinned source SHA `d3afef2605838b38d98e8a0b60c0f8e63f1b0065` completed successfully; all five jobs passed (Go, client, site, server E2E, desktop/web/voice E2E). This run predates the locally added Docker image and Docker voice jobs; their equivalent image health checks and voice scenario passed locally as recorded above. |

The baseline required test-only portability fixes in E2E setup: macOS uses `$HOME/Library/Application Support` and ignores `XDG_CONFIG_HOME`, and macOS `wc` pads counts. These changes isolate temporary profiles under each test's temp directory and normalize numeric output. `client/e2e/fixtures.ts` also sets a temporary `HOME` for spawned processes. No product behavior was changed. The amd64 voice test image builds from ARM with Buildx cross-compilation; running Chromium via QEMU failed, while the separate VZ+Rosetta profile passed the complete scenario.

The source commit has no tags, so the local `git describe` version is the abbreviated commit ID. npm installation reported 10 advisories (8 moderate, 2 high) in the locked client dependency graph; no automated dependency upgrades were applied. See the lockfile SBOM and review the affected package paths before release.
