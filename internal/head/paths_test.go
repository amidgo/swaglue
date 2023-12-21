package head_test

import (
	"bytes"
	_ "embed"
	"io"
	"testing"

	"github.com/amidgo/swaglue/internal/head"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed testdata/paths/get.yaml
	getPathData []byte
	//go:embed testdata/paths/post.yaml
	postPathData []byte

	//go:embed testdata/test_paths_expected_swagger.yaml
	pathsExpectedData []byte
)

func TestPaths(t *testing.T) {
	head, err := head.ParseHeadFromFile("testdata/swagger.yaml")
	assert.NoError(t, err, "failed open swagger.yaml")
	head.SetPaths(map[string]io.Reader{
		"#/paths/get":  bytes.NewReader(getPathData),
		"#/paths/post": bytes.NewReader(postPathData),
	})
	buf := &bytes.Buffer{}
	err = head.SaveTo(buf)
	assert.NoError(t, err, "failed save file")
	assert.Equal(t, pathsExpectedData, buf.Bytes())
}
