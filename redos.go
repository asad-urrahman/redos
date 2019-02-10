package redos

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"regexp"
	"strings"
	"time"
)

const (
	goStandardRegexPkgName = "regexp"
)

// Struct for for regex expresion
type regex struct {
	expression string
	pos        token.Pos
}

var (
	regExpresions []regex
)

// ScanDir scan pass directory recursively for regex expresions
func ScanDir(dirName string) {
	// parse directory
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirName, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	// Extract all regex expresions
	for _, pkg := range pkgs {
		for _, astFile := range pkg.Files {
			// Check if regex pkg is imported
			if !isRegexpImported(astFile) {
				continue
			}
			// Inspect
			ast.Inspect(astFile, extractRegexExpression)
		}
	}

	fuzzRegix(regExpresions)
}

func isRegexpImported(file *ast.File) bool {
	for _, pkg := range file.Imports {
		c := strings.Trim(pkg.Path.Value, "\"")
		if c == goStandardRegexPkgName {
			return true
		}
	}
	return false
}

func extractRegexExpression(node ast.Node) bool {
	// Find Functions call
	fcall, ok := node.(*ast.CallExpr)
	if ok {
		// Retrive called packge name
		fcalPkg, ok := fcall.Fun.(*ast.SelectorExpr).X.(*ast.Ident)
		if ok {
			// if called packge is regex package
			// TODO also check func name
			if fcalPkg.Name == goStandardRegexPkgName {
				// get only first argument "patern"
				if len(fcall.Args) > 0 {
					re := regex{
						expression: fcall.Args[0].(*ast.BasicLit).Value,
						pos:        fcall.Pos(),
					}
					regExpresions = append(regExpresions, re)
				}
			}
		}
		return true
	}
	return true

}

func fuzzRegix(re []regex) error {

	for _, r := range re {
		testRegex, err := regexp.Compile(r.expression)
		if err != nil {
			return err
		}

		ch := make(chan bool, 1)
		defer close(ch)

		// start timer
		timer := time.NewTimer(5 * time.Second)
		defer timer.Stop()

		go func() {
			fuzzString := "aaaaaaaaaaaaaaaaaaaaaaaa!"
			testRegex.FindAllSubmatch([]byte(fuzzString), -1)
			ch <- true
		}()

		select {
		case <-ch:
			// TODO If verbose, Print info
			continue
		case <-timer.C:
			// Timeout
			fmt.Printf("EVIL REGEX: %v \n", r.expression)
		}
	}

	return nil
}
