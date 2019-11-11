package main

import (
	"fmt"
	"github.com/godcong/chronos"
	"github.com/godcong/fate"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

const programName = `fate`

func main() {

	app := cli.NewApp()
	app.Name = programName
	app.Usage = "CLI for IPFS Cluster"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "last, l",
			Value: "",
			Usage: "set the last name",
		},
		&cli.StringFlag{
			Name:  "born, b",
			Value: "",
			Usage: "set the born date Format(2006/01/02 03:04)",
		},
		&cli.StringFlag{
			Name:  "database",
			Usage: "set the database address",
			Value: `"localhost:3306", "root", "111111"`,
		},
	}
	var f fate.Fate
	app.Before = func(c *cli.Context) error {
		db := strings.Split(c.String("database"), ",")
		fmt.Println("database:", db)
		eng := fate.InitMysql(db[0], db[1], db[2])
		born := c.String("born")
		fmt.Println("born:", born)
		chr := chronos.New(born)
		//fate.DefaultStrokeMin = 3
		//fate.DefaultStrokeMax = 10
		//fate.HardMode = false
		last := c.String("last")
		fmt.Println("last", last)
		f = fate.NewFate(last, chr.Solar().Time(), fate.Database(eng), fate.BaGuaFilter(), fate.ZodiacFilter(), fate.SupplyFilter())
		return nil
	}
	app.Action = func(context *cli.Context) error {
		e := f.MakeName()
		if e != nil {
			return e
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
