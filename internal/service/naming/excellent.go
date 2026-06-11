// Package naming provides the name generation service pipeline.
package naming

import (
	"container/heap"
	"math/rand"
	"sort"
	"sync"
)

const excellentTableCapacity = 10000
const maxShownNames = 100

// ExcellentEntry represents a named candidate in the excellent table.
type ExcellentEntry struct {
	Char1     string  `json:"char1"`
	Char2     string  `json:"char2"`
	Score     float64 `json:"score"`
	Grade     string  `json:"grade"`
	WuXing1   string  `json:"wu_xing1"`
	WuXing2   string  `json:"wu_xing2"`
	HasPoetry bool    `json:"has_poetry"`
}

type excellentMinHeap []ExcellentEntry

func (h excellentMinHeap) Len() int           { return len(h) }
func (h excellentMinHeap) Less(i, j int) bool { return h[i].Score < h[j].Score }
func (h excellentMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *excellentMinHeap) Push(x any)        { *h = append(*h, x.(ExcellentEntry)) }
func (h *excellentMinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// ExcellentTable is a thread-safe, heap-backed collection of top-scored name candidates.
type ExcellentTable struct {
	mu      sync.RWMutex
	h       excellentMinHeap
	entries []ExcellentEntry
	index   map[string]int
	shown   map[string]bool
	seen    map[string]bool
	cap     int
}

// NewExcellentTable creates a new empty excellent table.
func NewExcellentTable() *ExcellentTable {
	return &ExcellentTable{
		h:     make(excellentMinHeap, 0, excellentTableCapacity),
		shown: make(map[string]bool),
		seen:  make(map[string]bool),
		cap:   excellentTableCapacity,
	}
}

// TryPush attempts to insert an entry into the table.
// Duplicate entries (same Char1+Char2) are silently ignored.
func (t *ExcellentTable) TryPush(entry ExcellentEntry) {
	t.mu.Lock()
	defer t.mu.Unlock()

	key := entry.Char1 + entry.Char2
	if t.seen[key] {
		return
	}
	t.seen[key] = true

	if len(t.h) < t.cap {
		heap.Push(&t.h, entry)
		return
	}

	if entry.Score > t.h[0].Score {
		heap.Pop(&t.h)
		heap.Push(&t.h, entry)
	}
}

// Finalize locks the table, sorts entries by score descending, and frees heap memory.
func (t *ExcellentTable) Finalize() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.entries = make([]ExcellentEntry, len(t.h))
	copy(t.entries, t.h)
	sort.Slice(t.entries, func(i, j int) bool {
		return t.entries[i].Score > t.entries[j].Score
	})
	t.index = make(map[string]int, len(t.entries))
	for i, e := range t.entries {
		t.index[e.Char1+e.Char2] = i
	}
	t.h = nil
	t.seen = nil
}

// HeapLen returns the current number of entries in the heap (before finalization).
func (t *ExcellentTable) HeapLen() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.h)
}

// TopN returns the top n entries by score.
func (t *ExcellentTable) TopN(n int) []ExcellentEntry {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if n > len(t.entries) {
		n = len(t.entries)
	}
	if n <= 0 {
		return nil
	}
	result := make([]ExcellentEntry, n)
	copy(result, t.entries[:n])
	return result
}

// Explore returns a randomized sample of unshown eligible entries.
func (t *ExcellentTable) Explore(count int, filter func(ExcellentEntry) bool) []ExcellentEntry {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.shown) >= maxShownNames {
		return nil
	}

	var eligible []int
	for i, e := range t.entries {
		key := e.Char1 + e.Char2
		if t.shown[key] {
			continue
		}
		if filter != nil && !filter(e) {
			continue
		}
		eligible = append(eligible, i)
	}

	rand.Shuffle(len(eligible), func(i, j int) {
		eligible[i], eligible[j] = eligible[j], eligible[i]
	})

	remaining := maxShownNames - len(t.shown)
	if count > remaining {
		count = remaining
	}
	if count > len(eligible) {
		count = len(eligible)
	}

	result := make([]ExcellentEntry, count)
	for i := 0; i < count; i++ {
		result[i] = t.entries[eligible[i]]
		t.shown[t.entries[eligible[i]].Char1+t.entries[eligible[i]].Char2] = true
	}
	return result
}

// Entries returns a slice of entries in [offset, offset+limit).
func (t *ExcellentTable) Entries(offset, limit int) []ExcellentEntry {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if offset >= len(t.entries) {
		return nil
	}
	end := offset + limit
	if end > len(t.entries) {
		end = len(t.entries)
	}
	result := make([]ExcellentEntry, end-offset)
	copy(result, t.entries[offset:end])
	return result
}

// MarkShown marks an entry as shown (for Explore tracking).
func (t *ExcellentTable) MarkShown(char1, char2 string) {
	t.mu.Lock()
	t.shown[char1+char2] = true
	t.mu.Unlock()
}

// AllEntries returns all finalized entries (before finalization this returns nil).
func (t *ExcellentTable) AllEntries() []ExcellentEntry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.entries
}

// Len returns the total number of finalized entries.
func (t *ExcellentTable) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.entries)
}

// ShownCount returns the number of entries marked as shown.
func (t *ExcellentTable) ShownCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.shown)
}

// FindEntry looks up an entry by its two characters.
func (t *ExcellentTable) FindEntry(char1, char2 string) *ExcellentEntry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if i, ok := t.index[char1+char2]; ok {
		return &t.entries[i]
	}
	return nil
}
