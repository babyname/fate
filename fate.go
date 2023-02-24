package fate

import (
	"strings"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/database"
	"github.com/babyname/fate/model"
	"golang.org/x/net/context"
)

type Session interface {
	Start(ctx context.Context, input *Input) error
}

// Fate ...
type Fate interface {
	NewSession(Properties *Property) Session
}

type fateImpl struct {
	cfg *config.Config
	db  *model.Model
}

func (f *fateImpl) NewSession(props *Property) Session {
	return &session{
		props: props,
		db:    f.db,
	}
}

func (f *fateImpl) Query() *model.Model {
	return f.db
}

func filterSex(lucky *WuGeLucky) bool {
	return lucky.ZongSex == true
}

func isLucky(s string) bool {
	if strings.Compare(s, "吉") == 0 || strings.Compare(s, "半吉") == 0 {
		return true
	}
	return false
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
