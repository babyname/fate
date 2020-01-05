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
			fmt.Printf("config will output to %s/config.json\n", path)
			if ext := filepath.Ext(path); ext != "" {
				log.Fatal("path cannot have a file name")
			}
			config.DefaultJSONPath = path

			e := config.OutputConfig(config.DefaultConfig())
			if e != nil {
				log.Fatal(e)
			}
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", ".", "set the output path")
	return cmd
}
