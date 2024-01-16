package gluer

import "github.com/amidgo/swaglue/pkg/logger"

//go:generate mockgen -source routes.go -destination mocks/routes_mocks.gen.go -package gluermocks

type RoutesParser interface{}

type RoutesAppender interface{}

type RoutesGluer struct {
	log      logger.Logger
	parser   RoutesParser
	appender RoutesAppender
}

func NewRoutesGluer(log logger.Logger, parser RoutesParser, appender RoutesAppender) *RoutesGluer {
	return &RoutesGluer{
		log:      log,
		parser:   parser,
		appender: appender,
	}
}

func (g *RoutesGluer) Glue() error {
	return nil
}
