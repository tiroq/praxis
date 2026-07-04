#!/usr/bin/env sh

set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT_DIR"

mkdir -p build

# Configure storage for the smoke test
export PRAXIS_STORAGE_BACKEND=sqlite
export PRAXIS_SQLITE_PATH=build/smoke-praxis.db
rm -f "$PRAXIS_SQLITE_PATH" # Clean slate for each smoke test run

cleanup() {
	if [ -n "${WATCH_PID:-}" ] && kill -0 "$WATCH_PID" 2>/dev/null; then
		kill "$WATCH_PID" 2>/dev/null || true
		wait "$WATCH_PID" 2>/dev/null || true
	fi
	if [ -n "${WORKER_PID:-}" ] && kill -0 "$WORKER_PID" 2>/dev/null; then
		kill "$WORKER_PID" 2>/dev/null || true
		wait "$WORKER_PID" 2>/dev/null || true
	fi
	task stop:nats || true
}

trap cleanup EXIT INT TERM

task build

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

./build/praxis watch --max-messages 1 --poll-timeout 500ms > build/smoke-praxis-watch.json 2> build/smoke-praxis-watch.log &
WATCH_PID=$!

sleep 0.25

PUBLISH_OUTPUT=$(./build/praxis publish --text "urgent review: buy tickets to Shanghai")
printf '%s\n' "$PUBLISH_OUTPUT" > build/smoke-praxis-publish.log

wait "$WATCH_PID"
unset WATCH_PID

python3 - <<'PY'
import json
import re
from pathlib import Path

pub = Path("build/smoke-praxis-publish.log").read_text().strip()
match_id = re.search(r"message_id=([^\s]+)", pub)
match_corr = re.search(r"correlation_id=([^\s]+)", pub)
if not match_id or not match_corr:
    raise SystemExit("failed to parse publish output")
message_id = match_id.group(1)
correlation_id = match_corr.group(1)

watch = json.loads(Path("build/smoke-praxis-watch.json").read_text())
if watch.get("status") != "ok":
    raise SystemExit(f"unexpected output status: {watch.get('status')}")
if watch.get("input_event_id") != message_id:
    raise SystemExit(f"input_event_id mismatch: expected {message_id}, got {watch.get('input_event_id')}")

report = {
    "worker_flow_ok": True,
    "message": "validated praxis publish/watch flow over NATS JetStream",
    "published": {
        "message_id": message_id,
        "correlation_id": correlation_id,
    },
    "received": watch,
}
Path("build/smoke-praxis.json").write_text(json.dumps(report, indent=2) + "\n")
print(json.dumps(report, indent=2))
PY

# Verify that the worker persisted kernel.pipeline.completed events to storage
if [ ! -f "$PRAXIS_SQLITE_PATH" ]; then
	echo "storage database not found: $PRAXIS_SQLITE_PATH" >&2
	exit 1
fi

./build/verify-storage "$PRAXIS_SQLITE_PATH" > build/smoke-praxis-storage.json

STORAGE_VERIFICATION=$(python3 -c "
import json
try:
    d = json.load(open('build/smoke-praxis-storage.json'))
    if d.get('verification') == 'passed':
        print('ok')
    else:
        print('failed')
except Exception as e:
    print('error')
")

if [ "$STORAGE_VERIFICATION" != "ok" ]; then
	echo "storage verification failed" >&2
	cat build/smoke-praxis-storage.json >&2 || true
	exit 1
fi

echo "storage verification passed: kernel.pipeline.completed events persisted"
cat build/smoke-praxis-storage.json
