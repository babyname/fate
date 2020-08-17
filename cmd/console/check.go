package main

import (
	"fmt"
	"github.com/goextension/log"
	"github.com/spf13/cobra"
)

func cmdCheck() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "check the env does correct",
		Run: func(cmd *cobra.Command, args []string) {
			if e := zoneCheck(); e != nil {
				log.Fatalw("zoneinfo check failed", "error", e)
			}
			fmt.Println("check all done!!!")
		},
	}
	return cmd
}
