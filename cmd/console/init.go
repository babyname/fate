package main

import (
	"fmt"
	"github.com/godcong/fate/config"
	"github.com/spf13/cobra"
)

func cmdInit() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "output the init config",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("config will output to %s/config.json\n", path)
			e := config.OutputConfig(path, config.DefaultConfig())
			if e != nil {
				panic(e)
			}
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", ".", "set the output path")
	return cmd
}
