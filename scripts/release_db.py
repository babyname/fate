#!/usr/bin/env python3
import os
import sys
import json
import zipfile
import hashlib
import argparse
from datetime import datetime, timezone, timedelta

def main():
    parser = argparse.ArgumentParser(description="Package fate database for release")
    parser.add_argument("--version", help="Override version number")
    args = parser.parse_args()

    script_dir = os.path.dirname(__file__)
    root_dir = os.path.join(script_dir, "..")

    version_path = os.path.join(root_dir, "VERSION")
    if args.version:
        version = args.version
    elif os.path.exists(version_path):
        with open(version_path, "r", encoding="utf-8") as f:
            version = f.read().strip()
    else:
        print("ERROR: VERSION file not found and --version not specified")
        sys.exit(1)

    db_path = os.path.join(root_dir, "fate.db.gz")
    if not os.path.exists(db_path):
        print(f"ERROR: database not found at {db_path}")
        print("Please run 'go run ./cmd/dbinit' first to generate the database")
        sys.exit(1)

    releases_dir = os.path.join(root_dir, "releases")
    os.makedirs(releases_dir, exist_ok=True)

    zip_name = f"fate-v{version}.zip"
    zip_path = os.path.join(releases_dir, zip_name)

    print(f"Packaging {db_path} ({os.path.getsize(db_path):,} bytes)...")
    with zipfile.ZipFile(zip_path, "w", zipfile.ZIP_DEFLATED) as zf:
        zf.write(db_path, "fate.db.gz")

    zip_size = os.path.getsize(zip_path)
    print(f"  Zipped: {zip_name} ({zip_size:,} bytes)")

    sha256 = hashlib.sha256()
    with open(zip_path, "rb") as f:
        while True:
            chunk = f.read(8192)
            if not chunk:
                break
            sha256.update(chunk)
    sha256_hex = sha256.hexdigest()

    sha256_path = zip_path + ".sha256"
    with open(sha256_path, "w", encoding="utf-8") as f:
        f.write(f"{sha256_hex}  {zip_name}\n")
    print(f"  SHA256: {sha256_hex}")

    tz = timezone(timedelta(hours=8))
    metadata = {
        "version": version,
        "file": zip_name,
        "size": zip_size,
        "sha256": sha256_hex,
        "created": datetime.now(tz).isoformat()
    }

    latest_path = os.path.join(releases_dir, "latest.json")
    with open(latest_path, "w", encoding="utf-8") as f:
        json.dump(metadata, f, indent=2)

    print(f"\nRelease complete:")
    print(f"  File:   {zip_path}")
    print(f"  Size:   {zip_size:,} bytes ({zip_size / (1024*1024):.1f} MB)")
    print(f"  SHA256: {sha256_hex}")
    print(f"\nTo upload to GitHub Releases:")
    print(f"  gh release create v{version} {zip_path} {sha256_path} --title 'fate database v{version}'")

if __name__ == "__main__":
    main()
