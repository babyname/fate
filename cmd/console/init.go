package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func cmdInit() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "output the init config",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("config will output to %s/config.json\n", path)
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", ".", "set the output path")
	return cmd
}
