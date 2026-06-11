# Phase 3A: Fate Service Layer Extraction

**Author:** QClaw  
**Date:** 2026-06-11  
**Status:** DESIGN — ready for implementation

## Goal

Extract business logic from `internal/session/session.go` (the God Object) into a clean `internal/service/naming/` package, while maintaining full backward compatibility with the public `fate.Session` API used by Qiming.

## Current Architecture (Problem)

```
fate.go (public API)
  └── session/session.go  ← 400+ lines, 8 imports
       ├── State machine (Start/Stop/Wait/Err)
       ├── Name generation pipeline (generate)
       │   ├── Fate calculation (chronosfate)
       │   ├── Lucky wuge combo lookup
       │   ├── Character preloading (DB)
       │   ├── Poetry set query (DB)
       │   ├── 3-tier filter fallback
       │   └── Output assembly (charMap + top results)
       ├── Scoring loop (scoreAllCandidates)  ← 60 lines, 4 nested loops
       ├── Character preloading (preloadChars)
       ├── Poetry query (queryPoetrySet)
       └── Stroke helpers (getLastStrokeFromBasic, strokeGetFromFilterType)
```

**Why it's bad:**
- Untestable in isolation: session ties DB access to scoring logic
- Scoring loop locked inside session: can't reuse for different modes
- Mixes state management with domain logic
- Filter fallback logic scattered across generate() instead of in scoring

## Target Architecture

```
internal/service/naming/        ← NEW package
  ├── pipeline.go               — Generate(ctx, req) → result
  ├── scorer.go                 — runScoring(candidates, rater, table)
  └── preloader.go              — preloadChars(db, filter, lucky)
       ↓ imports
       repository/ (DB), wuge/ (combos), rating/ (Rater),
       session/ (ExcellentTable — data structure kept here)

session/session.go             ← SIMPLIFIED
  ├── State machine (unchanged)
  ├── generate() — 15 lines: calls service.Generate(), builds Output
  └── Removed: scoreAllCandidates, preloadChars, queryPoetrySet,
      getLastStrokeFromBasic, strokeGetFromFilterType

fate.go (public API)           ← UNCHANGED
  Session, Input, Output, ExcellentTable — all identical
```

## Dependency Graph

```
service/naming
  ├── repository   (DB queries)
  ├── session      (ExcellentTable — data structure)
  ├── rating       (Rater)
  ├── wuge         (lucky combos)
  ├── chronosfate  (fate calc)
  ├── filter       (Filter interface)
  ├── naming       (NameBasic)
  ├── ent          (Character)
  └── log          (logging)

session            ← NO LONGER imports
  ├── repository   (only for Session.NewSession factory — unchanged)
  ├── service/naming ← NEW import
  ├── naming       (NameBasic)
  ├── filter       (Filter)
  ├── chronosfate  (FateData) — still needed for fate calc, OR delegated
  └── log
```

## Detailed Design

### 1. `internal/service/naming/types.go`

Internal DTOs used by the pipeline:

```go
package naming

// PipelineRequest carries all inputs for name generation.
type PipelineRequest struct {
    LastName  [2]*ent.Character
    Born      time.Time
    Sex       naming.Sex
    Filter    filter.Filter
}

// PipelineResult is the output of a generation run.
type PipelineResult struct {
    TopNames       []analysis.NameResult
    ExcellentTable *session.ExcellentTable
    CharMap        map[string]*ent.Character
    FateData       *chronosfate.FateData
}
```

### 2. `internal/service/naming/pipeline.go`

The main orchestrator. Replaces `session.generate()`.

```go
type Pipeline struct {
    db *repository.Repository
}

func NewPipeline(db *repository.Repository) *Pipeline

// Generate runs the full name generation pipeline.
// It is the functional equivalent of session.generate().
func (p *Pipeline) Generate(ctx context.Context, req PipelineRequest) (*PipelineResult, error)
```

**Pipeline steps (identical to current generate() logic):**

1. Calculate strokes from last name via filter
2. Get lucky wuge combinations: `wuge.GetLuckyByLastName(s1, s2)`
3. Preload characters: `p.preload(ctx, lucky, req.Filter)`
4. Query poetry set: `p.db.QueryPoetryChars(ctx)`
5. Create rater: `rating.NewRaterWithStrokes(fateData, s1, s2)`
6. Score candidates via 3-tier fallback:
   ```
   runScoring(lucky, rater, table, poetrySet, chars, filter)
   if heap empty && strictness != "relaxed" → moderate → re-score
   if heap empty && strictness != "relaxed" → relaxed → re-score
   ```
7. Finalize table
8. Build charMap
9. Build top results from heap entries

### 3. `internal/service/naming/scorer.go`

Extracted from `session.scoreAllCandidates()`.

```go
// runScoring scores all candidate name pairs against the given lucky combos,
// pushing qualifying entries into the excellent table.
func runScoring(
    lucky []wuge.WuGeResult,
    rater *rating.Rater,
    table *session.ExcellentTable,
    poetrySet map[string]bool,
    chars map[int][]*ent.Character,
    flt filter.Filter,
)
```

Key behaviors preserved:
- 4 nested loops: lucky → stroke filter → char1 → char2
- Context cancellation check on inner loop
- Poetry mode 2 filtering (both chars must have poetry)
- `table.TryPush()` for each candidate pair
- Logging `heap_size` after completion

### 4. `internal/service/naming/preloader.go`

Extracted from `session.preloadChars()`.

```go
// preloadChars loads characters for all stroke values needed by the lucky
// wuge combinations, respecting the filter's stroke scope and regular checks.
func preloadChars(
    ctx context.Context,
    db *repository.Repository,
    lucky []wuge.WuGeResult,
    flt filter.Filter,
) (map[int][]*ent.Character, error)
```

Key behaviors preserved:
- Collect unique stroke values from lucky combos
- Skip strokes that fail `CheckSkipStrokeNumberScope`
- Cache check: only query DB for missing stroke values
- Uses `db.GetCharactersCached()` for query

### 5. `session/session.go` changes

**Removed functions:**
- `scoreAllCandidates()` → moved to `service/naming`
- `preloadChars()` → moved to `service/naming`
- `queryPoetrySet()` → inlined into pipeline
- `getLastStrokeFromBasic()` → moved to `service/naming`
- `strokeGetFromFilterType()` → moved to `service/naming`
- `buildCharMap()` → moved to `service/naming`

**New `session` struct field:**
```go
type session struct {
    // ... existing fields ...
    pipeline *naming.Pipeline  // NEW
}
```

**New `generate()` (approximately 15 lines):**
```go
func (s *session) generate() error {
    defer s.close()

    result, err := s.pipeline.Generate(s.ctx, naming.PipelineRequest{
        LastName: s.output.Basic().LastName,
        Born:     /* from input */,
        Sex:      /* from input */,
        Filter:   s.filter,
    })

    s.output.SetTopNames(result.TopNames)
    s.output.SetExcellentTable(result.ExcellentTable)
    s.output.SetCharMap(result.CharMap)
    s.output.SetFateData(result.FateData)

    s.SetState(SessionStateFinish)
    return nil
}
```

### 6. `fate.go` changes

**None.** The public API surface is unchanged. Qiming continues to import `fate.Session`, `fate.Input`, etc.

## Files Changed

| File | Change | Type |
|------|--------|------|
| `internal/service/naming/pipeline.go` | New | Add |
| `internal/service/naming/scorer.go` | New | Add |
| `internal/service/naming/preloader.go` | New | Add |
| `internal/session/session.go` | ~300 lines removed, ~15 lines new | Simplify |
| `docs/PHASE3A_DESIGN.md` | New | Doc |

## Verification

1. `go build ./...` — must pass
2. `go test -v -run TestNameAnalysisOutput -timeout 120s` — must produce 10000 names, score >80
3. `go test -v -run TestNameGenerationPerformance -timeout 120s` — 4 surnames, all >0 names
4. `go test ./... -timeout 180s` — all tests green

## Risks & Mitigations

| Risk | Severity | Mitigation |
|------|----------|------------|
| Circular imports (service/naming → session for ExcellentTable) | Medium | `service/naming` imports `session` for `ExcellentTable` only. `session` imports `service/naming` for `Pipeline`. Circular? NO — `session` → `service/naming` → `session` is a circle. **Fix:** Move `ExcellentTable` to a standalone data package, OR keep it in session and pass it as parameter. |
| Scoring behavior change | Low | `runScoring` is a pure extraction — same loops, same filter checks |
| Performance regression | Low | Same number of DB queries, same heap operations |

### Circular Import Fix

The pipeline creates an `ExcellentTable` and returns it. This means `service/naming` must import `session` for `ExcellentTable`. But `session` imports `service/naming` for `Pipeline`. This is circular.

**Solution: Move `ExcellentTable` to `internal/service/naming/`.**

Since `ExcellentTable` is already re-exported via `fate.ExcellentTable` (type alias), moving it doesn't break the public API. The internal `session` package and `http/handler.go` import paths change, but those are internal.

**Revised dependency graph after fix:**
```
service/naming/           ← ExcellentTable lives here now
  ├── repository
  ├── rating
  ├── wuge
  ├── chronosfate
  ├── filter
  ├── ent
  └── log

session/                  ← IMPORTS service/naming for Pipeline
  ├── service/naming
  ├── repository
  ├── filter
  └── naming

service/naming NO LONGER imports session — circularity resolved.
```

`session/excellent.go` becomes a thin re-export wrapper OR is deleted and all internal callers in `http/handler.go` + tests update their imports.

## Implementation Order

1. Create `internal/service/naming/` package
2. Move `ExcellentTable` + `ExcellentEntry` → `service/naming/excellent.go`
3. Implement `scorer.go` (pure extraction)
4. Implement `preloader.go` (pure extraction)
5. Implement `pipeline.go` (orchestrator, uses steps 2-4)
6. Rewrite `session/generate()` to use pipeline
7. Update `session/excellent.go` → re-export from `service/naming`
8. Update `http/handler.go` imports (session.ExcellentTable → naming.ExcellentTable)
9. Update `fate.go` type aliases
10. Delete dead code in `session/session.go`
11. Build + full test suite

## Post-Phase 3A Impact on Phase 3B (Qiming)

With this extraction:
- Qiming's `internal/api/echo_handlers_core.go` **still works unchanged** (imports `fate.Session`, not `session.ExcellentTable` directly)
- Phase 3B can now extract Qiming's handler logic into a shared handler pattern that consumes `service/naming.Pipeline` directly
- The `service/naming.Pipeline` becomes the canonical naming engine — Fate and Qiming both use it via different paths (Fate: through Session adapter, Qiming: potentially directly)
