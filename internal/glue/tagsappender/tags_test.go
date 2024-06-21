package tagsappender_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/amidgo/node/yaml"
	"github.com/amidgo/swaglue/internal/glue/tagsappender"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/model"
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
	tagsExpectedSwagger []byte
	//go:embed testdata/tags/empty_tags_expected_swagger.yaml
	emptyTagsExpectedSwagger []byte
)

func Test_AppendTags_Success(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/swagger.yaml", new(yaml.Decoder))
	require.NoError(t, err, "open swagger.yaml")

	appender := tagsappender.New(hd, new(yaml.Decoder))

	err = appender.AppendTags([]model.Item{
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
	require.NoError(t, err, "append components")

	buf := &bytes.Buffer{}
	err = hd.SaveTo(buf, &yaml.Encoder{Indent: 2})
	require.NoError(t, err, "save file")

	assert.Equal(t, tagsExpectedSwagger, buf.Bytes())
}

func Test_AppendTags_EmptyTags(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/empty.swagger.yaml", new(yaml.Decoder))
	require.NoError(t, err, "open swagger.yaml")

	appender := tagsappender.New(hd, new(yaml.Decoder))

	err = appender.AppendTags([]model.Item{
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
	require.NoError(t, err, "append components")

	buf := &bytes.Buffer{}
	err = hd.SaveTo(buf, &yaml.Encoder{Indent: 2})
	require.NoError(t, err, "save file")

	assert.Equal(t, emptyTagsExpectedSwagger, buf.Bytes())
}

func Test_AppendTags_Wrong_Name(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/swagger.yaml", new(yaml.Decoder))
	require.NoError(t, err, "open swagger.yaml")

	appender := tagsappender.New(hd, new(yaml.Decoder))

	err = appender.AppendTags([]model.Item{
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
	require.ErrorIs(t, err, tagsappender.ErrInvalidContentName)
}

func Test_AppendTags_TagNameExists(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/swagger.yaml", new(yaml.Decoder))
	require.NoError(t, err, "open swagger.yaml")

	appender := tagsappender.New(hd, new(yaml.Decoder))

	err = appender.AppendTags([]model.Item{
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
	require.ErrorIs(t, err, tagsappender.ErrTagNameExists)
}
