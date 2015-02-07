package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func tester(i int) {}

var (
	filename = flag.String("f", "", "filename to search")
	dir      = flag.String("d", "", "directory to recursively search")
)

var usage = `tas - type assert searcher

tas is looking for type asserts that might cause panics
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	if *filename != "" && *dir != "" {
		fmt.Println("specify filename or directory")
		return
	}
	if *filename != "" {
		err := parseFile(*filename)
		if err != nil {
			panic(err)
		}
	}
	if *dir != "" {
		err := filepath.Walk(*dir, func(path string, info os.FileInfo, _ error) error {
			if filepath.Ext(info.Name()) == ".go" {
				err := parseFile(path)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
}

func parseFile(filename string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return err
	}

	//Testing hack - go run main -f main.go should find this type assert
	var test interface{} = 42
	tester(test.(int))

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.AssignStmt:
			return false
		case *ast.TypeAssertExpr:
			if x.Type != nil {
				position := fset.Position(x.Pos())
				fmt.Println(position) // Found some naughtyness
			}
		}
		return true
	})
	return nil
}
