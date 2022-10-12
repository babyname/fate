package transfer

import (
	"context"
	"fmt"

	"github.com/babyname/fate/ent"
)

type Transfer interface {
	Start(ctx context.Context) error
}

type transferDatabase struct {
	Source *ent.Client
	Target *ent.Client
	Tables []string
}

func (t transferDatabase) Start(ctx context.Context) error {
	//todo
	return nil
}

func newTransfer(c *DatabaseConfig) (*transferDatabase, error) {
	source, err := c.Source.GetDatabase().BuildClient()
	if err != nil {
		return nil, fmt.Errorf("could not open source database: %v", err)
	}
	target, err := c.Target.GetDatabase().BuildClient()
	if err != nil {
		return nil, fmt.Errorf("could not open target database: %v", err)
	}
	return &transferDatabase{
		Tables: c.Tables,
		Source: source,
		Target: target,
	}, nil
}

func NewTransfer(config *DatabaseConfig) (Transfer, error) {
	return newTransfer(config)
}

var _ Transfer = (*transferDatabase)(nil)
