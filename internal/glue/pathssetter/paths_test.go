package pathssetter_test

import (
	"bytes"
	_ "embed"
	"io"
	"testing"

	"github.com/amidgo/node/yaml"
	"github.com/amidgo/swaglue/internal/glue/pathssetter"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/paths/get.yaml
	getPathData []byte
	//go:embed testdata/paths/post.yaml
	postPathData []byte

	//go:embed testdata/paths/paths_expected_swagger.yaml
	pathsExpectedData []byte
)

func TestPaths(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/swagger.yaml", new(yaml.Decoder))
	require.NoError(t, err, "failed open swagger.yaml")

	pathsSetter := pathssetter.New(hd, new(yaml.Decoder))

	err = pathsSetter.SetPaths(map[string]io.Reader{
		"#/paths/get":  bytes.NewReader(getPathData),
		"#/paths/post": bytes.NewReader(postPathData),
	})
	require.NoError(t, err, "failed set paths")

	buf := &bytes.Buffer{}
	err = hd.SaveTo(buf, &yaml.Encoder{Indent: 2})
	require.NoError(t, err, "failed save file")

	assert.Equal(t, pathsExpectedData, buf.Bytes())
}

func TestPaths_InvalidRef(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/swagger.yaml", new(yaml.Decoder))
	require.NoError(t, err, "failed open swagger.yaml")

	pathsSetter := pathssetter.New(hd, new(yaml.Decoder))

	err = pathsSetter.SetPaths(map[string]io.Reader{
		"#/paths/get": bytes.NewReader(getPathData),
	})

	require.ErrorIs(t, err, pathssetter.ErrInvalidRef)
}
