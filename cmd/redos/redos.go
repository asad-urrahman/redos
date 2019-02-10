package main

import (
	"flag"
	"log"

	"github.com/asad-urrahman/redos"
)

func main() {
	// usage: redos ./code-dir
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		log.Fatalf("Directory path is missing")
	}

	redos.ScanDir(args[0])

}
