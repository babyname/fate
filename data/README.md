# Data Architecture

## Principle

**Large external datasets are NOT tracked in git.** They are downloaded at build time or on first run.

## Directory Layout

```
data/
├── README.md              # This file
├── DATA_SOURCES.md         # Attribution: source, license, version
├── .gitignore              # Ignore all downloaded/generated data
├── download/               # Download scripts (idempotent)
│   ├── chinese_poetry.ps1  # Clone/pull chinese-poetry dataset
│   ├── unihan.ps1          # Download Unihan data
│   └── xinhua.ps1          # Download xinhua dictionary word.json
└── seed/                   # Curated seed data (tracked in git)
    └── .gitkeep
```

## Data Flow

```
External Source          Processing Tool       Project Artifact
───────────────────────────────────────────────────────────────
chinese-poetry (354MB) → dbinit import-poetry → SQLite DB (poems, poem_chars)
chinese-poetry (354MB) → dictctl import-poetry → poem_entries.json (runtime)
Unihan (49MB)          → dictctl import-unihan → resources/character.json
Xinhua (26MB)          → dictctl fill-xinhua   → resources/character.json
```

## Runtime Dependencies

The `data/` directory is NOT required for the core naming engine to run.
Only character.json (in `resources/`) is essential.

### Poetry Feature (Optional)

- Requires `data/chinese-poetry/` directory
- Server gracefully degrades if absent (poetry search returns empty)
- Run `./data/download/chinese_poetry.ps1` to download the dataset
