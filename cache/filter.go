package cache

import (
	"github.com/babyname/fate/ent"
)

type PutFilter struct {
	filterList map[string][]*Element[[2]*ent.Character]
	list       *List[[2]*ent.Character]
}

func NewPutFilter() *PutFilter {
	return &PutFilter{
		filterList: make(map[string][]*Element[[2]*ent.Character], 128),
		list:       New[[2]*ent.Character](),
	}
}

func (f *PutFilter) Put(names ...[2]*ent.Character) {
	for i := range names {
		e := f.list.PushBack(names[i])
		f.filterList[names[i][0].Ch] = append(f.filterList[names[i][0].Ch], e)
		f.filterList[names[i][1].Ch] = append(f.filterList[names[i][1].Ch], e)
	}
}

func (f *PutFilter) Len() int {
	return f.list.Len()
}

func (f *PutFilter) front() *Element[[2]*ent.Character] {
	return f.list.Front()
}

func (f *PutFilter) getElement(idx int) *Element[[2]*ent.Character] {
	if idx >= f.list.Len() || idx < 0 {
		return nil
	}
	for i, e := 0, f.list.Front(); e != nil; i, e = i+1, e.Next() {
		if i >= idx {
			return e
		}
		continue
	}
	return nil
}

func (f *PutFilter) gc() {
	if f.list.Len() != 0 {
		for s, i := range f.filterList {
			if v := removeDeletedElement(i); len(v) != 0 {
				f.filterList[s] = v
			} else {
				delete(f.filterList, s)
			}
		}
	} else {
		f.filterList = make(map[string][]*Element[[2]*ent.Character], 128)
	}
}

func (f *PutFilter) Filter(s string) [][2]*ent.Character {
	var filtered [][2]*ent.Character
	if v, ok := f.filterList[s]; ok {
		for _, e := range v {
			if e.Prev() == nil && e.Next() == nil {
				continue
			}

			filtered = append(filtered, f.list.Remove(e))
		}
	}
	return filtered
}

func removeDeletedElement(es []*Element[[2]*ent.Character]) []*Element[[2]*ent.Character] {
	ret := make([]*Element[[2]*ent.Character], len(es))
	count := 0
	for i, e := range es {
		if e.Prev() == nil && e.Next() == nil {
			continue
		}
		ret[count] = es[i]
		count++
	}
	return ret[:count]
}
