package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/babyname/fate/v4"
	"github.com/babyname/fate/v4/config"
	fatehttp "github.com/babyname/fate/v4/internal/http"
	"github.com/babyname/fate/v4/resources"
)

func main() {
	addr := flag.String("addr", ":18080", "listen address")
	dbDriver := flag.String("db-driver", "", "database driver: sqlite3 (default), mysql")
	dbFile := flag.String("db-file", "", "sqlite3 database file path (default: ./fate.db)")
	initMode := flag.String("init-mode", "", "database init mode: auto (default), db, json")
	noDownload := flag.Bool("no-download", false, "disable auto-download of database")
	downloadURL := flag.String("download-url", "", "custom database download URL")
	configPath := flag.String("config", "", "config file path")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	if *dbDriver != "" {
		cfg.Database.Driver = *dbDriver
	}
	if *dbFile != "" {
		cfg.Database.DBFile = *dbFile
	}
	if *initMode != "" {
		cfg.Database.InitMode = *initMode
	}
	if *noDownload {
		cfg.Database.AutoDownload = false
	}
	if *downloadURL != "" {
		cfg.Database.DownloadURL = *downloadURL
	}

	f, err := fate.New(cfg)
	if err != nil {
		log.Fatal("failed to initialize fate: ", err)
	}

	handler := fatehttp.NewHandler(f, resources.StaticSub, f.Repo())

	fmt.Printf("Fate server listening on %s\n", *addr)
	fmt.Println("Endpoints:")
	fmt.Println("  GET  /health           - health check")
	fmt.Println("  POST /api/generate     - generate names")
	fmt.Println("  GET  /api/name-detail  - get name detail")
	fmt.Println("  GET  /                 - web UI (embedded)")

	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal("server error: ", err)
	}
}
