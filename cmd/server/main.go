package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/internal/api"
)

func main() {
	port := flag.Int("port", 8080, "server port")
	configPath := flag.String("config", "", "config file path")
	flag.Parse()

	var cfg *config.Config
	if *configPath != "" {
		var err error
		cfg, err = config.LoadConfig(*configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "load config: %v\n", err)
			os.Exit(1)
		}
	} else {
		cfg = config.DefaultConfig()
	}

	server, err := api.NewServer(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create server: %v\n", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Fate 起名系统 Web 服务启动于 http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
