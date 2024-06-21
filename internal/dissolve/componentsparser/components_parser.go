package componentsparser

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/amidgo/node"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/model"
)

const componenetsTag = "components"

var (
	ErrNoComponentsTag   = errors.New("no 'components' tag")
	ErrEmptyComponents   = errors.New("empty components")
	ErrComponentNotFound = errors.New("component not found")
	ErrEncodeContentNode = errors.New("encode content node")
)

type ComponentsParser struct {
	headNode      *head.Head
	componentName string
	components    []*model.Item

	encoder node.EncoderTo
}

func NewComponentsParser(headNode *head.Head, componentName string, encoder node.EncoderTo) *ComponentsParser {
	return &ComponentsParser{
		headNode:      headNode,
		componentName: componentName,
		encoder:       encoder,
		components:    make([]*model.Item, 0),
	}
}

func (c *ComponentsParser) Parse() error {
	index := node.MapSearchByStringKey(c.headNode.Node(), componenetsTag)
	if index == -1 {
		return ErrNoComponentsTag
	}

	componentsNode := c.headNode.Node().Content()[index]
	if len(componentsNode.Content()) == 0 {
		return ErrEmptyComponents
	}

	specificComponentNode, err := c.searchSpecificComponent(componentsNode)
	if err != nil {
		return err
	}

	err = c.parseNodeItems(specificComponentNode)
	if err != nil {
		return err
	}

	return nil
}

func (c *ComponentsParser) searchSpecificComponent(nd node.Node) (node.Node, error) {
	iterator := node.MakeMapNodeIterator(nd.Content())
	for iterator.HasNext() {
		key, content := iterator.Next()
		if key.Value() == c.componentName {
			return content, nil
		}
	}

	return nil, fmt.Errorf("%q %w", c.componentName, ErrComponentNotFound)
}

func (c *ComponentsParser) parseNodeItems(nd node.Node) error {
	iterator := node.MakeMapNodeIterator(nd.Content())
	for iterator.HasNext() {
		key, content := iterator.Next()

		err := c.addItem(key, content)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ComponentsParser) addItem(key, content node.Node) error {
	itemContent := new(bytes.Buffer)
	item := model.Item{
		Name:    key.Value(),
		Content: itemContent,
	}

	err := c.encoder.EncodeTo(itemContent, content)
	if err != nil {
		return fmt.Errorf("%w, %w", ErrEncodeContentNode, err)
	}

	c.components = append(c.components, &item)

	return nil
}

func (c *ComponentsParser) Components() []*model.Item {
	return c.components
}
