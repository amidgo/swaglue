package gluer

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/amidgo/swaglue/internal/route"
	"github.com/amidgo/swaglue/pkg/logger"
)

var ErrGlueRoutes = errors.New("glue routes")

type RoutesParser interface {
	Parse() error
	Routes() []route.Route
}

type RoutesAppender interface {
	AppendRoutes(routes []route.Route) error
}

type RoutesGluer struct {
	log      logger.DebugLogger
	parser   RoutesParser
	appender RoutesAppender
}

func NewRoutesGluer(log logger.DebugLogger, parser RoutesParser, appender RoutesAppender) *RoutesGluer {
	return &RoutesGluer{
		log:      log,
		parser:   parser,
		appender: appender,
	}
}

func (g *RoutesGluer) error(err error) error {
	return fmt.Errorf("%w, %w", ErrGlueRoutes, err)
}

func (g *RoutesGluer) Glue() error {
	err := g.parser.Parse()
	if err != nil {
		return g.error(fmt.Errorf("parse routes, %w", err))
	}

	routes := g.parser.Routes()
	g.log.Debug("gluer routes", slog.Any("routes", routes))

	err = g.appender.AppendRoutes(routes)
	if err != nil {
		return g.error(fmt.Errorf("append routes, %w", err))
	}

	return nil
}
