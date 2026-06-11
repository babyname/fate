#!/usr/bin/env python3
"""Download chinese-poetry dataset.

Required for: server poetry index, dbinit import-poetry
Source: https://github.com/chinese-poetry/chinese-poetry (MIT)
Size: ~354 MB (shallow clone)
"""

import argparse
import subprocess
import sys
from pathlib import Path

REPO_URL = "https://github.com/chinese-poetry/chinese-poetry.git"
REPO_DIR = Path(__file__).resolve().parent.parent / "chinese-poetry"


def main() -> None:
    parser = argparse.ArgumentParser(description="Download chinese-poetry dataset")
    parser.add_argument("--target", default=str(REPO_DIR), help="Target directory")
    args = parser.parse_args()

    target = Path(args.target)

    if (target / ".git").is_dir():
        print("[*] Updating existing chinese-poetry checkout...")
        subprocess.run(["git", "-C", str(target), "pull", "--depth=1"], check=True)
    else:
        print("[*] Cloning chinese-poetry (shallow, ~354 MB)...")
        if target.exists():
            import shutil

            shutil.rmtree(target)
        subprocess.run(
            ["git", "clone", "--depth=1", "--filter=blob:none", REPO_URL, str(target)],
            check=True,
        )

    print(f"[✓] Dataset at: {target}")
    print("[>] Next: go run ./cmd/dbinit import-poetry")
    print("[>]   or: go run ./cmd/dictctl import-poetry")


if __name__ == "__main__":
    main()
