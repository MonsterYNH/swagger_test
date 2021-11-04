package parse

import (
	"fmt"
	"sort"
	"strings"
)

type SwaggerAble interface {
	SwaggerSummary() (string, error)
	SwaggerDescription() (string, error)
	SwaggerAccept() (string, error)
	SwaggerProduce() (string, error)
	SwaggerParams() (string, error)
	SwaggerSuccess() (string, error)
	// SwaggerFailures() (string, error)
}

type GinSwagger struct {
	*FunctionDesc
}

func (desc *GinSwagger) SwaggerSummary() string {
	// @Summary 测试SayHello
	return fmt.Sprintf("// @Summary %s", desc.Name)
}

func (desc *GinSwagger) SwaggerDescription() string {
	// @Description 向你说Hello
	return fmt.Sprintf("// @Description %s", desc.Name)
}

func (desc *GinSwagger) SwaggerAccept() string {
	// @Accept json
	return "// @Accept json"
}

func (desc *GinSwagger) SwaggerProduce() string {
	// @Produce json
	return "// @Produce json"
}

func (desc *GinSwagger) SwaggerParams() []string {
	params := []string{}

	selectors := []string{}
	for selector := range desc.Exprs {
		selectors = append(selectors, selector)
	}

	sort.Strings(selectors)

	for _, selector := range selectors {
		switch selector {
		case "*github.com/gin-gonic/gin.Context.ShouldBindJSON":
			arg := desc.Exprs[selector].Args[0]
			argType := strings.Split(arg.Type, "/")
			params = append(params, fmt.Sprintf("// @Params %s %s %s %t %s", arg.Name, "body", argType[len(argType)-1], true, arg.Name))
		}
	}

	return params
}

func (desc *GinSwagger) SwaggerSuccess() string {
	result := ""
	for selector, callExpr := range desc.Exprs {
		switch selector {
		case "*github.com/gin-gonic/gin.Context.JSON":
			// @Success 200 {object} GreatingResponse
			argType := strings.Split(callExpr.Args[1].Type, "/")
			result = fmt.Sprintf("// @Success 200 {object} %s", argType[len(argType)-1])
		}
	}
	return result
}
