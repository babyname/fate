package main

import (
	"fmt"
	"os"

	"github.com/babyname/fate/v4"
	"github.com/babyname/fate/v4/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	_ "github.com/sqlite3ent/sqlite3"
)

const (
	programName = `fate`

	helpContent = "正在使用 fate 生成姓名列表，如遇到问题请访问项目地址：https://github.com/babyname/fate获取帮助!"
)

var (
	flagConfigPath  string
	flagDBDriver    string
	flagDBFile      string
	flagInitMode    string
	flagNoDownload  bool
	flagDownloadURL string
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
		var err error
		cfg, err = config.LoadConfig(flagConfigPath)
		if err != nil {
			fmt.Println("load config error:", err)
			return
		}

		if flagDBDriver != "" {
			cfg.Database.Driver = flagDBDriver
		}
		if flagDBFile != "" {
			cfg.Database.DBFile = flagDBFile
		}
		if flagInitMode != "" {
			cfg.Database.InitMode = flagInitMode
		}
		if flagNoDownload {
			cfg.Database.AutoDownload = false
		}
		if flagDownloadURL != "" {
			cfg.Database.DownloadURL = flagDownloadURL
		}

		fmt.Printf("Config file: %+v\n", cfg)
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
	rootCmd.PersistentFlags().StringVar(&flagDBDriver, "db-driver", "", "database driver: sqlite3 (default), mysql")
	rootCmd.PersistentFlags().StringVar(&flagDBFile, "db-file", "", "sqlite3 database file path (default: ./fate.db)")
	rootCmd.PersistentFlags().StringVar(&flagInitMode, "init-mode", "", "database init mode: auto (default), db, json")
	rootCmd.PersistentFlags().BoolVar(&flagNoDownload, "no-download", false, "disable auto-download of database")
	rootCmd.PersistentFlags().StringVar(&flagDownloadURL, "download-url", "", "custom database download URL")

	rootCmd.AddCommand(cmdInit(), cmdName())
	e := rootCmd.Execute()
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
