package test_test

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/amidgo/swaglue/internal/glue/gluer"
	gluermocks "github.com/amidgo/swaglue/internal/glue/gluer/mocks"
	loggermocks "github.com/amidgo/swaglue/pkg/logger/mocks"
	"github.com/stretchr/testify/require"
)

func Test_PathsGluer_FailedParsePaths(t *testing.T) {
	pathsParseErr := os.ErrNotExist
	pathsGluerTester := NewPathsGluerTester(t, pathsParseErr)

	pathsGluerTester.ExpectPathsParse(pathsParseErr)

	pathsGluerTester.Test(t)
}

func Test_PathsGluer_FailedSetPaths(t *testing.T) {
	paths := map[string]io.Reader{
		"example-path": strings.NewReader("paths data"),
	}
	pathsSetErr := io.ErrUnexpectedEOF
	pathsGluerTester := NewPathsGluerTester(t, pathsSetErr)

	pathsGluerTester.ExpectPathsParse(nil)
	pathsGluerTester.ExpectPaths(paths)
	pathsGluerTester.ExpectPathsDebugLog(paths)
	pathsGluerTester.ExpectPathsSet(paths, pathsSetErr)

	pathsGluerTester.Test(t)
}

func Test_PathsGluer_Success(t *testing.T) {
	paths := map[string]io.Reader{
		"example-path": strings.NewReader("paths data"),
	}
	pathsGluerTester := NewPathsGluerTester(t, nil)

	pathsGluerTester.ExpectPathsParse(nil)
	pathsGluerTester.ExpectPaths(paths)
	pathsGluerTester.ExpectPathsDebugLog(paths)
	pathsGluerTester.ExpectPathsSet(paths, nil)

	pathsGluerTester.Test(t)
}

type PathsGluerTester struct {
	debugLogger *loggermocks.DebugLogger
	pathsParser *gluermocks.PathsParser
	pathsSetter *gluermocks.PathsSetter

	ExpectErr error
}

func NewPathsGluerTester(t *testing.T, expectErr error) *PathsGluerTester {
	t.Helper()

	return &PathsGluerTester{
		debugLogger: loggermocks.NewDebugLogger(t),
		pathsParser: gluermocks.NewPathsParser(t),
		pathsSetter: gluermocks.NewPathsSetter(t),
		ExpectErr:   expectErr,
	}
}

func (pt *PathsGluerTester) ExpectPathsParse(outErr error) {
	pt.pathsParser.EXPECT().Parse().Return(outErr).Once()
}

func (pt *PathsGluerTester) ExpectPaths(outPaths map[string]io.Reader) {
	pt.pathsParser.EXPECT().Paths().Return(outPaths).Once()
}

func (pt *PathsGluerTester) ExpectPathsDebugLog(paths map[string]io.Reader) {
	pt.debugLogger.EXPECT().Debug("gluer paths", slog.Any("paths", paths)).Once()
}

func (pt *PathsGluerTester) ExpectPathsSet(paths map[string]io.Reader, outErr error) {
	pt.pathsSetter.EXPECT().SetPaths(paths).Return(outErr).Once()
}

func (pt *PathsGluerTester) Test(t *testing.T) {
	pathsGluer := gluer.NewPathsGluer(pt.debugLogger, pt.pathsParser, pt.pathsSetter)

	actualErr := pathsGluer.Glue()

	if pt.ExpectErr == nil {
		require.NoError(t, actualErr)
	} else {
		require.ErrorIs(t, actualErr, gluer.ErrFailedGluePaths)
		require.ErrorIs(t, actualErr, pt.ExpectErr)
	}
}
