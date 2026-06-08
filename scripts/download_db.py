#!/usr/bin/env python3
import os
import sys
import json
import zipfile
import hashlib
import argparse
import urllib.request

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

def download_file(url, output_path):
    """Download a file with progress."""
    print(f"Downloading from {url}...")
    try:
        urllib.request.urlretrieve(url, output_path)
        return True
    except Exception as e:
        print(f"ERROR: Failed to download: {e}")
        return False

def main():
    parser = argparse.ArgumentParser(description="Download and extract fate database")
    parser.add_argument("--version", help="Specific version to download (default: latest from DB_VERSION)")
    parser.add_argument("--local", help="Use local release file instead of GitHub", action="store_true")
    args = parser.parse_args()

    script_dir = os.path.dirname(__file__)
    root_dir = os.path.join(script_dir, "..")

    default_version = load_db_version(root_dir)

    zip_name = None
    zip_path = None
    releases_dir = os.path.join(root_dir, "releases")

    if args.local:
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
                print(f"Using local latest release: v{latest['version']}")
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
            print("ERROR: Need DB_VERSION/VERSION file or --version")
            sys.exit(1)
        
        os.makedirs(releases_dir, exist_ok=True)
        zip_path = os.path.join(releases_dir, zip_name)
        
        url = f"https://github.com/babyname/fate/releases/download/{tag}/{zip_name}"
        
        if not download_file(url, zip_path):
            print(f"\nFailed to download v{default_version if not args.version else args.version}")
            sys.exit(1)

    # Verify SHA256
    sha256_path = zip_path + ".sha256"
    sha256_expected = None
    
    # Try to download SHA256 file if not local and not exists
    if not args.local and not os.path.exists(sha256_path):
        sha256_url = f"https://github.com/babyname/fate/releases/download/{tag}/{zip_name}.sha256"
        print(f"Downloading SHA256 checksum...")
        try:
            urllib.request.urlretrieve(sha256_url, sha256_path)
        except Exception:
            print(f"  Note: Could not download SHA256 checksum file, verification will be skipped.")
    
    if os.path.exists(sha256_path):
        with open(sha256_path, "r", encoding="utf-8") as f:
            line = f.read().strip()
            if line:
                sha256_expected = line.split()[0]

    print(f"Verifying {zip_name}...")
    actual_sha256 = calculate_sha256(zip_path)

    if sha256_expected and actual_sha256 != sha256_expected:
        print(f"ERROR: SHA256 mismatch!")
        print(f"  Expected: {sha256_expected}")
        print(f"  Actual:   {actual_sha256}")
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
