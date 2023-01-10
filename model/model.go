package model

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/babyname/fate/ent"
)

type Model struct {
	*ent.Client
}

func (m Model) Initialize() error {
	return nil
}

// ID ...
func ID(name string) string {
	sum := md5.Sum([]byte(name))
	return hex.EncodeToString(sum[:])
}

// New ...
// @param *ent.Client
// @return *Model
func New(client *ent.Client) *Model {
	return &Model{Client: client}
}
