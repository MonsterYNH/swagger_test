package main

import (
	"fmt"
	"net/http"
	"swag_init/parse"

	"github.com/gin-gonic/gin"
)

func main() {
	_, err := parse.ParseDir("example/gin", "func(*gin.Context)")
	if err != nil {
		panic(err)
	}

	// build.Print(true)

	app := gin.Default()

	app.GET("/", test)

	for _, info := range app.Routes() {
		fmt.Println(info.Method, info.Path, info.Handler, info.HandlerFunc)
	}
}

func test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"as": "asd",
	})
}
