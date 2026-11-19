package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/babyname/fate/v4"
	"github.com/babyname/fate/v4/internal/analysis"
	"github.com/babyname/fate/v4/internal/service/naming"
)

func TestRecommendationPagesAndDetailsUseFullCatalogue(t *testing.T) {
	h := NewHandler(nil, nil, nil)
	result := newPopulatedResult(22)
	names := make([]fate.NameResult, 22)
	for i := range names {
		names[i] = fate.NameResult{
			Rank: i + 1, Char1: analysis.CharInfo{Char: fmt.Sprintf("字%03d", i)}, Char2: analysis.CharInfo{Char: "远"},
			Ranking: &analysis.RecommendationRanking{Traditional: 70, SoundAndMeaning: 90, Total: 82},
		}
	}
	result.TopNames = names[:10]
	h.store.Set("recommendation", &taskEntry{result: result, candidates: names})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/names?task_id=recommendation&page=2", nil))
	var page struct {
		Page  int
		Total int
		Names []fate.NameResult
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK || page.Page != 2 || page.Total != 22 || len(page.Names) != 2 || page.Names[0].Rank != 21 {
		t.Fatalf("wrong catalogue page: HTTP %d: %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	query := url.Values{"task_id": {"recommendation"}, "char1": {names[21].Char1.Char}, "char2": {"远"}}
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/name_detail?"+query.Encode(), nil))
	var detail struct {
		Name fate.NameResult `json:"name_result"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &detail); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK || detail.Name.Rank != 22 || detail.Name.Ranking == nil || detail.Name.Ranking.Total != 82 {
		t.Fatalf("wrong candidate detail: HTTP %d: %s", rec.Code, rec.Body.String())
	}
}

// newPopulatedResult builds a finalized table with n distinct candidates.
func newPopulatedResult(n int) *fate.GenerateResult {
	result := &fate.GenerateResult{ExcellentTable: naming.NewExcellentTable()}
	for i := 0; i < n; i++ {
		result.ExcellentTable.TryPush(fate.ExcellentEntry{
			Char1: fmt.Sprintf("字%03d", i),
			Char2: "远",
			Score: float64(n - i),
		})
	}
	result.ExcellentTable.Finalize()
	return result
}

func newExploreHandler(t *testing.T, result *fate.GenerateResult) *Handler {
	t.Helper()
	entry := &taskEntry{}
	entry.storeResult(result)
	h := NewHandler(nil, nil, nil)
	h.store.Set("example", entry)
	return h
}

// Regression: once maxShownNames entries are marked, ExcellentTable.Explore
// returns a nil slice, which json.Marshal encodes as JSON null. The handler
// must normalize it to an empty array so clients can always iterate `names`.
func TestExploreReturnsArrayWhenPoolExhausted(t *testing.T) {
	result := newPopulatedResult(101)
	for _, e := range result.ExcellentTable.TopN(100) {
		result.ExcellentTable.MarkShown(e.Char1, e.Char2)
	}
	h := newExploreHandler(t, result)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/explore?task_id=example&count=3", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if got := string(body["names"]); got != "[]" {
		t.Fatalf("exhausted explore must return a JSON array, got %s (body: %s)", got, rec.Body.String())
	}
	if got := string(body["shown_count"]); got != "100" {
		t.Fatalf("shown_count = %s, want 100 (body: %s)", got, rec.Body.String())
	}
}

// Explore must advance the shown cursor instead of replaying the same names.
func TestExploreAdvancesShownCursor(t *testing.T) {
	h := newExploreHandler(t, newPopulatedResult(1))

	for _, want := range []int{1, 0} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/explore?task_id=example", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
		}
		var body struct {
			Names []fate.ExcellentEntry `json:"names"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body.Names == nil || len(body.Names) != want {
			t.Fatalf("expected array of length %d, got %s", want, rec.Body.String())
		}
	}
}
