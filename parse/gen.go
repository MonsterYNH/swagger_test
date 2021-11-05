package parse

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	commentRegExp = regexp.MustCompile("@(Summary|Description|Accept|Produce|Params|Success)")
)

type GinSwagger struct {
	*FunctionDesc
	comments []string
}

func (desc *GinSwagger) genSwaggerComment(routeInfos map[string]gin.RouteInfo) {

	desc.comments = append(desc.comments, desc.SwaggerSummary())
	desc.comments = append(desc.comments, desc.SwaggerAccept())
	desc.comments = append(desc.comments, desc.SwaggerProduce())
	desc.comments = append(desc.comments, desc.SwaggerSuccess())
	desc.comments = append(desc.comments, desc.SwaggerParams()...)
	for key := range routeInfos {
		if strings.HasPrefix(key, desc.PackageName) {
			desc.comments = append(desc.comments, desc.SwaggerRoute(routeInfos[key]))
		}
	}
}

func (desc *GinSwagger) mergeComments(comments []string) []string {
	commentMap := map[string][]string{}
	for _, comment := range desc.comments {
		if _, exist := commentMap[comment]; exist {
			log.Println("[WARNING] comment %s already exist")
			continue
		}
		commentMap[comment] = struct{}{}
	}

	newComments := make([]string, 0)
	for _, comment := range comments {
		if strings.HasPrefix("// @Summary")
	}
}

func (desc *GinSwagger) SwaggerSummary() string {
	// @Summary 测试SayHello
	return fmt.Sprintf("// @Summary %s", desc.Name)
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

func (desc *GinSwagger) SwaggerRoute(routeInfo gin.RouteInfo) string {
	bytes, _ := json.Marshal(routeInfo)
	fmt.Println(string(bytes))
	// @Router /greating [post]
	return fmt.Sprintf("// @Router %s [%s]", routeInfo.Path, strings.ToLower(routeInfo.Method))
}
