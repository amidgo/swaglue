package componentsappender

import (
	"errors"
	"fmt"

	"github.com/amidgo/node"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/model"
)

const componentsTag = "components"

var ErrNoComponentsTag = errors.New("'components' tag not found")

type HeadComponentAppender struct {
	head    *head.Head
	decoder node.DecoderFrom
}

func New(head *head.Head, decoder node.DecoderFrom) *HeadComponentAppender {
	return &HeadComponentAppender{
		head:    head,
		decoder: decoder,
	}
}

func (h *HeadComponentAppender) AppendComponent(componentName string, componentItems []model.Item) error {
	index := node.MapSearchByStringKey(h.head.Node(), componentsTag)
	if index == -1 {
		return ErrNoComponentsTag
	}

	componentTag := h.head.Node().Content()[index]

	appender := ComponentAppender{
		Decoder:        h.decoder,
		nd:             node.MakeMapNode(componentTag.Content()...),
		ComponentName:  componentName,
		ComponentItems: componentItems,
	}

	nd, err := appender.Node()
	if err != nil {
		return fmt.Errorf("append component, %w", err)
	}

	h.head.Node().Content()[index] = nd

	return nil
}

type ComponentAppender struct {
	Decoder        node.DecoderFrom
	ComponentName  string
	ComponentItems []model.Item
	nd             node.MapNode
}

func (a *ComponentAppender) Node() (nd node.Node, err error) {
	itemNode, index := a.searchComponentNode()
	if index != -1 {
		nd, err = a.appendExistComponent(itemNode)
		if err != nil {
			return nil, err
		}

		a.nd.Content()[index] = nd

		return a.nd, nil
	}

	err = a.appendNewComponent()
	if err != nil {
		return nil, err
	}

	return a.nd, nil
}

func (a *ComponentAppender) searchComponentNode() (itemNode node.MapNode, index int) {
	iter := node.NewIndexedIterator(node.MakeMapNodeIterator(a.nd.Content()))

	for iter.HasNext() {
		key, value := iter.Next()

		if key.Value() != a.ComponentName {
			continue
		}

		if iter.Index() != len(a.nd.Content())-1 {
			return node.MakeMapNode(value.Content()...), iter.Index() + 1
		}

		a.nd = node.MapRound(a.nd, value)

		return node.MakeMapNode(value.Content()...), iter.Index() + 1
	}

	return node.MakeMapNode(), -1
}

func (a *ComponentAppender) appendNewComponent() error {
	keyNode, contentNode, err := new(ComponentNodeBuilder).
		SetDecoder(a.Decoder).
		SetName(a.ComponentName).
		SetItems(a.ComponentItems).
		Build()
	if err != nil {
		return err
	}

	a.nd = node.MapAppend(a.nd, keyNode, contentNode)

	return nil
}

func (a *ComponentAppender) appendExistComponent(itemNode node.MapNode) (node.Node, error) {
	err := a.validateItemNode(itemNode)
	if err != nil {
		return nil, err
	}

	_, contentNode, err := new(ComponentNodeBuilder).
		SetDecoder(a.Decoder).
		SetItemNode(itemNode).
		AppendItems(a.ComponentItems).
		Build()
	if err != nil {
		return nil, err
	}

	return contentNode, nil
}

var ErrComponentItemNameExists = errors.New("component item name exists")

func (a *ComponentAppender) validateItemNode(itemNode node.Node) error {
	existsNames := make(map[string]struct{}, len(a.ComponentItems))
	for _, item := range a.ComponentItems {
		existsNames[item.Name] = struct{}{}
	}

	for i := 0; i < len(itemNode.Content()); i += 2 {
		nd := itemNode.Content()[i]
		if nd.Kind() != node.String {
			continue
		}

		_, exists := existsNames[nd.Value()]
		if exists {
			return fmt.Errorf("%w, name %s", ErrComponentItemNameExists, nd.Value())
		}
	}

	return nil
}

type ComponentNodeBuilder struct {
	decoder   node.DecoderFrom
	namedNode node.Node
	itemsNode node.MapNode
	err       error
}

func (c *ComponentNodeBuilder) SetDecoder(dec node.DecoderFrom) *ComponentNodeBuilder {
	if c.err != nil {
		return c
	}

	c.decoder = dec

	return c
}

func (c *ComponentNodeBuilder) SetName(name string) *ComponentNodeBuilder {
	if c.err != nil {
		return c
	}

	c.namedNode = node.MakeStringNode(name)

	return c
}

func (c *ComponentNodeBuilder) SetItemNode(itemsNode node.MapNode) *ComponentNodeBuilder {
	if c.err != nil {
		return c
	}

	c.itemsNode = itemsNode

	return c
}

func (c *ComponentNodeBuilder) AppendItems(items []model.Item) *ComponentNodeBuilder {
	if c.err != nil {
		return c
	}

	c.appendComponents(items)

	return c
}

func (c *ComponentNodeBuilder) SetItems(items []model.Item) *ComponentNodeBuilder {
	if c.err != nil {
		return c
	}
	content := make([]node.Node, 0, len(items))

	c.itemsNode = node.MakeMapNode(content...)
	c.appendComponents(items)

	return c
}

func (c *ComponentNodeBuilder) appendComponents(components []model.Item) {
	for _, component := range components {
		err := c.appendComponent(component)
		if err != nil {
			c.err = err

			return
		}
	}
}

func (c *ComponentNodeBuilder) appendComponent(component model.Item) error {
	namedNode := node.MakeStringNode(component.Name)

	itemNode, err := c.decoder.DecodeFrom(component.Content)
	if err != nil {
		return fmt.Errorf("%w, for %s, %w", head.ErrDecodeFile, component.Name, err)
	}

	c.itemsNode = node.MapAppend(c.itemsNode, namedNode, itemNode)

	return nil
}

func (c *ComponentNodeBuilder) Build() (keyNode, contentNode node.Node, err error) {
	return c.namedNode, c.itemsNode, c.err
}

func (c *ComponentNodeBuilder) Err() error {
	return c.err
}
