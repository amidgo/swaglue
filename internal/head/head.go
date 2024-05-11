package head

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/amidgo/node"
)

var ErrDecodeFile = errors.New("decode file to .yaml")

const DefaultYamlIndent = 2

type FieldNotFoundError struct {
	Field string
}

func (e FieldNotFoundError) Error() string {
	return fmt.Sprintf("field %s not found", e.Field)
}

type Head struct {
	node.Node
}

func ParseHeadFromFile(filePath string, decoder node.Decoder) (*Head, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed parse head from file: failed open file: %w", err)
	}

	head, err := ParseHead(file, decoder)
	if err != nil {
		return nil, fmt.Errorf("failed parse head from file: failed parse head: %w", err)
	}

	return head, nil
}

func ParseHead(r io.Reader, decoder node.Decoder) (*Head, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed read from source, %w", err)
	}

	node, err := decoder.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal data, %w", err)
	}

	return &Head{
		Node: node,
	}, nil
}

func (h *Head) SaveToFile(filePath string, flag int, mode os.FileMode, encoder node.Encoder) error {
	f, err := os.OpenFile(filePath, flag, mode)
	if err != nil {
		return fmt.Errorf("failed open file, %w", err)
	}

	err = h.SaveTo(f, encoder)
	if err != nil {
		return fmt.Errorf("failed save head to file, %w", err)
	}

	return nil
}

func (h *Head) SaveTo(w io.Writer, encoder node.Encoder) error {
	data, err := encoder.Encode(h.Node)
	if err != nil {
		return fmt.Errorf("failed save head to writer, %w", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("failed write to dest, %w", err)
	}

	return nil
}

func (h *Head) SearchRootTag(tag string) (int, bool) {
	if h.Type() != node.Map {
		return 0, false
	}

	return searchInContent(h.Node, tag)
}

func searchInContent(node node.Node, tag string) (int, bool) {
	nodes := node.Content()
	for i := range nodes {
		if i == len(nodes)-1 {
			return 0, false
		}

		if nodes[i].Value() == tag {
			return i + 1, true
		}
	}

	return 0, false
}

var ErrWrongYamlDocumentNodeFormat = errors.New("wrong yaml document node format")

func DecodeNodeFrom(src io.Reader, decoder node.Decoder) (node.Node, error) {
	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("failed read, %w", err)
	}

	node, err := decoder.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed decode, %w", err)
	}

	return node, nil
}
