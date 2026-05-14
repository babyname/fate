package repository

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
)

type Repository struct {
	*ent.Client
	cache *ModelCache
}

func (m *Repository) Initialize(ctx context.Context) error {
	return m.Schema.Create(ctx)
}

func (m *Repository) QueryLastName(ctx context.Context, last [2]string) (lastName [2]*ent.Character, err error) {
	lastName[0], err = m.Character.Query().Where(character.CharEQ(last[0])).First(ctx)
	if err != nil {
		return lastName, fmt.Errorf("query last name 0:%v", err)
	}
	if last[1] != "" {
		lastName[1], err = m.Character.Query().Where(character.CharEQ(last[1])).First(ctx)
		if err != nil {
			return lastName, fmt.Errorf("query last name 1:%v", err)
		}
	}
	return lastName, nil
}

func ID(name string) string {
	sum := md5.Sum([]byte(name))
	return hex.EncodeToString(sum[:])
}

func New(client *ent.Client) *Repository {
	Logger("Repository")
	return &Repository{Client: client, cache: NewModelCache()}
}
