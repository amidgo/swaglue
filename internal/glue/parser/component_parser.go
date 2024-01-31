package parser

import (
	"io"
	"path"
	"strings"

	"github.com/amidgo/swaglue/internal/glue/model"
)

type SwaggerComponentParser struct {
	*swaggerComponentFileHandler
	*fileParser
}

func NewSwaggerComponentParser(basePackage string, targetFileFormat FileFormat) *SwaggerComponentParser {
	swaggerComponentFileHandler := &swaggerComponentFileHandler{
		files: make([]*model.Item, 0),
	}
	parser := &fileParser{
		basePackage:      basePackage,
		fileHandler:      swaggerComponentFileHandler,
		targetFileFormat: targetFileFormat,
	}

	return &SwaggerComponentParser{
		swaggerComponentFileHandler: swaggerComponentFileHandler,
		fileParser:                  parser,
	}
}

type swaggerComponentFileHandler struct {
	files []*model.Item
}

func (s *swaggerComponentFileHandler) ComponentItems() []*model.Item {
	return s.files
}

func (s *swaggerComponentFileHandler) HandleFile(filePath string, file io.Reader) {
	filePath = strings.TrimSuffix(filePath, ".yaml")
	_, fileName := path.Split(filePath)

	s.files = append(s.files, &model.Item{
		Name:    fileName,
		Content: file,
	})
}
