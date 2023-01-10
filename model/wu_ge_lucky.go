package model

const (
	WuGeLuckyMax = 81
	PerInitStep  = 10000
)

func WuGeLuckyID(l1, l2, f1, f2 int) int {
	var id int
	id = id | l1<<24
	id = id | l2<<16
	id = id | f1<<8
	id = id | f2
	return id
}
