package fate

import (
	"github.com/babyname/fate/model"
	"github.com/godcong/chronos/v2"
	"golang.org/x/net/context"
)

type Session interface {
	Start(input *Input) error
}

type session struct {
	ctx    context.Context
	cancel context.CancelFunc
	db     *model.Model
	filter Filter
	base   NameBase
	detail []NameDetail
	name   chan Name
	err    error
}

func (s *session) Start(input *Input) error {
	log.Info("start", "input", input)
	s.ctx, s.cancel = context.WithCancel(context.Background())

	var err error
	s.base.Sex = input.Sex
	s.base.Born = chronos.ParseTime(input.Born)
	s.base.LastName, err = s.db.QueryLastName(s.Context(), input.Last)
	if err != nil {
		return err
	}

	return nil
}

func (s *session) Output() <-chan Name {
	return s.name
}

func (s *session) Err() error {
	return s.err
}

func (s *session) Stop() error {
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
	return nil
}

func (s *session) Context() context.Context {
	return s.ctx
}
