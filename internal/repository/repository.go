package repository

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
)

// Repository 封装 ent 客户端并提供字符数据的增删改查能力。
type Repository struct {
	*ent.Client
	cache *ModelCache
}

// Initialize 初始化数据库表结构。
func (m *Repository) Initialize(ctx context.Context) error {
	return m.Schema.Create(ctx)
}

// QueryLastName 根据姓氏字符串查询对应的字符记录。
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

// ID 根据名称生成 MD5 哈希作为唯一标识。
func ID(name string) string {
	sum := md5.Sum([]byte(name))
	return hex.EncodeToString(sum[:])
}

// New 创建 Repository 实例并初始化缓存。
func New(client *ent.Client) *Repository {
	Logger("Repository")
	return &Repository{Client: client, cache: NewModelCache()}
}
