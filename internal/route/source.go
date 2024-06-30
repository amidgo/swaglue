package route

import (
	"bytes"
	"errors"
	"slices"

	"github.com/amidgo/node"
	"github.com/amidgo/swaglue/pkg/httpmethod"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Source interface {
	Routes() ([]Route, error)
}

type SliceSource []Route

func (s SliceSource) Routes() ([]Route, error) {
	return s, nil
}

var ErrJoinRouteSources = errors.New("join route sources")

type JoinSource struct {
	routesSources []Source
}

func NewJoinSource(routeSources ...Source) JoinSource {
	return JoinSource{routesSources: routeSources}
}

func (j JoinSource) Routes() ([]Route, error) {
	totalRoutes := make([][]Route, 0, len(j.routesSources))
	totalSize := 0

	for _, source := range j.routesSources {
		routes, err := source.Routes()
		if err != nil {
			return nil, errors.Join(ErrJoinRouteSources, err)
		}

		totalSize += len(routes)
		totalRoutes = append(totalRoutes, routes)
	}

	joinedRoutes := make([]Route, totalSize)

	offset := 0
	for _, routes := range totalRoutes {
		copy(joinedRoutes[offset:len(routes)+offset], routes)
		offset += len(routes)
	}

	return joinedRoutes, nil
}

var ErrMergeRouteDuplicates = errors.New("merge route duplicates")

type RemoveDuplicatesSource struct {
	source Source
}

func NewRemoveDuplicatesSource(source Source) RemoveDuplicatesSource {
	return RemoveDuplicatesSource{source: source}
}

func (m RemoveDuplicatesSource) Routes() ([]Route, error) {
	routes, err := m.source.Routes()
	if err != nil {
		return nil, errors.Join(ErrMergeRouteDuplicates, err)
	}

	existsRoutes := orderedmap.New[string, Route](
		orderedmap.WithCapacity[string, Route](len(routes)),
	)

	for _, rt := range routes {
		existsRoute, ok := existsRoutes.Get(rt.Name)
		if !ok {
			existsRoutes.Set(rt.Name, rt)

			continue
		}

		existsRoutes.Set(
			existsRoute.Name,
			Route{
				Name:    existsRoute.Name,
				Methods: mergeRouteMethods(existsRoute, rt),
			},
		)
	}

	mergedRoutes := make([]Route, 0, existsRoutes.Len())

	for pair := existsRoutes.Oldest(); pair != nil; pair = pair.Next() {
		mergedRoutes = append(mergedRoutes, pair.Value)
	}

	return mergedRoutes, nil
}

func mergeRouteMethods(sourceRoute, inputRoute Route) []Method {
	sourceRouteMethods := routeMethods(sourceRoute)

	for _, method := range inputRoute.Methods {
		if !slices.Contains(sourceRouteMethods, method.Method) {
			sourceRoute.Methods = append(sourceRoute.Methods, method)
		}
	}

	return sourceRoute.Methods
}

func routeMethods(route Route) []string {
	methods := make([]string, 0, len(route.Methods))

	for _, method := range route.Methods {
		methods = append(methods, method.Method)
	}

	return methods
}

var (
	ErrEncodeNode      = errors.New("encode node")
	ErrInvalidNodeKind = errors.New("invalid node kind")
)

type NodeRouteSource struct {
	nd       node.Node
	encoder  node.EncoderTo
	validate node.Validate
}

func NewNodeRouteSource(
	nd node.Node,
	encoder node.EncoderTo,
) NodeRouteSource {
	return NodeRouteSource{
		nd:       nd,
		encoder:  encoder,
		validate: node.NewKindValidate(node.Map, node.Empty),
	}
}

func (n NodeRouteSource) Routes() ([]Route, error) {
	err := n.validate.Validate(n.nd)
	if err != nil {
		return nil, err
	}

	routes := make([]Route, 0)

	iter := node.MakeMapNodeIterator(n.nd.Content())
	for iter.HasNext() {
		path, methods := iter.Next()

		routeMethods, err := n.routeExistsMethods(methods)
		if err != nil {
			return nil, err
		}

		route := Route{
			Name:    path.Value(),
			Methods: routeMethods,
		}

		routes = append(routes, route)
	}

	return routes, nil
}

func (n NodeRouteSource) routeExistsMethods(methods node.Node) ([]Method, error) {
	routeMethods := make([]Method, 0)

	iter := node.MakeMapNodeIterator(methods.Content())
	for iter.HasNext() {
		key, content := iter.Next()
		if !httpmethod.Valid(key.Value()) {
			return nil, httpmethod.ErrInvalidMethod
		}

		buf := &bytes.Buffer{}

		err := n.encoder.EncodeTo(buf, content)
		if err != nil {
			return nil, errors.Join(ErrEncodeNode, err)
		}

		method := Method{
			Method:  key.Value(),
			Content: buf,
		}

		routeMethods = append(routeMethods, method)
	}

	return routeMethods, nil
}
