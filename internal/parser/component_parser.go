package parser

import (
	"io"
	"path"
	"strings"

	"github.com/amidgo/swaglue/internal/model"
)

type SwaggerComponentParser struct {
	*swaggerComponentFileHandler
	*yamlFileParser
}

func NewSwaggerComponentParser(basePackage string) *SwaggerComponentParser {
	swaggerComponentFileHandler := &swaggerComponentFileHandler{
		files: make([]*model.SwaggerComponentItem, 0),
	}
	parser := &yamlFileParser{
		basePackage: basePackage,
		fileHandler: swaggerComponentFileHandler,
	}

	return &SwaggerComponentParser{
		swaggerComponentFileHandler: swaggerComponentFileHandler,
		yamlFileParser:              parser,
	}
}

type swaggerComponentFileHandler struct {
	files []*model.SwaggerComponentItem
}

func (s *swaggerComponentFileHandler) Files() []*model.SwaggerComponentItem {
	return s.files
}

func (s *swaggerComponentFileHandler) HandleFile(filePath string, file io.Reader) {
	filePath = strings.TrimSuffix(filePath, ".yaml")
	_, fileName := path.Split(filePath)
	s.files = append(s.files, &model.SwaggerComponentItem{
		Name:    fileName,
		Content: file,
	})
}
