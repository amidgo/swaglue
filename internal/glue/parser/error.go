package parser

import "fmt"

type FailedReadFileError struct {
	FilePath string
	Err      error
}

func (e *FailedReadFileError) Error() string {
	return fmt.Sprintf("failed read file: %s, err: %s", e.FilePath, e.Err)
}

func (e *FailedReadFileError) Unwrap() error {
	return e.Err
}

type FailedReadDirectoryError struct {
	DirectoryPath string
	Err           error
}

func (e *FailedReadDirectoryError) Error() string {
	return fmt.Sprintf("failed read directory: %s, err: %s", e.DirectoryPath, e.Err)
}

func (e *FailedReadDirectoryError) Unwrap() error {
	return e.Err
}
