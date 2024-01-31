package head

import (
	"errors"
	"fmt"

	"github.com/amidgo/swaglue/internal/glue/model"
	"github.com/amidgo/swaglue/pkg/httpmethod"
	"gopkg.in/yaml.v3"
)

type RouteExistsMethods struct {
	routeMethods      map[string]PathMethods
	routeContentNodes map[string]*yaml.Node
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

func (m *RouteExistsMethods) RouteContentNode(routeName string) (route RouteContentNode, err error) {
	contentNode, ok := m.routeContentNodes[routeName]
	if !ok {
		return route, ErrContentNodeNotExists
	}

	return RouteContentNode{Node: contentNode}, nil
}

func (m *RouteExistsMethods) ScanNode(node *yaml.Node) error {
	const routeKindCount = 2

	expectRoutesCount := len(node.Content) / routeKindCount

	m.routeMethods = make(map[string]PathMethods, expectRoutesCount)
	m.routeContentNodes = make(map[string]*yaml.Node, expectRoutesCount)

	for i, route := range node.Content {
		if i == len(node.Content)-1 {
			return nil
		}

		if route.Kind != yaml.ScalarNode {
			continue
		}

		err := m.scanRoute(node.Content[i], node.Content[i+1])
		if err != nil {
			return fmt.Errorf("failed scan route methods, err: %w, line %d", err, route.Line)
		}
	}

	return nil
}

var (
	ErrInvalidKeyNodeKind     = errors.New("invalid key node kind")
	ErrInvalidMappingNodeKind = errors.New("invalid mapping node kind")
)

func (m *RouteExistsMethods) scanRoute(routeNameNode *yaml.Node, routeContentNode *yaml.Node) error {
	if routeNameNode.Kind != yaml.ScalarNode {
		return ErrInvalidKeyNodeKind
	}

	if routeContentNode.Kind != yaml.MappingNode {
		return ErrInvalidMappingNodeKind
	}

	routeName := routeNameNode.Value

	m.routeContentNodes[routeName] = routeContentNode

	err := m.scanRouteMethods(routeName, routeContentNode)
	if err != nil {
		return fmt.Errorf("failed scan route methods, err %w, line %d", err, routeContentNode.Line)
	}

	return nil
}

func (m *RouteExistsMethods) scanRouteMethods(routeName string, routeContentNode *yaml.Node) error {
	for i, method := range routeContentNode.Content {
		if i == len(routeContentNode.Content)-1 {
			return nil
		}

		if method.Kind != yaml.ScalarNode {
			continue
		}

		if !httpmethod.Valid(method.Value) {
			return fmt.Errorf("%w, line %d", httpmethod.ErrInvalidMethod, method.Line)
		}

		m.addRouteMethod(routeName, method.Value)
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
	*yaml.Node
}

func (n RouteContentNode) AddMethod(method *model.RouteMethod) error {
	methodNameNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: method.Method,
	}

	methodContentNode, err := DecodeYamlNode(yaml.NewDecoder(method.Content))
	if err != nil {
		return fmt.Errorf("failed decode yaml node, err: %w", err)
	}

	n.Node.Content = append(n.Node.Content, methodNameNode, methodContentNode)

	return nil
}
