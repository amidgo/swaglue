package head

import (
	"errors"
	"fmt"
	"strings"

	"github.com/amidgo/swaglue/internal/glue/model"
	"gopkg.in/yaml.v3"
)

var ErrInvalidTypePathsNode = errors.New("invalid type of paths node")

func (h *Head) AppendRoutes(routes []*model.Route) (err error) {
	pathNode := h.SearchTag(pathsTag)
	if pathNode == nil {
		return ErrNoPathTag
	}

	if len(pathNode.Content) == 0 {
		pathNode.Kind = yaml.MappingNode
		pathNode.Tag = ""
	}

	var existsRouteMethods RouteExistsMethods

	err = existsRouteMethods.ScanNode(pathNode)
	if err != nil {
		return err
	}

	routePathNode := RouteAppender{
		ExistsRouteMethods: existsRouteMethods,
		Node:               pathNode,
	}

	err = routePathNode.AppendRoutes(routes)
	if err != nil {
		return fmt.Errorf("failed append routes to path node, err: %w", err)
	}

	return nil
}

type RouteAppender struct {
	ExistsRouteMethods RouteExistsMethods
	Node               *yaml.Node
}

func (r *RouteAppender) AppendRoutes(routes []*model.Route) error {
	for _, route := range routes {
		if r.ExistsRouteMethods.RouteNameExists(route.Name) {
			err := r.appendRouteMethodsToExistsRoute(route)
			if err != nil {
				return fmt.Errorf("failed append route to exists routes, err: %w", err)
			}
		} else {
			err := r.appendRoute(route)
			if err != nil {
				return fmt.Errorf("failed append route, err: %w", err)
			}
		}
	}

	return nil
}

var ErrDetectRouteMethodsConflicts = errors.New("detect route methods conflicts")

func (r *RouteAppender) appendRouteMethodsToExistsRoute(route *model.Route) error {
	err := r.detectRouteMethodsConflits(route)
	if err != nil {
		return err
	}

	err = r.addRouteMethods(route)
	if err != nil {
		return fmt.Errorf("failed append route methods, err: %w", err)
	}

	return nil
}

func (r *RouteAppender) detectRouteMethodsConflits(route *model.Route) error {
	methods := make([]string, 0, len(route.Methods))
	for _, routeMethod := range route.Methods {
		methods = append(methods, routeMethod.Method)
	}

	if existsMethods, exists := r.ExistsRouteMethods.AnyRouteMethodExists(route.Name, methods); exists {
		methodsString := strings.Join(existsMethods, ",")

		return fmt.Errorf(
			"%w, route %s, route methods %s already exists in head file",
			ErrDetectRouteMethodsConflicts,
			route.Name,
			methodsString,
		)
	}

	return nil
}

func (r *RouteAppender) addRouteMethods(route *model.Route) error {
	routeMethodContentNode, err := r.ExistsRouteMethods.RouteContentNode(route.Name)
	if err != nil {
		return err
	}

	for _, method := range route.Methods {
		err := routeMethodContentNode.AddMethod(method)
		if err != nil {
			return fmt.Errorf("failed add method, err: %w", err)
		}
	}

	return nil
}

func (r *RouteAppender) appendRoute(route *model.Route) error {
	pathNameRoute := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: route.Name,
	}
	pathMethodsRoute := RouteContentNode{
		Node: &yaml.Node{
			Kind:    yaml.MappingNode,
			Content: make([]*yaml.Node, 0, len(route.Methods)),
		},
	}

	for _, method := range route.Methods {
		err := pathMethodsRoute.AddMethod(method)
		if err != nil {
			return fmt.Errorf("failed add method, err: %w", err)
		}
	}

	r.Node.Content = append(r.Node.Content, pathNameRoute, pathMethodsRoute.Node)

	return nil
}
