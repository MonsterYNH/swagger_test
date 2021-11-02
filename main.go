package main

import (
	"log"
	"net/http"
	"swag_init/api"
	"swag_init/swagger"

	"github.com/gorilla/mux"
)

func main() {
	m := mux.NewRouter()
	if err := swagger.NewSwaggerRouter(m, &swagger.Routeronfig{
		SearchDirs: []string{"."},
		APIFile:    "main.go",
		BasePath:   "/",
		Host:       "http://localhost:1323",
	}); err != nil {
		panic(err)
	}

	m.HandleFunc("/greating", api.Greating)

	log.Fatal(http.ListenAndServe(":1323", m))

}
