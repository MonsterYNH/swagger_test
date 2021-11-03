package parse

import "fmt"

type CommentConfig struct {
	swaggerTag string
}

type SwaggerAble interface {
	SwaggerSummary() (string, error)
	SwaggerDescription() (string, error)
	SwaggerAccept() (string, error)
	SwaggerProduce() (string, error)
	SwaggerParams() (string, error)
	SwaggerSuccess() (string, error)
	SwaggerFailures() (string, error)
}

type GinSwagger struct {
	*FunctionDesc
}

func (desc *GinSwagger) SwaggerSummary() (string, error) {
	return desc.Name, nil
}

func (desc *GinSwagger) SwaggerDescription() (string, error) {
	return desc.Name, nil
}

func (desc *GinSwagger) SwaggerAccept() (string, error) {
	return "json", nil
}

func (desc *GinSwagger) SwaggerProduce() (string, error) {
	return "json", nil
}

func (desc *GinSwagger) SwaggerParams() ([]string, error) {
	params := []string{}
	for _, item := range desc.Exprs {
		switch item.CallName {
		case "c.ShouldBindJSON":
			if len(item.Args) != 1 {
				return nil, fmt.Errorf("[ERROR] wrong func format %s", desc.Name)
			}
			params = append(params, fmt.Sprintf("// @Params %s %s %s %t %s", desc.Name, "body", "models.GreatingRequest", true, item.Args[0].Name))
		}
	}

	return params, nil
}
