package head

import (
	"errors"
	"fmt"
	"io"
	"log"

	"gopkg.in/yaml.v3"
)

var (
	NoPathTag     = errors.New("'paths' tag not found")
	WrongPathKind = errors.New("wrong path kind, expected map")
	InvalidRef    = errors.New("invalid ref value")
)

type HttpMethod string

const (
	MethodGet     HttpMethod = "get"
	MethodHead    HttpMethod = "head"
	MethodPost    HttpMethod = "post"
	MethodPut     HttpMethod = "put"
	MethodPatch   HttpMethod = "patch"
	MethodDelete  HttpMethod = "delete"
	MethodConnect HttpMethod = "connect"
	MethodOptions HttpMethod = "options"
	MethodTrace   HttpMethod = "trace"
)

func (h HttpMethod) Valid() bool {
	for _, method := range []HttpMethod{
		MethodGet,
		MethodHead,
		MethodPost,
		MethodPut,
		MethodPatch,
		MethodDelete,
		MethodConnect,
		MethodOptions,
		MethodTrace,
	} {
		if h == method {
			return true
		}
	}
	return false
}

func (h Head) SetPaths(paths map[string]io.Reader) error {
	pathNode := h.SearchTag("paths")
	if err := validatePathNode(pathNode); err != nil {
		return err
	}
	log.Println(pathNode.Value)
	pathChilds := pathNode.Content
	for i := range pathChilds {
		route := PathRoute{Node: pathChilds[i]}
		err := route.Handle(paths)
		if err != nil {
			return err
		}
	}
	return nil
}

func validatePathNode(pathNode *yaml.Node) error {
	if pathNode == nil {
		return NoPathTag
	}
	if pathNode.Kind != yaml.MappingNode {
		return fmt.Errorf("%w, actual %s", WrongPathKind, pathNode.Tag)
	}
	return nil
}

type PathRoute struct {
	Node *yaml.Node
}

type PathRouteMethod [2]*yaml.Node

func (p PathRouteMethod) Valid() bool {
	return p[0].Value == "$ref" && p[1].Value != ""
}

func (p PathRouteMethod) Ref() string {
	return p[1].Value
}

func (p PathRoute) Handle(paths map[string]io.Reader) error {
	if !(p.Node.Kind == yaml.MappingNode || p.Node.Kind == yaml.SequenceNode) {
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
			return fmt.Errorf("%w, ref %s not found, line %d", InvalidRef, ref, next.Line)
		}
		node := Decoder{yaml.NewDecoder(r)}.Node()
		if node == nil {
			return fmt.Errorf("%w, for ref %s", FailedDecodeFile, ref)
		}
		nodes[i+1] = node
	}
	return nil
}
