#!/usr/bin/env python3
"""Download xinhua dictionary dataset.

Required for: dictctl fill-xinhua (fills character.json with Chinese meanings)
Source: https://github.com/pwxcoo/chinese-xinhua
License: Open data
Size: ~26 MB
"""

import argparse
import sys
from pathlib import Path
from urllib.request import urlretrieve

XINHUA_URL = (
    "https://raw.githubusercontent.com/pwxcoo/chinese-xinhua/master/data/word.json"
)
XINHUA_DIR = Path(__file__).resolve().parent.parent / "xinhua"


def main() -> None:
    parser = argparse.ArgumentParser(description="Download xinhua dictionary")
    parser.add_argument("--target", default=str(XINHUA_DIR), help="Target directory")
    args = parser.parse_args()

    target = Path(args.target)
    target.mkdir(parents=True, exist_ok=True)

    out = target / "word.json"
    print(f"[*] Downloading xinhua word.json (~26 MB)...")
    urlretrieve(XINHUA_URL, out)
    print(f"[✓] Xinhua data at: {out}")
    print(f"[>] Next: go run ./cmd/dictctl fill-xinhua resources/character.json {out}")


if __name__ == "__main__":
    main()
