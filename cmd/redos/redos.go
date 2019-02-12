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
	regex := flag.String("regex", "", "Regex expresion")
	flag.Parse()
	args := flag.Args()

	opts := redos.Options{
		Verbose:  *verbose,
		FuzzFile: *fuzFile,
		Timeout:  *timeout,
		Regex:    *regex,
	}

	if len(args) == 0 {
		log.Fatalf("Directory path is missing")
	}

	if opts.Regex != "" {
		// / TODO handle error
		_ = redos.CompileRegix(opts.Regex, opts, "Regex: \""+opts.Regex+"\"")
		return
	}

	redos.ScanDir(args[0], opts)

}
