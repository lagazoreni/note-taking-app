#!/usr/bin/env sh
set -eu
IMAGE="${1:-noted:local}"
NAME="noted-smoke-$$"
DATA_DIR="$(mktemp -d)"
cleanup() { docker rm -f "$NAME" >/dev/null 2>&1 || true; rm -rf "$DATA_DIR"; }
trap cleanup EXIT
docker run -d --name "$NAME" -p 18081:8080 -v "$DATA_DIR:/data" "$IMAGE" >/dev/null
for _ in $(seq 1 30); do curl -fsS http://127.0.0.1:18081/readyz >/dev/null && break; sleep 1; done
curl -fsS http://127.0.0.1:18081/healthz >/dev/null
curl -fsS http://127.0.0.1:18081/readyz >/dev/null
curl -fsS -X POST http://127.0.0.1:18081/api/v1/workspaces -H 'content-type: application/json' -d '{"name":"Smoke"}' >/dev/null
if [ "$(docker exec "$NAME" id -u)" = "0" ]; then echo 'container unexpectedly runs as root' >&2; exit 1; fi
