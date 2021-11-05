package main

import (
	"fmt"
	"net/http"
	"swag_init/example/gin/models"
	"swag_init/swagger"

	"github.com/gin-gonic/gin"
)

type T struct{}

// @Summary 测试SayHello
// @Description 向你说Hello
// @Tags 测试
// @Accept json
// @Produce json
// @Param request body GreatingRequest true "asdasd"
// @Success 200 {object} GreatingResponse
// @Router /greating [post]
func (t *T) SayHello(ctx *gin.Context) {
	var (
		request  models.GreatingRequest
		response models.GreatingResponse
		err      error
	)

	defer func() {
		ctx.JSON(http.StatusOK, response)
	}()

	if err = ctx.ShouldBindJSON(&request); err != nil {
		return
	}

	response = models.GreatingResponse{
		Greating: fmt.Sprintf("Hello %s ~", request.Name),
	}

}

func main() {
	app := gin.Default()

	swagger.NewGinSwaggerRouter(app, &swagger.Routeronfig{
		SearchDirs: []string{"."},
		APIFile:    "main.go",
		BasePath:   "/",
		Host:       "http://localhost:1323",
	})

	t := T{}

	app.POST("/greating", t.SayHello)

	app.Run(":1323")
}
