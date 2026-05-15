package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/babyname/fate"
	"github.com/babyname/fate/config"
	"github.com/babyname/fate/internal/analysis"
	"github.com/babyname/fate/internal/database"
	filterpkg "github.com/babyname/fate/internal/filter"
	"github.com/babyname/fate/internal/poetry"
	"github.com/babyname/fate/internal/repository"
	v2 "github.com/godcong/chronos/v2"
)

var sessions sync.Map

type Server struct {
	cfg  *config.Config
	repo *repository.Repository
	fate fate.Fate
	mux  *http.ServeMux
}

func NewServer(cfg *config.Config) (*Server, error) {
	b := database.New(cfg.Database)
	client, err := b.Client()
	if err != nil {
		return nil, err
	}
	repo := repository.New(client)
	f, err := fate.New(cfg)
	if err != nil {
		return nil, err
	}

	s := &Server{
		cfg:  cfg,
		repo: repo,
		fate: f,
		mux:  http.NewServeMux(),
	}
	s.routes()
	return s, nil
}

func (s *Server) routes() {
	s.mux.HandleFunc("POST /api/generate", s.handleGenerate)
	s.mux.HandleFunc("GET /api/generate/{taskID}/status", s.handleStatus)
	s.mux.HandleFunc("GET /api/generate/{taskID}/names", s.handleNames)
	s.mux.HandleFunc("GET /api/name/detail", s.handleNameDetail)
	s.mux.HandleFunc("GET /api/poetry/search", s.handlePoetrySearch)
	s.mux.HandleFunc("GET /api/poetry/random", s.handlePoetryRandom)
	s.mux.Handle("GET /", http.FileServer(http.FS(webSub)))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

type GenerateRequest struct {
	Surname         string `json:"surname"`
	Born            string `json:"born"`
	Sex             string `json:"sex"`
	PoetryMode      int    `json:"poetry_mode"`
	XiYongMethod    string `json:"xiyong_method"`
	FilterStrictness string `json:"filter_strictness"`
}

type GenerateResponse struct {
	TaskID   string                `json:"task_id"`
	State    string                `json:"state"`
	TopNames []analysis.NameResult `json:"top_names"`
	Total    int                   `json:"total"`
}

type NameEntry struct {
	Surname string  `json:"surname"`
	Char1   string  `json:"char1"`
	Char2   string  `json:"char2"`
	Score   float64 `json:"score"`
	Grade   string  `json:"grade"`
}

func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	var req GenerateRequest
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

	filterOpt := filterpkg.NewDefaultFilterOption()
	filterOpt.DaYanFilter = true
	filterOpt.WuXingFilter = true
	filterOpt.PoetryMode = req.PoetryMode
	filterOpt.XiYongMethod = req.XiYongMethod
	filterOpt.FilterStrictness = req.FilterStrictness

	sess := s.fate.NewSessionWithFilter(filterpkg.NewFilter(filterOpt))

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

	taskID := fmt.Sprintf("%d", time.Now().UnixNano())
	sessions.Store(taskID, output)

	writeJSON(w, 200, GenerateResponse{
		TaskID:   taskID,
		State:    "finished",
		TopNames: topNames,
		Total:    len(allNames),
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskID")
	if _, ok := sessions.Load(taskID); !ok {
		writeError(w, 404, "task not found")
		return
	}
	writeJSON(w, 200, map[string]string{"state": "finished"})
}

func (s *Server) handleNames(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskID")
	outputI, ok := sessions.Load(taskID)
	if !ok {
		writeError(w, 404, "task not found")
		return
	}
	output := outputI.(*fate.Output)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	allNames := output.AllNames()
	start := (page - 1) * size
	end := start + size
	if start >= len(allNames) {
		start = len(allNames)
	}
	if end >= len(allNames) {
		end = len(allNames)
	}

	surname := ""
	basic := output.Basic()
	if basic != nil && basic.LastName[0] != nil {
		surname = basic.LastName[0].Char
	}

	entries := make([]NameEntry, 0, end-start)
	for i := start; i < end; i++ {
		sn := allNames[i]
		c1 := ""
		c2 := ""
		if sn.Name[0] != nil {
			c1 = sn.Name[0].Char
		}
		if sn.Name[1] != nil {
			c2 = sn.Name[1].Char
		}
		entries = append(entries, NameEntry{
			Surname: surname,
			Char1:   c1,
			Char2:   c2,
			Score:   sn.Score,
			Grade:   sn.Grade,
		})
	}

	writeJSON(w, 200, map[string]any{
		"names": entries,
		"total": len(allNames),
		"page":  page,
		"size":  size,
	})
}

func (s *Server) handleNameDetail(w http.ResponseWriter, r *http.Request) {
	surname := r.URL.Query().Get("surname")
	char1 := r.URL.Query().Get("char1")
	char2 := r.URL.Query().Get("char2")
	bornStr := r.URL.Query().Get("born")
	sexStr := r.URL.Query().Get("sex")

	if surname == "" || char1 == "" || char2 == "" {
		writeError(w, 400, "surname, char1, char2 are required")
		return
	}

	born, _ := time.Parse("2006/01/02 15:04", bornStr)
	sx := 1
	if sexStr == "girl" {
		sx = 0
	}

	fateData, _ := v2.GetFateData(&v2.FateInput{
		BirthDate: born,
		Gender:    sx,
		Surname:   surname,
	})

	c1, err := s.repo.GetCharacter(r.Context(), repository.Char(char1))
	if err != nil {
		writeError(w, 404, "character "+char1+" not found")
		return
	}
	c2, err := s.repo.GetCharacter(r.Context(), repository.Char(char2))
	if err != nil {
		writeError(w, 404, "character "+char2+" not found")
		return
	}

	l := [2]string{}
	runes := []rune(surname)
	if len(runes) >= 1 {
		l[0] = string(runes[0])
	}
	if len(runes) >= 2 {
		l[1] = string(runes[1])
	}

	l1, l2 := 0, 0
	ln, err := s.repo.QueryLastName(r.Context(), l)
	if err == nil && ln[0] != nil {
		l1 = ln[0].ScienceStroke
		if ln[1] != nil {
			l2 = ln[1].ScienceStroke
		}
	}

	nr := analysis.BuildNameResult(0, surname, c1, c2, l1, l2, fateData)

	pn := poetry.NewNamer(s.repo.Client)
	sources, _ := pn.FindCharPoetry(r.Context(), char1)
	if len(sources) > 0 {
		nr.PoetrySource = &analysis.PoetrySourceInfo{
			Title:    sources[0].Title,
			Author:   sources[0].Author,
			Dynasty:  sources[0].Dynasty,
			Sentence: sources[0].Sentence,
			Type:     sources[0].Type,
		}
	}

	writeJSON(w, 200, map[string]any{"name_result": nr})
}

func (s *Server) handlePoetrySearch(w http.ResponseWriter, r *http.Request) {
	char := r.URL.Query().Get("char")
	if char == "" {
		writeError(w, 400, "char parameter is required")
		return
	}

	pn := poetry.NewNamer(s.repo.Client)
	sources, err := pn.FindCharPoetry(r.Context(), char)
	if err != nil {
		writeError(w, 500, "query failed: "+err.Error())
		return
	}

	writeJSON(w, 200, map[string]any{"results": sources})
}

func (s *Server) handlePoetryRandom(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.URL.Query().Get("count"))
	if count < 1 || count > 20 {
		count = 5
	}

	poems, err := s.repo.Poem.Query().Limit(count).All(r.Context())
	if err != nil {
		writeError(w, 500, "query failed: "+err.Error())
		return
	}

	writeJSON(w, 200, map[string]any{"poems": poems})
}
