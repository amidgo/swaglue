package gluer

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/amidgo/swaglue/internal/model"
	"github.com/amidgo/swaglue/pkg/logger"
)

var ErrFailedGlueTags = errors.New("failed glue tags")

type TagsAppender interface {
	AppendTags(items []*model.Item) error
}

type TagsGluer struct {
	log      logger.DebugLogger
	parser   ComponentsParser
	appender TagsAppender
}

func NewTagsGluer(log logger.DebugLogger, parser ComponentsParser, appender TagsAppender) *TagsGluer {
	return &TagsGluer{
		log:      log,
		parser:   parser,
		appender: appender,
	}
}

func (g *TagsGluer) error(err error) error {
	return fmt.Errorf("%w, %w", ErrFailedGlueTags, err)
}

func (g *TagsGluer) Glue() error {
	err := g.parser.Parse()
	if err != nil {
		return g.error(fmt.Errorf("failed parse tags component items, err: %w", err))
	}

	componentItems := g.parser.ComponentItems()
	g.log.Debug(
		"tags component items",
		slog.Any("componentItems", componentItems),
	)

	err = g.appender.AppendTags(componentItems)
	if err != nil {
		return g.error(fmt.Errorf("failed append tags, err: %w", err))
	}

	return nil
}
