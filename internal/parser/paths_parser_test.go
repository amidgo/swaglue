package parser_test

import (
	_ "embed"
	"io"
	"testing"

	"github.com/amidgo/swaglue/internal/parser"
	"github.com/stretchr/testify/assert"
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
	parser := parser.NewSwaggerPathParser(basePackage, keyPrefix)

	err := parser.Parse()
	assert.NoError(t, err, "failed parse")

	files := parser.Files()

	assert.Equal(t, len(files), 2, "wrong len")

	postContent, _ := io.ReadAll(files["#/paths/groups/post"])
	assert.Equal(t, postContent, groupsPostFileContent)

	putContent, _ := io.ReadAll(files["#/paths/users/put"])
	assert.Equal(t, putContent, usersPutFileContent)
}
