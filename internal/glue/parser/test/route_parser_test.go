package parser_test

import (
	"bytes"
	_ "embed"
	"io"
	"log"
	"slices"
	"testing"

	"github.com/amidgo/swaglue/internal/glue/parser"
	"github.com/amidgo/swaglue/internal/route"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/routes/@login@vk/get.yaml
	loginVKGet []byte
	//go:embed testdata/routes/@login@vk/post.yaml
	loginVKPost []byte

	//go:embed testdata/routes/user/@user@{id}/get.yaml
	userGet []byte

	//go:embed testdata/routes/user/@group@all/get.yaml
	groupAllGet []byte
)

func TestRouteParser(t *testing.T) {
	const basePackage = "./testdata/routes"

	parser := parser.NewRouteParser(basePackage)

	err := parser.Parse()
	require.NoError(t, err, "parse routes")

	expectedRoutes := []route.Route{
		{
			Name: "/login/vk",
			Methods: []route.Method{
				{
					Method:  "get",
					Content: bytes.NewReader(loginVKGet),
				},
				{
					Method:  "post",
					Content: bytes.NewReader(loginVKPost),
				},
			},
		},
		{
			Name: "/group/all",
			Methods: []route.Method{
				{
					Method:  "get",
					Content: bytes.NewReader(groupAllGet),
				},
			},
		},
		{
			Name: "/user/{id}",
			Methods: []route.Method{
				{
					Method:  "get",
					Content: bytes.NewReader(userGet),
				},
			},
		},
	}
	parserRoutes := parser.Routes()

	assert.True(t, RoutesEqual(expectedRoutes, parserRoutes), "routes not equal")
}

func RoutesEqual(routes1, routes2 []route.Route) bool {
	return slices.EqualFunc(routes1, routes2, RouteEqual)
}

func RouteEqual(route1, route2 route.Route) bool {
	if route1.Name != route2.Name {
		return false
	}

	return slices.EqualFunc(route1.Methods, route2.Methods, RouteMethodEqual)
}

func RouteMethodEqual(routeMethod1, routeMethod2 route.Method) bool {
	if routeMethod1.Method != routeMethod2.Method {
		return false
	}

	content1, err := io.ReadAll(routeMethod1.Content)
	if err != nil {
		log.Fatalf("read routemethod1 content, %s", err)
	}

	content2, err := io.ReadAll(routeMethod2.Content)
	if err != nil {
		log.Fatalf("read routemethod2 content, %s", err)
	}

	return slices.Equal(content1, content2)
}
