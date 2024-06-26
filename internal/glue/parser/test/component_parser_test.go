package parser_test

import (
	_ "embed"
	"io"
	"slices"
	"testing"

	"github.com/amidgo/swaglue/internal/fileformats"
	"github.com/amidgo/swaglue/internal/glue/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/schemas/AnimalName.yaml
	animalNameFileContent []byte
	//go:embed testdata/schemas/User.yaml
	userFileContent []byte
)

func TestSwaggerComponentParser(t *testing.T) {
	const basePackage = "./testdata/schemas"

	parser := parser.NewSwaggerComponentParser(basePackage, fileformats.YAML())
	err := parser.Parse()
	require.NoError(t, err, "parse")

	files := parser.ComponentItems()
	assert.Len(t, files, 2, "wrong files length")

	animalFileContentEqual := readerEqualContent(files[0].Content, animalNameFileContent)
	assert.True(t, animalFileContentEqual, "animal file content not equal")

	userFileContentEqual := readerEqualContent(files[1].Content, userFileContent)
	assert.True(t, userFileContentEqual, "user file content not equal")
}

func readerEqualContent(r io.Reader, content []byte) bool {
	data, err := io.ReadAll(r)
	if err != nil {
		return false
	}

	return slices.Equal(data, content)
}
