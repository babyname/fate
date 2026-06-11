#!/usr/bin/env python3
"""Download Unihan database.

Required for: dictctl import-unihan
Source: https://www.unicode.org/Public/UCD/latest/ucd/Unihan.zip
License: Unicode License (free)
Size: ~2 MB compressed, ~49 MB extracted
"""

import argparse
import sys
import zipfile
from pathlib import Path
from tempfile import NamedTemporaryFile
from urllib.request import urlretrieve

UNIHAN_URL = "https://www.unicode.org/Public/UCD/latest/ucd/Unihan.zip"
UNIHAN_DIR = Path(__file__).resolve().parent.parent / "raw" / "unihan"


def main() -> None:
    parser = argparse.ArgumentParser(description="Download Unihan database")
    parser.add_argument("--target", default=str(UNIHAN_DIR), help="Target directory")
    args = parser.parse_args()

    target = Path(args.target)

    print(f"[*] Downloading Unihan.zip...")
    tmp = NamedTemporaryFile(suffix=".zip", delete=False)
    try:
        tmp.close()
        urlretrieve(UNIHAN_URL, tmp.name)
        print(f"[*] Extracting to {target}...")
        if target.exists():
            import shutil

            shutil.rmtree(target)
        target.mkdir(parents=True)
        with zipfile.ZipFile(tmp.name, "r") as zf:
            zf.extractall(target)
    finally:
        Path(tmp.name).unlink(missing_ok=True)

    print(f"[✓] Unihan data at: {target}")
    irg = target / "Unihan_IRGSources.txt"
    if irg.exists():
        print(f"[>] Next: go run ./cmd/dictctl import-unihan {irg}")


if __name__ == "__main__":
    main()
