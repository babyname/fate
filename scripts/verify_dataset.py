#!/usr/bin/env python3
"""Verify character.json against the engine's data contract.

The engine (fate) is the consumer, so the consumer schema is the single source
of truth for this contract. Any producer that supplies character data must emit
a file that passes this checker. No knowledge of the producer is encoded here.

Usage:
    python scripts/verify_dataset.py                 # check resources/character.json
    python scripts/verify_dataset.py --self-test     # detector must-hit / must-miss cases
"""
from __future__ import annotations

import argparse
import json
import re
import sys
from collections import Counter
from pathlib import Path
from typing import Any, Callable

DATASET = Path(__file__).resolve().parents[1] / "resources" / "character.json"

# ---- Contract: single source of truth, derived from the consumer schema ----
EXPECTED_COUNT = 30060
FIVE_ELEMENTS = {"木", "火", "水", "金", "土"}
UNICODE_RE = re.compile(r"^U\+[0-9A-F]{4,6}$")

# field -> (runtime type, required). "required" means present and non-null on
# every record; optional fields may be absent or null (sparse in real data).
FIELD_SPEC: dict[str, tuple[type, bool]] = {
    "char": (str, True),
    "unicode": (str, True),
    "radical": (str, True),
    "radical_stroke": (int, True),
    "simplified_stroke": (int, True),
    "traditional_stroke": (int, True),
    "kangxi_stroke": (int, False),
    "science_stroke": (int, True),
    "wu_xing": (str, True),
    "source": (str, True),
    "is_simplified": (bool, False),
    "is_traditional": (bool, False),
    "is_kangxi": (bool, False),
    "is_variant": (bool, False),
    "regular": (bool, False),
    "nameable": (bool, False),
    "pinyin": (list, False),
    "meaning": (str, False),
    "simplified_of_char": (str, False),
    "variant_of_char": (str, False),
}


def _is_int(value: Any) -> bool:
    # bool is a subclass of int in Python; the contract treats them as distinct.
    return isinstance(value, int) and not isinstance(value, bool)


def _type_ok(value: Any, expected: type) -> bool:
    if expected is int:
        return _is_int(value)
    if expected is bool:
        return isinstance(value, bool)
    return isinstance(value, expected)


def validate(records: list[dict[str, Any]]) -> list[str]:
    """Return a list of contract violations (empty == pass)."""
    errors: list[str] = []

    if len(records) != EXPECTED_COUNT:
        errors.append(f"count: expected {EXPECTED_COUNT}, got {len(records)}")

    seen: set[str] = set()
    for idx, rec in enumerate(records):
        if not isinstance(rec, dict):
            errors.append(f"[{idx}] record is {type(rec).__name__}, expected dict")
            continue

        for field, (expected, required) in FIELD_SPEC.items():
            if field not in rec or rec[field] is None:
                if required:
                    errors.append(f"[{idx}] required field missing/null: {field}")
                continue
            if not _type_ok(rec[field], expected):
                errors.append(
                    f"[{idx}] field {field}: expected {expected.__name__}, "
                    f"got {type(rec[field]).__name__}"
                )

        char = rec.get("char")
        if isinstance(char, str):
            if not char:
                errors.append(f"[{idx}] char is empty")
            elif char in seen:
                errors.append(f"[{idx}] duplicate char: {char}")
            else:
                seen.add(char)

        uni = rec.get("unicode")
        if uni is not None and (not isinstance(uni, str) or not UNICODE_RE.match(uni)):
            errors.append(f"[{idx}] unicode malformed: {uni!r}")

        wx = rec.get("wu_xing")
        if wx is not None and wx not in FIVE_ELEMENTS:
            errors.append(f"[{idx}] wu_xing out of enum: {wx!r}")

        pinyin = rec.get("pinyin")
        if pinyin is not None and isinstance(pinyin, list):
            if any(not isinstance(p, str) for p in pinyin):
                errors.append(f"[{idx}] pinyin contains non-string")

    return errors


def _load() -> list[dict[str, Any]]:
    with DATASET.open(encoding="utf-8") as fh:
        data = json.load(fh)
    if not isinstance(data, list):
        raise SystemExit(f"dataset root is {type(data).__name__}, expected list")
    return data


def self_test() -> int:
    """Prove each detector fires on a bad case and stays quiet on a good case."""
    good: dict[str, Any] = {
        "char": "木", "unicode": "U+6728", "wu_xing": "木",
        "radical": "木", "radical_stroke": 0, "simplified_stroke": 4,
        "traditional_stroke": 4, "kangxi_stroke": 4, "science_stroke": 4,
        "source": "test", "pinyin": ["mù"],
    }
    cases: list[tuple[str, Callable[[list[dict[str, Any]]], bool], str]] = [
        ("count mismatch", lambda rs: any("count:" in e for e in validate(rs)), "bad"),
        ("required missing", lambda rs: any("required field" in e for e in validate(rs)), "bad"),
        ("type mismatch", lambda rs: any("expected" in e for e in validate(rs)), "bad"),
        ("duplicate char", lambda rs: any("duplicate char" in e for e in validate(rs)), "bad"),
        ("empty char", lambda rs: any("char is empty" in e for e in validate(rs)), "bad"),
        ("unicode malformed", lambda rs: any("unicode malformed" in e for e in validate(rs)), "bad"),
        ("wu_xing enum", lambda rs: any("wu_xing out of enum" in e for e in validate(rs)), "bad"),
        ("pinyin element", lambda rs: any("pinyin contains" in e for e in validate(rs)), "bad"),
    ]

    def mutate(fn: Callable[[dict[str, Any]], None]) -> list[dict[str, Any]]:
        r = dict(good)
        fn(r)
        return [r]

    bad_inputs: list[list[dict[str, Any]]] = [
        mutate(lambda r: None),  # count == 1 != 30060
        mutate(lambda r: r.pop("char")),
        mutate(lambda r: r.__setitem__("radical_stroke", "four")),
        [dict(good), dict(good)],  # duplicate char (count also wrong)
        mutate(lambda r: r.__setitem__("char", "")),
        mutate(lambda r: r.__setitem__("unicode", "6728")),
        mutate(lambda r: r.__setitem__("wu_xing", "风")),
        mutate(lambda r: r.__setitem__("pinyin", ["mù", 3])),
    ]

    failures = 0
    for (name, detector, _), bad in zip(cases, bad_inputs):
        hit = detector(bad)
        print(f"  {'PASS' if hit else 'FAIL'}  must-hit: {name}")
        if not hit:
            failures += 1

    # must-miss: a single well-formed record should not trip field/value
    # detectors (the count detector is expected to fire on a 1-element list).
    clean = validate([dict(good)])
    field_errors = [e for e in clean if not e.startswith("count:")]
    ok = not field_errors
    print(f"  {'PASS' if ok else 'FAIL'}  must-miss: well-formed record")
    if not ok:
        failures += 1
        print("        " + "; ".join(field_errors))

    print(f"== self-test: {'PASS' if failures == 0 else str(failures) + ' failures'} ==")
    return 1 if failures else 0


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--self-test", action="store_true", help="run detector must-hit/must-miss cases")
    ap.add_argument("--file", type=Path, default=DATASET, help="dataset path override")
    args = ap.parse_args()

    if args.self_test:
        return self_test()

    records: list[dict[str, Any]]
    try:
        with args.file.open(encoding="utf-8") as fh:
            records = json.load(fh)
    except (OSError, json.JSONDecodeError) as exc:
        print(f"FAIL  load {args.file}: {exc}")
        return 1

    errors = validate(records)
    fields = set().union(*(r.keys() for r in records)) if records else set()
    print(f"dataset : {args.file}")
    print(f"records : {len(records)}")
    print(f"fields  : {len(fields)}")
    if errors:
        print(f"FAIL    : {len(errors)} violation(s)")
        for e in errors[:20]:
            print(f"  - {e}")
        if len(errors) > 20:
            print(f"  ... and {len(errors) - 20} more")
        return 1

    print("PASS    : dataset satisfies the contract")
    return 0


if __name__ == "__main__":
    sys.exit(main())
