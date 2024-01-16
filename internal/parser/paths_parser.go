package parser

import (
	"io"
	"strings"
)

type SwaggerPathParser struct {
	*swaggerPathFileHandler
	*yamlFileParser
}

func NewSwaggerPathParser(basePackage string, keyPrefix string) *SwaggerPathParser {
	pathFileHandler := &swaggerPathFileHandler{
		keyPrefix: keyPrefix,
		paths:     make(map[string]io.Reader),
	}
	parser := &yamlFileParser{
		basePackage: basePackage,
		fileHandler: pathFileHandler,
	}

	return &SwaggerPathParser{
		swaggerPathFileHandler: pathFileHandler,
		yamlFileParser:         parser,
	}
}

type swaggerPathFileHandler struct {
	keyPrefix string
	paths     map[string]io.Reader
}

func (p *swaggerPathFileHandler) Paths() map[string]io.Reader {
	return p.paths
}

func (p *swaggerPathFileHandler) HandleFile(relativeFilePath string, file io.Reader) {
	relativeFilePath = strings.TrimSuffix(relativeFilePath, ".yaml")
	key := p.fileKey(relativeFilePath)
	p.paths[key] = file
}

func (p *swaggerPathFileHandler) fileKey(filePath string) string {
	return p.keyPrefix + filePath
}
