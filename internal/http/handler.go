package http

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/babyname/fate"
	"github.com/babyname/fate/internal/analysis"
	"github.com/babyname/fate/internal/session"
)

type Handler struct {
	fate      fate.Fate
	mux       *http.ServeMux
	store     *taskStore
	staticFS  fs.FS
	fileServer http.Handler
}

type taskEntry struct {
	session fate.Session
	input   *fate.Input
	created time.Time
}

type taskStore struct {
	mu     sync.RWMutex
	tasks  map[string]*taskEntry
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

func NewHandler(f fate.Fate, staticFS fs.FS) *Handler {
	h := &Handler{
		fate:       f,
		mux:        http.NewServeMux(),
		store:      newTaskStore(),
		staticFS:   staticFS,
		fileServer: nil,
	}

	if staticFS != nil {
		h.fileServer = http.FileServerFS(staticFS)
	}

	h.mux.HandleFunc("GET /health", h.handleHealth)
	h.mux.HandleFunc("POST /api/generate", h.handleGenerate)
	h.mux.HandleFunc("GET /api/generate/status", h.handleStatus)
	h.mux.HandleFunc("GET /api/generate/result", h.handleResult)
	h.mux.HandleFunc("GET /api/generate/explore", h.handleExplore)
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

	taskID := fmt.Sprintf("%d", time.Now().UnixNano())
	h.store.Set(taskID, &taskEntry{
		session: sess,
		input:   input,
		created: time.Now(),
	})

	writeJSON(w, 200, map[string]any{
		"task_id": taskID,
		"state":   "computing",
	})
}

func (h *Handler) handleStatus(w http.ResponseWriter, r *http.Request) {
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

	state := entry.session.State()
	stateStr := ""
	switch state {
	case fate.SessionStateGenerating:
		stateStr = "computing"
	case fate.SessionStateFinish:
		stateStr = "done"
	case fate.SessionStateFailed:
		stateStr = "failed"
	case fate.SessionStateCanceled:
		stateStr = "canceled"
	default:
		stateStr = "waiting"
	}

	resp := map[string]any{
		"task_id": taskID,
		"state":   stateStr,
	}

	if state == fate.SessionStateFinish {
		output := entry.input.Output()
		table := output.ExcellentTable()
		if table != nil {
			resp["total"] = table.Len()
		}
	}

	if state == fate.SessionStateFailed {
		if err := entry.session.Err(); err != nil {
			resp["error"] = err.Error()
		}
	}

	writeJSON(w, 200, resp)
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

	if entry.session.State() != fate.SessionStateFinish {
		writeError(w, 409, "task not ready")
		return
	}

	output := entry.input.Output()
	topNames := output.TopNames()
	table := output.ExcellentTable()

	top3 := make([]fate.NameResult, 0, 3)
	for i := range topNames {
		if i >= 3 {
			break
		}
		top3 = append(top3, topNames[i])
	}

	top10Entries := make([]fate.ExcellentEntry, 0, 10)
	if table != nil {
		top10Entries = table.TopN(10)
	}

	writeJSON(w, 200, map[string]any{
		"task_id":   taskID,
		"top_names": topNames,
		"top3":      top3,
		"top10":     top10Entries,
		"total":     output.Total(),
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

	if entry.session.State() != fate.SessionStateFinish {
		writeError(w, 409, "task not ready")
		return
	}

	output := entry.input.Output()
	table := output.ExcellentTable()
	if table == nil {
		writeError(w, 500, "excellent table not available")
		return
	}

	count := 10
	if c := r.URL.Query().Get("count"); c != "" {
		fmt.Sscanf(c, "%d", &count)
		if count <= 0 || count > 20 {
			count = 10
		}
	}

	hasPoetry := r.URL.Query().Get("has_poetry")
	wuxing := r.URL.Query().Get("wuxing")

	var filter func(session.ExcellentEntry) bool
	if hasPoetry == "true" || wuxing != "" {
		filter = func(e session.ExcellentEntry) bool {
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
		if ok && entry.session.State() == fate.SessionStateFinish {
			output := entry.input.Output()
			topNames := output.TopNames()
			for i := range topNames {
				if topNames[i].Char1.Char == char1 && topNames[i].Char2.Char == char2 {
					writeJSON(w, 200, map[string]any{"name_result": topNames[i]})
					return
				}
			}

			charMap := output.CharMap()
			table := output.ExcellentTable()
			if charMap != nil && table != nil {
				c1 := charMap[char1]
				c2 := charMap[char2]
				if c1 != nil && c2 != nil {
					basic := output.Basic()
					surname := ""
					l1, l2 := 0, 0
					if basic.LastName[0] != nil {
						surname = basic.LastName[0].Char
						l1 = basic.LastName[0].ScienceStroke
					}
					if basic.LastName[1] != nil {
						surname += basic.LastName[1].Char
						l2 = basic.LastName[1].ScienceStroke
					}
					fateData := output.FateData()
					if fateData != nil {
						nr := analysis.BuildNameResult(0, surname, c1, c2, l1, l2, fateData)
						writeJSON(w, 200, map[string]any{"name_result": nr})
						return
					}
				}
			}
		}
	}

	writeError(w, 404, "name not found")
}
