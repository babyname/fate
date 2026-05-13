package main

import (
	"fmt"
	"log"
	"os"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/internal/seeddb"
	"github.com/spf13/cobra"
)

var (
	seedDir    string
	dbConfig   config.DBConfig
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "seeddb",
	Short: "Database seed and migration tool for fate",
	Long: `A database-agnostic tool to export, seed, and manage fate data.

Supports SQLite3 and MySQL via ent ORM. Raw data is stored as JSON
(database-independent), and can be imported into any supported database.

Workflow:
  1. seeddb export   — Export old SQLite3 database to JSON (one-time)
  2. seeddb init     — Create schema + import JSON seed data into target DB
  3. seeddb report   — Generate data quality report from JSON files`,
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export old SQLite3 database to JSON files",
	Long: `Read from an old SQLite3 database and write to JSON files.
This is a one-time operation to extract data from the legacy database.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runExport(); err != nil {
			log.Fatalf("Export failed: %v", err)
		}
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize database with seed data",
	Long: `Create database schema and import JSON seed data.
Supports SQLite3 and MySQL drivers via ent ORM.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runInit(); err != nil {
			log.Fatalf("Init failed: %v", err)
		}
	},
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate data quality report",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runReport(); err != nil {
			log.Fatalf("Report failed: %v", err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&seedDir, "dir", "d", "data/seed", "Directory for JSON seed files")

	exportCmd.Flags().StringVarP(&dbConfig.Name, "input", "i", "fate_old.db", "Path to old SQLite3 database file")

	initCmd.Flags().StringVar(&dbConfig.Driver, "driver", "sqlite3", "Database driver: sqlite3|mysql")
	initCmd.Flags().StringVar(&dbConfig.Name, "name", "fate", "Database name (or file path for sqlite3)")
	initCmd.Flags().StringVar(&dbConfig.Host, "host", "127.0.0.1", "Database host (mysql)")
	initCmd.Flags().StringVar(&dbConfig.Port, "port", "3306", "Database port (mysql)")
	initCmd.Flags().StringVar(&dbConfig.User, "user", "root", "Database user (mysql)")
	initCmd.Flags().StringVar(&dbConfig.Pwd, "password", "", "Database password (mysql)")
	initCmd.Flags().StringVar(&dbConfig.DSN, "dsn", "", "Full DSN (overrides individual fields)")
	initCmd.Flags().StringVarP(&configPath, "config", "c", "", "YAML config file (overrides flags)")

	rootCmd.AddCommand(exportCmd, initCmd, reportCmd)
}

func resolveDBConfig() seeddb.DBConfig {
	if configPath != "" {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		return seeddb.DBConfig{
			Driver: cfg.Database.Driver,
			DSN:    cfg.Database.DSN,
			Host:   cfg.Database.Host,
			Port:   cfg.Database.Port,
			User:   cfg.Database.User,
			Pwd:    cfg.Database.Pwd,
			Name:   cfg.Database.Name,
		}
	}
	return seeddb.DBConfig{
		Driver: dbConfig.Driver,
		DSN:    dbConfig.DSN,
		Host:   dbConfig.Host,
		Port:   dbConfig.Port,
		User:   dbConfig.User,
		Pwd:    dbConfig.Pwd,
		Name:   dbConfig.Name,
	}
}

func runExport() error {
	if err := os.MkdirAll(seedDir, 0755); err != nil {
		return fmt.Errorf("create seed dir: %w", err)
	}
	exporter := seeddb.NewExporter(dbConfig.Name, seedDir)
	return exporter.Export()
}

func runInit() error {
	cfg := resolveDBConfig()
	if err := os.MkdirAll(seedDir, 0755); err != nil {
		return fmt.Errorf("create seed dir: %w", err)
	}
	importer := seeddb.NewImporter(seedDir, cfg)
	return importer.Import()
}

func runReport() error {
	reporter := seeddb.NewReporter(seedDir)
	return reporter.Generate()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
