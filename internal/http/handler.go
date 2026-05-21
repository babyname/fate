package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/babyname/fate"
)

type Handler struct {
	fate fate.Fate
	mux  *http.ServeMux
}

func NewHandler(f fate.Fate) *Handler {
	h := &Handler{
		fate: f,
		mux:  http.NewServeMux(),
	}
	h.mux.HandleFunc("GET /health", h.handleHealth)
	h.mux.HandleFunc("POST /api/generate", h.handleGenerate)
	h.mux.HandleFunc("GET /api/names", h.handleNames)
	h.mux.HandleFunc("GET /api/name-detail", h.handleNameDetail)
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

type generateRequest struct {
	Surname      string   `json:"surname"`
	Born         string   `json:"born"`
	Sex          string   `json:"sex"`
	PoetryMode   int      `json:"poetry_mode"`
	XiYongMethod string   `json:"xiyong_method"`
	AvoidChars   []string `json:"avoid_chars"`
	RequireChars []string `json:"require_chars"`
	FilterType   uint     `json:"filter_type"`
}

type generateResponse struct {
	TopNames []fate.NameResult `json:"top_names"`
	Total    int               `json:"total"`
}

func (h *Handler) handleGenerate(w http.ResponseWriter, r *http.Request) {
	var req generateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}
	if req.Surname == "" {
		writeError(w, 400, "surname is required")
		return
	}

	born, err := time.Parse("2006/01/02 15:04", req.Born)
	if err != nil {
		writeError(w, 400, "invalid born date format, use 2006/01/02 15:04")
		return
	}

	sex := fate.SexBoy
	if req.Sex == "girl" {
		sex = fate.SexGirl
	}

	l := [2]string{}
	runes := []rune(req.Surname)
	if len(runes) >= 1 {
		l[0] = string(runes[0])
	}
	if len(runes) >= 2 {
		l[1] = string(runes[1])
	}

	filterOpt := fate.NewFilter(fate.FilterOption{
		CharacterFilter:     true,
		RegularFilter:       true,
		WuXingFilter:        true,
		PoetryMode:          req.PoetryMode,
		XiYongMethod:        req.XiYongMethod,
		AvoidCharacters:     req.AvoidChars,
		RequireCharacters:   req.RequireChars,
		CharacterFilterType: fate.CharacterFilterType(req.FilterType),
	})

	sess := h.fate.NewSessionWithFilter(filterOpt)
	input := &fate.Input{
		Last: l,
		Born: born,
		Sex:  sex,
	}

	if err := sess.Start(input); err != nil {
		writeError(w, 500, "failed to start: "+err.Error())
		return
	}
	sess.Wait()

	output := input.Output()
	topNames := output.TopNames()
	allNames := output.AllNames()

	writeJSON(w, 200, generateResponse{
		TopNames: topNames,
		Total:    len(allNames),
	})
}

func (h *Handler) handleNames(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		writeError(w, 400, "task_id is required")
		return
	}
	writeError(w, 404, "task not found (stateless server)")
}

func (h *Handler) handleNameDetail(w http.ResponseWriter, r *http.Request) {
	surname := r.URL.Query().Get("surname")
	char1 := r.URL.Query().Get("char1")
	char2 := r.URL.Query().Get("char2")
	bornStr := r.URL.Query().Get("born")
	sexStr := r.URL.Query().Get("sex")

	if surname == "" || char1 == "" || char2 == "" {
		writeError(w, 400, "surname, char1, char2 are required")
		return
	}

	born, err := time.Parse("2006/01/02 15:04", bornStr)
	if err != nil {
		born, err = time.Parse("2006-01-02 15:04", bornStr)
	}
	if err != nil {
		born, err = time.Parse("2006/01/02", bornStr)
	}
	if err != nil {
		writeError(w, 400, "invalid born date format, use 2006/01/02 15:04 or 2006/01/02")
		return
	}
	sex := fate.SexBoy
	if sexStr == "girl" {
		sex = fate.SexGirl
	}

	l := [2]string{}
	runes := []rune(surname)
	if len(runes) >= 1 {
		l[0] = string(runes[0])
	}
	if len(runes) >= 2 {
		l[1] = string(runes[1])
	}

	filterOpt := fate.NewFilter(fate.FilterOption{
		CharacterFilter: true,
		RegularFilter:   true,
		DaYanFilter:     false,
		WuXingFilter:    true,
	})
	sess := h.fate.NewSessionWithFilter(filterOpt)
	input := &fate.Input{
		Last: l,
		Born: born,
		Sex:  sex,
	}
	if err := sess.Start(input); err != nil {
		writeError(w, 500, "failed to start session: "+err.Error())
		return
	}
	sess.Wait()

	output := input.Output()
	topNames := output.TopNames()

	for _, nr := range topNames {
		if nr.Char1.Char == char1 && nr.Char2.Char == char2 {
			writeJSON(w, 200, map[string]any{"name_result": nr})
			return
		}
	}

	writeError(w, 404, "name not found in results")
}
