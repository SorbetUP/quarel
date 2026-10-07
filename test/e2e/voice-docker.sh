#!/usr/bin/env bash
# Build and run the unchanged voice.sh scenario in an isolated Linux container.
# Usage: test/e2e/voice-docker.sh [linux/amd64]
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../.." && pwd)
PLATFORM=${1:-linux/amd64}
if [[ "$PLATFORM" != "linux/amd64" ]]; then
  echo "Unsupported platform: $PLATFORM (the pinned LiveKit image is linux/amd64)" >&2
  exit 2
fi

SUFFIX="${UID:-$(id -u)}-$$"
IMAGE="quarel-e2e-voice:$SUFFIX"
CONTAINER="quarel-e2e-voice-$SUFFIX"
cleanup() {
  docker rm -f "$CONTAINER" >/dev/null 2>&1 || true
  docker image rm "$IMAGE" >/dev/null 2>&1 || true
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

docker buildx build --load --platform "$PLATFORM" -f "$ROOT/Dockerfile.e2e" -t "$IMAGE" "$ROOT"
docker image inspect "$IMAGE" --format 'Image {{.Id}}; platform {{.Os}}/{{.Architecture}}'
docker run --rm --platform "$PLATFORM" --entrypoint sha256sum "$IMAGE" /usr/local/bin/livekit-server
docker create --platform "$PLATFORM" --name "$CONTAINER" --network bridge --init "$IMAGE" >/dev/null
docker start --attach "$CONTAINER"
