package route_test

import (
	"bytes"
	"testing"

	"github.com/amidgo/node"
	"github.com/amidgo/node/yaml"
	"github.com/amidgo/swaglue/internal/route"
	"github.com/amidgo/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "embed"
)

var (
	//go:embed testdata/routes/@login@vk/get.yaml
	loginVKGet []byte
	//go:embed testdata/routes/@login@vk/post.yaml
	loginVKPost []byte
	//go:embed testdata/routes/user/@group@all/get.yaml
	groupAllGet []byte
	//go:embed testdata/routes/user/@user@{id}/get.yaml
	UserIDGet []byte
)

func routes() []route.Route {
	return []route.Route{
		{
			Name: "/login/vk",
			Methods: []route.Method{
				{
					Method:  "get",
					Content: bytes.NewReader(loginVKGet),
				},
			},
		},
		{
			Name: "/login/vk",
			Methods: []route.Method{
				{
					Method:  "post",
					Content: bytes.NewReader(loginVKPost),
				},
			},
		},
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
					Content: bytes.NewReader(UserIDGet),
				},
			},
		},
	}
}

type RemoveDuplicatesRouteSourceTest struct {
	CaseName       string
	Routes         []route.Route
	ExpectedRoutes []route.Route
	ExpectedErr    error
}

func (r *RemoveDuplicatesRouteSourceTest) Name() string {
	return r.CaseName
}

func (r *RemoveDuplicatesRouteSourceTest) Test(t *testing.T) {
	routeSource := route.NewRemoveDuplicatesSource(
		route.SliceSource(r.Routes),
	)

	routes, err := routeSource.Routes()
	require.ErrorIs(t, err, r.ExpectedErr)
	assert.Equal(t, r.ExpectedRoutes, routes)
}

func Test_RemoveDuplicatesRouteSource(t *testing.T) {
	tester.RunNamedTesters(t,
		&RemoveDuplicatesRouteSourceTest{
			CaseName: "routes",
			Routes:   routes(),
			ExpectedRoutes: []route.Route{
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
							Content: bytes.NewReader(UserIDGet),
						},
					},
				},
			},
		},
	)
}

type NodeRouteSourceTest struct {
	CaseName       string
	Node           node.Node
	Encoder        node.EncoderTo
	ExpectedRoutes []route.Route
	ExpectedErr    error
}

func (n *NodeRouteSourceTest) Name() string {
	return n.CaseName
}

func (n *NodeRouteSourceTest) Test(t *testing.T) {
	routeSource := route.NewNodeRouteSource(n.Node, n.Encoder)

	routes, err := routeSource.Routes()
	require.ErrorIs(t, err, n.ExpectedErr)
	assert.Equal(t, n.ExpectedRoutes, routes)
}

//go:embed testdata/noderoutesource/node.yaml
var routeSourceNode []byte

func Test_NodeRouteSource(t *testing.T) {
	nd, err := new(yaml.Decoder).DecodeFrom(bytes.NewReader(routeSourceNode))
	require.NoError(t, err)

	tester.RunNamedTesters(t,
		&NodeRouteSourceTest{
			CaseName: "node source",
			Node:     nd,
			Encoder:  &yaml.Encoder{Indent: 2},
			ExpectedRoutes: []route.Route{
				{
					Name: "/user/{id}",
					Methods: []route.Method{
						{
							Method:  "get",
							Content: bytes.NewBuffer(UserIDGet),
						},
					},
				},
				{
					Name: "/login/vk",
					Methods: []route.Method{
						{
							Method:  "get",
							Content: bytes.NewBuffer(loginVKGet),
						},
						{
							Method:  "post",
							Content: bytes.NewBuffer(loginVKPost),
						},
					},
				},
			},
		},
	)
}

type JoinRouteSourceTest struct {
	CaseName       string
	RouteSources   []route.Source
	ExpectedRoutes []route.Route
	ExpectedError  error
}

func (j *JoinRouteSourceTest) Name() string {
	return j.CaseName
}

func (j *JoinRouteSourceTest) Test(t *testing.T) {
	routeSource := route.NewJoinSource(j.RouteSources...)

	routes, err := routeSource.Routes()
	require.ErrorIs(t, err, j.ExpectedError)
	assert.Equal(t, j.ExpectedRoutes, routes)
}

func Test_JoinRouteSource(t *testing.T) {
	tester.RunNamedTesters(t,
		&JoinRouteSourceTest{
			CaseName: "join single slice",
			RouteSources: []route.Source{
				route.SliceSource(routes()),
			},
			ExpectedRoutes: routes(),
		},
		&JoinRouteSourceTest{
			CaseName: "join two slices",
			RouteSources: []route.Source{
				route.SliceSource(routes()),
				route.SliceSource(routes()),
			},
			ExpectedRoutes: append(routes(), routes()...),
		},
	)
}
