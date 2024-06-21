package routesappender

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

type HeadRoutesAppender struct {
	head    *head.Head
	decoder node.DecoderFrom
	encoder node.EncoderTo
}

func New(head *head.Head, decoder node.DecoderFrom, encoder node.EncoderTo) *HeadRoutesAppender {
	return &HeadRoutesAppender{
		head:    head,
		decoder: decoder,
		encoder: encoder,
	}
}

func (h *HeadRoutesAppender) AppendRoutes(routes []route.Route) (err error) {
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
