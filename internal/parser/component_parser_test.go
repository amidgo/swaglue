package parser_test

import (
	_ "embed"
	"io"
	"slices"
	"testing"

	"github.com/amidgo/swaglue/internal/parser"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed testdata/schemas/AnimalName.yaml
	animalNameFileContent []byte
	//go:embed testdata/schemas/User.yaml
	userFileContent []byte
)

func TestSwaggerComponentParser(t *testing.T) {
	const basePackage = "./testdata/schemas"

	parser := parser.NewSwaggerComponentParser(basePackage)
	err := parser.Parse()
	assert.NoError(t, err, "failed parse")

	files := parser.Files()
	assert.Equal(t, len(files), 2, "wrong files length")

	animalFileContentEqual := readerEqualContent(files[0].Content, animalNameFileContent)
	assert.True(t, animalFileContentEqual, "animal file content not equal")

	userFileContentEqual := readerEqualContent(files[1].Content, userFileContent)
	assert.True(t, userFileContentEqual, "user file content not equal")
}

func readerEqualContent(r io.Reader, content []byte) bool {
	data, _ := io.ReadAll(r)
	return slices.Equal(data, content)
}
