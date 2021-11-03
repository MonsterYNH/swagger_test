package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// @Summary 测试SayHello
// @Description 向你说Hello
// @Tags 测试
// @Accept json
// @Produce json
// @Param request body []GreatingRequest true "asdasd"
// @Success 200 {object} GreatingResponse
// @Router /greating [post]
func Greating(rw http.ResponseWriter, r *http.Request) (int, string) {
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

	bytes, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err = json.Unmarshal(bytes, &request); err != nil {
		return 1, ""
	}

	response.Greating = fmt.Sprintf("Hello %s ~", request.Name)

	return 2, ""
}

type GreatingRequest struct {
	Name string `json:"name"`
}

type GreatingResponse struct {
	Greating string `json:"greating"`
}
