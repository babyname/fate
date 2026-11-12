package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/babyname/fate/v4"
	"github.com/babyname/fate/v4/internal/service/naming"
)

func TestExploreReturnsArrayWhenNoNamesRemain(t *testing.T) {
	entry := &taskEntry{}
	result := &fate.GenerateResult{ExcellentTable: *naming.NewExcellentTable()}
	result.ExcellentTable.TryPush(fate.ExcellentEntry{Char1: "知", Char2: "远", Score: 80})
	result.ExcellentTable.Finalize()
	entry.storeResult(result)
	handler := NewHandler(nil, nil, nil)
	handler.store.Set("example", entry)

	for _, want := range []int{1, 0} {
		request := httptest.NewRequest(http.MethodGet, "/api/v1/explore?task_id=example", nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("HTTP %d: %s", recorder.Code, recorder.Body.String())
		}
		var response struct {
			Names []fate.ExcellentEntry `json:"names"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if response.Names == nil || len(response.Names) != want {
			t.Fatalf("expected array of length %d, got %s", want, recorder.Body.String())
		}
	}
}
