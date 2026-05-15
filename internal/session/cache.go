package session

import (
	"github.com/babyname/fate/ent"
)

// CacheFilter 定义名字缓存的过滤与访问接口。
type CacheFilter interface {
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

// FilterCache 定义带过滤器的名字缓存接口，支持设置过滤器并按字符过滤名字。
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
