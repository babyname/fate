package main

import (
	"fmt"
	"github.com/godcong/fate/config"
	"github.com/spf13/cobra"
	"log"
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
				log.Fatalf("wrong path with:%v", e)
			}
			fmt.Printf("config will output to %s\n", filepath.Join(absPath, config.JSONName))
			config.DefaultJSONPath = path

			e = config.OutputConfig(config.DefaultConfig())
			if e != nil {
				log.Fatal(e)
			}
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "set the output path")
	return cmd
}
