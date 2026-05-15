// Package session 提供命名会话管理，负责控制命名生成的生命周期和并发调度。
package session

import (
	"context"
	"sort"
	"sync/atomic"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/internal/analysis"
	filterpkg "github.com/babyname/fate/internal/filter"
	"github.com/babyname/fate/internal/naming"
	"github.com/babyname/fate/internal/rating"
	"github.com/babyname/fate/internal/repository"
	"github.com/babyname/fate/internal/wuge"
	"github.com/babyname/fate/log"
	v2 "github.com/godcong/chronos/v2"
	"golang.org/x/sync/errgroup"
)

// SessionState 表示会话的运行状态。
//
//nolint:revive // stutter name is intentional: SessionState is clearer than State in a session package
type SessionState int32

const (
	// SessionStateWaiting 会话等待启动。
	SessionStateWaiting SessionState = iota
	// SessionStateGenerating 会话正在生成名字。
	SessionStateGenerating
	// SessionStateFinish 会话已完成生成。
	SessionStateFinish
	// SessionStateCanceled 会话已被取消。
	SessionStateCanceled
	// SessionStateFailed 会话生成失败。
	SessionStateFailed
)

// Session 定义命名会话的接口，提供启动、停止和等待等操作。
type Session interface {
	Context() context.Context
	Start(input *Input) error
	Stop() error
	Err() error
	Wait()
}

type session struct {
	ctx        context.Context
	cancel     context.CancelFunc
	db         *repository.Repository
	group      errgroup.Group
	state      int32
	filter     filterpkg.Filter
	outputDone chan struct{}
	fateData   *v2.FateData

	chars map[int][]*ent.Character

	name   chan naming.FirstName
	output *Output
}

// NewSession 创建一个新的命名会话。
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
	s.name = make(chan naming.FirstName, 1024)

	var err error
	s.output = input.Output()
	ln, err := s.db.QueryLastName(s.Context(), input.Last)
	if err != nil {
		return err
	}
	s.output.SetLastName(ln)

	fateData, err := v2.GetFateData(&v2.FateInput{
		BirthDate: input.Born,
		Gender:    int(input.Sex),
		Surname:   input.Last[0],
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
	outputDone := make(chan struct{})
	go func() {
		s.startOutput()
		close(outputDone)
	}()
	s.outputDone = outputDone
	return nil
}

func (s *session) startOutput() {
	put := NewPutFilter()
	defer s.output.SetCacheFilter(put)
	for {
		select {
		case name, ok := <-s.name:
			if !ok {
				return
			}
			put.Put(name)
		case <-s.Context().Done():
			for {
				select {
				case name, ok := <-s.name:
					if !ok {
						return
					}
					put.Put(name)
				default:
					return
				}
			}
		}
	}
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
	if s.outputDone != nil {
		<-s.outputDone
	}
}

func (s *session) generate() error {
	defer close(s.name)
	defer s.close()

	basic := s.output.Basic()
	strokes := getLastStrokeFromBasic(s.filter, basic)
	lucky := wuge.GetLuckyByLastName(strokes[0], strokes[1])

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
				s.SetState(SessionStateFailed)
				return err
			}
			s.chars[stroke] = cs
		}
	}

	var candidates []naming.FirstName
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
			for i2 := range f2s {
				select {
				case <-s.Context().Done():
					s.SetState(SessionStateCanceled)
					return nil
				default:
					candidates = append(candidates, naming.FirstName{f1s[i1], f2s[i2]})
				}
			}
		}
	}

	type scoredEntry struct {
		name  naming.FirstName
		score float64
		grade string
	}

	rater := rating.NewRater(s.fateData)
	scored := make([]scoredEntry, 0, len(candidates))
	for _, fn := range candidates {
		nr := rater.RateName("", fn[0], fn[1])
		scored = append(scored, scoredEntry{name: fn, score: nr.TotalScore, grade: nr.Grade})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	surname := ""
	if basic.LastName[0] != nil {
		surname = basic.LastName[0].Char
	}
	l1 := strokes[0]
	l2 := strokes[1]

	topN := 10
	if topN > len(scored) {
		topN = len(scored)
	}

	topResults := make([]analysis.NameResult, 0, topN)
	for i := 0; i < topN; i++ {
		nr := analysis.BuildNameResult(i+1, surname, scored[i].name[0], scored[i].name[1], l1, l2, s.fateData)
		topResults = append(topResults, nr)
	}

	allScored := make([]ScoredName, 0, len(scored))
	for _, se := range scored {
		allScored = append(allScored, ScoredName{Name: se.name, Score: se.score, Grade: se.grade})
	}
	s.output.SetTopNames(topResults)
	s.output.SetAllNames(allScored)

	for _, se := range scored {
		select {
		case s.name <- se.name:
		default:
		}
	}

	s.SetState(SessionStateFinish)
	return nil
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
	switch ct {
	case filterpkg.CharacterFilterTypeChs:
		return func(c *ent.Character) int {
			return c.SimplifiedStroke
		}
	case filterpkg.CharacterFilterTypeCht:
		return func(c *ent.Character) int {
			return c.TraditionalStroke
		}
	case filterpkg.CharacterFilterTypeKangxi:
		return func(c *ent.Character) int {
			return c.KangxiStroke
		}
	case filterpkg.CharacterFilterTypeDefault:
		fallthrough
	default:
		return func(c *ent.Character) int {
			return c.ScienceStroke
		}
	}
}
