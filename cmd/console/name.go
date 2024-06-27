package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/godcong/chronos"
	"github.com/goextension/log"
	"github.com/spf13/cobra"

	"github.com/babyname/fate"
	"github.com/babyname/fate/config"
)

func cmdName() *cobra.Command {
	path := ""
	born := ""
	last := ""
	output := ""
	sex := false
	cmd := &cobra.Command{
		Use:   "name",
		Short: "output the name",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("start", time.Now().String())
			config.DefaultJSONPath = path
			cfg := config.LoadConfig()

			if cfg == nil {
				log.Warnw("config file not found,use default config")
				cfg = config.DefaultConfig()
			}
			log.Infow("show config info", "data", *cfg)

			if output != "" {
				cfg.FileOutput.Path = filepath.Join(output, cfg.FileOutput.Path)
			}

			bornTime, e := time.Parse(chronos.DateFormat, born)
			if e != nil {
				log.Fatalw("parse born failed", "error", e)
			}
			f := fate.NewFate(last, bornTime, fate.ConfigOption(cfg), fate.SexOption(fate.Sex(sex)))

			e = f.MakeName(context.Background())
			if e != nil {
				log.Fatalw("makename failed", "error", e)
			}

			fmt.Println("end", time.Now().String())
		},
	}
	cmd.Flags().StringVarP(&last, "last", "l", "", "set lastname")
	cmd.Flags().StringVarP(&born, "born", "b", time.Now().Format(chronos.DateFormat), "set birth as 2016/01/02 15:04")
	cmd.Flags().StringVarP(&path, "path", "p", ".", "set the input config path")
	cmd.Flags().StringVarP(&output, "outout", "o", "", "set the output path")
	cmd.Flags().BoolVarP(&sex, "sex", "s", false, "set sex of the baby")
	return cmd
}

func getCurrentPath() string {
	getwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return getwd
}
