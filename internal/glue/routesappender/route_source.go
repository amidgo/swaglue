package routesappender

import (
	"errors"

	"github.com/amidgo/node"

	"github.com/amidgo/swaglue/internal/route"
)

type PathsNode struct {
	routeSource route.Source
	decoder     node.DecoderFrom
}

func NewPathsNode(routeSource route.Source, decoder node.DecoderFrom) PathsNode {
	return PathsNode{routeSource: routeSource, decoder: decoder}
}

var ErrGetRoutesFromSource = errors.New("get node from source")

func (p PathsNode) Node() (node.Node, error) {
	routes, err := p.routeSource.Routes()
	if err != nil {
		return nil, errors.Join(ErrGetRoutesFromSource, err)
	}

	pathsNode := node.MakeMapNode()

	for _, route := range routes {
		keyNode := node.MakeStringNode(route.Name)

		routeMethodsNode, err := p.routeMethodsNode(route.Methods)
		if err != nil {
			return nil, err
		}

		pathsNode = node.MapAppend(pathsNode, keyNode, routeMethodsNode)
	}

	return pathsNode, nil
}

var ErrDecode = errors.New("err decode")

func (p PathsNode) routeMethodsNode(methods []route.Method) (node.Node, error) {
	routeMethodsNode := node.MakeMapNode()

	for _, method := range methods {
		keyNode := node.MakeStringNode(method.Method)

		contentNode, err := p.decoder.DecodeFrom(method.Content)
		if err != nil {
			return nil, errors.Join(ErrDecode, err)
		}

		routeMethodsNode = node.MapAppend(routeMethodsNode, keyNode, contentNode)
	}

	return routeMethodsNode, nil
}
