# Phase 3B: Post-Session Cleanup

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

Phase 3B covers remaining areas:

| # | Area | Priority | Effort |
|---|------|----------|--------|
| 1 | Pipeline unit tests | ✅ DONE | `5d8355d` |
| 2 | Legacy tool cleanup | ✅ DONE | Local deletion |
| 3 | `constraint failed` performance test warnings | Low | Trivial |

---

## 1. Pipeline Unit Tests

**Status: ✅ Completed** (commit `5d8355d`)

Six integration-level tests covering smoke, pre-calculated FateData, multi-surname, stability, and context cancellation. See `internal/service/naming/pipeline_test.go`.

---

## 3. `constraint failed` Warnings

**Status: Low priority — known issue**

In performance tests, 7 `constraint failed` warnings appear during character linking. Root cause: 4,270+ traditional-simplified character connection conflicts where traditional and simplified forms of the same character have duplicate entries. These are logged but handled gracefully — no functional impact.

**Fix options:**
- Option A: Add `ON CONFLICT DO NOTHING` in link query
- Option B: Pre-deduplicate characters during dbinit
- Option C: Suppress log level to DEBUG


