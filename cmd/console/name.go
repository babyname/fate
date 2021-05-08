package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"
	_ "time/tzdata"

	"github.com/godcong/chronos"
	"github.com/godcong/fate"
	"github.com/godcong/fate/config"
	"github.com/godcong/yi"
	"github.com/goextension/log"
	"github.com/spf13/cobra"
)

func cmdName() *cobra.Command {
	cfg := ""
	born := ""
	last := ""
	output := ""
	sexStr := ""
	xiyong := ""
	cmd := &cobra.Command{
		Use:   "name",
		Short: "output the name",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("start", time.Now().String())
			config.DefaultJSONPath = cfg
			cfg := config.LoadConfig()

			if cfg == nil {
				panic("config file not found,use default config")
				//cfg = config.DefaultConfig()
			}
			log.Infow("show config info", "data", *cfg)

			if output != "" {
				cfg.FileOutput.Path = filepath.Join(output, cfg.FileOutput.Path)
			}

			bornTime, e := time.Parse(chronos.DateFormat, born)
			if e != nil {
				log.Fatalw("parse born failed", "error", e)
			}
			var sex yi.Sex
			if sexStr == "男" {
				sex = yi.SexBoy
			} else if sexStr == "女" {
				sex = yi.SexGirl
			} else {
				panic(fmt.Sprintf("性别不对劲哦:%s", sexStr))
			}
			f := fate.NewFate(last, bornTime, fate.ConfigOption(cfg), fate.SexOption(sex), fate.XiYongOption(xiyong))

			e = f.MakeName(context.Background())
			if e != nil {
				log.Fatalw("makename failed", "error", e)
			}

			fmt.Println("end", time.Now().String())
		},
	}
	cmd.Flags().StringVarP(&last, "last", "l", "", "set lastname")
	cmd.Flags().StringVarP(&born, "born", "b", time.Now().Format(chronos.DateFormat), "set birth format as 2016/01/02 15:04")
	cmd.Flags().StringVarP(&cfg, "cfg", "f", ".", "set the input config path")
	cmd.Flags().StringVarP(&output, "outout", "o", "", "set the output path")
	cmd.Flags().StringVarP(&sexStr, "sex", "s", "", "set sex of the baby")
	cmd.Flags().StringVarP(&xiyong, "xiyong", "x", "", "set xiyong of the baby")

	return cmd
}

func getCurrentPath() string {
	getwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return getwd
}
