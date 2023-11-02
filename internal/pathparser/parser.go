package pathparser

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
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
	keyPrefix   string
	files       map[string]io.Reader
	usePrefix   bool
}

func NewParser(basePackage string, keyPrefix string) *Parser {
	return &Parser{basePackage: basePackage, keyPrefix: keyPrefix, files: make(map[string]io.Reader)}
}

func (p *Parser) Files() map[string]io.Reader {
	return p.files
}

func (p *Parser) setFile(filePath string, file io.Reader) {
	filePath = strings.TrimSuffix(filePath, ".yaml")
	var key string
	key = p.keyPrefix + filePath
	p.files[key] = file
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
