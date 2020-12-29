package main

import (
	"github.com/godcong/fate"
	"github.com/godcong/fate/config"
	"github.com/godcong/fate/regular"
)

func main() {
	c := config.LoadConfig()
	db := fate.InitDatabaseWithConfig(*c)
	r := regular.New(db)
	r.Run()
}
