package main

import (
	"log"
	"net/http"
	"swag_init/swagger"
)

func main() {
	mux, err := swagger.NewSwaggerRouter(&swagger.Routeronfig{
		SearchDirs: []string{"."},
		APIFile:    "main.go",
		BasePath:   "/",
		Host:       "http://localhost:1323",
	})
	if err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(":1323", mux))

}
