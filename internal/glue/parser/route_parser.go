package parser

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/amidgo/swaglue/internal/route"
	"github.com/amidgo/swaglue/pkg/httpmethod"
	"github.com/amidgo/swaglue/pkg/routes"
)

type RouteParser struct {
	basePackage string
	routes      []route.Route
}

func NewRouteParser(basePackage string) *RouteParser {
	return &RouteParser{
		basePackage: basePackage,
		routes:      make([]route.Route, 0),
	}
}

func (p *RouteParser) Parse() error {
	entries, err := os.ReadDir(p.basePackage)
	if err != nil {
		return fmt.Errorf("read base package, %w", err)
	}

	err = p.handleDirEntries(entries, p.basePackage)
	if err != nil {
		return fmt.Errorf("handle base package dir entries, %w", err)
	}

	return nil
}

func (p *RouteParser) handleDirEntries(entries []os.DirEntry, pathPrefix string) error {
	for _, entry := range entries {
		err := p.handleDirEntry(entry, pathPrefix)
		if err != nil {
			return fmt.Errorf("handle dir entry, %w", err)
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
			return fmt.Errorf("read dir by path %s, %w", dirPath, err)
		}

		return p.handleDirEntries(entries, dirPath)
	}

	return nil
}

func (p *RouteParser) handleRouteEntry(entry os.DirEntry, pathPrefix string) error {
	dirPath := path.Join(pathPrefix, entry.Name())

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("read route entry %s, %w", dirPath, err)
	}

	route := route.Route{
		Name:    routeEntryName(entry),
		Methods: make([]route.Method, 0),
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

func routeEntryMethod(entry os.DirEntry, pathPrefix string) (routeMethod route.Method, ok bool) {
	if entry.IsDir() {
		return routeMethod, false
	}

	method := strings.TrimSuffix(entry.Name(), ".yaml")
	if !httpmethod.Valid(method) {
		return routeMethod, false
	}

	filePath := path.Join(pathPrefix, entry.Name())

	f, err := os.Open(filePath)
	if err != nil {
		return routeMethod, false
	}

	return route.Method{
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

func (p *RouteParser) Routes() []route.Route {
	return p.routes
}
