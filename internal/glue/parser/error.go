package parser

type ReadFileError struct {
	FilePath string
}

func (e *ReadFileError) Error() string {
	return "read file: " + e.FilePath
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
	return "read directory: " + e.DirectoryPath
}

func (e *ReadDirectoryError) Is(target error) bool {
	dirErr, ok := target.(*ReadDirectoryError)
	if ok {
		return e.DirectoryPath == dirErr.DirectoryPath
	}

	return false
}
