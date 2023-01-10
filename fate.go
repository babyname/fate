package fate

import (
	"context"
	"fmt"
	"strings"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/model"
)

type Session interface{}

// Fate ...
type Fate interface {
	NewSession(Properties Property) Session
}

type fateImpl struct {
	cfg *config.Config
	db  *model.Model
}

func (f *fateImpl) NewSession(Properties Property) Session {
	//TODO implement me
	panic("implement me")
}

func (f *fateImpl) Query() *model.Model {
	return f.db
}

func (f *fateImpl) Initialize(ctx context.Context) error {
	panic("implement me")
}

func InitDayanLuckyTable(ctx context.Context, model *model.Model) error {
	err := model.Schema.Create(ctx)
	if err != nil {
		return err
	}
	lucky := make(chan *WuGeLucky)
	go initWuGe(lucky)
	for la := range lucky {
		//todo
		fmt.Println("la", la)
	}
	return nil
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

func New(cfg *config.Config) Fate {
	return &fateImpl{
		cfg: cfg,
	}
}

var _ Fate = (*fateImpl)(nil)
