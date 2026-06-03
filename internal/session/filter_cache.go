package session

import (
	"github.com/babyname/fate/v4/ent"
)

// PutFilter 名字过滤器，使用链表和字符索引实现名字的存储、过滤和移除。
type PutFilter struct {
	filterList map[string][]*Element[[2]*ent.Character]
	list       *List[[2]*ent.Character]
}

// NewPutFilter 创建一个新的名字过滤器实例。
func NewPutFilter() *PutFilter {
	return &PutFilter{
		filterList: make(map[string][]*Element[[2]*ent.Character], 128),
		list:       New[[2]*ent.Character](),
	}
}

// Put 将名字对添加到过滤器中，同时建立字符到链表节点的索引。
func (f *PutFilter) Put(names ...[2]*ent.Character) {
	for i := range names {
		e := f.list.PushBack(names[i])
		f.filterList[names[i][0].Char] = append(f.filterList[names[i][0].Char], e)
		f.filterList[names[i][1].Char] = append(f.filterList[names[i][1].Char], e)
	}
}

// Len 返回过滤器中名字的数量。
func (f *PutFilter) Len() int {
	return f.list.Len()
}

// ForEach 遍历过滤器中的所有名字对，对每个名字调用 fn。
func (f *PutFilter) ForEach(fn func([2]*ent.Character)) {
	for e := f.list.Front(); e != nil; e = e.Next() {
		fn(e.Value)
	}
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

// Filter 按指定字符过滤名字，移除并返回包含该字符的所有名字对。
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
