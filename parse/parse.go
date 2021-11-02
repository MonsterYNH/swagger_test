package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

func ParseDir(path string) (*Build, error) {
	var fileSet token.FileSet

	packages, err := parser.ParseDir(&fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// create new build for the file set
	build := NewBuild()

	// iterate over all packages in the directory
	for _, pkg := range packages {
		// iterate over all files within the package
		for name, astTree := range pkg.Files {
			baseName := filepath.Base(name)

			fileAST, err := ParseFileAST(baseName, astTree)
			if err != nil {
				return nil, err
			}

			ast.Print(&fileSet, astTree)

			build.AddFile(baseName, fileAST)
		}
	}

	return build, nil
}

func ParseFileAST(name string, tree *ast.File) (*File, error) {
	file := NewFile(name, tree)

	fmt.Println(tree.Name.Name)
	for _, declaration := range tree.Decls {
		switch decValue := declaration.(type) {
		case *ast.FuncDecl:
			fmt.Println(decValue.Name.Name)
			for _, stmt := range decValue.Body.List {
				fmt.Println(stmt)
			}

		}
	}

	return file, nil
}
