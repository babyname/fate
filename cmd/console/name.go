package main

import (
	"fmt"
	"github.com/godcong/fate/config"
	"github.com/spf13/cobra"
)

func cmdName() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "name",
		Short: "output the name",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("config will output to %s/config.json\n", path)
			config.DefaultJSONPath = path
			cfg := config.LoadConfig()
			fmt.Printf("config loaded: %+v", cfg)
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", ".", "set the input path")
	return cmd
}
