package main

import (
	"fmt"
	"net/http"
	"swag_init/swagger"

	"github.com/gin-gonic/gin"
)

// @Summary 测试SayHello
// @Description 向你说Hello
// @Tags 测试
// @Accept json
// @Produce json
// @Param request body GreatingRequest true "asdasd"
// @Success 200 {object} GreatingResponse
// @Router /greating [post]
func SayHello(c *gin.Context) {
	var (
		request  GreatingRequest
		response GreatingResponse
		err      error
	)

	defer func() {
		c.JSON(http.StatusOK, response)
	}()

	if err = c.ShouldBindJSON(&request); err != nil {
		return
	}

	response = GreatingResponse{
		Greating: fmt.Sprintf("Hello %s ~", request.Name),
	}

}

type GreatingRequest struct {
	Name string `json:"name"`
}

type GreatingResponse struct {
	Greating string `json:"greating"`
}

func main() {
	app := gin.Default()

	swagger.NewGinSwaggerRouter(app, &swagger.Routeronfig{
		SearchDirs: []string{"."},
		APIFile:    "main.go",
		BasePath:   "/",
		Host:       "http://localhost:1323",
	})

	app.POST("/greating", SayHello)

	app.Run(":1323")
}
