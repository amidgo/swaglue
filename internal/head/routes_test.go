package head_test

import (
	"bytes"
	_ "embed"
	"log"
	"testing"

	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/model"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed testdata/routes/@login@vk/get.yaml
	login_vk_get []byte
	//go:embed testdata/routes/@login@vk/post.yaml
	login_vk_post []byte
	//go:embed testdata/routes/user/@group@all/get.yaml
	group_all_get []byte
	//go:embed testdata/routes/user/@user@{id}/get.yaml
	user_id_get []byte

	//go:embed testdata/test_routes_with_empty_paths_expected.yaml
	expectedEmptyPaths []byte
	//go:embed testdata/test_routes_with_exists_paths_expect.yaml
	expectedExistsPaths []byte
)

func routes() []*model.Route {
	return []*model.Route{
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
					Content: bytes.NewReader(user_id_get),
				},
			},
		},
	}
}

func TestAppendRoutes_EmptyPaths(t *testing.T) {
	head, err := head.ParseHeadFromFile("testdata/routes_with_empty_paths.yaml")
	assert.NoError(t, err, "failed open routes_with_empty_paths.yaml")
	err = head.AppendRoutes(routes())
	assert.Nil(t, err, "failed append routes")
	buf := &bytes.Buffer{}
	err = head.SaveTo(buf)
	assert.NoError(t, err, "failed save file")
	log.Println(buf.String())
	assert.Equal(t, expectedEmptyPaths, buf.Bytes())
}

func TestAppendRoutes_ExistsPaths(t *testing.T) {
	head, err := head.ParseHeadFromFile("testdata/routes_with_exists_paths.yaml")
	assert.NoError(t, err, "failed open routes_with_exists_paths.yaml")
	err = head.AppendRoutes(routes())
	assert.Nil(t, err, "failed append routes")
	buf := &bytes.Buffer{}
	err = head.SaveTo(buf)
	assert.NoError(t, err, "failed save file")
	log.Println(buf.String())
	assert.Equal(t, expectedExistsPaths, buf.Bytes())
}
