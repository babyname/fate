package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/babyname/fate/v4"
	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/internal/repository"
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
	mu      sync.RWMutex
	result  *fate.GenerateResult
	request generateRequest
	created time.Time
}

func (e *taskEntry) loadResult() *fate.GenerateResult {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.result
}

func (e *taskEntry) storeResult(r *fate.GenerateResult) {
	e.mu.Lock()
	e.result = r
	e.mu.Unlock()
}

type taskStore struct {
	mu    sync.RWMutex
	tasks map[string]*taskEntry
}

type generateRequest struct {
	Surname       string `json:"surname"`
	Born          string `json:"born"`
	Sex           string `json:"sex"`
	PoetryMode    int    `json:"poetry_mode"`
	XiYongMethod  string `json:"xi_yong_method"`
	AvoidChars    string `json:"avoid_chars"`
	RequireChars  string `json:"require_chars"`
	FilterType    int    `json:"filter_type"`
}

// NewHandler creates an initialized Handler with all routes registered.
func NewHandler(f fate.Fate, staticFS fs.FS, repo *repository.Repository) *Handler {
	h := &Handler{
		fate:     f,
		mux:      http.NewServeMux(),
		store:    newTaskStore(),
		staticFS: staticFS,
		repo:     repo,
	}
	h.fileServer = http.FileServer(http.FS(h.staticFS))
	h.registerRoutes()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
	if r.Method == http.MethodOptions {
		w.WriteHeader(204)
		return
	}
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) registerRoutes() {
	h.mux.HandleFunc("POST /api/v1/generate", h.handleGenerate)
	h.mux.HandleFunc("GET /api/v1/status", h.handleStatus)
	h.mux.HandleFunc("GET /api/v1/result", h.handleResult)
	h.mux.HandleFunc("GET /api/v1/explore", h.handleExplore)
	h.mux.HandleFunc("GET /api/v1/names", h.handleNames)
	h.mux.HandleFunc("GET /api/v1/name_detail", h.handleNameDetail)
	h.mux.HandleFunc("GET /api/v1/characters", h.handleCharacters)
	h.mux.HandleFunc("GET /", h.handleStatic)
}

func newTaskStore() *taskStore {
	return &taskStore{tasks: make(map[string]*taskEntry)}
}

func (s *taskStore) Get(key string) (*taskEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.tasks[key]
	return entry, ok
}

func (s *taskStore) Set(key string, entry *taskEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[key] = entry
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
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

	// Parse optional comma-separated char lists.
	avoidChars := strings.Split(req.AvoidChars, ",")
	requireChars := strings.Split(req.RequireChars, ",")
	if len(avoidChars) == 1 && avoidChars[0] == "" {
		avoidChars = nil
	}
	if len(requireChars) == 1 && requireChars[0] == "" {
		requireChars = nil
	}

	// Generate synchronously.
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
			AvoidCharacters:     avoidChars,
			RequireCharacters:   requireChars,
			CharacterFilterType: fate.CharacterFilterType(req.FilterType),
		},
	})
	if err != nil {
		writeError(w, 500, "failed to generate: "+err.Error())
		return
	}

	taskID := fmt.Sprintf("%d", time.Now().UnixNano())
	entry := &taskEntry{
		result:  result,
		request: req,
		created: time.Now(),
	}
	h.store.Set(taskID, entry)

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

	result := entry.loadResult()
	if result == nil {
		writeError(w, 500, "result not available")
		return
	}

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

	result := entry.loadResult()
	if result == nil {
		writeError(w, 500, "result not available")
		return
	}

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

	result := entry.loadResult()
	if result == nil {
		writeError(w, 500, "result not available")
		return
	}

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
			result := entry.loadResult()
			if result != nil {
				// Search top names for a match.
				for i := range result.TopNames {
					if result.TopNames[i].Char1.Char == char1 && result.TopNames[i].Char2.Char == char2 {
						writeJSON(w, 200, map[string]any{"name_result": result.TopNames[i]})
						return
					}
				}

				// Return partial info from CharMap.
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
	}

	writeError(w, 404, "name not found")
}

func (h *Handler) handleCharacters(w http.ResponseWriter, r *http.Request) {
	char := r.URL.Query().Get("char")
	if char == "" {
		writeError(w, 400, "character is required")
		return
	}

	ctx := r.Context()
	results, err := h.repo.GetCharacters(ctx, repository.Char(char))
	if err != nil {
		writeError(w, 500, "query failed: "+err.Error())
		return
	}

	if results == nil {
		results = []*ent.Character{}
	}

	writeJSON(w, 200, map[string]any{"characters": results})
}

func (h *Handler) handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" || path == "" {
		path = "/index.html"
	}

	// If staticFS is nil, fallback gracefully.
	if h.staticFS == nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("Fate API Server\n\nAvailable endpoints:\n  POST /api/v1/generate\n  GET  /api/v1/status\n  GET  /api/v1/result\n  GET  /api/v1/explore\n  GET  /api/v1/characters\n"))
		return
	}

	// Try the requested path.
	cleaned := strings.TrimPrefix(path, "/")
	f, err := h.staticFS.Open(cleaned)
	if err != nil {
		// Fallback to index.html for SPA routes.
		f, err = h.staticFS.Open("index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}
	defer f.Close()

	// Read file content and serve.
	data, err := io.ReadAll(f)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Set content type based on extension.
	switch {
	case strings.HasSuffix(path, ".html"):
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case strings.HasSuffix(path, ".css"):
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case strings.HasSuffix(path, ".js"):
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case strings.HasSuffix(path, ".png"):
		w.Header().Set("Content-Type", "image/png")
	case strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg"):
		w.Header().Set("Content-Type", "image/jpeg")
	case strings.HasSuffix(path, ".svg"):
		w.Header().Set("Content-Type", "image/svg+xml")
	case strings.HasSuffix(path, ".ico"):
		w.Header().Set("Content-Type", "image/x-icon")
	case strings.HasSuffix(path, ".json"):
		w.Header().Set("Content-Type", "application/json")
	case strings.HasSuffix(path, ".woff2"):
		w.Header().Set("Content-Type", "font/woff2")
	case strings.HasSuffix(path, ".ttf"):
		w.Header().Set("Content-Type", "font/ttf")
	}

	w.WriteHeader(200)
	_, _ = w.Write(data)
}
