#!/usr/bin/env python3
import os
import sys
import json
import zipfile
import hashlib
import argparse
from datetime import datetime, timezone, timedelta

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

def load_latest_metadata(releases_dir):
    """Load the latest release metadata if it exists."""
    latest_path = os.path.join(releases_dir, "latest.json")
    if os.path.exists(latest_path):
        with open(latest_path, "r", encoding="utf-8") as f:
            return json.load(f)
    return None

def load_db_version(root_dir):
    """Load DB_VERSION file if it exists, fallback to VERSION."""
    db_version_path = os.path.join(root_dir, "DB_VERSION")
    if os.path.exists(db_version_path):
        with open(db_version_path, "r", encoding="utf-8") as f:
            return f.read().strip()
    
    version_path = os.path.join(root_dir, "VERSION")
    if os.path.exists(version_path):
        with open(version_path, "r", encoding="utf-8") as f:
            return f.read().strip()
    
    return None

def main():
    parser = argparse.ArgumentParser(description="Package fate database for release")
    parser.add_argument("--version", help="Override version number")
    parser.add_argument("--force", action="store_true", help="Force release even if database unchanged")
    args = parser.parse_args()

    script_dir = os.path.dirname(__file__)
    root_dir = os.path.join(script_dir, "..")

    if args.version:
        version = args.version
    else:
        version = load_db_version(root_dir)
        if not version:
            print("ERROR: DB_VERSION or VERSION file not found and --version not specified")
            sys.exit(1)

    db_path = os.path.join(root_dir, "fate.db.gz")
    if not os.path.exists(db_path):
        print(f"ERROR: database not found at {db_path}")
        print("Please run 'go run ./cmd/dbinit' first to generate the database")
        sys.exit(1)

    releases_dir = os.path.join(root_dir, "releases")
    os.makedirs(releases_dir, exist_ok=True)

    current_db_sha256 = calculate_sha256(db_path)
    latest_metadata = load_latest_metadata(releases_dir)
    
    if latest_metadata and not args.force:
        print(f"Checking if database has changed...")
        
        # Check if the latest release has the same SHA256 hash (database not changed)
        if "sha256" in latest_metadata:
            latest_sha256 = latest_metadata["sha256"]
            # We need to check if it's the same database content
            # To do that, we need to look at the actual zip file's SHA256
            latest_zip_path = os.path.join(releases_dir, latest_metadata["file"])
            if os.path.exists(latest_zip_path):
                latest_zip_sha256 = calculate_sha256(latest_zip_path)
                if latest_zip_sha256 == current_db_sha256:
                    print(f"Database unchanged from v{latest_metadata['version']}")
                    print(f"Latest release is still v{latest_metadata['version']}")
                    print("\nUse --force to force a new release anyway.")
                    return

    zip_name = f"fate-v{version}.zip"
    zip_path = os.path.join(releases_dir, zip_name)

    print(f"Packaging {db_path} ({os.path.getsize(db_path):,} bytes)...")
    with zipfile.ZipFile(zip_path, "w", zipfile.ZIP_DEFLATED) as zf:
        zf.write(db_path, "fate.db.gz")

    zip_size = os.path.getsize(zip_path)
    print(f"  Zipped: {zip_name} ({zip_size:,} bytes)")

    zip_sha256 = calculate_sha256(zip_path)

    sha256_path = zip_path + ".sha256"
    with open(sha256_path, "w", encoding="utf-8") as f:
        f.write(f"{zip_sha256}  {zip_name}\n")
    print(f"  SHA256: {zip_sha256}")

    tz = timezone(timedelta(hours=8))
    metadata = {
        "version": version,
        "file": zip_name,
        "size": zip_size,
        "sha256": zip_sha256,
        "db_sha256": current_db_sha256,
        "created": datetime.now(tz).isoformat()
    }

    latest_path = os.path.join(releases_dir, "latest.json")
    with open(latest_path, "w", encoding="utf-8") as f:
        json.dump(metadata, f, indent=2)

    print(f"\nRelease complete:")
    print(f"  File:   {zip_path}")
    print(f"  Size:   {zip_size:,} bytes ({zip_size / (1024*1024):.1f} MB)")
    print(f"  SHA256: {zip_sha256}")
    print(f"\nTo upload to GitHub Releases:")
    print(f"  gh release create v{version} {zip_path} {sha256_path} --title 'fate database v{version}'")

if __name__ == "__main__":
    main()
