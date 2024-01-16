package gluer

import (
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/amidgo/swaglue/pkg/logger"
)

//go:generate mockgen -source paths.go -destination mocks/paths_mocks.gen.go -package gluermocks

var ErrFailedGluePaths = errors.New("failed glue paths")

type PathsParser interface {
	Parse() error
	Paths() map[string]io.Reader
}

type PathsSetter interface {
	SetPaths(paths map[string]io.Reader) error
}

type PathsGluer struct {
	log    logger.Logger
	parser PathsParser
	setter PathsSetter
}

func NewPathsGluer(log logger.Logger, parser PathsParser, setter PathsSetter) *PathsGluer {
	return &PathsGluer{
		log:    log,
		parser: parser,
		setter: setter,
	}
}

func (g *PathsGluer) Glue() error {
	err := g.parser.Parse()
	if err != nil {
		return fmt.Errorf("%w, failed parse paths, err: %w", ErrFailedGluePaths, err)
	}

	paths := g.parser.Paths()
	g.log.Debug("gluer paths", slog.Any("paths", paths))

	err = g.setter.SetPaths(paths)
	if err != nil {
		return fmt.Errorf("%w, failed set paths, err: %w", ErrFailedGluePaths, err)
	}

	return nil
}
