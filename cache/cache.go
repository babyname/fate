package cache

import (
	"github.com/babyname/fate/ent"
)

type Filter interface {
	Put(...[2]*ent.Character)
	Next() ([2]*ent.Character, bool)
	Count() int
	SetCount(int)
	GetOne(idx int) ([2]*ent.Character, bool)
	GetList(sta, limit int) [][2]*ent.Character
	Filter(s string) [][2]*ent.Character
	Len() int
	Reset()
	Free()
}

type FilterCache interface {
	SetFilter(cache *PutFilter)
	Put(...[2]*ent.Character)
	Next() ([2]*ent.Character, bool)
	Count() int
	SetCount(int)
	GetOne(idx int) ([2]*ent.Character, bool)
	GetList(sta, limit int) [][2]*ent.Character
	Filter(s string) [][2]*ent.Character
	Len() int
	Reset()
	Free()
}
