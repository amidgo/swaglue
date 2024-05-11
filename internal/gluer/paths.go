package gluer

import (
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/amidgo/swaglue/pkg/logger"
)

var ErrGluePaths = errors.New("glue paths")

type PathsParser interface {
	Parse() error
	Paths() map[string]io.Reader
}

type PathsSetter interface {
	SetPaths(paths map[string]io.Reader) error
}

type PathsGluer struct {
	log    logger.DebugLogger
	parser PathsParser
	setter PathsSetter
}

func NewPathsGluer(log logger.DebugLogger, parser PathsParser, setter PathsSetter) *PathsGluer {
	return &PathsGluer{
		log:    log,
		parser: parser,
		setter: setter,
	}
}

func (g *PathsGluer) error(err error) error {
	return fmt.Errorf("%w, %w", ErrGluePaths, err)
}

func (g *PathsGluer) Glue() error {
	err := g.parser.Parse()
	if err != nil {
		return g.error(fmt.Errorf("parse paths, %w", err))
	}

	paths := g.parser.Paths()
	g.log.Debug("gluer paths", slog.Any("paths", paths))

	err = g.setter.SetPaths(paths)
	if err != nil {
		return g.error(fmt.Errorf("set paths, %w", err))
	}

	return nil
}
