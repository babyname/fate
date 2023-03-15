package fate

import (
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/model"
	"golang.org/x/net/context"
)

type Session interface {
	Start(input *Input) error
	Stop() error
	Err() error
}

type session struct {
	ctx    context.Context
	cancel context.CancelFunc
	db     *model.Model
	chars  map[int][]*ent.Character
	filter Filter

	err    error
	name   chan FirstName
	output *Output
}

func (s *session) Start(input *Input) error {
	log.Info("start", "input", input)
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
	go s.generate()
	go s.startOutput()
	return nil
}

func (s *session) startOutput() {
	for {
		select {
		case <-s.Context().Done():
			return
		case name, ok := <-s.name:
			if !ok {
				return
			}
			s.output.Put(name)
		}
	}
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

func (s *session) generate() {
	defer close(s.name)
	lucky, err := s.db.GetWuGeLucky(s.Context(), s.getLastStroke())
	if err != nil {
		log.Error("get wuge lucky", err)
		s.err = err
		return
	}
	log.Info("wuge lucky list", "size", len(lucky))
	var tmp *ent.WuGeLucky
	for i := range lucky {
		tmp = lucky[i]
		log.Info("current lukcy", "lucky", tmp)
		if s.filter.SexFilter(tmp) {
			continue
		}
		log.Info("current lukcy sex filtered", "lucky", tmp, "dayan", s.filter.DaYanFilter(tmp))
		if s.filter.DaYanFilter(tmp) {
			continue
		}
		log.Info("current lukcy dayan filterd", "lucky", tmp)
		if s.filter.WuXingFilter(tmp.TianGe, tmp.RenGe, tmp.DiGe) {
			continue
		}
		log.Info("current lukcy get chars", "lucky", tmp)
		var f1s []*ent.Character

		if cs, ok := s.chars[tmp.FirstStroke1]; !ok {
			f1s, err = s.db.GetCharacters(s.Context(), s.filter.StrokeFilter(tmp.FirstStroke1), s.filter.RegularFilter)
			if err != nil {
				log.Error("get first1 name", err)
				s.err = err
				return
			}
		} else {
			f1s = cs
		}

		var f2s []*ent.Character
		if cs, ok := s.chars[tmp.FirstStroke2]; !ok {
			f2s, err = s.db.GetCharacters(s.Context(), s.filter.StrokeFilter(tmp.FirstStroke2), s.filter.RegularFilter)
			if err != nil {
				log.Error("get first2 name", err)
				s.err = err
				return
			}
		} else {
			f2s = cs
		}

		//make first name
		for i1 := range f1s {
			for i2 := range f2s {
				select {
				case <-s.Context().Done():
					return
				default:
					s.name <- FirstName{
						f1s[i1],
						f2s[i2],
					}
				}
			}
		}
	}
	return
}
