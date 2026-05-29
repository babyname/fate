package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/babyname/fate"
	"github.com/babyname/fate/config"
	fatehttp "github.com/babyname/fate/internal/http"
	"github.com/babyname/fate/resources"
)

func main() {
	addr := flag.String("addr", ":18080", "listen address")
	flag.Parse()

	cfg := config.DefaultConfig()
	f, err := fate.New(cfg)
	if err != nil {
		log.Fatal("failed to initialize fate: ", err)
	}

	handler := fatehttp.NewHandler(f, resources.StaticSub)

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
