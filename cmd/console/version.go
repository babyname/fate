package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/babyname/fate"
)

func versionCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print current version to screen",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("version:", fate.Version)
		},
	}
	return cmd
}
