package main

import (
	"github.com/urfave/cli/v2"
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
			Name:  "born",
			Value: "",
			Usage: "set the born date Format(2006/01/02 03:04)",
		},
	}

	app.Before = func(c *cli.Context) error {
		return nil
	}
}
