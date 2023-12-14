package fate

import (
	"time"

	"github.com/babyname/fate/cache"
	"github.com/babyname/fate/ent"
)

type Sex int //girl:0,boy:1

type Input struct {
	Last   [2]string
	Born   time.Time
	Sex    Sex
	output *Output
}

func (i *Input) Output() *Output {
	if i.output == nil {
		b := parseNameBasicFromInput(i)
		i.output = &Output{
			basic: b,
			cache: cache.NewCache(),
			name:  make(chan FirstName, 128),
		}
	}
	return i.output
}

type Output struct {
	basic *BasicInfo
	cache cache.FilterCache
	name  chan FirstName
}

func (o *Output) Basic() *BasicInfo {
	return o.basic
}

func (o *Output) SetLastName(ln [2]*ent.Character) {
	o.basic.LastName = ln
}

func (o *Output) ResetNextName() {
	o.cache.SetCount(0)
}

func (o *Output) NextName() (Name, bool) {
	fn, ok := o.cache.Next()
	if ok {
		return Name{
			BasicInfo: o.basic,
			FirstName: fn,
		}, true
	}
	return Name{}, false
}

// Filter 过滤文字
func (o *Output) Filter(s string) int {
	return len(o.cache.Filter(s))
}

func (o *Output) Total() int {
	return o.cache.Len()
}

func (o *Output) SetCacheFilter(filterCache *cache.PutFilter) {
	o.cache.SetFilter(filterCache)
}
