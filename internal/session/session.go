package session

import (
	"sync"
	"sync/atomic"

	"github.com/babyname/fate/ent"
	filterpkg "github.com/babyname/fate/internal/filter"
	"github.com/babyname/fate/internal/naming"
	"github.com/babyname/fate/internal/repository"
	"github.com/babyname/fate/log"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type SessionState int32

const (
	SessionStateWaiting SessionState = iota
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
}

type session struct {
	ctx    context.Context
	cancel context.CancelFunc
	db     *repository.Repository
	group  errgroup.Group
	state  int32
	filter filterpkg.Filter

	chars map[int][]*ent.Character

	name   chan naming.FirstName
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
	s.name = make(chan naming.FirstName, 1024)

	var err error
	s.output = input.Output()
	ln, err := s.db.QueryLastName(s.Context(), input.Last)
	if err != nil {
		return err
	}
	s.output.SetLastName(ln)
	log.Info("generate", "base", s.output.Basic())
	s.SetState(SessionStateGenerating)

	s.group.Go(s.generate)
	go s.startOutput()
	return nil
}

func (s *session) startOutput() {
	put := NewPutFilter()
	defer s.output.SetCacheFilter(put)
	for {
		select {
		case <-s.Context().Done():
			return
		case name, ok := <-s.name:
			if !ok {
				return
			}
			put.Put(name)
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

func (s *session) generate() error {
	defer close(s.name)
	defer s.close()
	lucky, err := s.db.GetWuGeLuckyCached(s.Context(), getLastStrokeFromBasic(s.filter, s.output.Basic()))
	if err != nil {
		log.Error("get wuge lucky", err)
		s.SetState(SessionStateFailed)
		return err
	}
	log.Info("wuge lucky list", "size", len(lucky))

	strokesNeeded := make(map[int]struct{})
	for i := range lucky {
		tmp := lucky[i]
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

	wg := sync.WaitGroup{}
	for i := range lucky {
		tmp := lucky[i]
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

		wg.Add(1)
		go func(wg *sync.WaitGroup, f1s, f2s []*ent.Character) {
			defer wg.Done()
			for i1 := range f1s {
				for i2 := range f2s {
					select {
					case <-s.Context().Done():
						s.SetState(SessionStateCanceled)
						return
					default:
						s.name <- naming.FirstName{
							f1s[i1],
							f2s[i2],
						}
					}
				}
			}
		}(&wg, f1s, f2s)
	}
	wg.Wait()
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
