package parse

import (
	"encoding/json"
	"fmt"
	"go/ast"
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

	Comments    []string
	Name        string
	PackageName string
	Params      []FuncItem
	Results     []FuncItem

	Vars  map[string]FuncItem
	Exprs map[string]ExprItem
}

type FuncItem struct {
	Name string
	Type string
}

type ExprItem struct {
	Receiver string
	Name     string
	Args     []ExprArgItem
}

type ExprArgItem struct {
	Name string
	Type string
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

func (build *Build) GetFiles() []*ast.File {
	files := []*ast.File{}

	for _, file := range build.files {
		files = append(files, file.source)
	}
	return files
}
