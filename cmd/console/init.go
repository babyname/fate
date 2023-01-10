package main

import (
	"fmt"
	"github.com/babyname/fate/database"
	"github.com/babyname/fate/model"
	"github.com/spf13/cobra"
)

func cmdInit() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "output the init config",
		Run: func(cmd *cobra.Command, args []string) {
			b := database.New(cfg.Database)
			cli, err := b.Client()
			if err != nil {
				fmt.Println("building database error:", err)
			}
			m := model.New(cli)
			m.Initialize()
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "set the output path")
	return cmd
}
