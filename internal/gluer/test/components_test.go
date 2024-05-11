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

func Test_ComponentsGluer_FailedParseComponentItems(t *testing.T) {
	parseErr := os.ErrPermission
	componentName := "schema"
	gluerTester := NewComponentsGluerTester(t, componentName, parseErr)

	gluerTester.ExpectComponentsParse(parseErr)

	gluerTester.Test(t)
}

func Test_ComponentsGluer_FailedAppendComponentItems(t *testing.T) {
	appendErr := io.ErrNoProgress
	componentName := "component name"
	componentItems := []*model.Item{
		{
			Name:    "any component",
			Content: strings.NewReader("contentik"),
		},
	}
	gluerTester := NewComponentsGluerTester(t, componentName, appendErr)

	gluerTester.ExpectComponentsParse(nil)
	gluerTester.ExpectComponentItems(componentItems)
	gluerTester.ExpectComponentItemsDebugLog(componentItems)
	gluerTester.ExpectComponentItemsAppend(componentItems, appendErr)

	gluerTester.Test(t)
}

func Test_ComponentsGluer_Success(t *testing.T) {
	componentName := "component name"
	componentItems := []*model.Item{
		{
			Name:    "any component",
			Content: strings.NewReader("contentik"),
		},
	}
	gluerTester := NewComponentsGluerTester(t, componentName, nil)

	gluerTester.ExpectComponentsParse(nil)
	gluerTester.ExpectComponentItems(componentItems)
	gluerTester.ExpectComponentItemsDebugLog(componentItems)
	gluerTester.ExpectComponentItemsAppend(componentItems, nil)

	gluerTester.Test(t)
}

type ComponentsGluerTester struct {
	debugLogger        *loggermocks.DebugLogger
	componentsParser   *gluermocks.ComponentsParser
	componentsAppender *gluermocks.ComponentsAppender

	ComponentName string
	ExpectErr     error
}

func NewComponentsGluerTester(t *testing.T, componentName string, expectErr error) *ComponentsGluerTester {
	t.Helper()

	return &ComponentsGluerTester{
		debugLogger:        loggermocks.NewDebugLogger(t),
		componentsParser:   gluermocks.NewComponentsParser(t),
		componentsAppender: gluermocks.NewComponentsAppender(t),
		ComponentName:      componentName,
		ExpectErr:          expectErr,
	}
}

func (ct *ComponentsGluerTester) ExpectComponentsParse(outErr error) {
	ct.componentsParser.EXPECT().Parse().Return(outErr).Once()
}

func (ct *ComponentsGluerTester) ExpectComponentItemsDebugLog(componentItems []*model.Item) {
	ct.debugLogger.EXPECT().Debug(
		"component items",
		slog.String("componentName", ct.ComponentName),
		slog.Any("componentItems", componentItems),
	).Once()
}

func (ct *ComponentsGluerTester) ExpectComponentItems(outItems []*model.Item) {
	ct.componentsParser.EXPECT().ComponentItems().Return(outItems).Once()
}

func (ct *ComponentsGluerTester) ExpectComponentItemsAppend(items []*model.Item, outErr error) {
	ct.componentsAppender.EXPECT().AppendComponent(ct.ComponentName, items).Return(outErr).Once()
}

func (ct *ComponentsGluerTester) Test(t *testing.T) {
	componentsGluer := gluer.NewComponentsGluer(
		ct.ComponentName,
		ct.debugLogger,
		ct.componentsParser,
		ct.componentsAppender,
	)
	actualErr := componentsGluer.Glue()

	if ct.ExpectErr == nil {
		require.NoError(t, actualErr)
	} else {
		require.ErrorIs(t, actualErr, gluer.ErrFailedGlueComponent)
		require.ErrorIs(t, actualErr, ct.ExpectErr)
	}
}
