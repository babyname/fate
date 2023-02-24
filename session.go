package fate

import (
	"github.com/babyname/fate/model"
	"golang.org/x/net/context"
)

type session struct {
	db    *model.Model
	props *Property
	name  map[string]string
}

func (s *session) Start(ctx context.Context, input *Input) error {
	return nil
}
