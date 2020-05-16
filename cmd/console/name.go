package main

import (
	"context"
	"fmt"
	"github.com/godcong/chronos"
	"github.com/godcong/fate"
	"github.com/godcong/fate/config"
	"github.com/goextension/log"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func cmdName() *cobra.Command {
	path := ""
	born := ""
	last := ""
	sex := false
	cmd := &cobra.Command{
		Use:   "name",
		Short: "output the name",
		Run: func(cmd *cobra.Command, args []string) {
			os.Setenv("ZONEINFO", "zoneinfo.zip")
			fmt.Println("start", time.Now().String())
			config.DefaultJSONPath = path
			cfg := config.LoadConfig()
			fmt.Printf("config loaded: %+v", cfg)
			bornTime, e := time.Parse(chronos.DateFormat, born)
			if e != nil {
				log.Fatalw("parseborn", "error", e)
			}
			f := fate.NewFate(last, bornTime, fate.ConfigOption(cfg), fate.SexOption(fate.Sex(sex)))

			e = f.MakeName(context.Background())
			if e != nil {
				log.Fatalw("makename", "error", e)
			}

			fmt.Println("end", time.Now().String())
		},
	}
	cmd.Flags().StringVarP(&last, "last", "l", "", "set lastname")
	cmd.Flags().StringVarP(&born, "born", "b", time.Now().Format(chronos.DateFormat), "set birth as 2016/01/02 15:04")
	cmd.Flags().StringVarP(&path, "path", "p", ".", "set the input path")
	cmd.Flags().BoolVarP(&sex, "sex", "s", false, "set sex of the baby")
	return cmd
}
