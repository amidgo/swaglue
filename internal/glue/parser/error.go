package parser

import (
	"fmt"
)

type ReadFileError struct {
	FilePath string
}

func (e *ReadFileError) Error() string {
	return fmt.Sprintf("read file: %s", e.FilePath)
}

func (e *ReadFileError) Is(target error) bool {
	dirErr, ok := target.(*ReadFileError)
	if ok {
		return e.FilePath == dirErr.FilePath
	}

	return false
}

type ReadDirectoryError struct {
	DirectoryPath string
}

func (e *ReadDirectoryError) Error() string {
	return fmt.Sprintf("read directory: %s", e.DirectoryPath)
}

func (e *ReadDirectoryError) Is(target error) bool {
	dirErr, ok := target.(*ReadDirectoryError)
	if ok {
		return e.DirectoryPath == dirErr.DirectoryPath
	}

	return false
}
