package fate

import (
	"sync"
	"sync/atomic"

	"github.com/babyname/fate/cache"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/model"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

// SessionState ...
// ENUM(waiting,generating,finish,canceled,failed)
type SessionState int32

type Session interface {
	Context() context.Context
	Start(input *Input) error
	Stop() error
	Err() error
}

type session struct {
	ctx    context.Context
	cancel context.CancelFunc
	db     *model.Model
	group  errgroup.Group
	state  int32
	filter Filter

	chars map[int][]*ent.Character

	name   chan FirstName
	output *Output
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
	s.name = make(chan FirstName, 1024)

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
	put := cache.NewPutFilter()
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
	lucky, err := s.db.GetWuGeLucky(s.Context(), getLastStrokeFromBasic(s.filter, s.output.Basic()))
	if err != nil {
		log.Error("get wuge lucky", err)
		s.SetState(SessionStateFailed)
		return err
	}
	log.Info("wuge lucky list", "size", len(lucky))
	wg := sync.WaitGroup{}
	var tmp *ent.WuGeLucky
	for i := range lucky {
		tmp = lucky[i]
		if s.filter.CheckStrokeNumber(tmp.FirstStroke1) || s.filter.CheckStrokeNumber(tmp.FirstStroke2) {
			log.Info("current lukcy", "lucky", tmp, "blocked", "stroke")
			continue
		}
		if s.filter.CheckSexFilter(tmp) {
			log.Info("current lukcy", "lucky", tmp, "blocked", "sex")
			continue
		}
		if s.filter.CheckDaYanFilter(tmp) {
			log.Info("current lukcy", "lucky", tmp, "blocked", "dayan")
			continue
		}
		if s.filter.CheckWuXingFilter(tmp.TianGe, tmp.RenGe, tmp.DiGe) {
			log.Info("current lukcy", "lucky", tmp, "blocked", "wuxing")
			continue
		}
		log.Info("lucky get chars", "lucky", tmp)
		var f1s []*ent.Character

		if cs, ok := s.chars[tmp.FirstStroke1]; !ok {
			f1s, err = s.db.GetCharacters(s.Context(), s.filter.QueryStrokeFilter(tmp.FirstStroke1), s.filter.QueryRegularFilter)
			if err != nil {
				log.Error("get first1 name", err)
				s.SetState(SessionStateFailed)
				return err
			}
		} else {
			f1s = cs
		}

		var f2s []*ent.Character
		if cs, ok := s.chars[tmp.FirstStroke2]; !ok {
			f2s, err = s.db.GetCharacters(s.Context(), s.filter.QueryStrokeFilter(tmp.FirstStroke2), s.filter.QueryRegularFilter)
			if err != nil {
				log.Error("get first2 name", err)
				s.SetState(SessionStateFailed)
				return err
			}
		} else {
			f2s = cs
		}

		//make first name
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
						s.name <- FirstName{
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

func getLastStrokeFromBasic(filter Filter, basic *NameBasic) [2]int {
	var strokes [2]int
	strokes[0] = filter.GetCharacterStroke(basic.LastName[0])
	if basic.LastName[1] != nil {
		strokes[1] = filter.GetCharacterStroke(basic.LastName[1])
	}
	return strokes
}
