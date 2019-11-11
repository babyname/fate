package main

import (
	"github.com/godcong/chronos"
	"github.com/godcong/fate"
	"github.com/urfave/cli/v2"
	"os"
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
		&cli.StringSliceFlag{
			Name:  "database, db",
			Usage: "set the database address",
			Value: cli.NewStringSlice("localhost:3306", "root", "111111"),
		},
	}
	var f fate.Fate
	app.Before = func(c *cli.Context) error {
		db := c.StringSlice("database")
		eng := fate.InitMysql(db[0], db[1], db[2])
		chronos := chronos.New(c.String("born"))
		//fate.DefaultStrokeMin = 3
		//fate.DefaultStrokeMax = 10
		//fate.HardMode = false
		f = fate.NewFate(c.String("last"), chronos.Solar().Time(), fate.Database(eng), fate.BaGuaFilter(), fate.ZodiacFilter(), fate.SupplyFilter())
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
