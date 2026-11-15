#!/usr/bin/env python3
"""FATE repository gate -- executable rules for AGENTS.md §6.

Why this exists
---------------
Text-only conventions rot. This gate turns the rules in AGENTS.md and
docs/ACCEPTANCE.md into deterministic checks, and ships a --self-test mode
that PROVES every detector can actually fire (positive control) and does not
fire on legitimate input (negative control). A gate that silently cannot
match its own pattern is worse than no gate.

Modes
-----
  python scripts/check_repo.py                 full static gate
  python scripts/check_repo.py --staged       pre-commit gate (staged set only)
  python scripts/check_repo.py --ci           CI gate (adds HEAD~1 range)
  python scripts/check_repo.py --self-test    detector positive/negative controls
  python scripts/check_repo.py --commits A..B  also audit a commit range

Exit code 0 = PASS, 1 = at least one FAIL, 2 = usage/environment error.
"""
from __future__ import annotations

import argparse
import re
import subprocess
import sys
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
FAILURES = 0
RESULTS: list[tuple[str, str, str]] = []

# Commit-message rules take effect AFTER this commit. Older messages are
# historical debt (216 CR bytes in the first 20 messages alone); blocking on
# them would make CI permanently red and turn the gate into noise.
BASELINE = "1f0cfa3"

TEMP_SUFFIXES = (".tmp", ".bak", ".orig", ".swp", ".swo", ".log")
AI_PREFIXES = (
    ".trae/", ".claude/", ".cursor/", ".workbuddy/", ".codex/", ".aider/",
    ".augment/", ".windsurf/", ".agents/", ".continue/",
)
AI_EXACT = {"CLAUDE.md", ".cursorrules", ".cursorignore", "skills-lock.json"}

# Machine-absolute paths are non-portable. System browser lookup paths
# (C:/Program Files/...) are legitimate and deliberately not matched here.
MACHINE_PATH = re.compile(
    r"(?:^|[^A-Za-z])[A-Za-z]:[\\/](?:Users|workspace|project)[\\/]|/[cCdD]/Users/"
)
CJK = re.compile(r"[\u3000-\u303f\u3400-\u4dbf\u4e00-\u9fff\uf900-\ufaff\uff00-\uffef]")
CONVENTIONAL = re.compile(
    r"^(?:feat|fix|refactor|docs|chore|test|style|perf|ci|build|release)"
    r"(?:\([a-z0-9_.-]+\))?!?: .+"
)


# --------------------------------------------------------------------------
# Pure detectors (all self-tested below)
# --------------------------------------------------------------------------
def is_temp_name(path: str) -> bool:
    p = path.replace("\\", "/")
    # Evidence paths are handled by is_evidence_path (one clear message);
    # a bare `.log` suffix alone is still a temp-name signal.
    if p.startswith("tests/") and "/evidence/" in p:
        return False
    return p.endswith(TEMP_SUFFIXES) or p.endswith("~")


def is_ai_artifact(path: str) -> bool:
    p = path.replace("\\", "/")
    return p.startswith(AI_PREFIXES) or p in AI_EXACT or "/.workbuddy/" in p


def is_evidence_log(path: str) -> bool:
    # A successful `go build` / `go vet` produces no output, so its evidence
    # log is legitimately empty. The summary report carries the exit status.
    p = path.replace("\\", "/")
    return p.startswith("tests/acceptance/") and "/evidence/" in p


def is_evidence_path(path: str) -> bool:
    """Acceptance evidence must never be committed.

    Evidence is generated per-run (HTTP bodies, logs, probe files) and often
    empty on success. 83 such files (506K) once polluted main because a
    `.gitignore` negation re-included them. This detector is the gate side of
    that fix: any tracked or staged evidence path is pollution.
    """
    p = path.replace("\\", "/")
    return p.startswith("tests/") and "/evidence/" in p


def has_machine_path(data: bytes) -> bool:
    return bool(MACHINE_PATH.search(data.decode("utf-8", "replace")))


def has_bom(data: bytes) -> bool:
    return data.startswith(b"\xef\xbb\xbf")


def has_cjk(text: str) -> bool:
    return bool(CJK.search(text))


def has_cr(data: bytes) -> bool:
    return b"\r" in data


def is_conventional(subject: str) -> bool:
    return bool(CONVENTIONAL.match(subject)) and not has_cjk(subject)


def local_refs(html: str) -> set[str]:
    """Local /js/... and /css/... paths referenced by index.html."""
    return {
        m.lstrip("/")
        for m in re.findall(r'(?:src|href)="(/[^"]+)"', html)
        if not m.startswith("//")
    }


# --------------------------------------------------------------------------
# Infrastructure
# --------------------------------------------------------------------------
def git(*args: str, text: bool = True) -> str | bytes:
    r = subprocess.run(["git", *args], cwd=ROOT, capture_output=True)
    out: str | bytes = r.stdout.decode("utf-8", "replace") if text else r.stdout
    return out


def read(rel: str) -> str:
    return (ROOT / rel).read_text(encoding="utf-8-sig")


def read_bytes(rel: str) -> bytes | None:
    p = ROOT / rel
    return p.read_bytes() if p.is_file() else None


def tracked(path: str | None = None) -> list[str]:
    args = ["ls-files"] + ([path] if path else [])
    return [x for x in git(*args).splitlines() if x]


def changed(paths_from: str) -> list[str]:
    out = git("diff", "--name-only", "--diff-filter=ACMR", paths_from)
    out += git("diff", "--cached", "--name-only", "--diff-filter=ACMR")
    out += git("ls-files", "--others", "--exclude-standard")
    return sorted({x for x in out.splitlines() if x})


def staged() -> list[str]:
    return [x for x in git("diff", "--cached", "--name-only", "--diff-filter=ACMR").splitlines() if x]


def untracked() -> list[str]:
    # These are exactly the files a careless `git add -A` would commit.
    return [x for x in git("ls-files", "--others", "--exclude-standard").splitlines() if x]


def record(name: str, ok: bool, detail: str = "") -> None:
    global FAILURES
    RESULTS.append((name, "PASS" if ok else "FAIL", detail))
    if not ok:
        FAILURES += 1


def run_cmd(name: str, argv: list[str]) -> None:
    r = subprocess.run(argv, cwd=ROOT, capture_output=True, text=True)
    tail = (r.stdout + r.stderr).strip().splitlines()[-6:]
    record(name, r.returncode == 0, "; ".join(tail))


def gofmt_clean(data: bytes) -> bool:
    """Compare gofmt output byte-for-byte, ignoring only line endings.

    `gofmt -l` on a Windows checkout reports CRLF-only differences and BOMs,
    which have their own checks. Feeding stdin makes the verdict identical on
    Windows, Linux and CI, and avoids the `grep -q .` false-negative trap
    where a failed gofmt invocation produces empty output and passes.
    """
    norm = data.replace(b"\r\n", b"\n").removeprefix(b"\xef\xbb\xbf")
    r = subprocess.run(["gofmt"], input=norm, capture_output=True)
    return r.returncode == 0 and r.stdout == norm


def check_gofmt(files: list[str]) -> None:
    # 29 historical files are not gofmt-clean (23 carry a UTF-8 BOM, 6 have
    # stale struct alignment). Blocking on them would make the gate red
    # forever, so this rule applies to CHANGED files: new code is zero-tolerance.
    targets = [f for f in files if f.endswith(".go") and not f.startswith("ent/")]
    dirty = [f for f in targets if not gofmt_clean(read_bytes(f) or b"")]
    record(f"code: gofmt (changed Go files, n={len(targets)})", not dirty, "; ".join(dirty))


# --------------------------------------------------------------------------
# Checks
# --------------------------------------------------------------------------
SCRIPT_SUFFIXES = (".sh", ".py", ".mjs", ".cjs", ".js", ".ps1", ".bat", ".cmd")


def check_pollution(files: list[str], label: str) -> None:
    bad: list[str] = []
    for f in files:
        if is_evidence_path(f):
            bad.append(f"acceptance evidence committed (test pollution): {f}")
        if is_temp_name(f):
            bad.append(f"temp/binary artifact: {f}")
        if is_ai_artifact(f):
            bad.append(f"AI tool artifact: {f}")
        data = read_bytes(f)
        # Zero-byte files are probe/scratch leftovers (e.g. a stray `status`
        # or `explore.` created by a shell redirect). No source file is empty.
        if data == b"" and not is_evidence_log(f):
            bad.append(f"empty probe file: {f}")
        if f.endswith(SCRIPT_SUFFIXES) and data and has_machine_path(data):
            bad.append(f"machine-absolute path in script: {f}")
    record(f"pollution [{label}]", not bad, "; ".join(bad))


def check_evidence_untracked() -> None:
    """Head-level guard: evidence under tests/**/evidence/ must not be tracked.

    check_pollution only sees the staged/changed set; this check walks the
    whole tracked index so pollution already sitting in HEAD is caught too
    (it fired on the 83 files from the v4.3.0 commits).
    """
    bad = sorted(f for f in tracked() if is_evidence_path(f))
    record("evidence: no acceptance evidence tracked in git", not bad,
           f"{len(bad)} tracked evidence files; run: git rm -r --cached "
           f"tests/acceptance/*/evidence/" + ("; " + "; ".join(bad[:5]) if bad else ""))


def check_static_assets() -> None:
    """The root cause of the recurring resources/static dirt.

    Rsbuild writes CONTENT-HASHED bundles into the tracked resources/static
    directory. Before cleanDistPath was enabled, every build ADDED new hashed
    files without deleting the old ones (js files: 7 -> 12 -> 13 -> 19 -> 29
    -> 36), so `git add -A` kept committing obsolete bundles. This check makes
    that state unrepresentable: every tracked bundle must be referenced by
    index.html (or be its .LICENSE.txt sidecar).
    """
    html_path = "resources/static/index.html"
    html = read(html_path)
    refs = local_refs(html)
    prefix = "resources/static/"
    tracked_rel = {p[len(prefix):] for p in tracked("resources/static") if p.startswith(prefix)}

    missing = sorted(r for r in refs if r not in tracked_rel)
    sidecars = {r + ".LICENSE.txt" for r in refs}
    extra = sorted(tracked_rel - {"index.html"} - refs - sidecars)
    record("static: every referenced bundle is tracked", not missing, "; ".join(missing))
    record("static: no unreferenced bundles", not extra, "; ".join(extra))

    on_disk = {str(p.relative_to(ROOT / "resources/static")).replace("\\", "/")
               for p in (ROOT / "resources/static").rglob("*") if p.is_file()}
    stray = sorted(on_disk - {"index.html"} - refs - sidecars)
    record("static: no stray files on disk", not stray, "; ".join(stray))

    cfg = read("web/rsbuild.config.ts")
    record("static: cleanDistPath enabled", "cleanDistPath: true" in cfg)
    record("static: bun.lock is tracked", "web/bun.lock" in tracked())


def check_contracts() -> None:
    handler = read("internal/http/handler.go")
    routes = set(re.findall(r'HandleFunc\("(?:GET|POST) (/api/v1/[^"]+)"', handler))
    dev = read("docs/DEVELOPMENT.md")
    undocumented = sorted(r for r in routes if r not in dev)
    record("contract: all /api/v1 routes documented", len(routes) >= 7 and not undocumented,
           f"routes={len(routes)} undocumented={undocumented}")

    legacy = [f for f in ("README.md", "README_EN.md", "docs/DEVELOPMENT.md")
              if "/api/generate" in read(f) or "/api/name-detail" in read(f)]
    record("contract: no legacy /api/* URLs in docs", not legacy, "; ".join(legacy))

    actual = {p.name for p in (ROOT / "examples").iterdir()
              if p.is_dir() and (p / "main.go").is_file()}
    mentioned = set(re.findall(r"`([a-z-]+)/`", read("examples/README.md")))
    listed = {m for m in mentioned if (ROOT / "examples" / m).is_dir()}
    record("contract: examples match README", actual == listed and len(actual) == 5,
           f"actual={sorted(actual)} listed={sorted(listed)}")


def check_versions() -> None:
    v = read("VERSION").strip()
    db = read("DB_VERSION").strip()
    m = re.search(r'const Version = "v([0-9.]+)"', read("version.go"))
    code = m.group(1) if m else ""
    record("versions: VERSION/DB_VERSION/version.go agree",
           bool(v) and v == db == code, f"VERSION={v} DB_VERSION={db} version.go={code}")


def check_bom(files: list[str], label: str) -> None:
    bad = [f for f in files if f.endswith(".go") and (read_bytes(f) or b"")[:3] == b"\xef\xbb\xbf"]
    record(f"bom: no BOM in changed Go files [{label}]", not bad, "; ".join(bad))


def check_acceptance_assets() -> None:
    required = [
        "AGENTS.md", "docs/ACCEPTANCE.md", "scripts/verify-release.sh",
        "scripts/verify-examples.sh", "scripts/verify-release-assets.sh",
        "scripts/check_repo.py", ".github/workflows/ci.yml",
    ]
    present = set(tracked()) | {p for p in required if (ROOT / p).is_file()}
    missing = [p for p in required if p not in present]
    record("acceptance: gate + scripts present", not missing, "; ".join(missing))


def check_commits(rng: str | None) -> None:
    if not rng:
        rng = "HEAD~1..HEAD"
    shas = [s for s in git("rev-list", rng).split() if s]
    if not shas:
        record("commits: range has commits", False, f"empty range {rng}")
        return
    bad: list[str] = []
    for sha in shas:
        # Read the raw message bytes: text=True would translate CRLF and hide CR.
        msg = git("log", "-1", "--format=%B", sha, text=False)
        subject = msg.split(b"\n", 1)[0].decode("utf-8", "replace").strip()
        if subject.startswith("Merge ") or subject.startswith("Revert "):
            continue
        # Only commits AFTER the baseline are judged; older ones are debt.
        if BASELINE and subprocess.run(
            ["git", "merge-base", "--is-ancestor", sha, BASELINE], cwd=ROOT, capture_output=True
        ).returncode == 0:
            continue
        if has_cr(msg):
            bad.append(f"{sha[:7]}: CR in message")
        if not is_conventional(subject):
            bad.append(f"{sha[:7]}: non-conventional/non-English subject: {subject!r}")
    record(f"commits [{rng}]", not bad, "; ".join(bad))


# --------------------------------------------------------------------------
# Positive/negative controls -- a detector that cannot fire is a false gate
# --------------------------------------------------------------------------
def self_test() -> int:
    global FAILURES
    win_abs = b'const P = "' + b"D:" + b"\\" + b"Users" + b"\\" + b'someone"'
    cjk_msg = chr(0x4FEE) + chr(0x590D) + " gate"

    cases: list[tuple[str, object, object, object]] = [
        ("temp name", is_temp_name, "build/foo.tmp", "build/foo.go"),
        ("temp name (evidence exempt)", lambda p: not is_temp_name(p),
         "tests/acceptance/v4.3.0/evidence/a.log", "x.log"),
        ("AI artifact", is_ai_artifact, ".workbuddy/memory/x.md", "docs/rules/x.md"),
        ("machine path", has_machine_path, win_abs,
         b'const CHROME = "C:/Program Files/Google/Chrome/Application/chrome.exe"'),
        ("BOM", has_bom, b"\xef\xbb\xbfpackage main", b"package main"),
        ("CJK message", has_cjk, cjk_msg, "ci(gate): add checks"),
        ("CR message", has_cr, b"subject\r\nbody", b"subject\nbody"),
        ("conventional subject", is_conventional, "fix(ci): add checks", "random words"),
        ("local refs", lambda h: len(local_refs(h)) == 2,
         '<script src="/js/a.js"></script><link href="/css/b.css">',
         '<link href="https://fonts.example/x"><script src="/js/a.js"></script>'),
        ("gofmt detects bad format", lambda d: not gofmt_clean(d),
         b"package main\nfunc  main( ) { }\n", b"package main\n\nfunc main() {}\n"),
        ("empty probe file", lambda d: d == b"", b"", b"package main\n"),
        ("evidence log exempt (empty ok)", is_evidence_log,
         "tests/acceptance/v4.3.0/evidence/go-build.log", "status"),
        ("evidence path flagged", is_evidence_path,
         "tests/acceptance/v4.3.0/evidence/go-build.log", "docs/ACCEPTANCE.md"),
        ("evidence dir flagged", is_evidence_path,
         "tests\\acceptance\\v4.3.0\\evidence\\x.body", "tests/acceptance/v4.3.0/report.md"),
    ]
    print("== detector self-test (positive control must fire, negative must not) ==")
    for name, fn, positive, negative in cases:  # type: ignore[misc]
        fired = bool(fn(positive))  # type: ignore[operator]
        clean = not bool(fn(negative))  # type: ignore[operator]
        ok = fired and clean
        print(f"  {'PASS' if ok else 'FAIL'}  {name:32s} fired={fired} clean={clean}")
        if not ok:
            FAILURES += 1

    # End-to-end control: a fake static tree with an unreferenced bundle must
    # be detected by the same set arithmetic used in check_static_assets.
    with tempfile.TemporaryDirectory() as td:
        t = Path(td)
        (t / "index.html").write_text('<script src="/js/a.js"></script>', encoding="utf-8")
        (t / "js").mkdir()
        (t / "js" / "a.js").write_text("ok", encoding="utf-8")
        (t / "js" / "stale.DEADBEEF.js").write_text("old", encoding="utf-8")
        refs = local_refs((t / "index.html").read_text(encoding="utf-8"))
        on_disk = {str(p.relative_to(t)).replace("\\", "/") for p in t.rglob("*") if p.is_file()}
        detected = sorted(on_disk - {"index.html"} - refs - {r + ".LICENSE.txt" for r in refs})
        ok = detected == ["js/stale.DEADBEEF.js"]
        print(f"  {'PASS' if ok else 'FAIL'}  stale-bundle e2e control          detected={detected}")
        if not ok:
            FAILURES += 1

        # gofmt control: CRLF and BOM are normalised away (they have their own
        # checks), but a genuine formatting difference must still be caught.
        ok = (gofmt_clean(b"\xef\xbb\xbfpackage main\r\n\r\nfunc main() {}\r\n")
              and not gofmt_clean(b"package main\nfunc  main( ) { }\n"))
        print(f"  {'PASS' if ok else 'FAIL'}  gofmt normalisation control")
        if not ok:
            FAILURES += 1

    print(f"== self-test: {'PASS' if FAILURES == 0 else 'FAIL'} ({FAILURES} failures) ==")
    return 1 if FAILURES else 0


# --------------------------------------------------------------------------
def main() -> int:
    ap = argparse.ArgumentParser(description="FATE repository gate")
    ap.add_argument("--staged", action="store_true", help="pre-commit: staged set only")
    ap.add_argument("--ci", action="store_true", help="CI: use HEAD~1 range, run Go checks")
    ap.add_argument("--self-test", action="store_true", help="detector positive/negative controls")
    ap.add_argument("--commits", metavar="RANGE", help="also audit a commit range")
    args = ap.parse_args()

    if args.self_test:
        return self_test()

    print(f"FATE repository gate  mode={'ci' if args.ci else 'staged' if args.staged else 'full'}")
    print(f"root={ROOT}\n")

    if args.staged:
        files = staged()
        if not files:
            print("no staged files; nothing to check")
            return 0
        check_pollution(files, "staged")
        check_evidence_untracked()
        check_bom(files, "staged")
        check_static_assets()
        check_contracts()
        check_versions()
    else:
        check_pollution(tracked() + untracked(), "tracked+untracked")
        check_evidence_untracked()
        check_static_assets()
        check_contracts()
        check_versions()
        check_acceptance_assets()
        rng = "HEAD~1" if args.ci else "HEAD"
        changed_files = changed(rng)
        check_bom(changed_files, rng)
        if args.ci or args.commits:
            check_commits(args.commits)
        check_gofmt(changed_files)
        run_cmd("code: go build ./...", ["go", "build", "./..."])
        run_cmd("code: go vet ./...", ["go", "vet", "./..."])

    for name, status, detail in RESULTS:
        line = f"{status}  {name}"
        if detail and status == "FAIL":
            line += f"\n        {detail}"
        print(line)
    overall = FAILURES == 0
    print(f"\n**Overall: {'PASS' if overall else 'FAIL'}**  ({len(RESULTS)} checks, {FAILURES} failures)")
    return 0 if overall else 1


if __name__ == "__main__":
    sys.exit(main())
