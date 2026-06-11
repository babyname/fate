# Fate CE/EE Feature Boundary Analysis

**Date**: 2026-06-11  
**Project**: babyname/fate v4

## Current Problems

### 1. Poetry has leaked into CE codebase
- `PoetryRater` (15% weight) in `internal/naming/raters.go`
- `buildPoetryIndex()` in `internal/http/handler.go`
- Poetry search API (`GET /api/poetry/search`)
- `PoetryOrigin` in `analysis/types.go`
- `HasPoetry` boolean on `CharInfo`
- All depend on loading 354MB chinese-poetry dataset at server startup

### 2. Meaning/Interpretation is ambiguous
- `Character.Meaning` (e.g., "美好") — CE: basic metadata
- `generateInterpretation()` — produces "名字 - 评分: 等级" currently useless
- `PoetryOrigin` enrichment — EE: poem source for character meanings

### 3. No clear package boundary between CE and EE
- Everything lives in `internal/` with no CE/EE separation
- HTTP handler mixes CE naming + EE poetry in one file
- No build tags, no feature flags

---

## Proposed CE/EE Boundary

### CE (Community Edition) — Core Naming Engine

Open source, MIT. Everything needed for: birthday → ranked name list.

```
CE Scoring Dimensions (total weight = 100%):
├── WuXing/BaZi   (35%)  — internal/chronosfate/, bazi/
├── WuGe ShuLi    (25%)  — internal/wuge/, wuxing/
├── BiHua         (25%)  — internal/naming/raters.go
├── ZhouYi        (10%)  — internal/zhouyi/
└── Poetry Flag    (5%)  — boolean only, pre-computed in DB
```

**CE packages:**
```
internal/bazi/          BaZi (8 characters) calculation
internal/chronosfate/   Fate data, XiYong (balance/geju)
internal/wuxing/        WuXing aggregation, SanCai
internal/wuge/          WuGe (sky/person/earth/outer/total)
internal/zhouyi/        ZhouYi hexagram divination
internal/dict/           Kangxi stroke/radical data tools
internal/naming/         Name filtering + rating + recommendation
internal/rating/         5-dimension comprehensive scoring
internal/filter/         Character filter options
internal/session/        Async session, excellent table
internal/analysis/       Result assembly + ranking
internal/repository/     DB access layer
internal/seeddb/         Seed data import pipeline
ent/                     Ent ORM schema (character, poem_char, poem, wu_xing)
config/                  YAML configuration
cmd/console/             CLI naming tool
cmd/server               HTTP API server (CE endpoints only)
cmd/dbinit               Database init tool
resources/               Embedded character.json + static web assets
```

**CE: Poetry — Boolean Flag Only**
- `Character.HasPoetry` (bool) — set during `dbinit import-seed`
  - Sources: character.json `has_poetry` field OR `poem_char` join table
- `PoetryRater` stays: checks `HasPoetry` boolean, adjusts score ±5 pts
  - **No external dataset required** — works from DB alone
- Weight reduced from 15% → 5% (poetry is a soft signal, not core)

**CE: Meaning — Basic Chinese Definition**
- `Character.Meaning` (string) — stays in CE
- Source: character.json `explanation` field
- Used by `rateWenHua()`: characters with meaning get +4 pts
- **No xinhua download required** — meaning comes from curated character.json

### EE (Enterprise Edition) — Premium Features

Licensed. Built on top of CE via plugins/middleware.

```
EE Core Features:
├── Advanced Poetry       Full poetry index, search, origin enrichment
├── Rich Interpretation   AI-generated name interpretations
├── Advanced Analysis     Detailed PDF/HTML reports, batch analysis
├── Database Management   Admin dashboard, import/export, migration
└── Enterprise            Auth, rate limiting, multi-tenancy
```

#### EE-1: Advanced Poetry Integration

**EE packages to create:**
```
ee/
├── poetry/
│   ├── index.go          Full poetry index (character → poem source)
│   ├── search.go         Poetry search by keyword/character
│   └── origin.go          PoetryOrigin enrichment for NameResult
├── interpretation/
│   ├── generator.go      AI-powered name interpretation
│   └── template.go       Interpretation templates
├── analysis/
│   ├── report.go         PDF/HTML report generation
│   └── batch.go          Batch name generation
├── admin/
│   ├── dashboard.go      Data quality dashboard API
│   ├── import.go         Advanced import pipeline
│   └── export.go         Data export tools
└── enterprise/
    ├── auth.go           API key management
    ├── ratelimit.go      Rate limiting
    └── analytics.go      Usage analytics
```

**EE HTTP endpoints (on top of CE):**
```
GET  /api/poetry/search          Search poetry by keyword
GET  /api/poetry/origin/{char}    Get poem source for character
GET  /api/name/{id}/interpretation  AI interpretation
POST /api/analysis/report        Generate PDF report
POST /api/analysis/batch         Batch name generation
GET  /api/admin/dashboard        Data quality dashboard
POST /api/admin/import           Advanced import
GET  /api/admin/export           Data export
```

**Enrichment pattern**: EE wraps CE results
```go
// CE returns NameResult with HasPoetry=true/false
// EE enriches with PoetryOrigin (title, author, dynasty, sentence)
func (ee *EEHandler) enrichResult(nr *NameResult) {
    if nr.Char1.HasPoetry {
        nr.Char1.PoetryOrigin = ee.poetryIndex[nr.Char1.Char]
    }
    if nr.Char2.HasPoetry {
        nr.Char2.PoetryOrigin = ee.poetryIndex[nr.Char2.Char]
    }
}
```

#### EE-2: Rich Interpretation

Current `generateInterpretation()` is empty: `"名字 - 评分: 等级"`

EE replacement:
- Character meaning explanation (from poetry source)
- WuXing compatibility with BaZi
- WuGe interpretation
- Cultural/literary allusions
- Overall evaluation paragraph

#### EE-3: Database Management

Current problems:
1. `cmd/migratedb` — V3→V4 one-shot migration, fragile
2. No incremental migration strategy
3. No data quality monitoring
4. No backup/restore
5. Seed data scattered in multiple JSON files

EE solution:
- Admin dashboard (`GET /api/admin/dashboard`) showing:
  - Character count, coverage stats
  - Missing fields (wu_xing, pinyin, stroke)
  - Source distribution
- Import pipeline with validation
- Export for backup
- Incremental schema migrations

#### EE-4: Enterprise Features
- API key management
- Rate limiting per key
- Usage analytics dashboard
- Multi-tenant support (separate DB per tenant)

---

## Implementation Plan

### Phase 1: Clean CE Boundary (Immediate)

1. **Reduce PoetryRater weight**: 15% → 5%
2. **Remove `buildPoetryIndex()` from CE handler** — make `poetryDir` optional
3. **Make `HasPoetry` a DB field** — pre-computed during `dbinit`
4. **Remove poetry search API from CE** (`GET /api/poetry/search`)
5. **Remove `PoetryOrigin` from CE types** — `analysis/types.go`
6. **Clean up `Meaning` semantics**: ensure it comes from character.json only

### Phase 2: Create EE Package Structure

1. Create `ee/` top-level directory with build tags
2. Extract poetry enrichment to `ee/poetry/`
3. Create EE HTTP handler that wraps CE handler
4. Add EE-specific endpoints

### Phase 3: Build EE Features

1. Full poetry search with origin tracking
2. AI-powered name interpretation
3. Report generation
4. Admin dashboard

---

## Scoring Weight Rebalance

### Current (problematic)
| Dimension | Weight | Problem |
|-----------|--------|---------|
| WuXing    | 35%   | OK |
| BiHua     | 25%   | OK |
| YinYun    | 25%   | OK |
| Poetry    | 15%   | Too high for optional feature |

### Proposed CE
| Dimension | Weight | Notes |
|-----------|--------|-------|
| WuXing    | 35%   | Core: BaZi compatibility |
| WuGe      | 25%   | Core: five-grid math |
| BiHua     | 20%   | Core: stroke balance |
| YinYun    | 15%   | Core: phonetics |
| Poetry    |  5%   | Soft: boolean flag only |

### Proposed EE adds
| Dimension | Weight | Notes |
|-----------|--------|-------|
| Poetry Depth | +5-10% | Full poem origin matching |
| Cultural Score | +5-10% | Literary allusion depth |

---

## Package Dependency Direction

```
EE (ee/) ──depends on──▶ CE (internal/)
                           ▲
                           │
                    config/, ent/, model/
```

- CE never imports EE
- EE wraps CE, adds enrichment
- Build tags: `//go:build ee` for EE features
```
