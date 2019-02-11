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
	fuzFile := flag.String("fuzfile", "", "Input file source for fuzzer")
	timeout := flag.Int("timeout", 5, "Timeout time in secods for regex fuzzer")
	flag.Parse()
	args := flag.Args()

	opts := redos.Options{
		Verbose:  *verbose,
		FuzzFile: *fuzFile,
		Timeout:  *timeout,
	}

	if len(args) == 0 {
		log.Fatalf("Directory path is missing")
	}

	redos.ScanDir(args[0], opts)

}
