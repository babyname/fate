#!/usr/bin/env bash
# Reproducible acceptance gate for the examples/ tree.
#
#   bash scripts/verify-examples.sh [v4.3.0]
#
# Produces evidence under tests/acceptance/<version>/evidence/:
#   - examples-verification.md   overall report
#   - structure.log              directory + documentation contract
#   - go-build/vet log
#   - examples-server.log        real server log
#   - examples-http-*.body/.status/.json/.assert   real wire responses
#   - example-<name>.log         real stdout of each example
#
# Requires: Go, curl, python3. No external services. The gate starts and stops
# its own server on FATE_EXAMPLES_PORT (default 18082).
set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LABEL="${1:-v4.3.0}"
OUT="$ROOT/tests/acceptance/$LABEL/evidence"
TMP="$(mktemp -d)"
PORT="${FATE_EXAMPLES_PORT:-18082}"
BASE="http://127.0.0.1:$PORT"
RESULT=0
ROWS=()
SERVER_PID=""

export NO_PROXY="127.0.0.1,localhost${NO_PROXY:+,$NO_PROXY}"
export no_proxy="$NO_PROXY"
export PYTHONIOENCODING=utf-8

mkdir -p "$OUT"
cleanup() {
  if [ -n "$SERVER_PID" ]; then kill "$SERVER_PID" 2>/dev/null || true; fi
  rm -rf "$TMP"
}
trap cleanup EXIT

pass() { echo "PASS"; ROWS+=("| $1 | PASS |"); }
fail() { echo "FAIL"; ROWS+=("| $1 | FAIL |"); RESULT=1; }

step() {
  local name="$1"; shift
  local log="$OUT/examples-$name.log"
  printf '== %s == ' "$name"
  if "$@" > "$log" 2>&1; then pass "$name"; else fail "$name"; fi
}

check_structure() {
  local ok=1
  [ -d "$ROOT/examples" ] || { echo "missing examples/"; ok=0; }
  [ ! -d "$ROOT/example" ] || { echo "legacy example/ still present"; ok=0; }
  [ -f "$ROOT/examples/README.md" ] || { echo "missing examples/README.md"; ok=0; }
  for n in quickstart stream filters explore http-client; do
    [ -f "$ROOT/examples/$n/main.go" ] || { echo "missing examples/$n/main.go"; ok=0; }
  done
  if grep -Eq 'POST /api/generate|/api/task/:id|GetExcellentTable|create_a_name|example/' \
      "$ROOT/docs/DEVELOPMENT.md" "$ROOT/docs/architecture/02-directory-structure.md"; then
    echo "developer docs still advertise retired API paths or the old example/ tree"
    ok=0
  fi
  if [ "$ok" -eq 1 ]; then
    echo "examples/ tree complete; legacy example/ removed; docs use /api/v1/* contract"
  fi
  return $((1 - ok))
}

run_example() {
  local name="$1"; shift
  local logfile="$OUT/example-$name.log"
  printf '== example/%s == ' "$name"
  if (cd "$ROOT" && FATE_BASE_URL="$BASE" timeout 240 go run "./examples/$name") \
      > "$logfile" 2>&1; then
    local missing=""
    for pat in "$@"; do
      if ! grep -Eq "$pat" "$logfile"; then missing="$missing [$pat]"; fi
    done
    if [ -z "$missing" ]; then
      pass "example/$name (see example-$name.log)"
    else
      echo "FAIL (output assertions:$missing)"
      ROWS+=("| example/$name output assertions | FAIL |")
      RESULT=1
    fi
  else
    fail "example/$name exited non-zero (see example-$name.log)"
  fi
}

check_http() {
  local name="$1" method="$2" path="$3" data="${4:-}" want="${5:-200}"
  local bodyfile="$OUT/examples-http-$name.body"
  local jsonfile="$OUT/examples-http-$name.json"
  local status
  local args=(-sS -o "$bodyfile" -w '%{http_code}' -X "$method")
  if [ -n "$data" ]; then args+=(-H 'Content-Type: application/json' --data-binary @-); fi
  args+=("$BASE$path")

  printf '== http/%s == ' "$name"
  status="$(printf '%s' "$data" | curl "${args[@]}" 2> "$OUT/examples-http-$name.curl.log" || true)"
  printf 'status=%s expected=%s\n' "${status:-000}" "$want" \
    > "$OUT/examples-http-$name.status"

  MSYS2_ARG_CONV_EXCL='/api/v1/' python -c '
import json, sys
path, status, want, bodyfile = sys.argv[1:5]
raw = open(bodyfile, "rb").read()
try:
    body = json.loads(raw)
except Exception:
    body = raw.decode("utf-8", "replace")[:500]
print(json.dumps({"url": path, "status": int(status or 0), "expected": int(want), "body": body},
                 ensure_ascii=False))
' "$path" "${status:-0}" "$want" "$bodyfile" > "$jsonfile" 2> "$OUT/examples-http-$name.json.err"

  # Empty stderr captures are noise; keep only files with real diagnostics.
  for f in "$OUT/examples-http-$name.curl.log" "$OUT/examples-http-$name.json.err"; do
    [ -s "$f" ] || rm -f "$f"
  done

  if [ "$status" = "$want" ]; then pass "HTTP $method $path (see examples-http-$name.json)"; else
    fail "HTTP $method $path (see examples-http-$name.status)"
  fi
}

check_body() {
  local name="$1" expr="$2"
  printf '== body/%s == ' "$name"
  if python -c '
import json, sys
d = json.load(open(sys.argv[1], encoding="utf-8"))
b = d["body"]
assert eval(sys.argv[2]), (sys.argv[2], b)
print("assertion OK:", sys.argv[2])
' "$OUT/examples-http-$name.json" "$expr" > "$OUT/examples-http-$name.assert" 2>&1; then
    pass "body/$name (see examples-http-$name.assert)"
  else
    fail "body/$name (see examples-http-$name.assert)"
  fi
}

echo "=== examples/ acceptance ($LABEL) ==="
echo "evidence: $OUT"

# 1. Directory and documentation contract.
check_structure > "$OUT/structure.log" 2>&1
printf '== structure == '
if [ -s "$OUT/structure.log" ] && ! grep -qi 'missing\|still present\|still advertise' "$OUT/structure.log"; then
  pass "structure (see structure.log)"
else
  fail "structure (see structure.log)"
fi

# 2. Every example must compile and pass vet.
step go-build bash -c "cd '$ROOT' && go build ./..."
step go-vet   bash -c "cd '$ROOT' && go vet ./..."

# 3. The examples read the local SQLite file; make sure it exists.
if [ ! -f "$ROOT/fate" ]; then
  (cd "$ROOT" && go run ./cmd/dbinit) > "$OUT/examples-dbinit.log" 2>&1 || true
fi

# 4. Start a dedicated server and wait until it answers.
printf '== server == '
(cd "$ROOT" && go build -o "$TMP/fate-server" ./cmd/server) \
  > "$OUT/examples-server-build.log" 2>&1
"$TMP/fate-server" -addr ":$PORT" > "$OUT/examples-server.log" 2>&1 &
SERVER_PID=$!
ready=0
for _ in $(seq 1 60); do
  if curl -sS -o /dev/null "$BASE/" 2>/dev/null; then ready=1; break; fi
  sleep 0.5
done
if [ "$ready" -eq 1 ]; then pass "server on $PORT (see examples-server.log)"; else
  fail "server did not become ready (see examples-server.log)"
fi

# 5. Real wire evidence for the endpoints the http-client example uses.
GEN_BODY='{"surname":"张","born":"2024/06/15 10:30","sex":"boy","poetry_mode":0,"xi_yong_method":"balance","avoid_chars":"病,死","require_chars":"","filter_type":0}'
check_http generate POST "/api/v1/generate" "$GEN_BODY" 200
check_body generate 'isinstance(b, dict) and bool(b.get("task_id")) and b.get("state") == "done"'

TASK_ID="$(python -c '
import json, sys
print(json.load(open(sys.argv[1], encoding="utf-8"))["body"].get("task_id", ""))
' "$OUT/examples-http-generate.json" 2>/dev/null || true)"
echo "task_id=$TASK_ID" >> "$OUT/structure.log"

if [ -n "$TASK_ID" ]; then
  check_http status  GET "/api/v1/status?task_id=$TASK_ID" "" 200
  check_http result  GET "/api/v1/result?task_id=$TASK_ID" "" 200
  check_body result 'isinstance(b, dict) and bool(b.get("top_names")) and b["top_names"][0]["surname"] == "\u5f20"'
  check_http names   GET "/api/v1/names?task_id=$TASK_ID&page=1&size=5" "" 200
  check_body names 'isinstance(b, dict) and isinstance(b.get("names"), list)'
  check_http explore GET "/api/v1/explore?task_id=$TASK_ID&count=3" "" 200
  check_body explore 'isinstance(b, dict) and isinstance(b.get("names"), list)'
else
  fail "generate response had no task_id"
fi

# 6. Run every example for real and assert its observable output.
run_example quickstart \
  'generated [0-9]+ name results' \
  '^  #1 ' \
  'entries in the excellent table'
run_example stream \
  'batch 1: .*final=false' \
  'final=true' \
  'stream complete: [0-9]+ batches, final batch delivered last'
run_example filters \
  '^filter: avoid=' \
  'verified [0-9]+ results, none contains' \
  'first name:'
run_example explore \
  '^table: total=' \
  'after MarkShown' \
  'explore\(nil filter\):' \
  'explore complete:'
run_example http-client \
  'POST /api/v1/generate -> 200' \
  'GET /api/v1/result -> 200' \
  'names is a JSON array' \
  'http-client complete: 5 endpoints OK'

# 7. Repository gate and runtime checks for the public Go API.
step gate-self-test bash -c "cd '$ROOT' && python scripts/check_repo.py --self-test"
step gate-full      bash -c "cd '$ROOT' && python scripts/check_repo.py"
step nil-filter-integration bash -c "cd '$ROOT' && go test -tags integration . -run TestGenerateNilFilterCompletes -v -count=1 -timeout 120s"
step stream-weight-integration bash -c "cd '$ROOT' && go test -tags integration . -run TestGenerateStreamRecommendationWeight -v -count=1 -timeout 120s"

{
  echo "# examples/ 验收报告"
  echo
  echo "- 时间: $(date -Iseconds)"
  echo "- 版本标签: $LABEL"
  echo "- 服务端口: $PORT"
  echo "- 证据目录: tests/acceptance/$LABEL/evidence/"
  echo
  echo "| 检查 | 结果 |"
  echo "| :--- | :--- |"
  if [ "${#ROWS[@]}" -gt 0 ]; then printf '%s\n' "${ROWS[@]}"; fi
  echo
  if [ "$RESULT" -eq 0 ]; then echo "**Overall: PASS**"; else echo "**Overall: FAIL**"; fi
} > "$OUT/examples-verification.md"

echo
if [ "$RESULT" -eq 0 ]; then echo "**Overall: PASS**"; else echo "**Overall: FAIL**"; fi
echo "report: $OUT/examples-verification.md"
exit "$RESULT"
