package fileformats_test

import (
	"testing"

	"github.com/amidgo/swaglue/internal/fileformats"
	"github.com/amidgo/tester"
	"github.com/stretchr/testify/require"
)

func Test_DetectFileFormat(t *testing.T) {
	tester.RunNamedTesters(t,
		DetectFileFormatCase{
			Format:               fileformats.JSONFormat,
			ExpectFileExtensions: []string{".json"},
		},
		DetectFileFormatCase{
			Format:               fileformats.YAMLFormat,
			ExpectFileExtensions: []string{".yaml", ".yml"},
		},
		DetectFileFormatCase{
			Format:    "abracadabra",
			ExpectErr: fileformats.ErrDetectFileFormat,
		},
	)
}

type DetectFileFormatCase struct {
	Format               string
	ExpectFileExtensions []string
	ExpectErr            error
}

func (c DetectFileFormatCase) Name() string {
	return "detect file format " + c.Format
}

func (c DetectFileFormatCase) Test(t *testing.T) {
	fileFormat, err := fileformats.Detect(c.Format)

	require.ErrorIs(t, err, c.ExpectErr)
	require.Equal(t, c.ExpectFileExtensions, fileFormat.Extensions())
}
