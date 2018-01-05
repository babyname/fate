package fate

type FiveGrid struct {
	SkyGrid    int //天格
	LandGrid   int //地格
	PersonGrid int //人格
	OutGrid    int //外格
	AllGrid    int //总格
}

// NewFiveGrid 新建一个五格属性
// 五格是进行大衍计算的基础
func NewFiveGrid(sg, lg, pg int) FiveGrid {
	return FiveGrid{SkyGrid: sg, LandGrid: lg, PersonGrid: pg}
}

func outGrid() {

}

func allGrid() {

}
