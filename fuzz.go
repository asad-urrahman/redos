package redos

import (
	"bufio"
	"fmt"
	"go/token"
	"os"
	"regexp"
	"strings"
	"time"
)

var fuzzStrings = []string{
	// https://github.com/rust-lang/regex/blob/master/examples/shootout-regex-dna-cheat.rs
	"agggtaaa|tttaccct",
	"[cgt]gggtaaa|tttaccc[acg]",
	"a[act]ggtaaa|tttacc[agt]t",
	"ag[act]gtaaa|tttac[agt]ct",
	"agg[act]taaa|ttta[agt]cct",
	"aggg[acg]aaa|ttt[cgt]ccct",
	"agggt[cgt]aa|tt[acg]accct",
	"agggta[cgt]a|t[acg]taccct",
	"agggtaa[cgt]|[acg]ttaccct",
	// https://www.owasp.org/index.php/Regular_expression_Denial_of_Service_-_ReDoS
	// Examples of Evil Patterns:
	// (a+)+
	// ([a-zA-Z]+)*
	// (a|aa)+
	// (a|a?)+
	// (.*a){x} | for x > 10
	// ^([a-zA-Z0-9])(([\-.]|[_]+)?([a-zA-Z0-9]+))*(@){1}[a-z0-9]+[.]{1}(([a-z]{2,3})|([a-z]{2,3}[.]{1}[a-z]{2,3}))$
	// ^(([a-z])+.)+[A-Z]([a-z])+$
	// All the above are susceptible to the input aaaaaaaaaaaaaaaaaaaaaaaa!
	// (The minimum input length might change slightly, when using faster or slower machines).
	"aaaaaaaaaaaaaaaaaaaaaaaa!",
	"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa!",
}

func fuzzRegix(fset *token.FileSet, re []regex, opts Options) error {

	fuzSource, err := getFuzzSource(opts)
	if err != nil {
		return err
	}

	for _, r := range re {
		testRegex, err := regexp.Compile(r.expression)
		if err != nil {
			return err
		}
		fd := fset.File(r.pos)

		ch := make(chan bool, 1)
		defer close(ch)

		// start timer
		timer := time.NewTimer(time.Duration(opts.Timeout) * time.Second)
		defer timer.Stop()

		go func() {
			for fuzSource.Scan() {
				testRegex.FindAllSubmatch([]byte(fuzSource.Text()), -1)
			}
			ch <- true
		}()

		select {
		case <-ch:
			if opts.Verbose {
				fmt.Printf("GOOD REGEX at %v[%0.4d] Reg: %v \n", fd.Name(), fd.Line(r.pos), r.expression)
			}
			continue
		case <-timer.C:
			// Timeout
			fmt.Printf("EVIL REGEX at %v[%0.4d] Reg: %v \n", fd.Name(), fd.Line(r.pos), r.expression)
		}
	}

	return nil
}

func getFuzzSource(opts Options) (*bufio.Scanner, error) {
	var scanner *bufio.Scanner
	if opts.FuzzFile != "" { // File source
		file, err := os.Open(opts.FuzzFile)
		if err != nil {
			panic(err)
		}
		scanner = bufio.NewScanner(file)
		return scanner, nil
	}

	s := ""
	for _, v := range fuzzStrings {
		s += v
	}
	scanner = bufio.NewScanner(strings.NewReader(s))
	return scanner, nil
}
