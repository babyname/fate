package main

import (
	"github.com/spf13/cobra"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	transfer2 "github.com/babyname/fate/tools/transfer"
)

const programName = `transfer`
const version = `0.0.1`

var path string

var rootCmd = &cobra.Command{
	Use:     programName,
	Short:   "call command for transfer database to other database",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := transfer2.ReadTransferConfig(path)
		if err != nil {
			panic(err)
		}
		t, err := transfer2.NewTransfer(config)
		if err != nil {
			panic(err)
		}
		err = t.Start(cmd.Context())
		if err != nil {
			panic(err)
		}
	},
	DisableSuggestions:         false,
	SuggestionsMinimumDistance: 1,
}

func main() {
	rootCmd.Flags().StringVarP(&path, "path", "p", "transfer.cfg", "set database configuration file path")
	e := rootCmd.Execute()
	if e != nil {
		panic(e)
	}
}
