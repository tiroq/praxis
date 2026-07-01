#!/usr/bin/env sh

set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT_DIR"

mkdir -p build

cleanup() {
	if [ -n "${WORKER_PID:-}" ] && kill -0 "$WORKER_PID" 2>/dev/null; then
		kill "$WORKER_PID" 2>/dev/null || true
		wait "$WORKER_PID" 2>/dev/null || true
	fi
}

trap cleanup EXIT INT TERM

task run:nats

./build/worker > build/worker-smoke.log 2>&1 &
WORKER_PID=$!

i=0
until grep -q "nats subscriber started" build/worker-smoke.log 2>/dev/null; do
	i=$((i + 1))
	if [ "$i" -ge 40 ]; then
		echo "worker did not become ready" >&2
		cat build/worker-smoke.log >&2 || true
		exit 1
	fi
	if ! kill -0 "$WORKER_PID" 2>/dev/null; then
		echo "worker exited before becoming ready" >&2
		cat build/worker-smoke.log >&2 || true
		exit 1
	fi
	sleep 0.25
	done

./build/nats-smoke --out build/smoke-nats.json