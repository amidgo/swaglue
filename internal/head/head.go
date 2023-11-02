package head

import (
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

var FailedDecodeFile = errors.New("failed decode file to .yaml")

type HeadFieldNotFoundError struct {
	Field string
}

func (e HeadFieldNotFoundError) Error() string {
	return fmt.Sprintf("field %s not found", e.Field)
}

type Head yaml.Node

func ParseHeadFromFile(filePath string) (Head, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return Head{}, err
	}
	return ParseHead(file)
}

func ParseHead(r io.Reader) (Head, error) {
	var node yaml.Node
	decoder := yaml.NewDecoder(r)
	err := decoder.Decode(&node)
	if err != nil {
		return Head{}, fmt.Errorf("failed unmarshal data, %w", err)
	}
	return Head(node), nil
}

func (h Head) SaveTo(w io.Writer) error {
	encoder := yaml.NewEncoder(w)
	err := encoder.Encode((*yaml.Node)(&h))
	encoder.Close()
	if err != nil {
		return err
	}
	return nil
}

func (h Head) SaveToFile(filePath string, flag int, mode os.FileMode) error {
	f, err := os.OpenFile(filePath, flag, mode)
	if err != nil {
		return err
	}
	return h.SaveTo(f)
}

func (h Head) SearchTag(tag string) *yaml.Node {
	var node *yaml.Node = (*yaml.Node)(&h)
	for {
		switch node.Kind {
		case yaml.DocumentNode:
			node = node.Content[0]
			continue
		case yaml.MappingNode:
			return searchInContent(node, func(n *yaml.Node) bool { return n.Value == tag })
		case yaml.SequenceNode:
			return searchInContent(node, func(n *yaml.Node) bool { return n.Value == tag })
		}
		return nil
	}
}

func searchInContent(node *yaml.Node, f func(*yaml.Node) bool) *yaml.Node {
	nodes := node.Content
	for i := range nodes {
		if i == len(nodes)-1 {
			return nil
		}
		if f(nodes[i]) {
			return nodes[i+1]
		}
	}
	return nil
}

type Decoder struct{ *yaml.Decoder }

func (d Decoder) Node() *yaml.Node {
	var n yaml.Node
	err := d.Decode(&n)
	if err != nil {
		return nil
	}
	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) == 0 {
			return nil
		}
		return n.Content[0]
	default:
		return &n
	}
}
