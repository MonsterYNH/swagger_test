package swagger

import (
	"encoding/json"
	"fmt"
	"swag_init/parser"

	"github.com/gorilla/mux"
)

type Routeronfig struct {
	SearchDirs []string
	APIFile    string
	BasePath   string
	Host       string
}

func NewSwaggerRouter(config *Routeronfig) (*mux.Router, error) {
	p := parser.New()

	p.PropNamingStrategy = "camelcase"

	if err := p.ParseAPIMultiSearchDir(config.SearchDirs, config.APIFile, 100); err != nil {
		return nil, err
	}
	swagger := p.GetSwagger()

	swagger.BasePath = config.BasePath
	swagger.Host = config.Host

	bytes, _ := json.MarshalIndent(swagger, "", "    ")

	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(Handler(
		func(c *Config) {
			c.JsonData = bytes
			c.URL = fmt.Sprintf("%s/swagger/doc.json", config.Host)
			c.DeepLinking = true
			// c.DocExpansion = "none"
			c.DomID = "#swagger-ui"
		},
	))

	return r, nil
}
