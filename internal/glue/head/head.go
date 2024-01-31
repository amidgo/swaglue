package head

import (
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

var ErrFailedDecodeFile = errors.New("failed decode file to .yaml")

type FieldNotFoundError struct {
	Field string
}

func (e FieldNotFoundError) Error() string {
	return fmt.Sprintf("field %s not found", e.Field)
}

const defaultEncoderIndent = 2

type Head struct {
	encoderIndent int
	headNode      yaml.Node
}

func ParseHeadFromFile(filePath string) (*Head, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed parse head from file: failed open file: %w", err)
	}

	head, err := ParseHead(file)
	if err != nil {
		return nil, fmt.Errorf("failed parse head from file: failed parse head: %w", err)
	}

	return head, nil
}

func ParseHead(r io.Reader) (*Head, error) {
	var node yaml.Node

	decoder := yaml.NewDecoder(r)

	err := decoder.Decode(&node)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal data, %w", err)
	}

	return &Head{
		headNode:      node,
		encoderIndent: defaultEncoderIndent,
	}, nil
}

func (h *Head) SaveToFile(filePath string, flag int, mode os.FileMode) error {
	f, err := os.OpenFile(filePath, flag, mode)
	if err != nil {
		return fmt.Errorf("failed open file, %w", err)
	}

	err = h.SaveTo(f)
	if err != nil {
		return fmt.Errorf("failed save head to file, %w", err)
	}

	return nil
}

func (h *Head) SaveTo(w io.Writer) error {
	encoder := yaml.NewEncoder(w)
	defer encoder.Close()

	encoder.SetIndent(h.encoderIndent)

	err := encoder.Encode(&h.headNode)
	if err != nil {
		return fmt.Errorf("failed save head to writer, %w", err)
	}

	return nil
}

func (h *Head) SearchTag(tag string) *yaml.Node {
	node := &h.headNode

	for {
		switch node.Kind {
		case yaml.DocumentNode:
			node = node.Content[0]

			continue
		case yaml.MappingNode:
			return searchInContent(node, tag)
		case yaml.SequenceNode:
			return searchInContent(node, tag)
		}

		return nil
	}
}

func searchInContent(node *yaml.Node, tag string) *yaml.Node {
	nodes := node.Content
	for i := range nodes {
		if i == len(nodes)-1 {
			return nil
		}

		if nodes[i].Value == tag {
			return nodes[i+1]
		}
	}

	return nil
}

var ErrWrongYamlDocumentNodeFormat = errors.New("wrong yaml document node format")

func DecodeYamlNode(dec *yaml.Decoder) (*yaml.Node, error) {
	var n yaml.Node

	err := dec.Decode(&n)
	if err != nil {
		return nil, fmt.Errorf("failed decode yaml node, err: %w", err)
	}

	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) == 0 {
			return nil, ErrWrongYamlDocumentNodeFormat
		}

		return n.Content[0], nil
	default:
		return &n, nil
	}
}
