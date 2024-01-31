package head_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/amidgo/swaglue/internal/glue/head"
	"github.com/amidgo/swaglue/internal/glue/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/tags/name.yaml
	nameTagData []byte
	//go:embed testdata/tags/users.yaml
	usersTagData []byte
	//go:embed testdata/tags/haters.yaml
	hatersTagData []byte

	//go:embed testdata/tags/tags_expected_swagger.yaml
	tagExpectedSwagger []byte
)

func Test_AppendTags_Success(t *testing.T) {
	head, err := head.ParseHeadFromFile("testdata/swagger.yaml")
	require.NoError(t, err, "failed open swagger.yaml")

	err = head.AppendTags([]*model.Item{
		{
			Name:    "haters",
			Content: bytes.NewReader(hatersTagData),
		},
		{
			Name:    "name",
			Content: bytes.NewReader(nameTagData),
		},
		{
			Name:    "users",
			Content: bytes.NewReader(usersTagData),
		},
	})
	require.NoError(t, err, "failed append components")

	buf := &bytes.Buffer{}
	err = head.SaveTo(buf)
	require.NoError(t, err, "failed save file")

	assert.Equal(t, tagExpectedSwagger, buf.Bytes())
}

func Test_AppendTags_Wrong_Name(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/swagger.yaml")
	require.NoError(t, err, "failed open swagger.yaml")

	err = hd.AppendTags([]*model.Item{
		{
			Name:    "haters",
			Content: bytes.NewReader(hatersTagData),
		},
		{
			Name:    "name",
			Content: bytes.NewReader(nameTagData),
		},
		{
			Name:    "userss",
			Content: bytes.NewReader(usersTagData),
		},
	})
	require.ErrorIs(t, err, head.ErrInvalidContentName)
}

func Test_AppendTags_TagNameExists(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/swagger.yaml")
	require.NoError(t, err, "failed open swagger.yaml")

	err = hd.AppendTags([]*model.Item{
		{
			Name:    "haters",
			Content: bytes.NewReader(hatersTagData),
		},
		{
			Name:    "teachers",
			Content: bytes.NewReader(nameTagData),
		},
		{
			Name:    "userss",
			Content: bytes.NewReader(usersTagData),
		},
	})
	require.ErrorIs(t, err, head.ErrTagNameExists)
}
