package cache

import (
	"sync"
	"time"

	"github.com/babyname/fate/ent"
)

type filterCache struct {
	sync.RWMutex
	done   chan struct{}
	filter *PutFilter

	count        int
	countElement *Element[[2]*ent.Character]
}

func (f *filterCache) SetFilter(filter *PutFilter) {
	f.Lock()
	defer f.Unlock()
	f.filter = filter
}

func (f *filterCache) Next() ([2]*ent.Character, bool) {
	f.RLock()
	f.RUnlock()
	if f.count >= f.filter.Len() {
		return [2]*ent.Character{}, false
	}

	if f.count == 0 {
		f.countElement = f.filter.front()
	}

	if f.count > 0 {
		if f.countElement == nil || f.countElement.Next() == nil {
			f.countElement = f.filter.getElement(f.count)
		} else {
			f.countElement = f.countElement.Next()
		}
	}
	f.count++
	return f.countElement.Value, true
}

func (f *filterCache) Count() int {
	f.RLock()
	defer f.RUnlock()
	return f.count
}

func (f *filterCache) SetCount(i int) {
	f.Lock()
	defer f.Unlock()
	f.count = i
	f.countElement = f.filter.getElement(i)
}

func (f *filterCache) Filter(s string) [][2]*ent.Character {
	f.Lock()
	defer f.Unlock()
	return f.filter.Filter(s)
}

func (f *filterCache) Reset() {
	f.Lock()
	defer f.Unlock()
	close(f.done)
	f.filter = NewPutFilter()
	f.done = make(chan struct{})
	go f.gc()
}

func (f *filterCache) Free() {
	f.Lock()
	defer f.Unlock()
	f.filter = NewPutFilter()
}

func (f *filterCache) gc() {
	t := time.NewTimer(10 * time.Minute)
	for {
		select {
		case <-f.done:
			return
		default:
			if !f.TryLock() {
				t.Reset(3 * time.Second)
				continue
			}
			f.filter.gc()
			f.Unlock()
			t.Reset(10 * time.Minute)
		}
	}
}

func (f *filterCache) Put(names ...[2]*ent.Character) {
	f.Lock()
	defer f.Unlock()
	f.filter.Put(names...)
}

func (f *filterCache) Len() int {
	f.RLock()
	defer f.RUnlock()
	return f.filter.Len()
}

func (f *filterCache) GetList(sta, limit int) [][2]*ent.Character {
	f.RLock()
	defer f.RUnlock()
	if sta >= f.filter.Len() || sta < 0 {
		return nil
	}
	if limit >= f.filter.Len() || limit <= 0 {
		return nil
	}
	var result [][2]*ent.Character
	for i, e := 0, f.filter.front(); e != nil; i, e = i+1, e.Next() {
		if i >= sta && limit > 0 {
			limit--
			result = append(result, e.Value)
			continue
		}
		break
	}

	return result
}

func (f *filterCache) GetOne(idx int) ([2]*ent.Character, bool) {
	f.RLock()
	defer f.RUnlock()
	if idx >= f.filter.Len() || idx < 0 {
		return [2]*ent.Character{}, false
	}
	ret := f.filter.getElement(idx)
	if ret == nil {
		return [2]*ent.Character{}, false
	}
	return ret.Value, true
}

// NewCache ...
func NewCache() FilterCache {
	return &filterCache{
		done:   make(chan struct{}),
		filter: NewPutFilter(),
	}
}

// NewCacheWithPut ...
func NewCacheWithPut(filter *PutFilter) FilterCache {
	return &filterCache{
		done:   make(chan struct{}),
		filter: filter,
	}
}

var _ Filter = (*filterCache)(nil)
