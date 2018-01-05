package main

import (
	"flag"
	"log"
	"os"

	"github.com/godcong/fate/config"
)

func main() {
	log.Println(os.Args)
	log.Println("parsed? = ", flag.Parsed())
	log.Println(config.GetFlag("l"))
	log.Println(config.GetFlag("d"))
	log.Println(config.GetFlag("c"))
	flag.PrintDefaults()
}
