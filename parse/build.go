package parse

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
)

type Build struct {
	files map[string]*File
}

func (build Build) Print(format bool) {
	bytes, _ := json.MarshalIndent(build, "", "  ")
	fmt.Println(string(bytes))
}

type File struct {
	Name   string
	source *ast.File

	Functions []FunctionDesc
}

type FunctionDesc struct {
	source *ast.FuncDecl

	Name    string
	Params  []FuncItem
	Results []FuncItem

	Vars  map[string]FuncItem
	Exprs map[string]ExprItem
}

type FuncItem struct {
	Name string
	Type string
	Pos  string
}

type ExprItem struct {
	CallName string
	Args     []ExprArgItem
}

type ExprArgItem struct {
	Name string
	Pos  string
}

func NewBuild() *Build {
	return &Build{
		files: make(map[string]*File),
	}
}

func (build *Build) AddFile(name string, file *File) {
	build.files[name] = file
}

func NewFile(name string, source *ast.File) *File {
	return &File{
		Name:      name,
		source:    source,
		Functions: make([]FunctionDesc, 0),
	}
}

func parseFuncItem(fields *ast.FieldList, fset token.FileSet) []FuncItem {
	items := []FuncItem{}

	if fields == nil || fields.List == nil {
		return items
	}

	for _, field := range fields.List {
		for _, nameEntry := range field.Names {
			items = append(items, FuncItem{
				Name: nameEntry.Name,
				Type: ExprString(field.Type),
				Pos:  fset.Position(nameEntry.Pos()).String(),
			})
			// fmt.Println("params", nameEntry.Name, ExprString(field.Type))
		}
	}

	return items
}
