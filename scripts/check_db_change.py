#!/usr/bin/env python3
"""
Check if database has changed compared to previous versions.
Outputs:
- "changed" if database is different from latest release
- "unchanged" if database is the same as latest release
- "unknown" if cannot determine (no previous release)
"""
import os
import sys
import json
import hashlib
import argparse


def calculate_sha256(file_path):
    """Calculate SHA256 hash of a file."""
    sha256 = hashlib.sha256()
    with open(file_path, "rb") as f:
        while True:
            chunk = f.read(8192)
            if not chunk:
                break
            sha256.update(chunk)
    return sha256.hexdigest()


def main():
    parser = argparse.ArgumentParser(description="Check if database has changed")
    parser.add_argument("--db-path", default="fate.db.gz", help="Path to database file")
    parser.add_argument("--releases-dir", default="releases", help="Releases directory")
    args = parser.parse_args()

    if not os.path.exists(args.db_path):
        print("ERROR: Database file not found", file=sys.stderr)
        sys.exit(1)

    current_sha = calculate_sha256(args.db_path)
    print(f"Current database SHA256: {current_sha}", file=sys.stderr)

    latest_json = os.path.join(args.releases_dir, "latest.json")
    if os.path.exists(latest_json):
        with open(latest_json, "r", encoding="utf-8") as f:
            meta = json.load(f)
        if "db_sha256" in meta:
            if meta["db_sha256"] == current_sha:
                print(f"unchanged: {meta['version']}")
                return

    # Check GitHub releases using GitHub CLI if available
    # (this part will be handled by the CI workflow)

    print("unknown")


if __name__ == "__main__":
    main()
