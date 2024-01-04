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
		files:     make(map[string]io.Reader),
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
	files     map[string]io.Reader
}

func (p *swaggerPathFileHandler) Files() map[string]io.Reader {
	return p.files
}

func (p *swaggerPathFileHandler) HandleFile(relativeFilePath string, file io.Reader) {
	relativeFilePath = strings.TrimSuffix(relativeFilePath, ".yaml")
	key := p.fileKey(relativeFilePath)
	p.files[key] = file
}

func (p *swaggerPathFileHandler) fileKey(filePath string) string {
	return p.keyPrefix + filePath
}
