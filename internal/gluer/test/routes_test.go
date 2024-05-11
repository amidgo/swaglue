package test_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/amidgo/swaglue/internal/gluer"
	gluermocks "github.com/amidgo/swaglue/internal/gluer/mocks"
	"github.com/amidgo/swaglue/internal/model"
	loggermocks "github.com/amidgo/swaglue/pkg/logger/mocks"
	"github.com/stretchr/testify/require"
)

func Test_RoutesGluer_FailedParseRoutes(t *testing.T) {
	parseErr := os.ErrClosed
	routesGluerTester := NewRoutesGluerTester(t, parseErr)

	routesGluerTester.ExpectParseRoutes(parseErr)

	routesGluerTester.Test(t)
}

func Test_RoutestGluer_FailedAppendRoutes(t *testing.T) {
	appendErr := os.ErrInvalid
	routes := []*model.Route{
		{
			Name: "dflsdfkjdddlja",
			Methods: []*model.RouteMethod{
				{}, {}, {},
			},
		},
	}
	routesGluerTester := NewRoutesGluerTester(t, appendErr)

	routesGluerTester.ExpectParseRoutes(nil)
	routesGluerTester.ExpectRoutes(routes)
	routesGluerTester.ExpectRoutesDebugLog(routes)
	routesGluerTester.ExpectAppendRoutes(routes, appendErr)

	routesGluerTester.Test(t)
}

func Test_RoutesGluer_Success(t *testing.T) {
	routes := []*model.Route{
		{
			Name: "dflsdfkjdddlja",
			Methods: []*model.RouteMethod{
				{}, {}, {},
			},
		},
	}
	routesGluerTester := NewRoutesGluerTester(t, nil)

	routesGluerTester.ExpectParseRoutes(nil)
	routesGluerTester.ExpectRoutes(routes)
	routesGluerTester.ExpectRoutesDebugLog(routes)
	routesGluerTester.ExpectAppendRoutes(routes, nil)

	routesGluerTester.Test(t)
}

type RoutesGluerTester struct {
	debugLogger    *loggermocks.DebugLogger
	routesParser   *gluermocks.RoutesParser
	routesAppender *gluermocks.RoutesAppender

	ExpectErr error
}

func NewRoutesGluerTester(t *testing.T, expectErr error) *RoutesGluerTester {
	t.Helper()

	return &RoutesGluerTester{
		debugLogger:    loggermocks.NewDebugLogger(t),
		routesParser:   gluermocks.NewRoutesParser(t),
		routesAppender: gluermocks.NewRoutesAppender(t),
		ExpectErr:      expectErr,
	}
}

func (rt *RoutesGluerTester) ExpectParseRoutes(outErr error) {
	rt.routesParser.EXPECT().Parse().Return(outErr).Once()
}

func (rt *RoutesGluerTester) ExpectRoutes(routes []*model.Route) {
	rt.routesParser.EXPECT().Routes().Return(routes).Once()
}

func (rt *RoutesGluerTester) ExpectRoutesDebugLog(routes []*model.Route) {
	rt.debugLogger.EXPECT().Debug("gluer routes", slog.Any("routes", routes))
}

func (rt *RoutesGluerTester) ExpectAppendRoutes(routes []*model.Route, outErr error) {
	rt.routesAppender.EXPECT().AppendRoutes(routes).Return(outErr).Once()
}

func (rt *RoutesGluerTester) Test(t *testing.T) {
	routesGluer := gluer.NewRoutesGluer(
		rt.debugLogger,
		rt.routesParser,
		rt.routesAppender,
	)

	actualErr := routesGluer.Glue()

	if rt.ExpectErr == nil {
		require.NoError(t, actualErr)
	} else {
		require.ErrorIs(t, actualErr, gluer.ErrGlueRoutes)
		require.ErrorIs(t, actualErr, rt.ExpectErr)
	}
}
