package componentsappender

import (
	"errors"
	"fmt"

	"github.com/amidgo/node"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/model"
)

const componenetsTag = "components"

var ErrNoComponentsTag = errors.New("'components' tag not found")

type HeadComponentAppender struct {
	head    *head.Head
	decoder node.Decoder
}

func New(head *head.Head, decoder node.Decoder) *HeadComponentAppender {
	return &HeadComponentAppender{
		head:    head,
		decoder: decoder,
	}
}

func (h *HeadComponentAppender) AppendComponent(componentName string, componentItems []*model.Item) error {
	index, ok := h.head.SearchRootTag(componenetsTag)
	if !ok {
		return ErrNoComponentsTag
	}

	componentTag := h.head.Content()[index]

	if len(componentTag.Content()) == 0 {
		componentTag = node.MakeMapNode()
		h.head.Content()[index] = componentTag
	}

	appender := ComponentAppender{
		Decoder:        h.decoder,
		Node:           componentTag,
		ComponentName:  componentName,
		ComponentItems: componentItems,
	}

	err := appender.AppendComponent()
	if err != nil {
		return fmt.Errorf("failed append component, %w", err)
	}

	return nil
}

type ComponentAppender struct {
	Decoder        node.Decoder
	ComponentName  string
	ComponentItems []*model.Item
	Node           node.Node
}

func (a *ComponentAppender) AppendComponent() error {
	itemNode, exists := a.searchComponentNode()
	if exists {
		return a.appendExistComponent(itemNode)
	}

	return a.appendNewComponent()
}

func (a *ComponentAppender) searchComponentNode() (itemNode node.Node, exists bool) {
	for i := 0; i < len(a.Node.Content()); i += 2 {
		nd := a.Node.Content()[i]
		if nd.Kind() != node.String {
			continue
		}

		if nd.Value() != a.ComponentName {
			continue
		}

		if i == len(a.Node.Content())-1 {
			const nodesPerItem = 2

			return node.MakeMapNodeWithCap(len(a.ComponentItems) * nodesPerItem), true
		}

		return a.Node.Content()[i+1], true
	}

	return nil, false
}

func (a *ComponentAppender) appendNewComponent() error {
	nodes, err := new(ComponentNodeBuilder).
		SetDecoder(a.Decoder).
		SetName(a.ComponentName).
		SetItems(a.ComponentItems).
		Build()
	if err != nil {
		return err
	}

	a.Node.SetContent(append(a.Node.Content(), nodes...))

	return nil
}

func (a *ComponentAppender) appendExistComponent(itemNode node.Node) error {
	err := a.validateItemNode(itemNode)
	if err != nil {
		return err
	}

	err = new(ComponentNodeBuilder).
		SetDecoder(a.Decoder).
		SetItemNode(itemNode).
		AppendItems(a.ComponentItems).
		Err()
	if err != nil {
		return err
	}

	return nil
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
	decoder   node.Decoder
	namedNode node.Node
	itemsNode node.Node
	err       error
}

func (c *ComponentNodeBuilder) SetDecoder(dec node.Decoder) *ComponentNodeBuilder {
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

func (c *ComponentNodeBuilder) SetItemNode(itemsNode node.Node) *ComponentNodeBuilder {
	if c.err != nil {
		return c
	}

	c.itemsNode = itemsNode

	return c
}

func (c *ComponentNodeBuilder) AppendItems(items []*model.Item) *ComponentNodeBuilder {
	if c.err != nil {
		return c
	}

	c.appendComponents(items)

	return c
}

func (c *ComponentNodeBuilder) SetItems(items []*model.Item) *ComponentNodeBuilder {
	if c.err != nil {
		return c
	}

	c.itemsNode = node.MakeMapNodeWithCap(len(items))
	c.appendComponents(items)

	return c
}

func (c *ComponentNodeBuilder) appendComponents(components []*model.Item) {
	for _, component := range components {
		err := c.appendComponent(component)
		if err != nil {
			c.err = err

			return
		}
	}
}

func (c *ComponentNodeBuilder) appendComponent(component *model.Item) error {
	namedNode := node.MakeStringNode(component.Name)

	itemNode, err := head.DecodeNodeFrom(component.Content, c.decoder)
	if err != nil {
		return fmt.Errorf("%w, for %s, err: %w", head.ErrDecodeFile, component.Name, err)
	}

	c.appendItems(namedNode, itemNode)

	return nil
}

func (c *ComponentNodeBuilder) appendItems(items ...node.Node) {
	c.itemsNode.SetContent(append(c.itemsNode.Content(), items...))
}

func (c *ComponentNodeBuilder) Build() ([]node.Node, error) {
	return []node.Node{c.namedNode, c.itemsNode}, c.err
}

func (c *ComponentNodeBuilder) Err() error {
	return c.err
}
