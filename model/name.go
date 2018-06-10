package model

type Name struct {
	Base   `xorm:"extends"`
	Last1  string
	Last2  string
	First1 string
	First2 string
}
