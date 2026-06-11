package session

import (
	"context"
	"sync/atomic"

	v2 "github.com/godcong/chronos/v2"
	"github.com/babyname/fate/v4/internal/chronosfate"
	filterpkg "github.com/babyname/fate/v4/internal/filter"
	"github.com/babyname/fate/v4/internal/log"
	"github.com/babyname/fate/v4/internal/repository"
	namesvc "github.com/babyname/fate/v4/internal/service/naming"
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

	pipeline *namesvc.Pipeline
	output   *Output
}

func NewSession(db *repository.Repository, f filterpkg.Filter) Session {
	return &session{
		db:       db,
		filter:   f,
		pipeline: namesvc.NewPipeline(db),
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

	s.output = input.Output()
	ln, err := s.db.QueryLastName(s.Context(), input.Last)
	if err != nil {
		return err
	}
	s.output.SetLastName(ln)

	// Calculate fate data eagerly so both session and pipeline have it.
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

	result, err := s.pipeline.Generate(s.ctx, namesvc.GenerateRequest{
		LastName: basic.LastName,
		Sex:      basic.Sex,
		Filter:   s.filter,
		FateData: s.fateData,
	})
	if err != nil {
		log.Error("pipeline generate", err)
		s.SetState(SessionStateFailed)
		return err
	}

	s.output.SetTopNames(result.TopNames)
	s.output.SetExcellentTable(result.ExcellentTable)
	s.output.SetCharMap(result.CharMap)
	// Fate data was already set in Start().

	s.SetState(SessionStateFinish)
	return nil
}

func (s *session) close() {
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
}
