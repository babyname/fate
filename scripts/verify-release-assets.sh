#!/usr/bin/env bash
# Reproducible acceptance for a FATE GitHub Release.
#
#   bash scripts/verify-release-assets.sh [label]     # default: v4.3.0
#
# Captures real GitHub API responses, workflow conclusions, cross-build results
# and sha256 evidence into tests/acceptance/<label>/evidence/. Every assertion
# is made against observed data, never against a hard-coded "green" expectation.

set -uo pipefail

LABEL="${1:-v4.3.0}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO="${FATE_REPO:-babyname/fate}"
OUT="$ROOT/tests/acceptance/$LABEL/evidence"
TMP="$(mktemp -d "${TMPDIR:-/tmp}/fate-release-verify.XXXXXX")"
mkdir -p "$OUT" "$TMP/dl"

REPORT="$OUT/release-verification.md"
: > "$REPORT"
FAILURES=0

log()  { printf '%s\n' "$*" | tee -a "$REPORT"; }
pass() { log "- $1: PASS"; }
fail() { log "- $1: FAIL"; FAILURES=$((FAILURES + 1)); }

TARGETS="windows/amd64 darwin/amd64 darwin/arm64 linux/amd64 linux/arm64"
DB_VER="${LABEL#v}"

log "# Release acceptance: $LABEL"
log ""
log "repo=$REPO  evidence=tests/acceptance/$LABEL/evidence/"
log ""

# 1. Root-cause regression: every GoReleaser target must build with CGO off.
log "## 1. Pure-Go cross-build for every GoReleaser target"
for t in $TARGETS; do
  os=${t%/*}; arch=${t#*/}; out="$TMP/fate-$os-$arch"
  if (cd "$ROOT" && CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" \
      go build -o "$out" ./cmd/console) >"$OUT/crossbuild-$os-$arch.log" 2>&1; then
    size=$(stat -c%s "$out" 2>/dev/null || echo 0)
    pass "CGO_ENABLED=0 $os/$arch (binary $size bytes)"
  else
    fail "CGO_ENABLED=0 $os/$arch (see crossbuild-$os-$arch.log)"
  fi
done

# 2. Real GitHub Release API response.
log ""
log "## 2. GitHub Release API"
if gh release view "$LABEL" --repo "$REPO" \
    --json tagName,name,isDraft,isPrerelease,publishedAt,url,assets \
    >"$OUT/release.json" 2>"$OUT/gh-release-view.err"; then
  pass "gh release view $LABEL (response saved to release.json)"
else
  fail "gh release view $LABEL (see gh-release-view.err)"
fi

if [ -s "$OUT/release.json" ]; then
  if python - "$OUT/release.json" "$LABEL" >"$OUT/release-api.log" 2>&1 <<'PY'
import json, sys
data = json.load(open(sys.argv[1], encoding="utf-8"))
label = sys.argv[2]
assert data["tagName"] == label, data
assert data["isDraft"] is False, "release is still a draft"
print("tagName      =", data["tagName"])
print("name         =", data.get("name"))
print("publishedAt  =", data.get("publishedAt"))
print("isDraft      =", data["isDraft"])
print("isPrerelease =", data.get("isPrerelease"))
print("url          =", data.get("url"))
print("assets       =", len(data["assets"]))
PY
  then pass "release is published (not a draft), tag matches"; else fail "release payload assertions (see release-api.log)"; fi
else
  fail "release payload assertions (no response body)"
fi

# 3. Asset manifest: 5 platform archives, checksums, database bundle.
log ""
log "## 3. Release asset manifest"
if python - "$OUT/release.json" "$DB_VER" >"$OUT/asset-manifest.log" 2>&1 <<'PY'
import json, sys
data = json.load(open(sys.argv[1], encoding="utf-8"))
dbver = sys.argv[2]
names = {a["name"] for a in data["assets"]}
required = {
    "checksums.txt",
    "fate-windows-amd64.zip",
    "fate-darwin-amd64.tar.gz",
    "fate-darwin-arm64.tar.gz",
    "fate-linux-amd64.tar.gz",
    "fate-linux-arm64.tar.gz",
    f"fate-v{dbver}.zip",
    f"fate-v{dbver}.zip.sha256",
}
for a in sorted(data["assets"]):
    print(f"{a['name']:<40} {a.get('size', '?')} bytes  downloads={a.get('downloadCount', '?')}")
missing = sorted(required - names)
assert not missing, "missing assets: " + ", ".join(missing)
print("all", len(required), "required assets present")
PY
  then pass "all required assets present (see asset-manifest.log)"; else fail "asset manifest (see asset-manifest.log)"; fi
else
  fail "asset manifest (no response body)"
fi

# 4. Workflow conclusion for this tag.
log ""
log "## 4. Release workflow conclusion"
gh run list --repo "$REPO" --workflow=release.yml --limit 5 \
  --json databaseId,conclusion,status,headBranch,displayTitle,createdAt,event \
  >"$OUT/workflow-runs.json" 2>"$OUT/gh-run-list.err"
if python - "$OUT/workflow-runs.json" "$LABEL" >"$OUT/workflow-conclusion.log" 2>&1 <<'PY'
import json, sys
runs = json.load(open(sys.argv[1], encoding="utf-8"))
label = sys.argv[2]
matching = [r for r in runs if r.get("headBranch") == label]
assert matching, f"no release workflow run for tag {label}; runs={runs}"
for r in matching:
    print(f"run {r['databaseId']}  {r['createdAt']}  status={r['status']}  conclusion={r['conclusion']}")
latest = matching[0]
assert latest["status"] == "completed", latest
assert latest["conclusion"] == "success", "release workflow did not succeed: " + str(latest)
print("latest release workflow conclusion = success")
PY
  then pass "release workflow concluded successfully"; else fail "release workflow conclusion (see workflow-conclusion.log)"; fi

# 5. Download assets and verify checksums.txt + database sha256.
log ""
log "## 5. Checksum integrity of downloaded assets"
if gh release download "$LABEL" --repo "$REPO" --dir "$TMP/dl" \
    --pattern 'checksums.txt' --pattern 'fate-*.tar.gz' --pattern 'fate-*.zip' \
    --pattern 'fate-*.sha256' >"$OUT/gh-release-download.log" 2>&1; then
  pass "assets downloaded (see gh-release-download.log)"
else
  fail "asset download (see gh-release-download.log)"
fi

if python - "$TMP/dl" "$DB_VER" >"$OUT/checksums.log" 2>&1 <<'PY'
import hashlib, os, sys
d, dbver = sys.argv[1], sys.argv[2]

def sha256(path):
    h = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(1 << 20), b""):
            h.update(chunk)
    return h.hexdigest()

csum = os.path.join(d, "checksums.txt")
assert os.path.exists(csum), "checksums.txt was not downloaded"
checked = 0
for line in open(csum, encoding="utf-8"):
    line = line.strip()
    if not line:
        continue
    expected, name = line.split()[0], line.split()[-1]
    path = os.path.join(d, name)
    assert os.path.exists(path), "archive missing: " + name
    actual = sha256(path)
    assert actual == expected, f"{name}: {actual} != {expected}"
    print(f"sha256 OK  {expected[:16]}  {name}")
    checked += 1
assert checked >= 6, f"only {checked} archives listed in checksums.txt"

dbzip = os.path.join(d, f"fate-v{dbver}.zip")
dbsha = os.path.join(d, f"fate-v{dbver}.zip.sha256")
assert os.path.exists(dbzip) and os.path.exists(dbsha), "database bundle or its .sha256 missing"
expected = open(dbsha, encoding="utf-8").read().split()[0]
actual = sha256(dbzip)
assert actual == expected, f"database bundle: {actual} != {expected}"
print(f"sha256 OK  {expected[:16]}  fate-v{dbver}.zip")
print("verified", checked, "archives + 1 database bundle")
PY
  then pass "all downloaded assets match their published sha256"; else fail "checksum verification (see checksums.log)"; fi

log ""
if [ "$FAILURES" -eq 0 ]; then
  log "**Overall: PASS**"
else
  log "**Overall: FAIL ($FAILURES check(s))**"
fi

rm -rf "$TMP"
exit "$FAILURES"
