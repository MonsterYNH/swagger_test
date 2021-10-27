package main

import "swag_init/test/models"

type Helloworld interface {
	// @Summary 测试SayHello
	// @Description 向你说Hello
	// @Tags 测试
	// @Accept json
	// @Produce json
	// @Param request body GreatingRequest true "asdasd"
	// @Success 200 {object} GreatingResponse
	// @Router /greating [post]
	Greating(name, asd int, request *models.GreatingRequest) (response, dd *models.GreatingResponse)

	// @Summary 测试SayHello
	// @Description 向你说Hello
	// @Tags 测试
	// @Accept json
	// @Produce json
	// @Param request body GreatingRequest true "asdasd"
	// @Success 200 {object} GreatingResponse
	// @Router /greating [post]
	Hello(name, asd int, request *GreatingRequest) (response, dd *models.GreatingResponse)
}

type Hahahahaha interface {
	// @Summary 测试SayHello
	// @Description 向你说Hello
	// @Tags 测试
	// @Accept json
	// @Produce json
	// @Param request body GreatingRequest true "asdasd"
	// @Success 200 {object} GreatingResponse
	// @Router /greating [post]
	Greating(name, asd int, request *models.GreatingRequest) (response, dd *models.GreatingResponse)

	// @Summary 测试SayHello
	// @Description 向你说Hello
	// @Tags 测试
	// @Accept json
	// @Produce json
	// @Param request body GreatingRequest true "asdasd"
	// @Success 200 {object} GreatingResponse
	// @Router /greating [post]
	Hello(name, asd int, request *GreatingRequest) (response, dd *models.GreatingResponse)
}

type GreatingRequest struct {
	Name string `json:"name"`
}

type GreatingResponse struct {
	Greating string `json:"greating"`
}
