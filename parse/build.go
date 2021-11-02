package parse

import (
	"go/ast"
)

type Build struct {
	files map[string]*File
}

type File struct {
	name   string
	source *ast.File

	Functions []*ast.FuncDecl
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
		name:      name,
		source:    source,
		Functions: make([]*ast.FuncDecl, 0),
	}
}
