package main

import (
	"swag_init/parse"
)

func main() {
	build, err := parse.ParseDir("example/gin", "func(*gin.Context)")
	if err != nil {
		panic(err)
	}

	build.Print(true)
}
