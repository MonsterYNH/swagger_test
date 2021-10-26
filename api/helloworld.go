package api

import "fmt"

type HelloworldService interface {
	Greating(*GreatingRequest) *GreatingResponse
}

type HelloworldLogic struct{}

// @Summary 测试SayHello
// @Description 向你说Hello
// @Tags 测试
// @Accept json
// @Produce json
// @Param request body object GreatingRequest "wasd"
// @Success 200 {object} GreatingResponse
// @Router /greating [post]
func (logic *HelloworldLogic) Greating(request *GreatingRequest) *GreatingResponse {
	return &GreatingResponse{
		Greating: fmt.Sprintf("Hello %s ~", request.Name),
	}
}

type GreatingRequest struct {
	Name string `json:"name"`
}

type GreatingResponse struct {
	Greating string `json:"greating"`
	Error    error  `json:"-"`
}
