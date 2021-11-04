package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"swag_init/swagger"

	"github.com/gorilla/mux"
)

func SayHello(rw http.ResponseWriter, r *http.Request) {
	var (
		request  GreatingRequest
		response GreatingResponse
		err      error
	)

	defer func() {
		bytes, _ := json.Marshal(response)
		rw.Write(bytes)
		rw.WriteHeader(http.StatusOK)
	}()

	requestBytes, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err = json.Unmarshal(requestBytes, &request); err != nil {
		return
	}

	response = GreatingResponse{
		Greating: fmt.Sprintf("Hello %s ~!", request.Name),
	}
}

type GreatingRequest struct {
	Name string `json:"name"`
}

type GreatingResponse struct {
	Greating string `json:"greating"`
}

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

	m.HandleFunc("/greating", SayHello)

	log.Fatal(http.ListenAndServe(":1323", m))
}
