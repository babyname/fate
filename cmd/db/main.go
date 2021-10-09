package main

import (
	"github.com/spf13/cobra"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

const programName = `db`
const dbVersion = `0.0.1`

var path string

var rootCmd = &cobra.Command{
	Use:     programName,
	Short:   "call db command to transfer old database to new",
	Version: dbVersion,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runTransfer(path); err != nil {
			panic(err)
		}
	},
	DisableSuggestions:         false,
	SuggestionsMinimumDistance: 1,
}

func main() {
	rootCmd.Flags().StringVarP(&path, "path", "p", "db.config", "set db config path")
	e := rootCmd.Execute()
	if e != nil {
		panic(e)
	}
}
