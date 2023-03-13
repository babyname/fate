package model

import (
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/wugelucky"
	"golang.org/x/net/context"
)

const (
	WuGeLuckyMax = 81
	PerInitStep  = 1000
)

func WuGeLuckyID(l1, l2, f1, f2 int) int {
	var id int
	id = id | l1<<24
	id = id | l2<<16
	id = id | f1<<8
	id = id | f2
	return id
}

func (m Model) GetWuGeLucky(ctx context.Context, strokes [2]int) ([]*ent.WuGeLucky, error) {
	query := m.WuGeLucky.Query().Where(wugelucky.LastStroke1EQ(strokes[0])).
		Where(wugelucky.And(wugelucky.LastStroke2EQ(strokes[1]))).
		Where(wugelucky.And(wugelucky.ZongLuckyEQ(true)))
	return query.All(ctx)
}
