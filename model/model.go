package model

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/babyname/fate/ent"
)

type Model struct {
	*ent.Client
}

// ID ...
func ID(name string) string {
	sum := md5.Sum([]byte(name))
	return hex.EncodeToString(sum[:])
}
