package routesappender

import (
	"errors"
	"fmt"

	"github.com/amidgo/node"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/model"
)

const pathsTag = "paths"

var (
	ErrInvalidTypePathsNode = errors.New("invalid type of paths node")
	ErrAppendRoutesToNode   = errors.New("append routes to node")
	ErrNoPathTag            = errors.New("no path tag")
)

type HeadRoutesAppender struct {
	head    *head.Head
	decoder node.Decoder
}

func New(head *head.Head, decoder node.Decoder) *HeadRoutesAppender {
	return &HeadRoutesAppender{
		head:    head,
		decoder: decoder,
	}
}

func (h *HeadRoutesAppender) AppendRoutes(routes []*model.Route) (err error) {
	index, ok := h.head.SearchRootTag(pathsTag)
	if !ok {
		return ErrNoPathTag
	}

	pathNode := h.head.Content()[index]

	if len(pathNode.Content()) == 0 {
		pathNode = node.MakeMapNode()
		h.head.Content()[index] = pathNode
	}

	existsRouteMethods := RouteExistsMethods{decoder: h.decoder}

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
		return errors.Join(ErrAppendRoutesToNode, err)
	}

	return nil
}

type RouteAppender struct {
	ExistsRouteMethods RouteExistsMethods
	node.Node
}

func (r *RouteAppender) AppendRoutes(routes []*model.Route) error {
	for _, route := range routes {
		if r.ExistsRouteMethods.RouteNameExists(route.Name) {
			err := r.appendRouteMethodsToExistsRoute(route)
			if err != nil {
				return fmt.Errorf("append route to exists routes, %w", err)
			}
		} else {
			err := r.appendRoute(route)
			if err != nil {
				return fmt.Errorf("append route, %w", err)
			}
		}
	}

	return nil
}

var ErrDetectRouteMethodsConflicts = errors.New("detect route methods conflicts")

func (r *RouteAppender) appendRouteMethodsToExistsRoute(route *model.Route) error {
	r.ExistsRouteMethods.FilterExistsMethods(route)

	err := r.addRouteMethods(route)
	if err != nil {
		return fmt.Errorf("append route methods, %w", err)
	}

	r.ExistsRouteMethods.AddRouteMethods(route)

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
			return fmt.Errorf("add method, %w", err)
		}
	}

	return nil
}

func (r *RouteAppender) appendRoute(route *model.Route) error {
	pathNameRoute := node.MakeStringNode(route.Name)
	pathMethodsRoute := RouteContentNode{
		Decoder: r.ExistsRouteMethods.decoder,
		Node:    node.MakeMapNode(),
	}

	for _, method := range route.Methods {
		err := pathMethodsRoute.AddMethod(method)
		if err != nil {
			return fmt.Errorf("add method, %w", err)
		}
	}

	r.AppendNode(pathNameRoute)
	r.AppendNode(pathMethodsRoute.Node)

	r.ExistsRouteMethods.ScanRoute(pathNameRoute, pathMethodsRoute.Node)

	return nil
}
