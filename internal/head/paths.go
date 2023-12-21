package head

import (
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

const path_tag = "paths"

var (
	ErrNoPathTag     = errors.New("'paths' tag not found")
	ErrWrongPathKind = errors.New("wrong path kind, expected map")
	ErrInvalidRef    = errors.New("invalid ref value")
)

func (h Head) SetPaths(paths map[string]io.Reader) error {
	pathNode := h.SearchTag(path_tag)
	if err := validatePathNode(pathNode); err != nil {
		return err
	}
	pathChilds := pathNode.Content
	for i := range pathChilds {
		route := PathRoute{Node: pathChilds[i]}
		err := route.Handle(paths)
		if err != nil {
			return fmt.Errorf("failed handle paths")
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

type PathRoute struct {
	Node *yaml.Node
}

func (p *PathRoute) Handle(paths map[string]io.Reader) error {
	if !p.isContentableNodeKind() {
		return nil
	}
	nodes := p.Node.Content
	for i := range nodes {
		if i == len(nodes)-1 {
			return nil
		}
		if !HttpMethod(nodes[i].Value).Valid() {
			continue
		}
		next := nodes[i+1]
		pathRouteMethod := PathRouteMethod(next.Content)
		if !pathRouteMethod.Valid() {
			continue
		}
		ref := pathRouteMethod.Ref()
		r, ok := paths[ref]
		if !ok {
			return fmt.Errorf("%w, ref %s not found, line %d", ErrInvalidRef, ref, next.Line)
		}
		node := DecodeYamlNode(yaml.NewDecoder(r))
		if node == nil {
			return fmt.Errorf("%w, for ref %s", ErrFailedDecodeFile, ref)
		}
		nodes[i+1] = node
	}
	return nil
}

func (p *PathRoute) isContentableNodeKind() bool {
	return p.Node.Kind == yaml.MappingNode || p.Node.Kind == yaml.SequenceNode
}

type PathRouteMethod []*yaml.Node

func (p PathRouteMethod) Valid() bool {
	if len(p) != 2 {
		return false
	}
	return p[0].Value == "$ref" && p[1].Value != ""
}

func (p PathRouteMethod) Ref() string {
	return p[1].Value
}
