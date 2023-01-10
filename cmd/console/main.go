package main

import (
	"fmt"
	"github.com/babyname/fate/config"
	"github.com/spf13/cobra"
)

const (
	programName = `fate`
	fateVersion = `4.0.0`
)

var (
	flagConfigPath = rootCmd.Flags().String("config", "", "set a config file path")
)

var rootCmd = &cobra.Command{
	Use:     programName,
	Short:   "run fate command to generate some baby name",
	Version: fateVersion,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("arguments [command] was not inputted")
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Fate running")

		c := config.LoadConfig(*flagConfigPath)
		fmt.Println(c)
	},
	DisableSuggestions:         false,
	SuggestionsMinimumDistance: 1,
}

func main() {
	rootCmd.AddCommand(cmdInit(), cmdName(), cmdCheck(), versionCMD())
	e := rootCmd.Execute()
	if e != nil {
		panic(e)
	}
}
