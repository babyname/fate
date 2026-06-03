package session

import (
	"sync"
	"time"

	"github.com/babyname/fate/v4/internal/chronosfate"
	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/internal/analysis"
	"github.com/babyname/fate/v4/internal/naming"
)

type Input struct {
	Last   [2]string
	Born   time.Time
	Sex    naming.Sex
	output *Output
}

func (i *Input) Output() *Output {
	if i.output == nil {
		b := naming.ParseNameBasicFromInput(i.Last, i.Born, i.Sex)
		i.output = &Output{
			basic: b,
		}
	}
	return i.output
}

type ScoredName struct {
	Name  naming.FirstName
	Score float64
	Grade string
}

type Output struct {
	basic          *naming.NameBasic
	fateData       *chronosfate.FateData
	topNames       []analysis.NameResult
	excellentTable *ExcellentTable
	charMap        map[string]*ent.Character
	mu             sync.RWMutex
}

func (o *Output) Basic() *naming.NameBasic {
	return o.basic
}

func (o *Output) SetLastName(ln [2]*ent.Character) {
	o.basic.LastName = ln
}

func (o *Output) SetFateData(fd *chronosfate.FateData) {
	o.fateData = fd
}

func (o *Output) FateData() *chronosfate.FateData {
	return o.fateData
}

func (o *Output) SetTopNames(names []analysis.NameResult) {
	o.mu.Lock()
	o.topNames = names
	o.mu.Unlock()
}

func (o *Output) TopNames() []analysis.NameResult {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.topNames
}

func (o *Output) SetExcellentTable(table *ExcellentTable) {
	o.mu.Lock()
	o.excellentTable = table
	o.mu.Unlock()
}

func (o *Output) ExcellentTable() *ExcellentTable {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.excellentTable
}

func (o *Output) SetCharMap(cm map[string]*ent.Character) {
	o.mu.Lock()
	o.charMap = cm
	o.mu.Unlock()
}

func (o *Output) CharMap() map[string]*ent.Character {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.charMap
}

func (o *Output) Total() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if o.excellentTable != nil {
		return o.excellentTable.Len()
	}
	return len(o.topNames)
}

func (o *Output) AllNames() []ScoredName {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if o.excellentTable == nil {
		return nil
	}
	entries := o.excellentTable.entries
	result := make([]ScoredName, 0, len(entries))
	for _, e := range entries {
		var c1, c2 *ent.Character
		if o.charMap != nil {
			c1 = o.charMap[e.Char1]
			c2 = o.charMap[e.Char2]
		}
		if c1 == nil || c2 == nil {
			continue
		}
		result = append(result, ScoredName{
			Name:  naming.FirstName{c1, c2},
			Score: e.Score,
			Grade: e.Grade,
		})
	}
	return result
}

func (o *Output) BuildNameResult(char1, char2 string) *analysis.NameResult {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if o.charMap == nil || o.fateData == nil || o.basic == nil {
		return nil
	}

	c1 := o.charMap[char1]
	c2 := o.charMap[char2]
	if c1 == nil || c2 == nil {
		return nil
	}

	surname := ""
	l1, l2 := 0, 0
	if o.basic.LastName[0] != nil {
		surname = o.basic.LastName[0].Char
		l1 = o.basic.LastName[0].ScienceStroke
	}
	if o.basic.LastName[1] != nil {
		surname += o.basic.LastName[1].Char
		l2 = o.basic.LastName[1].ScienceStroke
	}

	nr := analysis.BuildNameResult(0, surname, c1, c2, l1, l2, o.fateData)
	return &nr
}
