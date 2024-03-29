package head

import (
	"errors"
	"fmt"
	"io"

	"github.com/amidgo/swaglue/pkg/httpmethod"
	"gopkg.in/yaml.v3"
)

const pathsTag = "paths"

var (
	ErrNoPathTag     = errors.New("'paths' tag not found")
	ErrWrongPathKind = errors.New("wrong path kind, expected map")
	ErrInvalidRef    = errors.New("invalid ref value")
)

func (h *Head) SetPaths(paths map[string]io.Reader) error {
	pathNode := h.SearchTag(pathsTag)

	err := validatePathNode(pathNode)
	if err != nil {
		return err
	}

	pathChilds := pathNode.Content
	for i := range pathChilds {
		route := PathsSetter{Node: pathChilds[i]}

		err := route.SetPathRefs(paths)
		if err != nil {
			return fmt.Errorf("failed handle paths, err: %w", err)
		}
	}

	return nil
}

func validatePathNode(pathNode *yaml.Node) error {
	if pathNode == nil {
		return ErrNoPathTag
	}

	if pathNode.Kind != yaml.MappingNode {
		return fmt.Errorf("%w, actual %s", ErrWrongPathKind, pathNode.Tag)
	}

	return nil
}

type PathsSetter struct {
	Node *yaml.Node
}

func (p *PathsSetter) SetPathRefs(paths map[string]io.Reader) error {
	if !p.isContentableNodeKind() {
		return nil
	}

	nodes := p.Node.Content
	for i := 0; i < len(nodes); i += 2 {
		if !httpmethod.Valid(nodes[i].Value) {
			continue
		}

		next := nodes[i+1]

		pathRouteMethod, ok := ParsePathRouteMethodFromContent(next.Content)
		if !ok {
			continue
		}

		ref := pathRouteMethod.RefValue()

		r, ok := paths[ref]
		if !ok {
			return fmt.Errorf("%w, ref %s not found, line %d", ErrInvalidRef, ref, next.Line)
		}

		node, err := DecodeYamlNode(yaml.NewDecoder(r))
		if err != nil {
			return fmt.Errorf("%w, for ref %s, err: %w", ErrFailedDecodeFile, ref, err)
		}

		nodes[i+1] = node
	}

	return nil
}

func (p *PathsSetter) isContentableNodeKind() bool {
	return p.Node.Kind == yaml.MappingNode || p.Node.Kind == yaml.SequenceNode
}

type PathRouteMethod string

func ParsePathRouteMethodFromContent(content []*yaml.Node) (method PathRouteMethod, ok bool) {
	const requireRouteMethodLength = 2
	if len(content) != requireRouteMethodLength {
		return
	}

	if content[0].Value != "$ref" || content[1].Value == "" {
		return
	}

	return PathRouteMethod(content[1].Value), true
}

func (p PathRouteMethod) RefValue() string {
	return string(p)
}
