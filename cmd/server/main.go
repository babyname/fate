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
	flag.Parse()

	cfg := config.DefaultConfig()
	f, err := fate.New(cfg)
	if err != nil {
		log.Fatal("failed to initialize fate: ", err)
	}

	handler := fatehttp.NewHandler(f, resources.StaticSub, f.Repo())

	fmt.Printf("Fate server listening on %s\n", *addr)
	fmt.Println("Endpoints:")
	fmt.Println("  POST /api/v1/generate     - generate names")
	fmt.Println("  GET  /api/v1/status       - task status")
	fmt.Println("  GET  /api/v1/result       - top names")
	fmt.Println("  GET  /api/v1/names        - paginated candidates")
	fmt.Println("  GET  /api/v1/explore      - explore more names")
	fmt.Println("  GET  /api/v1/name_detail  - name detail")
	fmt.Println("  GET  /api/v1/characters   - character lookup")
	fmt.Println("  GET  /                    - web UI (embedded)")

	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal("server error: ", err)
	}
}
