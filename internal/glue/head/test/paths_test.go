package head_test

import (
	"bytes"
	_ "embed"
	"io"
	"testing"

	"github.com/amidgo/swaglue/internal/glue/head"
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
	head, err := head.ParseHeadFromFile("testdata/swagger.yaml")
	require.NoError(t, err, "failed open swagger.yaml")

	err = head.SetPaths(map[string]io.Reader{
		"#/paths/get":  bytes.NewReader(getPathData),
		"#/paths/post": bytes.NewReader(postPathData),
	})
	require.NoError(t, err, "failed set paths")

	buf := &bytes.Buffer{}
	err = head.SaveTo(buf)
	require.NoError(t, err, "failed save file")

	assert.Equal(t, pathsExpectedData, buf.Bytes())
}

func TestPaths_InvalidRef(t *testing.T) {
	hd, err := head.ParseHeadFromFile("testdata/swagger.yaml")
	require.NoError(t, err, "failed open swagger.yaml")

	err = hd.SetPaths(map[string]io.Reader{
		"#/paths/get": bytes.NewReader(getPathData),
	})
	require.ErrorIs(t, err, head.ErrInvalidRef)
}
