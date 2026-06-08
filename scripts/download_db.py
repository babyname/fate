#!/usr/bin/env python3
import os
import sys
import json
import zipfile
import hashlib
import argparse
import urllib.request

def main():
    parser = argparse.ArgumentParser(description="Download and extract fate database")
    parser.add_argument("--version", help="Specific version to download (default: latest)")
    parser.add_argument("--local", help="Use local release file instead of GitHub", action="store_true")
    args = parser.parse_args()

    script_dir = os.path.dirname(__file__)
    root_dir = os.path.join(script_dir, "..")

    version_path = os.path.join(root_dir, "VERSION")
    if os.path.exists(version_path):
        with open(version_path, "r", encoding="utf-8") as f:
            default_version = f.read().strip()
    else:
        default_version = None

    if args.local:
        releases_dir = os.path.join(root_dir, "releases")
        if args.version:
            zip_name = f"fate-v{args.version}.zip"
            zip_path = os.path.join(releases_dir, zip_name)
        else:
            latest_path = os.path.join(releases_dir, "latest.json")
            if os.path.exists(latest_path):
                with open(latest_path, "r", encoding="utf-8") as f:
                    latest = json.load(f)
                zip_name = latest["file"]
                zip_path = os.path.join(releases_dir, zip_name)
            else:
                print(f"ERROR: latest.json not found at {latest_path}")
                sys.exit(1)
        
        if not os.path.exists(zip_path):
            print(f"ERROR: Release file not found at {zip_path}")
            sys.exit(1)
    else:
        if args.version:
            tag = f"v{args.version}"
            zip_name = f"fate-v{args.version}.zip"
        elif default_version:
            tag = f"v{default_version}"
            zip_name = f"fate-v{default_version}.zip"
        else:
            print("ERROR: Need VERSION file or --version")
            sys.exit(1)
        
        releases_dir = os.path.join(root_dir, "releases")
        os.makedirs(releases_dir, exist_ok=True)
        zip_path = os.path.join(releases_dir, zip_name)
        
        url = f"https://github.com/babyname/fate/releases/download/{tag}/{zip_name}"
        print(f"Downloading from {url}...")
        
        try:
            urllib.request.urlretrieve(url, zip_path)
        except Exception as e:
            print(f"ERROR: Failed to download: {e}")
            sys.exit(1)

    sha256_path = zip_path + ".sha256"
    sha256_expected = None
    if os.path.exists(sha256_path):
        with open(sha256_path, "r", encoding="utf-8") as f:
            line = f.read().strip()
            if line:
                sha256_expected = line.split()[0]

    print(f"Verifying {zip_name}...")
    sha256 = hashlib.sha256()
    with open(zip_path, "rb") as f:
        while True:
            chunk = f.read(8192)
            if not chunk:
                break
            sha256.update(chunk)
    actual = sha256.hexdigest()

    if sha256_expected and actual != sha256_expected:
        print(f"ERROR: SHA256 mismatch!")
        print(f"  Expected: {sha256_expected}")
        print(f"  Actual:   {actual}")
        sys.exit(3)
    print("  SHA256 OK")

    db_path = os.path.join(root_dir, "fate.db.gz")
    if os.path.exists(db_path):
        backup = db_path + ".backup"
        print(f"  Backing up current database to {backup}")
        os.replace(db_path, backup)

    print(f"Extracting {zip_name}...")
    with zipfile.ZipFile(zip_path, "r") as zf:
        zf.extract("fate.db.gz", root_dir)

    print(f"\nDone! Database extracted to {db_path}")

if __name__ == "__main__":
    main()
