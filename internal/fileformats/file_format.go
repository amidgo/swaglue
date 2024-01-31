package fileformats

import "errors"

const (
	YamlFileFormat = "yaml"
	JSONFileFormat = "json"
)

var ErrFailedDetectFileFormat = errors.New("failed detect file format")

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
	case YamlFileFormat:
		return YAML(), nil
	case JSONFileFormat:
		return JSON(), nil
	}

	return &SliceFileFormat{}, ErrFailedDetectFileFormat
}
