package main

import (
	"fmt"
	"github.com/babyname/fate"
	"github.com/babyname/fate/database"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/model"
	"github.com/spf13/cobra"
)

func cmdInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize database",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("init running...")
			b := database.New(cfg.Database)
			cli, err := b.Client()
			if err != nil {
				fmt.Println("building database error:", err)
			}
			m := model.New(cli)
			wuge := make(chan *ent.WuGeLucky)
			go fate.InitWuGe(wuge)
			fmt.Println("database initializing...")
			err = m.Initialize(cmd.Context(), wuge)
			if err != nil {
				fmt.Println("initialize database error:", err)
			}
		},
	}
	return cmd
}
