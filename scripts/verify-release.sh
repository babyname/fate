#!/usr/bin/env bash
# Reproducible acceptance gate for a fate release.
#
#   bash scripts/verify-release.sh [v4.3.0]
#
# Produces evidence under tests/acceptance/<version>/evidence/ :
#   build/vet/test/tsc/web-build logs, real HTTP status + response bodies,
#   server log, README image link check, and report.md.
# Requires: Go, Node 22+, curl, python3. No external services.
set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LABEL="${1:-v4.3.0}"
OUT="$ROOT/tests/acceptance/$LABEL/evidence"
TMP="$(mktemp -d)"
PORT="${FATE_PORT:-18081}"
BASE="http://127.0.0.1:$PORT"
RESULT=0
ROWS=()
SERVER_PID=""

mkdir -p "$OUT"
: > "$OUT/report.md"
cleanup() {
  if [ -n "$SERVER_PID" ]; then kill "$SERVER_PID" 2>/dev/null || true; fi
  rm -rf "$TMP"
}
trap cleanup EXIT

step() {
  local name="$1"; shift
  local log="$OUT/$name.log"
  printf '== %s == ' "$name"
  if "$@" > "$log" 2>&1; then
    printf '0\n' > "$OUT/$name.status"
    echo PASS
    ROWS+=("| $name | PASS | [$name.log]($name.log) |")
  else
    printf '%s\n' "$?" > "$OUT/$name.status"
    echo FAIL
    ROWS+=("| $name | FAIL | [$name.log]($name.log) |")
    RESULT=1
  fi
}

# The integration tests open fate_test.db, which is intentionally untracked.
if [ ! -f "$ROOT/fate_test.db" ]; then
  if [ -f "$ROOT/fate" ]; then
    cp "$ROOT/fate" "$ROOT/fate_test.db"
  else
    (cd "$ROOT" && go run ./cmd/dbinit -db fate_test) > "$OUT/dbinit.log" 2>&1
  fi
fi

step go-build   bash -c "cd '$ROOT' && go build ./..."
step go-vet     bash -c "cd '$ROOT' && go vet ./..."
step go-test    bash -c "cd '$ROOT' && go test ./... -count=1"
step tsc        bash -c "cd '$ROOT/web' && node node_modules/typescript/bin/tsc --noEmit"
step web-build  bash -c "cd '$ROOT/web' && npm run build"

# README links must point at real screenshot files.
check_assets() {
  cd "$ROOT" || return 1
  local missing=0 f
  while IFS= read -r f; do
    if [ -f "$f" ]; then echo "OK $f"; else echo "MISSING $f"; missing=1; fi
  done < <(grep -hoE 'docs/screenshots/[A-Za-z0-9._-]+\.(png|jpg|svg)' README.md README_EN.md | sort -u)
  return "$missing"
}
step readme-assets check_assets

# Start a dedicated server instance so the captured log belongs to this run.
(cd "$ROOT" && go build -o "$TMP/fate-server" ./cmd/server) > "$OUT/server-build.log" 2>&1
(cd "$ROOT" && "$TMP/fate-server" -addr "127.0.0.1:$PORT") > "$OUT/server.log" 2>&1 &
SERVER_PID=$!
for _ in $(seq 1 60); do
  if curl -s -o /dev/null "$BASE/"; then break; fi
  sleep 0.5
done

check_http() {
  local name="$1" method="$2" url="$3" body="${4:-}" want="${5:-200}"
  local bodyfile="$OUT/http-$name.body" jsonfile="$OUT/http-$name.json"
  local code
  if [ "$method" = POST ]; then
    code=$(printf '%s' "$body" | curl -sS -X POST -H 'Content-Type: application/json' --data-binary @- \
      -o "$bodyfile" -w '%{http_code}' "$BASE$url")
  else
    code=$(curl -sS -X GET -o "$bodyfile" -w '%{http_code}' "$BASE$url")
  fi
  python -c '
import json, sys
code, want, url = int(sys.argv[1]), int(sys.argv[2]), sys.argv[3]
raw = sys.stdin.read()
try:
    payload = json.loads(raw)
except Exception:
    payload = raw[:4000]
print(json.dumps({"url": url, "status": code, "expected": want, "body": payload},
                 ensure_ascii=False, indent=2))
' "$code" "$want" "$url" < "$bodyfile" > "$jsonfile"
  if [ "$code" = "$want" ]; then
    echo "http $name: $code PASS"
    ROWS+=("| http-$name | PASS ($code) | [http-$name.json](http-$name.json) |")
  else
    echo "http $name: $code FAIL (want $want)"
    ROWS+=("| http-$name | FAIL ($code, want $want) | [http-$name.json](http-$name.json) |")
    RESULT=1
  fi
}

check_http root GET "/"
GEN_BODY='{"surname":"张","born":"2024/06/15 10:30","sex":"boy","poetry_mode":0}'
check_http generate POST "/api/v1/generate" "$GEN_BODY"
TASK_ID="$(python -c 'import json,sys; print(json.load(sys.stdin)["task_id"])' < "$OUT/http-generate.body" 2>/dev/null || true)"
if [ -z "$TASK_ID" ]; then
  echo "generate returned no task_id"; RESULT=1
fi
check_http status      GET "/api/v1/status?task_id=$TASK_ID"
check_http result      GET "/api/v1/result?task_id=$TASK_ID"
check_http names       GET "/api/v1/names?task_id=$TASK_ID&page=1&size=20"
check_http names-page-2 GET "/api/v1/names?task_id=$TASK_ID&page=2&size=20"
check_http names-out-of-range GET "/api/v1/names?task_id=$TASK_ID&page=9223372036854775807&size=20"
DETAIL_QUERY="$(python -c '
import json, sys, urllib.parse
name = json.load(sys.stdin)["names"][0]
print(urllib.parse.urlencode({"task_id": sys.argv[1], "char1": name["char1"]["char"], "char2": name["char2"]["char"]}))
' "$TASK_ID" < "$OUT/http-names-page-2.body" 2>/dev/null || true)"
check_http name-detail GET "/api/v1/name_detail?$DETAIL_QUERY"
check_http explore     GET "/api/v1/explore?task_id=$TASK_ID&count=6"
check_http characters  GET "/api/v1/characters?char=%E5%BC%A0"
check_http spa-fallback GET "/api/v1/name-score" "" 200

# The explore response must always expose an array, never JSON null.
check_explore_array() {
  python -c 'import json,sys; b=json.load(sys.stdin); assert isinstance(b.get("names"),list), b; print("names is a JSON array, length",len(b["names"]))' < "$OUT/http-explore.body"
}
step explore-array check_explore_array

check_candidate_pages() {
  python - "$OUT" <<'PYTHON'
import json
import pathlib
import sys

root = pathlib.Path(sys.argv[1])
def body(name):
    return json.loads((root / f"http-{name}.body").read_text(encoding="utf-8"))

first, second = body("names"), body("names-page-2")
total = body("generate")["total"]
print(f"Expected total > 20 and matching page totals. Actual generate={total}, page 1={first['total']}, page 2={second['total']}")
assert total > 20 and first["total"] == second["total"] == total
for p, offset in [(first, 0), (second, 20)]:
    wanted = list(range(offset + 1, min(offset + 20, total) + 1))
    found = [name["rank"] for name in p["names"]]
    print(f"Expected ranks {wanted}, actual {found}")
    assert p["page"] == offset // 20 + 1 and p["size"] == 20 and found == wanted
    assert all(name["surname"] == "张" for name in p["names"])
names = {name["full_name"] for name in first["names"] + second["names"]}
print(f"Expected distinct names: {len(first['names']) + len(second['names'])}, actual: {len(names)}")
assert len(names) == len(first["names"]) + len(second["names"])
assert body("names-out-of-range") == first
assert body("result")["top_names"] == first["names"][:10]
detail = body("name-detail")["name_result"]
print(f"Expected detail rank 21 and five-grid data. Actual rank: {detail['rank']}, five-grid data: {bool(detail['wu_ge'])}")
assert detail == second["names"][0] and detail["wu_ge"]
PYTHON
}
step candidate-pages check_candidate_pages

check_http weight-invalid POST "/api/v1/generate" '{"traditional_weight":101}' 400
check_http weight-zero POST "/api/v1/generate" "{\"task_id\":\"$TASK_ID\",\"traditional_weight\":0}"
WEIGHTED_ID="$(python -c 'import json,sys; print(json.load(sys.stdin)["task_id"])' < "$OUT/http-weight-zero.body" 2>/dev/null || true)"
check_http weighted-result GET "/api/v1/result?task_id=$WEIGHTED_ID"
check_http weighted-names GET "/api/v1/names?task_id=$WEIGHTED_ID&page=1&size=20"
check_http weighted-page-2 GET "/api/v1/names?task_id=$WEIGHTED_ID&page=2&size=20"
WEIGHTED_DETAIL="$(python -c '
import json, sys, urllib.parse
name = json.load(sys.stdin)["names"][0]
print(urllib.parse.urlencode({"task_id": sys.argv[1], "char1": name["char1"]["char"], "char2": name["char2"]["char"]}))
' "$WEIGHTED_ID" < "$OUT/http-weighted-page-2.body" 2>/dev/null || true)"
check_http weighted-detail GET "/api/v1/name_detail?$WEIGHTED_DETAIL"
check_http weight-hundred POST "/api/v1/generate" "{\"task_id\":\"$WEIGHTED_ID\",\"traditional_weight\":100}"
TRADITIONAL_ID="$(python -c 'import json,sys; print(json.load(sys.stdin)["task_id"])' < "$OUT/http-weight-hundred.body" 2>/dev/null || true)"
check_http weighted-traditional GET "/api/v1/names?task_id=$TRADITIONAL_ID&page=1&size=20"

check_score_weights() {
  python - "$OUT" <<'PYTHON'
import json
import pathlib
import sys

root = pathlib.Path(sys.argv[1])
def body(name):
    return json.loads((root / f"http-{name}.body").read_text(encoding="utf-8"))

legacy = body("names")["names"]
assert legacy and all("ranking" not in name for name in legacy)
first = body("weighted-names")["names"]
second = body("weighted-page-2")["names"]
assert len(first) == len(second) == 20
weighted = first + second
assert [name["rank"] for name in weighted] == list(range(1, 41))
assert [name["score"] for name in weighted] == sorted((name["score"] for name in weighted), reverse=True)
assert body("weighted-result")["top_names"] == first[:10]
assert body("weighted-detail")["name_result"] == second[0]
for expected, key, names in [(0, "sound_and_meaning", weighted), (100, "traditional", body("weighted-traditional")["names"])]:
    for name in names:
        ranking = name["ranking"]
        assert name["surname"] == "张"
        assert ranking["traditional_weight"] == expected and name["score"] == ranking["total"] == ranking[key]
    print(f"Expected traditional weight {expected} and total equal to {key}. Actual first ranking: {names[0]['ranking']}")
assert body("weight-zero")["task_id"] != body("generate")["task_id"]
assert body("weight-hundred")["task_id"] != body("weight-zero")["task_id"]
PYTHON
}
step score-weights check_score_weights

# The two unimplemented features must be gone from the frontend, and the old
# URL must be served by the SPA fallback (HTML), never by an API handler.
check_removed_ui() {
  cd "$ROOT" || return 1
  if grep -R -nE 'nameScore|searchPoetry|/name-score|/poetry/search' web/src; then
    echo "dead API calls still present in web/src"
    return 1
  fi
  echo "no dead API calls in web/src"
  python -c 'import sys; s=sys.stdin.read(); assert "<div id=\"root\">" in s, "not SPA HTML"; print("unimplemented path is SPA fallback, not an API endpoint")' < "$OUT/http-spa-fallback.body"
}
step removed-ui-contract check_removed_ui

{
  echo "# Acceptance evidence — $LABEL"
  echo
  echo "Generated at: $(date -Iseconds)"
  echo "Server base: $BASE"
  echo
  echo "| Check | Result | Evidence |"
  echo "|---|---|---|"
  for row in "${ROWS[@]}"; do echo "$row"; done
  echo
  echo "## Notes"
  echo "- HTTP evidence files contain the real status code and full response body."
  echo "- server.log is the log of the dedicated server started by this script."
  echo "- Screenshots are captured separately by scripts/capture-screenshots.mjs."
  echo "- Unimplemented name-score / poetry-search paths are served by the SPA fallback (HTML); their frontend entry points were removed."
  echo
  if [ "$RESULT" = 0 ]; then echo "**Overall: PASS**"; else echo "**Overall: FAIL**"; fi
} > "$OUT/report.md"

echo "report: $OUT/report.md"
exit "$RESULT"
