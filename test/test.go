package main

import (
	"swag_init/test/models"

	"github.com/gin-gonic/gin"
)

type Data struct{}

func SayHello(d int, c *gin.Context) {
	var data Data
	request := models.GreatingRequest{}
	c.BindHeader(&data)
	c.BindJSON(&request)
	c.BindQuery(&data)
	// c.BindUri(&data)
	// c.BindXML(&data)
	// c.BindYAML(&data)
	// // c.BindWith(&data, binding.Default())
	// c.ShouldBind(&data)
	// c.ShouldBindHeader(&data)
	// c.ShouldBindJSON(&data)
	// c.ShouldBindQuery(&data)
	// c.ShouldBindUri(&data)
	// c.ShouldBindXML(&data)
	// c.ShouldBindYAML(&data)

	// c.Query("name")
	// c.QueryArray("name")
	// c.QueryMap("name")

	// c.DefaultPostForm("name", "bob")
	// c.DefaultPostForm("name", "bob")

	// c.FormFile("name")
	// c.MultipartForm()

	// c.PostForm("name")
	// c.PostFormArray("name")
	// c.PostFormMap("name")

	// c.Param("name")
}
