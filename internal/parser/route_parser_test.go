package parser_test

import (
	"bytes"
	_ "embed"
	"io"
	"log"
	"slices"
	"testing"

	"github.com/amidgo/swaglue/internal/model"
	"github.com/amidgo/swaglue/internal/parser"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed testdata/routes/@login@vk/get.yaml
	login_vk_get []byte
	//go:embed testdata/routes/@login@vk/post.yaml
	login_vk_post []byte

	//go:embed testdata/routes/user/@user@{id}/get.yaml
	user_get []byte

	//go:embed testdata/routes/user/@group@all/get.yaml
	group_all_get []byte
)

func TestRouteParser(t *testing.T) {
	const basePackage = "./testdata/routes"
	parser := parser.NewRouteParser(basePackage)
	err := parser.Parse()
	assert.Nil(t, err, "failed parse routes")

	expectedRoutes := []*model.Route{
		{
			Name: "/login/vk",
			Methods: []*model.RouteMethod{
				{
					Method:  "get",
					Content: bytes.NewReader(login_vk_get),
				},
				{
					Method:  "post",
					Content: bytes.NewReader(login_vk_post),
				},
			},
		},
		{
			Name: "/group/all",
			Methods: []*model.RouteMethod{
				{
					Method:  "get",
					Content: bytes.NewReader(group_all_get),
				},
			},
		},
		{
			Name: "/user/{id}",
			Methods: []*model.RouteMethod{
				{
					Method:  "get",
					Content: bytes.NewReader(user_get),
				},
			},
		},
	}
	parserRoutes := parser.Routes()

	assert.True(t, RoutesEqual(expectedRoutes, parserRoutes), "routes not equal")
}

func RoutesEqual(routes1, routes2 []*model.Route) bool {
	return slices.EqualFunc(routes1, routes2, RouteEqual)
}

func RouteEqual(route1, route2 *model.Route) bool {
	if route1.Name != route2.Name {
		return false
	}
	return slices.EqualFunc(route1.Methods, route2.Methods, RouteMethodEqual)
}

func RouteMethodEqual(routeMethod1, routeMethod2 *model.RouteMethod) bool {
	if routeMethod1.Method != routeMethod2.Method {
		return false
	}
	content1, err := io.ReadAll(routeMethod1.Content)
	if err != nil {
		log.Fatalf("failed read routemethod1 content, err: %s", err)
	}
	content2, err := io.ReadAll(routeMethod2.Content)
	if err != nil {
		log.Fatalf("failed read routemethod2 content, err: %s", err)
	}
	return slices.Equal(content1, content2)
}
