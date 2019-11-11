package main

import (
	"fmt"
	"github.com/godcong/chronos"
	"github.com/godcong/fate"
	"github.com/urfave/cli/v2"
	"os"
)

const programName = `fate`

func main() {

	app := cli.NewApp()
	app.Name = programName
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
			Value: "",
			Usage: "set the database address",
		},
	}
	var f fate.Fate
	app.Before = func(c *cli.Context) error {
		db := c.String("database")
		fmt.Println("database:", db)
		eng := fate.InitMysql(db, "root", "111111")
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
