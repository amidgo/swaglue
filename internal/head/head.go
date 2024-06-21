package head

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/amidgo/node"
)

var ErrDecodeFile = errors.New("decode file to .yaml")

const DefaultIndent = 2

type FieldNotFoundError struct {
	Field string
}

func (e FieldNotFoundError) Error() string {
	return fmt.Sprintf("field %s not found", e.Field)
}

type Head struct {
	nd node.MapNode
}

func (h *Head) Node() node.MapNode {
	return h.nd
}

func ParseHeadFromFile(filePath string, decoder node.DecoderFrom) (*Head, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("parse head from file: open file: %w", err)
	}

	head, err := ParseHead(file, decoder)
	if err != nil {
		return nil, fmt.Errorf("parse head from file: parse head: %w", err)
	}

	return head, nil
}

var ErrWrongHeadKind = errors.New("wrong head kind")

func ParseHead(r io.Reader, decoder node.DecoderFrom) (*Head, error) {
	nd, err := decoder.DecodeFrom(r)
	if err != nil {
		return nil, fmt.Errorf("unmarshal data, %w", err)
	}

	if nd.Kind() != node.Map {
		return nil, ErrWrongHeadKind
	}

	return &Head{
		nd: node.MakeMapNodeWithContent(nd.Content()...),
	}, nil
}

func (h *Head) SaveToFile(filePath string, flag int, mode os.FileMode, encoder node.EncoderTo) error {
	f, err := os.OpenFile(filePath, flag, mode)
	if err != nil {
		return fmt.Errorf("open file, %w", err)
	}

	err = h.SaveTo(f, encoder)
	if err != nil {
		return fmt.Errorf("save head to file, %w", err)
	}

	return nil
}

func (h *Head) SaveTo(w io.Writer, encoder node.EncoderTo) error {
	err := encoder.EncodeTo(w, h.nd)
	if err != nil {
		return fmt.Errorf("save head to writer, %w", err)
	}

	return nil
}
