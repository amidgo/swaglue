package routesappender

import (
	"errors"
	"fmt"

	"github.com/amidgo/node"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/model"
	"github.com/amidgo/swaglue/pkg/httpmethod"
)

type RouteExistsMethods struct {
	decoder           node.Decoder
	routeMethods      map[string]PathMethods
	routeContentNodes map[string]node.Node
}

func (m *RouteExistsMethods) RouteNameExists(routeName string) bool {
	_, ok := m.routeMethods[routeName]

	return ok
}

func (m *RouteExistsMethods) AnyRouteMethodExists(
	routeName string,
	methods []string,
) (existsMethods []string, exists bool) {
	pathNode, ok := m.routeMethods[routeName]
	if !ok {
		return
	}

	existsMethods = make([]string, 0)

	for _, m := range methods {
		if pathNode.MethodExists(m) {
			existsMethods = append(existsMethods, m)
		}
	}

	return existsMethods, len(existsMethods) != 0
}

var ErrContentNodeNotExists = errors.New("content node not exists")

func (m *RouteExistsMethods) RouteContentNode(routeName string) (RouteContentNode, error) {
	contentNode, ok := m.routeContentNodes[routeName]
	if !ok {
		return RouteContentNode{}, ErrContentNodeNotExists
	}

	return RouteContentNode{Node: contentNode, Decoder: m.decoder}, nil
}

func (m *RouteExistsMethods) ScanNode(nd node.Node) error {
	const routeKindCount = 2

	expectRoutesCount := len(nd.Content()) / routeKindCount

	m.routeMethods = make(map[string]PathMethods, expectRoutesCount)
	m.routeContentNodes = make(map[string]node.Node, expectRoutesCount)

	for i, route := range nd.Content() {
		if i == len(nd.Content())-1 {
			return nil
		}

		if route.Kind() != node.String {
			continue
		}

		err := m.scanRoute(nd.Content()[i], nd.Content()[i+1])
		if err != nil {
			return fmt.Errorf("scan route methods, %w", err)
		}
	}

	return nil
}

var (
	ErrInvalidKeyNode     = errors.New("invalid key node")
	ErrInvalidMappingNode = errors.New("invalid mapping node")
)

func (m *RouteExistsMethods) scanRoute(routeNameNode, routeContentNode node.Node) error {
	if routeNameNode.Kind() != node.String {
		return ErrInvalidKeyNode
	}

	if routeContentNode.Type() != node.Map {
		return ErrInvalidMappingNode
	}

	routeName := routeNameNode.Value()

	m.routeContentNodes[routeName] = routeContentNode

	err := m.scanRouteMethods(routeName, routeContentNode)
	if err != nil {
		return fmt.Errorf("scan route methods, err %w", err)
	}

	return nil
}

func (m *RouteExistsMethods) scanRouteMethods(routeName string, routeContentNode node.Node) error {
	for i, method := range routeContentNode.Content() {
		if i == len(routeContentNode.Content())-1 {
			return nil
		}

		if method.Kind() != node.String {
			continue
		}

		if !httpmethod.Valid(method.Value()) {
			return httpmethod.ErrInvalidMethod
		}

		m.addRouteMethod(routeName, method.Value())
	}

	return nil
}

func (m *RouteExistsMethods) addRouteMethod(routeName, methodName string) {
	pathNode, ok := m.routeMethods[routeName]
	if !ok {
		m.routeMethods[routeName] = MakePathMethods(methodName)

		return
	}

	pathNode.addMethod(methodName)
}

type PathMethods struct {
	methods map[string]struct{}
}

func MakePathMethods(initialMethod string) PathMethods {
	pathMethods := PathMethods{
		methods: make(map[string]struct{}, 1),
	}
	pathMethods.addMethod(initialMethod)

	return pathMethods
}

func (p PathMethods) MethodExists(method string) bool {
	_, ok := p.methods[method]

	return ok
}

func (p PathMethods) addMethod(method string) {
	p.methods[method] = struct{}{}
}

type RouteContentNode struct {
	node.Node
	Decoder node.Decoder
}

func (n *RouteContentNode) AddMethod(method *model.RouteMethod) error {
	methodNameNode := node.MakeStringNode(method.Method)

	methodContentNode, err := head.DecodeNodeFrom(method.Content, n.Decoder)
	if err != nil {
		return fmt.Errorf("decode yaml node, %w", err)
	}

	n.AppendNode(methodNameNode)
	n.AppendNode(methodContentNode)

	return nil
}
