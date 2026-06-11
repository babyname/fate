package fate

import (
	"context"

	"github.com/babyname/fate/v4/config"
	"github.com/babyname/fate/v4/internal/analysis"
	"github.com/babyname/fate/v4/internal/chronosfate"
	"github.com/babyname/fate/v4/internal/database"
	filterpkg "github.com/babyname/fate/v4/internal/filter"
	"github.com/babyname/fate/v4/internal/naming"
	"github.com/babyname/fate/v4/internal/repository"
	"github.com/babyname/fate/v4/internal/session"
)

type Fate interface {
	NewSession() Session
	NewSessionWithFilter(f Filter) Session
	Repo() *repository.Repository
	// Generate runs the full name generation pipeline synchronously.
	// This is a simpler alternative to the Session state machine.
	Generate(ctx context.Context, req GenerateRequest) (*GenerateResult, error)
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

func (f *fateImpl) Repo() *repository.Repository {
	return f.db
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
type StrokeMode = filterpkg.StrokeMode
type NameStyle = filterpkg.NameStyle

type Sex = naming.Sex
type Name = naming.Name
type NameBasic = naming.NameBasic
type FirstName = naming.FirstName

type ScoredName = session.ScoredName
type ExcellentEntry = session.ExcellentEntry
type ExcellentTable = session.ExcellentTable

type NameResult = analysis.NameResult
type ZhouYiBasic = analysis.ZhouYiBasic
type BaziBasic = analysis.BaziBasic
type WuXingSection = analysis.WuXingSection
type CharInfo = analysis.CharInfo
type GeItem = analysis.GeItem
type WuGeResult = analysis.WuGeResult
type ScoreDetail = analysis.ScoreDetail

// FateData is a type alias for the internal chronosfate.FateData.
type FateData = chronosfate.FateData

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
	CharacterFilterTypeDefault     CharacterFilterType = filterpkg.CharacterFilterTypeDefault
	CharacterFilterTypeChs         CharacterFilterType = filterpkg.CharacterFilterTypeChs
	CharacterFilterTypeCht         CharacterFilterType = filterpkg.CharacterFilterTypeCht
	CharacterFilterTypeKangxi      CharacterFilterType = filterpkg.CharacterFilterTypeKangxi
	CharacterFilterTypeNameScience CharacterFilterType = filterpkg.CharacterFilterTypeNameScience
	CharacterFilterTypeChsCht      CharacterFilterType = filterpkg.CharacterFilterTypeChsCht
	CharacterFilterTypeAll         CharacterFilterType = filterpkg.CharacterFilterTypeAll
)

const (
	StrokeModeScience     StrokeMode = filterpkg.StrokeModeScience
	StrokeModeSimplified  StrokeMode = filterpkg.StrokeModeSimplified
	StrokeModeTraditional StrokeMode = filterpkg.StrokeModeTraditional
	StrokeModeKangxi      StrokeMode = filterpkg.StrokeModeKangxi
)

const (
	NameStyleSimplified  NameStyle = filterpkg.NameStyleSimplified
	NameStyleTraditional NameStyle = filterpkg.NameStyleTraditional
	NameStyleAuto        NameStyle = filterpkg.NameStyleAuto
)

var _ Fate = (*fateImpl)(nil)
