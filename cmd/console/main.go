package main

import (
	"github.com/godcong/fate"
)

func main() {
	name := fate.NewName("曹")
	fate.FilterBest(name)
}
