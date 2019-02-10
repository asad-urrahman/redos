package main

import (
	"flag"
	"log"

	"github.com/asad-urrahman/redos"
)

func main() {
	// usage: redos [-v] ./code-dir
	// main operation modes
	verbose := flag.Bool("v", false, "Enable verbose output")
	flag.Parse()
	args := flag.Args()

	opts := redos.Options{
		Verbose: *verbose,
	}

	if len(args) == 0 {
		log.Fatalf("Directory path is missing")
	}

	redos.ScanDir(args[0], opts)

}
