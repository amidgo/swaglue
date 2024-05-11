package fileformats_test

import (
	"testing"

	"github.com/amidgo/swaglue/internal/fileformats"
	"github.com/amidgo/swaglue/pkg/tester"
	"github.com/stretchr/testify/require"
)

func Test_DetectFileFormat(t *testing.T) {
	var tester tester.NamedContainer

	tester.AddNamedTester(
		DetectFileFormatCase{
			Format:               fileformats.JSONFileFormat,
			ExpectFileExtensions: []string{".json"},
		},
	)

	tester.AddNamedTester(
		DetectFileFormatCase{
			Format:               fileformats.YamlFileFormat,
			ExpectFileExtensions: []string{".yaml", ".yml"},
		},
	)

	tester.AddNamedTester(
		DetectFileFormatCase{
			Format:    "abracadabra",
			ExpectErr: fileformats.ErrDetectFileFormat,
		},
	)

	tester.Test(t)
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
