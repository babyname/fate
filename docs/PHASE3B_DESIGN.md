# Phase 3B: Post-Session Cleanup & Qiming Alignment

**Author:** QClaw  
**Date:** 2026-06-11  
**Status:** DRAFT — design document for remaining Phase 3 work

## Phase 3A Recap

Phase 3A completed successfully (commit `622248d`):

- `internal/service/naming/` package created with 4 files: `pipeline.go` (172 lines), `scorer.go` (75 lines), `preloader.go` (40 lines), `excellent.go` (189 lines)
- `internal/session/session.go` reduced from 340 → 129 lines (−62%)
- `internal/session/excellent.go` reduced from 190 → 8 lines (re-export wrapper)
- Public API (`fate.Session`, `fate.Input`, `fate.Output`, `fate.ExcellentEntry`) **unchanged**
- Full test suite: 119s all green, output identical (张适抡 90.1)
- Design doc: `docs/PHASE3A_DESIGN.md`

## Phase 3B Scope

Phase 3B covers three remaining areas:

| # | Area | Priority | Effort |
|---|------|----------|--------|
| 1 | Pipeline unit tests | ✅ DONE | `5d8355d` |
| 2 | EE legacy tool cleanup | ✅ DONE | Local deletion |
| 3 | `constraint failed` performance test warnings | Low | Trivial |
| 4 | Qiming handler deduplication | High | Large |

---

## 1. Pipeline Unit Tests

**Status: ✅ Completed** (commit `5d8355d`)

Six integration-level tests covering smoke, pre-calculated FateData, multi-surname, stability, and context cancellation. See `internal/service/naming/pipeline_test.go`.

---

## 2. EE Legacy Tool Cleanup

**Status: ✅ Completed**

Three tools removed from `cmd/tools/`:

| Tool | Purpose | Why Removed |
|------|---------|-------------|
| `dbaudit/` | Database audit / data quality check | EE-only ops tool, hardcoded DSN |
| `inspectdb/` | Database inspector | EE-only, dead code |
| `inspectv4/` | v4 data migration inspection | EE-only, migration-era tool |

All three were untracked in git (never committed to the main branch). `go build ./...` confirmed no broken imports.

---

## 3. `constraint failed` Warnings

**Status: Low priority — known issue**

In performance tests, 7 `constraint failed` warnings appear during character linking. Root cause: 4,270+ traditional-simplified character connection conflicts where traditional and simplified forms of the same character have duplicate entries. These are logged but handled gracefully — no functional impact.

**Fix options:**
- Option A: Add `ON CONFLICT DO NOTHING` in link query
- Option B: Pre-deduplicate characters during dbinit
- Option C: Suppress log level to DEBUG

---

## 4. Qiming Handler Deduplication (Main Work)

### 4.1 Current State

**Fate handler** (`internal/http/handler.go`):
- ~400 lines
- Uses `fate.Session` → `fate.NewSessionWithFilter()` → `session.Start()` → poll state
- Web handler with SSE-style status polling

**Qiming handler** (`internal/api/echo_handlers_core.go`):
- ~500 lines
- Uses **same** `fate.Session` API via `fate.NewSessionWithFilter()`
- Echo-based HTTP handler with status polling (nearly identical to Fate)
- Adds Qiming-specific: membership check, pricing, favorites, poem retrieval

```go
// Fate handler                         // Qiming handler
sess := fate.NewSessionWithFilter(flt)  // same
input := &fate.Input{...}               // same pattern
sess.Start(ctx, input)                  // same
// poll State() in loop                 // same pattern
// wait for Finish/Failed               // same
output := sess.Output()                 // same
```

### 4.2 What's Actually Duplicated

| Component | Fate | Qiming | Overlap |
|-----------|------|--------|---------|
| Session lifecycle (create → start → poll → output) | handler.go | echo_handlers_core.go | 🔴 90% same |
| Filter option construction | handler.go | echo_handlers_core.go | 🟡 80% same |
| Name result JSON serialization | handler.go | echo_handlers_core.go | 🟡 70% same |
| Task state polling loop | inline | inline | 🔴 identical |
| Task store (in-flight sessions) | inline map | inline map + TTL | 🟡 different impl |
| Bazi/WuXing detail structs | `fate.BaziBasic` struct | inline construction | 🟡 70% same |

### 4.3 Core Problem

Qiming imports Fate as a **Go module** (`go.mod replace`):

```
// qiming/go.mod
replace github.com/babyname/fate/v4 => ../fate

// qiming/internal/api/echo_handlers_core.go
import "github.com/babyname/fate/v4"
```

This means Qiming can only access Fate's **public API** (`fate.Session`, `fate.Input`, `fate.FilterOption`, `fate.ExcellentEntry`, etc.). It cannot access `internal/` packages.

The `fate.Session` public API is an **asynchronous state machine** — a design that forces Qiming to replicate the polling loop. The handler duplication exists because:

1. **State machine abstraction leaks** — the caller must poll `State()` instead of getting a synchronous result
2. **No shared handler skeleton** — each project rewrites the same create-start-poll-output pattern
3. **Different web frameworks** — Fate uses standard `http.Handler`, Qiming uses Echo

### 4.4 Solution Options

#### Option A: Shared Public Handler SDK (Recommended)

Create a new public package in Fate that Qiming can import directly:

```go
// fate/namesvc.go  — NEW public package
package namesvc

// Generate is a synchronous convenience wrapper.
// Returns the same result as Session but without polling.
func Generate(ctx context.Context, req GenerateRequest) (*GenerateResult, error)

// GenerateResult mirrors the final session output.
type GenerateResult struct {
    TopNames       []fate.NameResult
    ExcellentTable *fate.ExcellentTable
    CharMap        map[string]*fate.CharInfo
    FateData       *fate.FateData
}
```

**Pros:**
- Qiming calls one function instead of 5-step async dance
- Eliminates handler duplication completely
- Synchronous API is easier to test
- No framework coupling

**Cons:**
- Minor API surface expansion (new public package)
- Must maintain backward compat with existing Session API

#### Option B: Handler Template Extraction

Extract the shared handler logic into a template/skeleton that both projects embed:

```go
// fate/handlerutil/  — shared handler utilities
func HandleNameGeneration(w http.ResponseWriter, r *http.Request, opts ...HandlerOption) error
```

**Pros:**
- No public API surface change
- Works with any HTTP framework via adapters

**Cons:**
- Less clean than Option A
- Still leaks HTTP concerns into the library

#### Option C: Status Quo (Not Recommended)

Leave both handlers as-is. Accept maintenance burden.

### 4.5 Implementation Plan (Recommended: Option A)

```
Phase 3B.4a — Create fate/namesvc.go public synchronous API
  ├── Define GenerateRequest / GenerateResult structs
  ├── Implement Generate(ctx, req) using internal Pipeline
  ├── Unit tests for synchronous Generate
  ├── Integration tests for Generate
  └── Verify Qiming can import without changes to go.mod

Phase 3B.4b — Refactor Qiming handler to use namesvc.Generate
  ├── Replace session create → start → poll → output with single call
  ├── Remove duplicate TaskStore (if namesvc.Generate is synchronous)
  ├── Keep Qiming-specific: membership check, pricing, favorites
  └── Expected reduction: ~300 lines removed from echo_handlers_core.go

Phase 3B.4c — (Optional) Refactor Fate handler to use same approach
  ├── Replace inline polling with namesvc.Generate
  ├── Keep handler-specific concerns (auth, logging)
  └── Expected reduction: ~200 lines removed from handler.go
```

### 4.6 Detailed Design: `fate/namesvc.go`

```go
package fate

import (
    "context"
    "time"
)

// GenerateRequest mirrors the input to the name generation pipeline.
// It is a synchronous alternative to the Session state machine.
type GenerateRequest struct {
    LastName   string     // 姓, e.g. "张"
    FirstName  string     // optional 名 hint
    Born       time.Time
    Sex        Sex
    Filter     *FilterOption
}

// GenerateResult is the synchronous output of name generation.
type GenerateResult struct {
    TopNames       []NameResult
    ExcellentTable ExcellentTable
    CharMap        map[string]CharInfo
    FateData       *FateData
}

// Generate runs the full name generation pipeline synchronously.
// It is a convenience wrapper around the internal Pipeline.
func Generate(ctx context.Context, req GenerateRequest) (*GenerateResult, error) {
    // 1. Lookup character for LastName
    // 2. Create filter from option
    // 3. Call internal Pipeline
    // 4. Convert to public types
    // 5. Return result
}
```

### 4.7 Qiming Handler After Refactor

```go
// echo_handlers_core.go — REFACTORED
func (s *Server) handleNameGenerate(c echo.Context) error {
    // Qiming-specific: membership check
    if !s.checkMembership(c) {
        return c.JSON(http.StatusPaymentRequired, ...)
    }

    req := fate.GenerateRequest{...}
    result, err := fate.Generate(c.Request().Context(), req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ...)
    }

    return c.JSON(http.StatusOK, result)
}
```

Reduction: from ~200 lines (create→start→poll→output) to ~30 lines (construct → call → return).

---

## Dependency Graph (Post-Phase 3B)

```
fate/namesvc.go          ← NEW: public synchronous API
  └── internal/service/naming/pipeline.go
       ├── repository/
       ├── rating/
       ├── wuge/
       ├── filter/
       └── chronosfate/

fate/session.go          ← UNCHANGED: backward compat (async state machine)
  └── internal/service/naming/pipeline.go

qiming/echo_handlers_core.go  ← SIMPLIFIED: uses fate.Generate()
  └── fate/namesvc.go

fate/internal/http/handler.go  ← OPTIONALLY simplified
  └── fate/namesvc.go or fate.Session (keeping all routes)
```

---

## Non-Goals (Out of Scope for Phase 3B)

- Refactoring Qiming's membership / pricing / favorites handlers
- Changing `fate.Session` public API (backward compat preserved)
- Merging Fate and Qiming into a monorepo
- UI changes

---

## Risks & Mitigations

| Risk | Severity | Mitigation |
|------|----------|------------|
| `namesvc.Generate` breaks existing Session-based code | Low — backwards compat guaranteed | Old API still works, new API is additive |
| Performance: synchronous call blocks while generating | Medium — same as current impl | Current Session also blocks; namesvc.Generate just removes the polling indirection |
| Qiming go.mod replace path breaks after refactor | Medium | Test with `go build` in Qiming after each change |
| Pipeline internal refactoring breaks public API | Low | Test both `fate.Generate` and `fate.Session` paths |

---

## Verification Plan

1. `go build ./...` in Fate — must pass
2. `go test ./...` in Fate — all green (119s)
3. `go build ./...` in Qiming — must pass (verifies module compatibility)
4. Manual: run Qiming server, trigger name generation, verify output matches current behavior
