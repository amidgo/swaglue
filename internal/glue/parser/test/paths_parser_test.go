package parser_test

import (
	_ "embed"
	"io"
	"testing"

	"github.com/amidgo/swaglue/internal/fileformats"
	"github.com/amidgo/swaglue/internal/glue/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/path/groups/post.yaml
	groupsPostFileContent []byte
	//go:embed testdata/path/users/put.yaml
	usersPutFileContent []byte
)

func TestSwaggerPathParser(t *testing.T) {
	const (
		basePackage = "./testdata/path"
		keyPrefix   = "#/paths/"
	)

	parser := parser.NewSwaggerPathsParser(basePackage, keyPrefix, fileformats.YAML())

	err := parser.Parse()
	require.NoError(t, err, "failed parse")

	files := parser.Paths()
	assert.Len(t, files, 2, "wrong len")

	postContent, _ := io.ReadAll(files["#/paths/groups/post"])
	assert.Equal(t, postContent, groupsPostFileContent)

	putContent, _ := io.ReadAll(files["#/paths/users/put"])
	assert.Equal(t, putContent, usersPutFileContent)
}
