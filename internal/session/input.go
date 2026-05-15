package session

import (
	"time"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/internal/naming"
)

// Input 命名会话的输入参数，包含姓氏、出生时间和性别信息。
type Input struct {
	Last   [2]string
	Born   time.Time
	Sex    naming.Sex
	output *Output
}

// Output 根据输入参数生成命名输出，提供名字的获取、过滤和统计功能。
func (i *Input) Output() *Output {
	if i.output == nil {
		b := naming.ParseNameBasicFromInput(i.Last, i.Born, i.Sex)
		i.output = &Output{
			basic: b,
			cache: NewCache(),
			name:  make(chan naming.FirstName, 128),
		}
	}
	return i.output
}

// Output 命名会话的输出结果，封装名字缓存和基础信息。
type Output struct {
	basic *naming.NameBasic
	cache FilterCache
	name  chan naming.FirstName
}

// Basic 返回命名的基础信息。
func (o *Output) Basic() *naming.NameBasic {
	return o.basic
}

// SetLastName 设置姓氏字符。
func (o *Output) SetLastName(ln [2]*ent.Character) {
	o.basic.LastName = ln
}

// ResetNextName 重置名字游标，使下次获取从头开始。
func (o *Output) ResetNextName() {
	o.cache.SetCount(0)
}

// NextName 获取下一个名字，返回名字和是否还有更多结果。
func (o *Output) NextName() (naming.Name, bool) {
	fn, ok := o.cache.Next()
	if ok {
		return naming.Name{
			NameBasic: o.basic,
			FirstName: fn,
		}, true
	}
	return naming.Name{}, false
}

// Filter 按指定字符过滤名字，返回被过滤掉的名字数量。
func (o *Output) Filter(s string) int {
	return len(o.cache.Filter(s))
}

// Total 返回当前缓存中的名字总数。
func (o *Output) Total() int {
	return o.cache.Len()
}

// SetCacheFilter 设置输出的名字过滤器缓存。
func (o *Output) SetCacheFilter(filterCache *PutFilter) {
	o.cache.SetFilter(filterCache)
}
