package main

import (
	"github.com/babyname/fate"
	"github.com/babyname/fate/config"
	"github.com/babyname/fate/regular"
)

func main() {
	c := config.LoadConfig()
	db := fate.InitDatabaseWithConfig(*c)
	r := regular.New(db)
	r.Run()
}
