package main

import (
	"fmt"
	"os"
	"time"

	"github.com/godcong/chronos"

	"github.com/spf13/cobra"
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
