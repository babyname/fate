package main

import (
	"fmt"
	"github.com/godcong/fate/config"
	"github.com/goextension/log"
	"github.com/spf13/cobra"
	"path/filepath"
)

func cmdInit() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "output the init config",
		Run: func(cmd *cobra.Command, args []string) {
			absPath, e := filepath.Abs(path)
			if e != nil {
				log.Fatalw("wrong path", "error", e, "path", path)
			}
			fmt.Printf("config will output to %s\n", filepath.Join(absPath, config.JSONName))
			config.DefaultJSONPath = path

			e = config.OutputConfig(config.DefaultConfig())
			if e != nil {
				log.Fatalw("config wrong", "error", e)
			}

		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "set the output path")
	return cmd
}
