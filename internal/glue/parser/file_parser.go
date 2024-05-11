package parser

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

//go:generate mockgen -source parser.go -destination parser_mocks.go -package parser

type FileFormat interface {
	Extensions() []string
}

type FileHandler interface {
	HandleFile(relativeFilePath string, file io.Reader)
}

type fileParser struct {
	basePackage      string
	fileHandler      FileHandler
	targetFileFormat FileFormat
}

func (p *fileParser) Parse() error {
	entries, err := os.ReadDir(p.basePackage)
	if err != nil {
		return errors.Join(
			&ReadDirectoryError{
				DirectoryPath: p.basePackage,
			},
			err,
		)
	}

	err = p.parseEntries(entries, "")
	if err != nil {
		return fmt.Errorf("parse entries: %w", err)
	}

	return nil
}

func (p *fileParser) parseEntries(entries []os.DirEntry, pathPrefix string) error {
	for _, entry := range entries {
		err := p.parseEntry(entry, pathPrefix)
		if err != nil {
			return fmt.Errorf("parse entry %s, %w", entry.Name(), err)
		}
	}

	return nil
}

func (p *fileParser) parseEntry(entry os.DirEntry, pathPrefix string) error {
	pathPrefix = path.Join(pathPrefix, entry.Name())

	switch {
	case entry.IsDir():
		err := p.parseDirEntry(pathPrefix)
		if err != nil {
			return fmt.Errorf("parse dir entry, pathprefix %s, %w", pathPrefix, err)
		}
	case p.isTargetFile(entry):
		err := p.parseFile(pathPrefix)
		if err != nil {
			return fmt.Errorf("parse yaml file, pathprefix %s, %w", pathPrefix, err)
		}
	}

	return nil
}

func (p *fileParser) parseDirEntry(pathPrefix string) error {
	path := p.path(pathPrefix)

	entries, err := os.ReadDir(path)
	if err != nil {
		return errors.Join(
			&ReadDirectoryError{
				DirectoryPath: path,
			},
			err,
		)
	}

	return p.parseEntries(entries, pathPrefix)
}

func (p *fileParser) parseFile(pathPrefix string) error {
	filePath := p.path(pathPrefix)

	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return errors.Join(
			&ReadFileError{
				FilePath: filePath,
			},
			err,
		)
	}

	p.fileHandler.HandleFile(pathPrefix, file)

	return nil
}

func (p *fileParser) isTargetFile(entry os.DirEntry) bool {
	fileExt := path.Ext(entry.Name())
	for _, ext := range p.targetFileFormat.Extensions() {
		if fileExt == ext {
			return true
		}
	}

	return false
}

func (p *fileParser) path(pathPrefix string) string {
	return path.Join(p.basePackage, pathPrefix)
}
