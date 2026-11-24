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
    --json tagName,name,isDraft,isPrerelease,publishedAt,url,assets,body \
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
for a in sorted(data["assets"], key=lambda x: x["name"]):
    print(f"{a['name']:<40} {a.get('size', '?')} bytes  downloads={a.get('downloadCount', '?')}")
missing = sorted(required - names)
assert not missing, "missing assets: " + ", ".join(missing)
print("all", len(required), "required assets present")
PY
  then pass "all required assets present (see asset-manifest.log)"; else fail "asset manifest (see asset-manifest.log)"; fi

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
# A transient network failure can leave a truncated archive on disk; never
# checksum a partial download. Clear the directory and retry a few times.
: > "$OUT/gh-release-download.log"
DL_OK=0
for attempt in 1 2 3; do
  rm -f "$TMP/dl"/*
  if gh release download "$LABEL" --repo "$REPO" --dir "$TMP/dl" \
      --pattern 'checksums.txt' --pattern 'fate-*.tar.gz' --pattern 'fate-*.zip' \
      --pattern 'fate-*.sha256' >>"$OUT/gh-release-download.log" 2>&1; then
    DL_OK=1
    pass "assets downloaded on attempt $attempt (see gh-release-download.log)"
    break
  fi
  echo "download attempt $attempt failed; retrying" >> "$OUT/gh-release-download.log"
  sleep 3
done

if [ "$DL_OK" -eq 0 ]; then
  fail "asset download after 3 attempts (see gh-release-download.log)"
  fail "checksum verification (skipped: download incomplete)"
else
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
assert checked >= 5, f"only {checked} archives listed in checksums.txt"

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
fi

# 6. Release notes must be produced automatically by GitHub's generation API.
#    No hand-written notes file, no full-history git changelog.
log ""
log "## 6. Automated release notes"
TAG_SHA="$(git -C "$ROOT" rev-parse "$LABEL^{commit}")"
if gh api "repos/$REPO/releases/generate-notes" \
    -f tag_name="$LABEL" -f target_commitish="$TAG_SHA" \
    >"$OUT/generated-notes.json" 2>"$OUT/generated-notes.err"; then
  if python - "$OUT/release.json" "$OUT/generated-notes.json" >"$OUT/release-notes.log" 2>&1 <<'PY'
import json, sys
release = json.load(open(sys.argv[1], encoding="utf-8"))
generated = json.load(open(sys.argv[2], encoding="utf-8"))
notes = (generated.get("body") or "").strip()
body = (release.get("body") or "").strip()
print("generated notes length =", len(notes))
print("published body length  =", len(body))
assert notes, "GitHub returned empty generated notes"
assert body == notes, "published release notes differ from generated-notes API output"
assert len(body) < 20000, "release notes look like an auto-generated full-history changelog"
assert "Initial commit" not in body, "release notes contain raw repository history"
print("published body equals GitHub-generated notes")
PY
  then pass "release notes equal GitHub-generated output"; else fail "release notes (see release-notes.log)"; fi
else
  fail "generate-notes API call (see generated-notes.err)"
fi

# 7. Automation wiring and cleaned history: no manual-notes residue, local tag
#    refs identical to origin, no orphan draft releases.
log ""
log "## 7. Automation wiring and cleaned release/tag history"
if python - "$ROOT" "$REPO" >"$OUT/automation-history.log" 2>&1 <<'PY'
import json, os, subprocess, sys
root, repo = sys.argv[1], sys.argv[2]

gr = open(os.path.join(root, ".goreleaser.yml"), encoding="utf-8").read()
wf = open(os.path.join(root, ".github/workflows/release.yml"), encoding="utf-8").read()
assert "use: github-native" in gr, "GoReleaser not configured for github-native changelog"
assert "Apply curated release notes" not in wf, "curated notes step still present"
assert "--notes-file" not in wf, "manual notes upload still present"
assert not os.path.exists(os.path.join(root, "docs/release-notes")), "manual notes directory still present"
print("config: github-native changelog; no manual notes path")

proxy_env = os.environ.get("FATE_GIT_PROXY")
proxy = ["-c", f"http.proxy={proxy_env}", "-c", f"https.proxy={proxy_env}"] if proxy_env else []
remote = subprocess.run(["git"] + proxy + ["ls-remote", "--tags", "origin"], cwd=root,
                        capture_output=True, text=True)
assert remote.returncode == 0, remote.stderr
tags = {}
for line in remote.stdout.splitlines():
    parts = line.split()
    if len(parts) < 2:
        continue
    sha, ref = parts[0], parts[1].replace("refs/tags/", "")
    if ref.endswith("^{}"):
        tags[ref[:-3]] = sha
    elif ref not in tags:
        tags[ref] = sha
local = subprocess.run(["git", "tag", "-l"], cwd=root, capture_output=True, text=True).stdout.split()
# `backup-*` tags are deliberately local-only (e.g. backup-evidence-v4.3.0,
# created before the evidence history rewrite). They must never be pushed.
# Whitelist ONLY that prefix; every other local-only/remote-only tag is still
# a real inconsistency and must fail the gate.
backup = sorted(t for t in local if t.startswith("backup-"))
publishable = [t for t in local if not t.startswith("backup-")]
print("local backup tags (intentionally not pushed):", ", ".join(backup) or "none")
assert not [t for t in backup if t in tags], (
    "backup tag leaked to origin: " + ", ".join(t for t in backup if t in tags))
assert set(publishable) == set(tags), (
    f"tag names differ: local-only={sorted(set(publishable) - set(tags))} "
    f"remote-only={sorted(set(tags) - set(publishable))}")
bad = []
for t in tags:
    got = subprocess.run(["git", "rev-parse", "-q", "--verify", f"refs/tags/{t}^{{commit}}"],
                         cwd=root, capture_output=True, text=True).stdout.strip()
    if got != tags[t]:
        bad.append((t, got[:10], tags[t][:10]))
print("local tags:", len(local), "remote tags:", len(tags))
assert not bad, "tag SHA mismatch: " + str(bad)
print("all local tags match origin")

rel = subprocess.run(["gh", "api", f"repos/{repo}/releases?per_page=100"], cwd=root,
                     capture_output=True, text=True)
assert rel.returncode == 0, rel.stderr
orphans = [x["id"] for x in json.loads(rel.stdout) if x["tag_name"] == "v4.2.0"]
assert not orphans, "orphan v4.2.0 draft releases remain: " + str(orphans)
print("no orphan v4.2.0 releases")
PY
  then pass "automation wiring clean, tags consistent, no orphan drafts"; else fail "automation/history (see automation-history.log)"; fi

# 8. Validate .goreleaser.yml with the same GoReleaser version CI uses.
#    Set GORELEASER_BIN to enable; without it the check is reported as SKIPPED,
#    never silently counted as a pass.
log ""
log "## 8. GoReleaser configuration validation"
if [ -n "${GORELEASER_BIN:-}" ] && [ -x "$GORELEASER_BIN" ]; then
  if "$GORELEASER_BIN" check --config "$ROOT/.goreleaser.yml" \
      >"$OUT/goreleaser-check.log" 2>&1; then
    pass "goreleaser check: configuration validated"
  else
    fail "goreleaser check (see goreleaser-check.log)"
  fi
else
  log "- goreleaser check: SKIPPED (set GORELEASER_BIN to enable)"
fi

log ""
if [ "$FAILURES" -eq 0 ]; then
  log "**Overall: PASS**"
else
  log "**Overall: FAIL ($FAILURES check(s))**"
fi

rm -rf "$TMP"
exit "$FAILURES"
