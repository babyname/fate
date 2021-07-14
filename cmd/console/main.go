package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const programName = `fate`
const fateVersion = `0.0.2`

var rootCmd = &cobra.Command{
	Use:     programName,
	Short:   "call fate command to make name",
	Version: fateVersion,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("arguments [command] was not inputted")
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
