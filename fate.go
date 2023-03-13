package fate

import (
	"github.com/babyname/fate/config"
	"github.com/babyname/fate/database"
	"github.com/babyname/fate/model"
)

// Fate ...
type Fate interface {
	NewSession() Session
	NewSessionWithFilter(f Filter) Session
}

type fateImpl struct {
	cfg *config.Config
	db  *model.Model
}

func (f *fateImpl) NewSessionWithFilter(filter Filter) Session {
	return &session{
		filter: filter,
		db:     f.db,
	}
}

func (f *fateImpl) NewSession() Session {
	return f.NewSessionWithFilter(DefaultFilter())
}

// New creates a new instance of Fate
// @param *config.Config
// @return Fate
// @return error
func New(cfg *config.Config) (Fate, error) {
	c := database.New(cfg.Database)
	client, err := c.Client()
	if err != nil {
		return nil, err
	}
	return &fateImpl{
		cfg: cfg,
		db:  model.New(client),
	}, nil
}

var _ Fate = (*fateImpl)(nil)
