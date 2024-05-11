package parser

import (
	"io"
	"path"
	"strings"

	"github.com/amidgo/swaglue/internal/model"
)

type ComponentParser struct {
	*swaggerComponentFileHandler
	*fileParser
}

func NewSwaggerComponentParser(basePackage string, targetFileFormat FileFormat) *ComponentParser {
	swaggerComponentFileHandler := &swaggerComponentFileHandler{
		files: make([]*model.Item, 0),
	}

	parser := &fileParser{
		basePackage:      basePackage,
		fileHandler:      swaggerComponentFileHandler,
		targetFileFormat: targetFileFormat,
	}

	return &ComponentParser{
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
	filePath = strings.TrimSuffix(filePath, ".yml")
	filePath = strings.TrimSuffix(filePath, ".json")

	_, fileName := path.Split(filePath)

	s.files = append(s.files, &model.Item{
		Name:    fileName,
		Content: file,
	})
}
