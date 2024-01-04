package parser

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

//go:generate mockgen -source parser.go -destination parser_mocks.go -package parser

type FileHandler interface {
	HandleFile(relativeFilePath string, file io.Reader)
}

type yamlFileParser struct {
	basePackage string
	fileHandler FileHandler
}

func (p *yamlFileParser) Parse() error {
	entries, err := os.ReadDir(p.basePackage)
	if err != nil {
		return &FailedReadDirectoryError{
			DirectoryPath: p.basePackage,
			Err:           err,
		}
	}
	err = p.parseEntries(entries, "")
	if err != nil {
		return fmt.Errorf("parse entries: %w", err)
	}
	return nil
}

func (p *yamlFileParser) parseEntries(entries []os.DirEntry, pathPrefix string) error {
	for _, entry := range entries {
		err := p.parseEntry(entry, pathPrefix)
		if err != nil {
			return fmt.Errorf("failed parse entry %s, err: %w", entry.Name(), err)
		}
	}
	return nil
}

func (p *yamlFileParser) parseEntry(entry os.DirEntry, pathPrefix string) error {
	pathPrefix = path.Join(pathPrefix, entry.Name())
	switch {
	case entry.IsDir():
		err := p.parseDirEntry(entry, pathPrefix)
		if err != nil {
			return fmt.Errorf("parse dir entry, pathprefix %s, err: %w", pathPrefix, err)
		}
	case isYamlFile(entry):
		err := p.parseYamlFile(pathPrefix)
		if err != nil {
			return fmt.Errorf("parse yaml file, pathprefix %s, err: %w", pathPrefix, err)
		}
	}
	return nil
}

func (p *yamlFileParser) parseDirEntry(entry os.DirEntry, pathPrefix string) error {
	path := p.path(pathPrefix)
	entries, err := os.ReadDir(path)
	if err != nil {
		return &FailedReadDirectoryError{
			DirectoryPath: path,
			Err:           err,
		}
	}
	return p.parseEntries(entries, pathPrefix)
}

func (p *yamlFileParser) parseYamlFile(pathPrefix string) error {
	filePath := p.path(pathPrefix)
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return &FailedReadFileError{
			FilePath: filePath,
			Err:      err,
		}
	}
	p.fileHandler.HandleFile(pathPrefix, file)
	return nil
}

func isYamlFile(entry os.DirEntry) bool {
	return strings.HasSuffix(entry.Name(), ".yaml")
}

func (p *yamlFileParser) path(pathPrefix string) string {
	return path.Join(p.basePackage, pathPrefix)
}
