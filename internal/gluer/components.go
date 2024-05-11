package gluer

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/amidgo/swaglue/internal/model"
	"github.com/amidgo/swaglue/pkg/logger"
)

var ErrGlueComponent = errors.New("glue component")

type ComponentsParser interface {
	Parse() error
	ComponentItems() []*model.Item
}

type ComponentsAppender interface {
	AppendComponent(componentName string, componentItems []*model.Item) error
}

type ComponentsGluer struct {
	ComponentName string
	log           logger.DebugLogger
	parser        ComponentsParser
	appender      ComponentsAppender
}

func NewComponentsGluer(
	componentName string,
	log logger.DebugLogger,
	parser ComponentsParser,
	appender ComponentsAppender,
) *ComponentsGluer {
	return &ComponentsGluer{
		ComponentName: componentName,
		log:           log,
		parser:        parser,
		appender:      appender,
	}
}

func (g *ComponentsGluer) error(err error) error {
	return fmt.Errorf("%w, componentName %s, %w", ErrGlueComponent, g.ComponentName, err)
}

func (g *ComponentsGluer) Glue() error {
	err := g.parser.Parse()
	if err != nil {
		return g.error(fmt.Errorf("parse component items, %w", err))
	}

	componentItems := g.parser.ComponentItems()
	g.log.Debug(
		"component items",
		slog.String("componentName", g.ComponentName),
		slog.Any("componentItems", componentItems),
	)

	err = g.appender.AppendComponent(g.ComponentName, componentItems)
	if err != nil {
		return g.error(fmt.Errorf("append component items, %w", err))
	}

	return nil
}
