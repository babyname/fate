package fate

import (
	"github.com/babyname/fate/config"
	"github.com/babyname/fate/internal/database"
	filterpkg "github.com/babyname/fate/internal/filter"
	"github.com/babyname/fate/internal/naming"
	"github.com/babyname/fate/internal/repository"
	"github.com/babyname/fate/internal/session"
)

type Fate interface {
	NewSession() Session
	NewSessionWithFilter(f Filter) Session
}

type fateImpl struct {
	cfg *config.Config
	db  *repository.Repository
}

func (f *fateImpl) NewSessionWithFilter(flt Filter) Session {
	return session.NewSession(f.db, flt)
}

func (f *fateImpl) NewSession() Session {
	return f.NewSessionWithFilter(filterpkg.DefaultFilter())
}

func New(cfg *config.Config) (Fate, error) {
	b := database.New(cfg.Database)
	client, err := b.Client()
	if err != nil {
		return nil, err
	}
	repo := repository.New(client)
	return &fateImpl{
		cfg: cfg,
		db:  repo,
	}, nil
}

func NewWithRepo(cfg *config.Config, repo *repository.Repository) Fate {
	return &fateImpl{
		cfg: cfg,
		db:  repo,
	}
}

type Session = session.Session
type Input = session.Input
type Output = session.Output
type SessionState = session.SessionState

type Filter = filterpkg.Filter
type FilterOption = filterpkg.FilterOption
type CharacterFilterType = filterpkg.CharacterFilterType

type Sex = naming.Sex
type Name = naming.Name
type NameBasic = naming.NameBasic
type FirstName = naming.FirstName

type ScoredName = session.ScoredName

const (
	SexBoy  Sex = naming.SexBoy
	SexGirl Sex = naming.SexGirl
)

var (
	NewFilter     = filterpkg.NewFilter
	DefaultFilter = filterpkg.DefaultFilter
)

const (
	SessionStateWaiting    SessionState = session.SessionStateWaiting
	SessionStateGenerating SessionState = session.SessionStateGenerating
	SessionStateFinish     SessionState = session.SessionStateFinish
	SessionStateCanceled   SessionState = session.SessionStateCanceled
	SessionStateFailed     SessionState = session.SessionStateFailed
)

const (
	CharacterFilterTypeDefault CharacterFilterType = filterpkg.CharacterFilterTypeDefault
	CharacterFilterTypeChs     CharacterFilterType = filterpkg.CharacterFilterTypeChs
	CharacterFilterTypeCht     CharacterFilterType = filterpkg.CharacterFilterTypeCht
	CharacterFilterTypeKangxi  CharacterFilterType = filterpkg.CharacterFilterTypeKangxi
)

var _ Fate = (*fateImpl)(nil)
