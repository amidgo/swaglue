package routesiterationstep_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/amidgo/node/yaml"
	"github.com/amidgo/swaglue/internal/glue/routesiterationstep"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/route"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	//go:embed testdata/routes/routes_with_empty_paths_expected.yaml
	expectedEmptyPaths []byte
	//go:embed testdata/routes/routes_with_exists_paths_expect.yaml
	expectedExistsPaths []byte
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

func TestAppendRoutes_EmptyPaths(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/routes/routes_with_empty_paths.yaml", new(yaml.Decoder))
	require.NoError(t, err, "open routes_with_empty_paths.yaml")

	appender := routesiterationstep.New(hd, new(yaml.Decoder), new(yaml.Encoder), route.SliceSource(routes()))

	err = appender.AppendRoutes(routes())
	require.NoError(t, err, "append routes")

	buf := &bytes.Buffer{}
	err = hd.SaveTo(buf, &yaml.Encoder{Indent: 2})
	require.NoError(t, err, "save file")

	assert.Equal(t, string(expectedEmptyPaths), buf.String())
}

func TestAppendRoutes_ExistsPaths(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/routes/routes_with_exists_paths.yaml", new(yaml.Decoder))
	require.NoError(t, err, "open routes_with_exists_paths.yaml")

	appender := routesiterationstep.New(hd, new(yaml.Decoder), &yaml.Encoder{Indent: 2}, route.SliceSource(routes()))

	err = appender.AppendRoutes(routes())
	require.NoError(t, err, "append routes")

	buf := &bytes.Buffer{}
	err = hd.SaveTo(buf, &yaml.Encoder{Indent: 2})
	require.NoError(t, err, "save file")

	assert.Equal(t, string(expectedExistsPaths), buf.String())
}
