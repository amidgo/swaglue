package componentparser

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/amidgo/swaglue/internal/model"
)

type FailedReadFileError struct {
	FilePath string
	Err      error
}

func (e FailedReadFileError) Error() string {
	return fmt.Sprintf("failed read file: %s, err: %s", e.FilePath, e.Err)
}

func (e FailedReadFileError) Unwrap() error {
	return e.Err
}

type FailedReadDirectoryError struct {
	DirectoryPath string
	Err           error
}

func (e FailedReadDirectoryError) Error() string {
	return fmt.Sprintf("failed read directory: %s, err: %s", e.DirectoryPath, e.Err)
}

func (e FailedReadDirectoryError) Unwrap() error {
	return e.Err
}

type Parser struct {
	basePackage string
	files       []*model.Component
}

func NewParser(basePackage string) *Parser {
	return &Parser{basePackage: basePackage, files: make([]*model.Component, 0)}
}

func (p *Parser) Files() []*model.Component {
	return p.files
}

func (p *Parser) setFile(filePath string, file io.Reader) {
	filePath = strings.TrimSuffix(filePath, ".yaml")
	_, fileName := path.Split(filePath)
	p.files = append(p.files, &model.Component{
		Name:    fileName,
		Content: file,
	})
}

func (p *Parser) path(pathPrefix string) string {
	return path.Join(p.basePackage, pathPrefix)
}

func (p *Parser) parseEntry(entry os.DirEntry, pathPrefix string) error {
	pathPrefix = path.Join(pathPrefix, entry.Name())
	if entry.IsDir() {
		path := p.path(pathPrefix)
		entries, err := os.ReadDir(path)
		if err != nil {
			return FailedReadDirectoryError{
				DirectoryPath: path,
				Err:           err,
			}
		}
		for _, entry := range entries {
			err := p.parseEntry(entry, pathPrefix)
			if err != nil {
				return err
			}
		}
	} else if strings.HasSuffix(entry.Name(), ".yaml") {
		filePath := p.path(pathPrefix)
		value, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return FailedReadFileError{
				FilePath: filePath,
				Err:      err,
			}
		}
		p.setFile(pathPrefix, value)
	}
	return nil
}

func (p *Parser) Parse() error {
	entries, err := os.ReadDir(p.basePackage)
	if err != nil {
		return FailedReadDirectoryError{
			DirectoryPath: p.basePackage,
			Err:           err,
		}
	}
	for _, entry := range entries {
		err := p.parseEntry(entry, "")
		if err != nil {
			return err
		}
	}
	return nil
}
