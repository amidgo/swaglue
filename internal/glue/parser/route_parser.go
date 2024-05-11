package parser

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/amidgo/swaglue/internal/model"
	"github.com/amidgo/swaglue/pkg/httpmethod"
	"github.com/amidgo/swaglue/pkg/routes"
)

type RouteParser struct {
	basePackage string
	routes      []*model.Route
}

func NewRouteParser(basePackage string) *RouteParser {
	return &RouteParser{
		basePackage: basePackage,
		routes:      make([]*model.Route, 0),
	}
}

func (p *RouteParser) Parse() error {
	entries, err := os.ReadDir(p.basePackage)
	if err != nil {
		return fmt.Errorf("read base package, err: %w", err)
	}

	err = p.handleDirEntries(entries, p.basePackage)
	if err != nil {
		return fmt.Errorf("handle base package dir entries, err: %w", err)
	}

	return nil
}

func (p *RouteParser) handleDirEntries(entries []os.DirEntry, pathPrefix string) error {
	for _, entry := range entries {
		err := p.handleDirEntry(entry, pathPrefix)
		if err != nil {
			return fmt.Errorf("failed handle dir entry, err: %w", err)
		}
	}

	return nil
}

func (p *RouteParser) handleDirEntry(entry os.DirEntry, pathPrefix string) error {
	if isRouteEntry(entry) {
		return p.handleRouteEntry(entry, pathPrefix)
	}

	if entry.IsDir() {
		dirPath := path.Join(pathPrefix, entry.Name())

		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return fmt.Errorf("failed read dir by path %s, err: %w", dirPath, err)
		}

		return p.handleDirEntries(entries, dirPath)
	}

	return nil
}

func (p *RouteParser) handleRouteEntry(entry os.DirEntry, pathPrefix string) error {
	dirPath := path.Join(pathPrefix, entry.Name())

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed read route entry %s, err: %w", dirPath, err)
	}

	route := &model.Route{
		Name:    routeEntryName(entry),
		Methods: make([]*model.RouteMethod, 0),
	}

	for _, entry := range entries {
		method, ok := routeEntryMethod(entry, dirPath)
		if ok {
			route.Methods = append(route.Methods, method)
		}
	}

	p.routes = append(p.routes, route)

	return nil
}

func routeEntryMethod(entry os.DirEntry, pathPrefix string) (routeMethod *model.RouteMethod, ok bool) {
	if entry.IsDir() {
		return nil, false
	}

	method := strings.TrimSuffix(entry.Name(), ".yaml")
	if !httpmethod.Valid(method) {
		return nil, false
	}

	filePath := path.Join(pathPrefix, entry.Name())

	f, err := os.Open(filePath)
	if err != nil {
		return nil, false
	}

	return &model.RouteMethod{
		Method:  method,
		Content: f,
	}, true
}

func routeEntryName(entry fs.DirEntry) string {
	return strings.ReplaceAll(entry.Name(), routes.Separator, "/")
}

func isRouteEntry(entry fs.DirEntry) bool {
	if !entry.IsDir() {
		return false
	}

	return strings.HasPrefix(entry.Name(), routes.Separator)
}

func (p *RouteParser) Routes() []*model.Route {
	return p.routes
}
