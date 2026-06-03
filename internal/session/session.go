package session

import (
	"context"
	"sync/atomic"

	v2 "github.com/godcong/chronos/v2"
	"github.com/babyname/fate/v4/internal/chronosfate"
	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/internal/analysis"
	filterpkg "github.com/babyname/fate/v4/internal/filter"
	"github.com/babyname/fate/v4/internal/log"
	"github.com/babyname/fate/v4/internal/naming"
	"github.com/babyname/fate/v4/internal/rating"
	"github.com/babyname/fate/v4/internal/repository"
	"github.com/babyname/fate/v4/internal/wuge"
	"golang.org/x/sync/errgroup"
)

type SessionState int32

const (
	SessionStateWaiting    SessionState = iota
	SessionStateGenerating
	SessionStateFinish
	SessionStateCanceled
	SessionStateFailed
)

type Session interface {
	Context() context.Context
	Start(input *Input) error
	Stop() error
	Err() error
	Wait()
	State() SessionState
}

type session struct {
	ctx      context.Context
	cancel   context.CancelFunc
	db       *repository.Repository
	group    errgroup.Group
	state    int32
	filter   filterpkg.Filter
	fateData *chronosfate.FateData

	chars map[int][]*ent.Character

	output *Output
}

func NewSession(db *repository.Repository, f filterpkg.Filter) Session {
	return &session{
		db:     db,
		chars:  make(map[int][]*ent.Character, 128),
		filter: f,
	}
}

func (s *session) State() SessionState {
	return SessionState(atomic.LoadInt32(&s.state))
}

func (s *session) SetState(state SessionState) {
	atomic.StoreInt32(&s.state, int32(state))
}

func (s *session) Start(input *Input) error {
	log.Info("start", "input", input)
	if s.State() != SessionStateWaiting {
		return nil
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	var err error
	s.output = input.Output()
	ln, err := s.db.QueryLastName(s.Context(), input.Last)
	if err != nil {
		return err
	}
	s.output.SetLastName(ln)

	method := chronosfate.XiYongMethodBalance
	if s.filter.XiYongMethod() == "geju" {
		method = chronosfate.XiYongMethodGeJu
	}
	fateData, err := chronosfate.GetFateData(chronosfate.FateInput{
		Calendar:     v2.NewSolarCalendar(input.Born),
		Gender:       int(input.Sex),
		XiYongMethod: method,
	})
	if err != nil {
		log.Error("get fate data", err)
	} else {
		s.fateData = fateData
		s.output.SetFateData(fateData)
	}

	log.Info("generate", "base", s.output.Basic())
	s.SetState(SessionStateGenerating)

	s.group.Go(s.generate)
	return nil
}

func (s *session) Err() error {
	if s.State() == SessionStateFailed {
		return s.group.Wait()
	}
	return nil
}

func (s *session) Stop() error {
	s.close()
	s.SetState(SessionStateWaiting)
	return nil
}

func (s *session) Context() context.Context {
	return s.ctx
}

func (s *session) Wait() {
	s.group.Wait()
}

func (s *session) generate() error {
	defer s.close()

	basic := s.output.Basic()
	strokes := getLastStrokeFromBasic(s.filter, basic)
	lucky := wuge.GetLuckyByLastName(strokes[0], strokes[1])

	s.preloadChars(lucky)

	poetrySet := s.queryPoetrySet()
	log.Info("poetry set loaded", "count", len(poetrySet))

	table := NewExcellentTable()
	rater := rating.NewRaterWithStrokes(s.fateData, strokes[0], strokes[1])

	s.scoreAllCandidates(lucky, rater, table, poetrySet)

	if table.HeapLen() == 0 && s.filter.FilterStrictness() != "relaxed" {
		s.filter.SetFilterStrictness("moderate")
		s.preloadChars(lucky)
		s.scoreAllCandidates(lucky, rater, table, poetrySet)
	}
	if table.HeapLen() == 0 && s.filter.FilterStrictness() != "relaxed" {
		s.filter.SetFilterStrictness("relaxed")
		s.preloadChars(lucky)
		s.scoreAllCandidates(lucky, rater, table, poetrySet)
	}

	table.Finalize()

	charMap := s.buildCharMap()

	surname := ""
	if basic.LastName[0] != nil {
		surname = basic.LastName[0].Char
	}
	if basic.LastName[1] != nil && basic.LastName[1].Char != "" {
		surname += basic.LastName[1].Char
	}

	topEntries := table.TopN(10)
	topResults := make([]analysis.NameResult, 0, len(topEntries))
	for i, e := range topEntries {
		c1 := charMap[e.Char1]
		c2 := charMap[e.Char2]
		if c1 == nil || c2 == nil {
			continue
		}
		nr := analysis.BuildNameResult(i+1, surname, c1, c2, strokes[0], strokes[1], s.fateData)
		topResults = append(topResults, nr)
	}

	s.output.SetTopNames(topResults)
	s.output.SetExcellentTable(table)
	s.output.SetCharMap(charMap)

	s.SetState(SessionStateFinish)
	return nil
}

func (s *session) scoreAllCandidates(lucky []wuge.WuGeResult, rater *rating.Rater, table *ExcellentTable, poetrySet map[string]bool) {
	poetryMode := s.filter.PoetryMode()

	for i := range lucky {
		tmp := &lucky[i]
		if s.filter.CheckSkipStrokeNumberScope(tmp.FirstStroke1, tmp.FirstStroke2) {
			continue
		}
		if s.filter.CheckSkipSexFilter(tmp) {
			continue
		}
		if s.filter.CheckSkipDaYanFilter(tmp) {
			continue
		}
		if s.filter.CheckSkipWuXingFilter(tmp.TianGe, tmp.RenGe, tmp.DiGe) {
			continue
		}

		f1s := s.chars[tmp.FirstStroke1]
		f2s := s.chars[tmp.FirstStroke2]

		for i1 := range f1s {
			c1 := f1s[i1]
			if poetryMode == 2 && !poetrySet[c1.Char] {
				continue
			}
			for i2 := range f2s {
				select {
				case <-s.Context().Done():
					return
				default:
				}

				c2 := f2s[i2]
				if poetryMode == 2 && !poetrySet[c2.Char] {
					continue
				}

				nr := rater.RateName("", c1, c2)

				hasPoetry := poetrySet[c1.Char] || poetrySet[c2.Char]

				table.TryPush(ExcellentEntry{
					Char1:     c1.Char,
					Char2:     c2.Char,
					Score:     nr.TotalScore,
					Grade:     nr.Grade,
					WuXing1:   c1.WuXing,
					WuXing2:   c2.WuXing,
					HasPoetry: hasPoetry,
				})
			}
		}
	}
	log.Info("score candidates done", "heap_size", table.HeapLen())
}

func (s *session) queryPoetrySet() map[string]bool {
	poetrySet := make(map[string]bool)
	chars, err := s.db.QueryPoetryChars(s.Context())
	if err != nil {
		log.Error("query poetry chars failed", err)
	} else {
		log.Info("query poetry chars", "count", len(chars))
		for _, ch := range chars {
			poetrySet[ch] = true
		}
	}
	return poetrySet
}

func (s *session) buildCharMap() map[string]*ent.Character {
	charMap := make(map[string]*ent.Character)
	for _, chars := range s.chars {
		for _, c := range chars {
			if _, ok := charMap[c.Char]; !ok {
				charMap[c.Char] = c
			}
		}
	}
	return charMap
}

func (s *session) preloadChars(lucky []wuge.WuGeResult) {
	strokesNeeded := make(map[int]struct{})
	for i := range lucky {
		tmp := &lucky[i]
		if !s.filter.CheckSkipStrokeNumberScope(tmp.FirstStroke1, tmp.FirstStroke2) &&
			!s.filter.CheckSkipSexFilter(tmp) &&
			!s.filter.CheckSkipDaYanFilter(tmp) &&
			!s.filter.CheckSkipWuXingFilter(tmp.TianGe, tmp.RenGe, tmp.DiGe) {
			strokesNeeded[tmp.FirstStroke1] = struct{}{}
			strokesNeeded[tmp.FirstStroke2] = struct{}{}
		}
	}

	for stroke := range strokesNeeded {
		if _, ok := s.chars[stroke]; !ok {
			cs, err := s.db.GetCharactersCached(s.Context(), stroke, s.filter.QueryStrokeFilter(stroke), s.filter.QueryRegularFilter)
			if err != nil {
				log.Error("preload characters", err)
				continue
			}
			s.chars[stroke] = cs
		}
	}
}

func (s *session) close() {
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
}

func getLastStrokeFromBasic(filter filterpkg.Filter, basic *naming.NameBasic) [2]int {
	var strokes [2]int
	sg := strokeGetFromFilterType(filter.FilterType())
	strokes[0] = sg(basic.LastName[0])
	if basic.LastName[1] != nil {
		strokes[1] = sg(basic.LastName[1])
	}
	return strokes
}

func strokeGetFromFilterType(ct filterpkg.CharacterFilterType) func(c *ent.Character) int {
	if ct == filterpkg.CharacterFilterTypeDefault {
		return func(c *ent.Character) int {
			return c.ScienceStroke
		}
	}
	if ct.HasType(filterpkg.CharacterFilterTypeChs) {
		return func(c *ent.Character) int {
			return c.SimplifiedStroke
		}
	}
	if ct.HasType(filterpkg.CharacterFilterTypeCht) {
		return func(c *ent.Character) int {
			return c.TraditionalStroke
		}
	}
	if ct.HasType(filterpkg.CharacterFilterTypeKangxi) {
		return func(c *ent.Character) int {
			return c.KangxiStroke
		}
	}
	return func(c *ent.Character) int {
		return c.ScienceStroke
	}
}
