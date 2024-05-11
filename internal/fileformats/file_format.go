package fileformats

import "errors"

const (
	YAMLFormat = "yaml"
	JSONFormat = "json"
)

var ErrDetectFileFormat = errors.New("detect file format")

type SliceFileFormat struct {
	extensions []string
}

func (s *SliceFileFormat) Extensions() []string {
	return s.extensions
}

func YAML() *SliceFileFormat {
	return &SliceFileFormat{
		extensions: []string{".yaml", ".yml"},
	}
}

func JSON() *SliceFileFormat {
	return &SliceFileFormat{
		extensions: []string{".json"},
	}
}

func Detect(format string) (*SliceFileFormat, error) {
	switch format {
	case YAMLFormat:
		return YAML(), nil
	case JSONFormat:
		return JSON(), nil
	}

	return &SliceFileFormat{}, ErrDetectFileFormat
}
