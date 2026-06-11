package http

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/babyname/fate/v4"
	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/ent/character"
	"github.com/babyname/fate/v4/internal/analysis"
	"github.com/babyname/fate/v4/internal/chronosfate"
	"github.com/babyname/fate/v4/internal/repository"
	"github.com/godcong/chronos/v2"
)

type Handler struct {
	fate       fate.Fate
	mux        *http.ServeMux
	store      *taskStore
	staticFS   fs.FS
	fileServer http.Handler
	repo       *repository.Repository
}

type taskEntry struct {
	result  *fate.GenerateResult
	request generateRequest
	created time.Time
}

type taskStore struct {
	mu    sync.RWMutex
	tasks map[string]*taskEntry
}

func newTaskStore() *taskStore {
	return &taskStore{
		tasks: make(map[string]*taskEntry),
	}
}

func (s *taskStore) Set(id string, entry *taskEntry) {
	s.mu.Lock()
	s.tasks[id] = entry
	s.mu.Unlock()
}

func (s *taskStore) Get(id string) (*taskEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.tasks[id]
	return e, ok
}

func NewHandler(f fate.Fate, staticFS fs.FS, repo *repository.Repository) *Handler {
	h := &Handler{
		fate:       f,
		mux:        http.NewServeMux(),
		store:      newTaskStore(),
		staticFS:   staticFS,
		fileServer: nil,
		repo:       repo,
	}

	if staticFS != nil {
		h.fileServer = http.FileServerFS(staticFS)
	}

	h.mux.HandleFunc("GET /health", h.handleHealth)
	h.mux.HandleFunc("POST /api/generate", h.handleGenerate)
	h.mux.HandleFunc("GET /api/generate/status", h.handleStatus)
	h.mux.HandleFunc("GET /api/generate/result", h.handleResult)
	h.mux.HandleFunc("GET /api/generate/explore", h.handleExplore)
	h.mux.HandleFunc("GET /api/generate/names", h.handleNames)
	h.mux.HandleFunc("POST /api/name-score", h.handleNameScore)
	h.mux.HandleFunc("GET /api/name-detail", h.handleNameDetail)
	h.mux.HandleFunc("/{path...}", h.handleDefault)

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) handleDefault(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	path := r.URL.Path

	if strings.HasPrefix(path, "/api/") {
		http.NotFound(w, r)
		return
	}

	if h.fileServer != nil {
		if path != "/" {
			if _, err := fs.Stat(h.staticFS, strings.TrimPrefix(path, "/")); err == nil {
				h.fileServer.ServeHTTP(w, r)
				return
			}
		}

		f, err := h.staticFS.Open("index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		stat, _ := f.Stat()
		if stat != nil {
			buf := make([]byte, stat.Size())
			f.Read(buf)
			w.Write(buf)
			return
		}
	}

	http.NotFound(w, r)
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

	// Synchronous generation using the new API.
	result, err := h.fate.Generate(r.Context(), fate.GenerateRequest{
		LastName: req.Surname,
		Born:     born,
		Sex:      sex,
		Filter: &fate.FilterOption{
			CharacterFilter:     true,
			RegularFilter:       true,
			WuXingFilter:        true,
			PoetryMode:          req.PoetryMode,
			XiYongMethod:        req.XiYongMethod,
			AvoidCharacters:     req.AvoidChars,
			RequireCharacters:   req.RequireChars,
			CharacterFilterType: fate.CharacterFilterType(req.FilterType),
		},
	})
	if err != nil {
		writeError(w, 500, "failed to generate: "+err.Error())
		return
	}

	taskID := fmt.Sprintf("%d", time.Now().UnixNano())
	h.store.Set(taskID, &taskEntry{
		result:  result,
		request: req,
		created: time.Now(),
	})

	writeJSON(w, 200, map[string]any{
		"task_id": taskID,
		"state":   "done",
		"total":   result.ExcellentTable.Len(),
	})
}

func (h *Handler) handleStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		writeError(w, 400, "task_id is required")
		return
	}

	_, ok := h.store.Get(taskID)
	if !ok {
		writeError(w, 404, "task not found")
		return
	}

	writeJSON(w, 200, map[string]any{
		"task_id": taskID,
		"state":   "done",
	})
}

func (h *Handler) handleResult(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		writeError(w, 400, "task_id is required")
		return
	}

	entry, ok := h.store.Get(taskID)
	if !ok {
		writeError(w, 404, "task not found")
		return
	}

	result := entry.result
	topNames := result.TopNames

	top3 := make([]fate.NameResult, 0, 3)
	for i := range topNames {
		if i >= 3 {
			break
		}
		top3 = append(top3, topNames[i])
	}

	top10Entries := result.ExcellentTable.TopN(10)
	for i := range top10Entries {
		result.ExcellentTable.MarkShown(top10Entries[i].Char1, top10Entries[i].Char2)
	}

	writeJSON(w, 200, map[string]any{
		"task_id":   taskID,
		"top_names": topNames,
		"top3":      top3,
		"top10":     top10Entries,
		"total":     result.ExcellentTable.Len(),
	})
}

func (h *Handler) handleExplore(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		writeError(w, 400, "task_id is required")
		return
	}

	entry, ok := h.store.Get(taskID)
	if !ok {
		writeError(w, 404, "task not found")
		return
	}

	result := entry.result
	table := result.ExcellentTable

	count := 10
	if c := r.URL.Query().Get("count"); c != "" {
		fmt.Sscanf(c, "%d", &count)
		if count <= 0 || count > 100 {
			count = 10
		}
	}

	hasPoetry := r.URL.Query().Get("has_poetry")
	wuxing := r.URL.Query().Get("wuxing")

	var filter func(fate.ExcellentEntry) bool
	if hasPoetry == "true" || wuxing != "" {
		filter = func(e fate.ExcellentEntry) bool {
			if hasPoetry == "true" && !e.HasPoetry {
				return false
			}
			if wuxing != "" && e.WuXing1 != wuxing && e.WuXing2 != wuxing {
				return false
			}
			return true
		}
	}

	names := table.Explore(count, filter)

	writeJSON(w, 200, map[string]any{
		"task_id":     taskID,
		"names":       names,
		"shown_count": table.ShownCount(),
		"total":       table.Len(),
	})
}

func (h *Handler) handleNames(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		writeError(w, 400, "task_id is required")
		return
	}

	entry, ok := h.store.Get(taskID)
	if !ok {
		writeError(w, 404, "task not found")
		return
	}

	result := entry.result
	topNames := result.TopNames

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
		if page < 1 {
			page = 1
		}
	}

	pageSize := 20
	if s := r.URL.Query().Get("size"); s != "" {
		fmt.Sscanf(s, "%d", &pageSize)
		if pageSize <= 0 || pageSize > 100 {
			pageSize = 20
		}
	}

	total := len(topNames)
	start := (page - 1) * pageSize
	if start >= total {
		start = 0
		page = 1
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	pagedNames := topNames[start:end]

	writeJSON(w, 200, map[string]any{
		"task_id": taskID,
		"names":   pagedNames,
		"page":    page,
		"size":    pageSize,
		"total":   total,
	})
}

func (h *Handler) handleNameDetail(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	char1 := r.URL.Query().Get("char1")
	char2 := r.URL.Query().Get("char2")

	if char1 == "" || char2 == "" {
		writeError(w, 400, "char1, char2 are required")
		return
	}

	if taskID != "" {
		entry, ok := h.store.Get(taskID)
		if ok {
			result := entry.result

			// Search top names for a match.
			for i := range result.TopNames {
				if result.TopNames[i].Char1.Char == char1 && result.TopNames[i].Char2.Char == char2 {
					writeJSON(w, 200, map[string]any{"name_result": result.TopNames[i]})
					return
				}
			}

			// Return partial info from CharMap if characters exist.
			c1, ok1 := result.CharMap[char1]
			c2, ok2 := result.CharMap[char2]
			if ok1 && ok2 {
				surname := ""
				if len(result.TopNames) > 0 {
					surname = result.TopNames[0].Surname
				}
				nr := fate.NameResult{
					FullName: surname + char1 + char2,
					Surname:  surname,
					Char1:    c1,
					Char2:    c2,
				}
				writeJSON(w, 200, map[string]any{"name_result": nr})
				return
			}
		}
	}

	writeError(w, 404, "name not found")
}

type nameScoreRequest struct {
	Surname string `json:"surname"`
	Name1   string `json:"name1"`
	Name2   string `json:"name2"`
	Born    string `json:"born"`
	Sex     string `json:"sex"`
}

func (h *Handler) handleNameScore(w http.ResponseWriter, r *http.Request) {
	var req nameScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request: "+err.Error())
		return
	}
	defer r.Body.Close()

	surnameRunes := []rune(req.Surname)
	if len(surnameRunes) == 0 {
		writeError(w, 400, "surname required")
		return
	}

	name1Runes := []rune(req.Name1)
	name2Runes := []rune(req.Name2)
	if len(name1Runes) == 0 || len(name2Runes) == 0 {
		writeError(w, 400, "both name characters required")
		return
	}

	ctx := r.Context()
	surnameChars, err := h.repo.QueryLastName(ctx, [2]string{string(surnameRunes[0]), ""})
	if err != nil {
		writeError(w, 500, "query surname: "+err.Error())
		return
	}

	c1, err := h.repo.Character.Query().Where(character.CharEQ(string(name1Runes[0]))).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			writeError(w, 404, "character not found: "+string(name1Runes[0]))
			return
		}
		writeError(w, 500, "query character: "+err.Error())
		return
	}

	c2, err := h.repo.Character.Query().Where(character.CharEQ(string(name2Runes[0]))).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			writeError(w, 404, "character not found: "+string(name2Runes[0]))
			return
		}
		writeError(w, 500, "query character: "+err.Error())
		return
	}

	l1 := 0
	if surnameChars[0] != nil {
		l1 = surnameChars[0].ScienceStroke
	}
	l2 := 0

	var fateData *chronosfate.FateData
	if req.Born != "" {
		sexInt := 1
		if req.Sex == "2" || req.Sex == "女" {
			sexInt = 2
		}
		input := chronosfate.FateInput{
			Calendar:     chronos.NewSolarCalendar(req.Born),
			Gender:       sexInt,
			XiYongMethod: chronosfate.XiYongMethodBalance,
		}
		fateData, err = chronosfate.GetFateData(input)
		if err != nil {
			writeError(w, 500, "calculate fate: "+err.Error())
			return
		}
	}

	nr := analysis.BuildNameResult(0, req.Surname, c1, c2, l1, l2, fateData)

	writeJSON(w, 200, map[string]any{
		"name_result": nr,
	})
}
