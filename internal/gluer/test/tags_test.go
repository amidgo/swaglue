package test_test

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/amidgo/swaglue/internal/gluer"
	gluermocks "github.com/amidgo/swaglue/internal/gluer/mocks"
	"github.com/amidgo/swaglue/internal/model"
	loggermocks "github.com/amidgo/swaglue/pkg/logger/mocks"
	"github.com/stretchr/testify/require"
)

func Test_TagsGluer_FailedParseComponentItems(t *testing.T) {
	parseErr := os.ErrPermission
	gluerTester := NewTagsGluerTester(t, parseErr)

	gluerTester.ExpectComponentsParse(parseErr)

	gluerTester.Test(t)
}

func Test_TagsGluer_FailedAppendTags(t *testing.T) {
	appendErr := io.ErrNoProgress
	componentItems := []*model.Item{
		{
			Name:    "any component",
			Content: strings.NewReader("contentik"),
		},
	}
	gluerTester := NewTagsGluerTester(t, appendErr)

	gluerTester.ExpectComponentsParse(nil)
	gluerTester.ExpectComponentItems(componentItems)
	gluerTester.ExpectComponentItemsDebugLog(componentItems)
	gluerTester.ExpectAppendTags(componentItems, appendErr)

	gluerTester.Test(t)
}

func Test_TagsGluer_Success(t *testing.T) {
	componentItems := []*model.Item{
		{
			Name:    "any component",
			Content: strings.NewReader("contentik"),
		},
	}
	gluerTester := NewTagsGluerTester(t, nil)

	gluerTester.ExpectComponentsParse(nil)
	gluerTester.ExpectComponentItems(componentItems)
	gluerTester.ExpectComponentItemsDebugLog(componentItems)
	gluerTester.ExpectAppendTags(componentItems, nil)

	gluerTester.Test(t)
}

type TagsGluerTester struct {
	debugLogger      *loggermocks.DebugLogger
	componentsParser *gluermocks.ComponentsParser
	tagsAppender     *gluermocks.TagsAppender

	ExpectErr error
}

func NewTagsGluerTester(t *testing.T, expectedErr error) *TagsGluerTester {
	t.Helper()

	return &TagsGluerTester{
		debugLogger:      loggermocks.NewDebugLogger(t),
		componentsParser: gluermocks.NewComponentsParser(t),
		tagsAppender:     gluermocks.NewTagsAppender(t),

		ExpectErr: expectedErr,
	}
}

func (gt *TagsGluerTester) ExpectComponentsParse(outErr error) {
	gt.componentsParser.EXPECT().Parse().Return(outErr).Once()
}

func (gt *TagsGluerTester) ExpectComponentItems(componentItems []*model.Item) {
	gt.componentsParser.EXPECT().ComponentItems().Return(componentItems).Once()
}

func (gt *TagsGluerTester) ExpectComponentItemsDebugLog(componentItems []*model.Item) {
	gt.debugLogger.EXPECT().Debug(
		"tags component items",
		slog.Any("componentItems", componentItems),
	).Once()
}

func (gt *TagsGluerTester) ExpectAppendTags(componentItems []*model.Item, outErr error) {
	gt.tagsAppender.EXPECT().AppendTags(componentItems).Return(outErr).Once()
}

func (gt *TagsGluerTester) Test(t *testing.T) {
	tagsGluer := gluer.NewTagsGluer(
		gt.debugLogger,
		gt.componentsParser,
		gt.tagsAppender,
	)
	actualErr := tagsGluer.Glue()

	if gt.ExpectErr == nil {
		require.NoError(t, actualErr)
	} else {
		require.ErrorIs(t, actualErr, gluer.ErrFailedGlueTags)
		require.ErrorIs(t, actualErr, gt.ExpectErr)
	}
}
