package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"swag_init/parse"
	"swag_init/parser"
	"swag_init/swagger"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	t := T{}
	app.GET("/order/:id", t.test)

	infos := parse.GetGinRouteInfos(app)

	build, err := parse.ParseDir(".", "func(*gin.Context)", infos)
	if err != nil {
		panic(err)
	}

	p := parser.New()

	for _, file := range build.GetFiles() {
		if err := p.GinSwagger("", "", file); err != nil {
			panic(err)
		}
	}

	sg := p.GetSwagger()

	sg.BasePath = "/"
	swaggerURL, err := url.Parse("http://localhost:1323")
	if err != nil {
		panic(err)
	}

	sg.Host = swaggerURL.Host

	bytes, _ := json.MarshalIndent(sg, "", "    ")

	app.GET("/swagger/*any", swagger.GinHandler(
		func(c *swagger.Config) {
			c.JsonData = bytes
			c.URL = fmt.Sprintf("%s/swagger/doc.json", "http://localhost:1323")
			c.DeepLinking = true
			c.DomID = "#swagger-ui"
		},
	))

	app.Run(":1323")
}

type T struct{}

func (t T) test(c *gin.Context) {
	var err error
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{})
		return
	}
	c.ShouldBindJSON(&t)
	c.Query("1")
	c.QueryArray("2")
	c.QueryMap("3")
	c.BindQuery(&t)
	c.DefaultQuery("4", "1")
	c.GetQuery("5")
	c.GetQueryArray("6")
	c.GetQueryMap("7")
	c.ShouldBindQuery("8")

	fmt.Println(c.Param("id"))

	// c.JSONP(http.StatusOK, map[string]interface{}{})
	// c.XML(http.StatusOK, map[string]interface{}{})
	// c.YAML(http.StatusOK, map[string]interface{}{})
	// c.ProtoBuf(http.StatusOK, map[string]interface{}{})
	c.JSON(http.StatusOK, map[string]interface{}{
		"asd": "asd",
	})
}
