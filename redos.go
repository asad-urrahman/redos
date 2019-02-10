package redos

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

const (
	goStandardRegexPkgName = "regexp"
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
					arg := fcall.Args[0]
					fmt.Printf(" Reg exp : %v\n", arg)
				}
			}
		}
		return true
	}
	return true
}
