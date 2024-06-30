package paths

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
	ErrNoPathTag          = errors.New("'paths' tag not found")
	ErrInvalidRef         = errors.New("invalid ref value")
	ErrInvalidNode        = errors.New("invalid node")
	ErrGetPathsFromSource = errors.New("get paths from source")
)

type Source interface {
	Paths() (map[string]io.Reader, error)
}

type MapSource map[string]io.Reader

func (m MapSource) Paths() (map[string]io.Reader, error) {
	return m, nil
}

type HeadPathSetter struct {
	head        *head.Head
	decoder     node.DecoderFrom
	validate    node.Validate
	pathsSource Source
}

func New(head *head.Head, decoder node.DecoderFrom, pathsSource Source) *HeadPathSetter {
	return &HeadPathSetter{
		head:        head,
		decoder:     decoder,
		validate:    node.NewKindValidate(node.Map),
		pathsSource: pathsSource,
	}
}

func (h *HeadPathSetter) SetPaths(paths map[string]io.Reader) error {
	index := node.MapSearchByStringKey(h.head.Node(), pathsTag)
	if index == -1 {
		return ErrNoPathTag
	}

	pathNode := h.head.Node().Content()[index]

	err := h.validate.Validate(pathNode)
	if err != nil {
		return errors.Join(ErrInvalidNode, err)
	}

	pathChilds := pathNode.Content()
	for i := range pathChilds {
		route := Setter{Node: pathChilds[i], Decoder: h.decoder}

		err := route.SetPathRefs(paths)
		if err != nil {
			return fmt.Errorf("handle paths, %w", err)
		}
	}

	return nil
}

func (h *HeadPathSetter) KeyValue(key, value node.Node) (resKey, resValue node.Node, err error) {
	if !node.StringEquals(key, pathsTag) {
		return key, value, nil
	}

	paths, err := h.pathsSource.Paths()
	if err != nil {
		return nil, nil, errors.Join(ErrGetPathsFromSource, err)
	}

	err = h.validate.Validate(value)
	if err != nil {
		return nil, nil, errors.Join(ErrInvalidNode, err)
	}

	pathChilds := value.Content()
	for i := range pathChilds {
		route := Setter{Node: pathChilds[i], Decoder: h.decoder}

		err := route.SetPathRefs(paths)
		if err != nil {
			return nil, nil, fmt.Errorf("handle paths, %w", err)
		}
	}

	return key, value, nil
}

type Setter struct {
	Node    node.Node
	Decoder node.DecoderFrom
}

func (p *Setter) SetPathRefs(paths map[string]io.Reader) error {
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

func (p *Setter) isContentableNodeKind() bool {
	return p.Node.Kind() == node.Map || p.Node.Kind() == node.Array
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
