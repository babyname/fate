package model

import (
	"context"

	"github.com/godcong/fate/ent"
)

func (m Model) GetCharacter(ctx context.Context, char string) (*ent.Character, error) {
	return m.Character.Query().First(ctx)
}
