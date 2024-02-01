package main

import (
	"context"
	"fmt"

	_ "github.com/sqlite3ent/sqlite3"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/database"
)

func main() {
	cli := database.New(config.DefaultConfig().Database)
	client, err := cli.Client()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("update table")
	err = client.Schema.Create(ctx)
	if err != nil {
		panic(err)
	}
	//total := 0
	//err = scripts.FindCharNeedFix(ctx, client, func(fix scripts.NeedFix) bool {
	//	total++
	//	fmt.Println("need fix", fix)
	//	return true
	//})
	//fmt.Println("all fixable count", total)

	//fmt.Println("update character total", total)

	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("char detail count", total)

	//total = 0

	//fmt.Println("kangxi char count", total)

}
