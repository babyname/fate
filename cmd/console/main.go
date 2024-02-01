package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/babyname/fate"
	"github.com/babyname/fate/config"
	"github.com/babyname/fate/log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/sqlite3ent/sqlite3"
)

const (
	programName = `fate`

	// helpContent ...
	helpContent = "正在使用Fate生成姓名列表，如遇到问题请访问项目地址：https://github.com/babyname/fate获取帮助!"
)

var (
	flagConfigPath = ""
)

var (
	cfg *config.Config
)

var rootCmd = &cobra.Command{
	Use:     programName,
	Short:   "生成姓名列表",
	Version: fate.Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(helpContent)
		err := cmd.Help()
		if err != nil {
			return
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if flagConfigPath != "" {
			fmt.Println("Loading config file from: ", flagConfigPath)
		}
		cfg = config.LoadConfig(flagConfigPath)
		fmt.Printf("Config file: %+v\n", cfg)
		err := log.SetGlobalLogger(cfg.Log)
		if err != nil {
			return
		}
		log.Info("logging config file", "path", cfg.Log.Path)
	},
	DisableSuggestions: false,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
		HiddenDefaultCmd:    true,
	},
	SuggestionsMinimumDistance: 1,
}

func main() {
	rootCmd.PersistentFlags().StringVarP(&flagConfigPath, "config", "c", "", "set a config file path")

	rootCmd.AddCommand(cmdInit(), cmdName())
	e := rootCmd.Execute()
	if e != nil {
		panic(e)
	}
}
