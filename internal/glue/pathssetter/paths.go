package pathssetter

import (
	"errors"
	"fmt"
	"io"

	"github.com/amidgo/node"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/pkg/httpmethod"
)

const pathsTag = "paths"

var (
	ErrNoPathTag     = errors.New("'paths' tag not found")
	ErrWrongPathKind = errors.New("wrong path type, expected map")
	ErrInvalidRef    = errors.New("invalid ref value")
)

type HeadPathSetter struct {
	head    *head.Head
	decoder node.DecoderFrom
}

func New(head *head.Head, decoder node.DecoderFrom) *HeadPathSetter {
	return &HeadPathSetter{
		head:    head,
		decoder: decoder,
	}
}

func (h *HeadPathSetter) SetPaths(paths map[string]io.Reader) error {
	index, ok := h.head.SearchRootTag(pathsTag)
	if !ok {
		return ErrNoPathTag
	}

	pathNode := h.head.Content()[index]

	err := validatePathNode(pathNode)
	if err != nil {
		return err
	}

	pathChilds := pathNode.Content()
	for i := range pathChilds {
		route := PathsSetter{Node: pathChilds[i], Decoder: h.decoder}

		err := route.SetPathRefs(paths)
		if err != nil {
			return fmt.Errorf("handle paths, %w", err)
		}
	}

	return nil
}

func validatePathNode(pathNode node.Node) error {
	if pathNode.Type() != node.Map {
		return fmt.Errorf("%w, actual %s", ErrWrongPathKind, pathNode.Type())
	}

	return nil
}

type PathsSetter struct {
	Node    node.Node
	Decoder node.DecoderFrom
}

func (p *PathsSetter) SetPathRefs(paths map[string]io.Reader) error {
	if !p.isContentableNodeKind() {
		return nil
	}

	nodes := p.Node.Content()
	for i := 0; i < len(nodes); i += 2 {
		if !httpmethod.Valid(nodes[i].Value()) {
			continue
		}

		next := nodes[i+1]

		pathRouteMethod, ok := ParsePathRouteMethodFromContent(next.Content())
		if !ok {
			continue
		}

		ref := pathRouteMethod.RefValue()

		r, ok := paths[ref]
		if !ok {
			return fmt.Errorf("%w, ref %s not found", ErrInvalidRef, ref)
		}

		node, err := p.Decoder.DecodeFrom(r)
		if err != nil {
			return fmt.Errorf("%w, for ref %s, %w", head.ErrDecodeFile, ref, err)
		}

		nodes[i+1] = node
	}

	return nil
}

func (p *PathsSetter) isContentableNodeKind() bool {
	return p.Node.Type() == node.Map || p.Node.Type() == node.Array
}

type PathRouteMethod string

func ParsePathRouteMethodFromContent(content []node.Node) (method PathRouteMethod, ok bool) {
	const requireRouteMethodLength = 2
	if len(content) != requireRouteMethodLength {
		return "", false
	}

	if content[0].Value() != "$ref" || content[1].Value() == "" {
		return "", false
	}

	return PathRouteMethod(content[1].Value()), true
}

func (p PathRouteMethod) RefValue() string {
	return string(p)
}
