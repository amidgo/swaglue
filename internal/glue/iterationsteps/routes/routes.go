package routes

import (
	"errors"

	"github.com/amidgo/node"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/route"
)

const pathsTag = "paths"

var (
	ErrInvalidTypePathsNode = errors.New("invalid type of paths node")
	ErrAppendRoutesToNode   = errors.New("append routes to node")
	ErrNoPathTag            = errors.New("no path tag")
)

type IterationStep struct {
	head        *head.Head
	decoder     node.DecoderFrom
	encoder     node.EncoderTo
	routeSource route.Source
}

func New(
	head *head.Head,
	decoder node.DecoderFrom,
	encoder node.EncoderTo,
	routeSource route.Source,
) *IterationStep {
	return &IterationStep{
		head:        head,
		decoder:     decoder,
		encoder:     encoder,
		routeSource: routeSource,
	}
}

func (h *IterationStep) AppendRoutes(routes []route.Route) (err error) {
	index := node.MapSearchByStringKey(h.head.Node(), pathsTag)
	if index == -1 {
		return ErrNoPathTag
	}

	pathsNode, err := NewPathsNode(
		route.NewRemoveDuplicatesSource(
			route.NewJoinSource(
				route.NewNodeRouteSource(h.head.Node().Content()[index], h.encoder),
				route.SliceSource(routes),
			),
		),
		h.decoder,
	).Node()
	if err != nil {
		return err
	}

	h.head.Node().Content()[index] = pathsNode

	return nil
}

func (h *IterationStep) KeyValue(key, value node.Node) (resKey, resValue node.Node, err error) {
	if key.Kind() != node.String || key.Value() != pathsTag {
		return key, value, nil
	}

	pathsNode, err := NewPathsNode(
		route.NewRemoveDuplicatesSource(
			route.NewJoinSource(
				route.NewNodeRouteSource(value, h.encoder),
				h.routeSource,
			),
		),
		h.decoder,
	).Node()
	if err != nil {
		return nil, nil, err
	}

	return key, pathsNode, nil
}
