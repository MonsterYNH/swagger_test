package swagger

import (
	"encoding/json"
	"fmt"
	"net/url"
	"swag_init/parser"

	"github.com/gorilla/mux"
)

type Routeronfig struct {
	SearchDirs []string
	APIFile    string
	BasePath   string
	Host       string
}

func NewSwaggerRouter(r *mux.Router, config *Routeronfig) error {
	p := parser.New()

	p.PropNamingStrategy = "camelcase"
	// p.ParseDependency = true

	if err := p.ParseAPIMultiSearchDir(config.SearchDirs, config.APIFile, 100); err != nil {
		return err
	}
	swagger := p.GetSwagger()

	swagger.BasePath = config.BasePath
	swaggerURL, err := url.Parse(config.Host)
	if err != nil {
		return err
	}
	swagger.Host = swaggerURL.Host

	bytes, _ := json.MarshalIndent(swagger, "", "    ")

	r.PathPrefix("/swagger/").Handler(Handler(
		func(c *Config) {
			c.JsonData = bytes
			c.URL = fmt.Sprintf("%s/swagger/doc.json", config.Host)
			c.DeepLinking = true
			// c.DocExpansion = "none"
			c.DomID = "#swagger-ui"
		},
	))

	return nil
}
