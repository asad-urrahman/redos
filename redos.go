package redos

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

const (
	goStandardRegexPkgName = "regexp"
)

// Command options
type Options struct {
	Verbose  bool
	FuzzFile string
	Timeout  int
	Regex    string
}

// Struct for for regex expresion
type regex struct {
	expression string
	pos        token.Pos
}

// ScanDir scan pass directory recursively for regex expresions
func ScanDir(dirName string, opts Options) {

	fset := token.NewFileSet()
	// parse directory
	pkgs, err := parser.ParseDir(fset, dirName, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	regExpresions := extractAllRegExpressions(pkgs)
	// TODO handle error
	_ = fuzzRegix(fset, regExpresions, opts)
}

// Extract all regex expresions
func extractAllRegExpressions(pkgs map[string]*ast.Package) []regex {
	var regExpresions []regex

	for _, pkg := range pkgs {
		for _, astFile := range pkg.Files {
			// Check if regex pkg is imported
			if !isRegexpImported(astFile) {
				continue
			}
			// Inspect
			ast.Inspect(astFile, func(node ast.Node) bool {
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
								return true
							}
						}
					}
				}
				return true
			})
		}
	}

	return regExpresions
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
