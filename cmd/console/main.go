package main

import "fmt"
import "github.com/spf13/cobra"

const programName = `fate`

var rootCmd = &cobra.Command{
	Use:     "fate [command]",
	Short:   "call fate command",
	Version: "v0.0.1",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("arguments [command] was not inputted")
	},
	DisableSuggestions:         false,
	SuggestionsMinimumDistance: 1,
}

func main() {
	rootCmd.AddCommand(cmdInit())
	e := rootCmd.Execute()
	if e != nil {
		panic(e)
	}
}
