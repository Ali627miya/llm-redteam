#!/usr/bin/env bash
# One-shot: build, start mock OpenAI server, run scan, stop server.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
export PATH="/opt/homebrew/bin:/usr/local/bin:${PATH}"

PORT="${PORT:-8765}"
if command -v lsof >/dev/null && lsof -ti ":$PORT" >/dev/null 2>&1; then
  echo "Freeing port $PORT …"
  lsof -ti ":$PORT" | xargs kill 2>/dev/null || true
  sleep 0.5
fi

echo "Building …"
go build -o bin/mocktarget ./cmd/mocktarget
go build -o bin/redteam ./cmd/redteam

CFG=$(mktemp)
trap 'rm -f "$CFG"; kill $MPID 2>/dev/null || true' EXIT
sed "s|127.0.0.1:8765|127.0.0.1:$PORT|g" examples/redteam-local-mockserver.yaml >"$CFG"

./bin/mocktarget -addr "127.0.0.1:$PORT" -persona vulnerable &
MPID=$!
sleep 0.8

echo "Scanning http://127.0.0.1:$PORT …"
./bin/redteam run --config "$CFG" --output report.html --format html
echo "Done. Open: $ROOT/report.html"
