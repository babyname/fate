package repository

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
	"github.com/google/uuid"
)

type Repository struct {
	*ent.Client
	cache *ModelCache
}

func (m *Repository) Initialize(ctx context.Context, luckies <-chan *ent.WuGeLucky) error {
	err := m.Schema.Create(ctx)
	if err != nil {
		return err
	}
	var tmp []*ent.WuGeLuckyCreate
	var count int
	for lucky := range luckies {
		_ = WuGeLuckyID(lucky.LastStroke1, lucky.LastStroke2, lucky.FirstStroke1, lucky.FirstStroke2)
		uid := uuid.Must(uuid.NewUUID())
		tmp = append(tmp, m.WuGeLucky.Create().SetWuGeLuckyWithOptional(lucky).SetID(uid))
		count++
		if len(tmp) >= PerInitStep {
			log.Info("insert into wugelucky", "count", count)
			_, err := m.insertWuGeLucky(ctx, tmp)
			if err != nil {
				return err
			}
			tmp = nil
		}
	}
	if len(tmp) != 0 {
		log.Info("insert into wugelucky", "count", count)
		_, err := m.insertWuGeLucky(ctx, tmp)
		if err != nil {
			return err
		}
		tmp = nil
	}
	return nil
}

func (m *Repository) insertWuGeLucky(ctx context.Context, tmp []*ent.WuGeLuckyCreate) ([]*ent.WuGeLucky, error) {
	return m.WuGeLucky.CreateBulk(tmp...).Save(ctx)
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
