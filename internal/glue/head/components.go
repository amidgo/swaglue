package head

import (
	"errors"
	"fmt"

	"github.com/amidgo/swaglue/internal/glue/model"
	"gopkg.in/yaml.v3"
)

const componenetsTag = "components"

var ErrNoComponentsTag = errors.New("'components' tag not found")

func (h *Head) AppendComponent(componentName string, componentItems []*model.Item) error {
	componentTag := h.SearchTag(componenetsTag)
	if componentTag == nil {
		return ErrNoComponentsTag
	}

	if len(componentTag.Content) == 0 {
		componentTag.Kind = yaml.MappingNode
		componentTag.Tag = ""
	}

	appender := ComponentAppender{
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
	ComponentName  string
	ComponentItems []*model.Item
	Node           *yaml.Node
}

func (a *ComponentAppender) AppendComponent() error {
	itemNode, exists := a.searchComponentNode()
	if exists {
		return a.appendExistComponent(itemNode)
	}

	return a.appendNewComponent()
}

func (a *ComponentAppender) searchComponentNode() (itemNode *yaml.Node, exists bool) {
	for i := 0; i < len(a.Node.Content); i += 2 {
		node := a.Node.Content[i]
		if node.Kind != yaml.ScalarNode {
			continue
		}

		if node.Value != a.ComponentName {
			continue
		}

		if i == len(a.Node.Content)-1 {
			const nodesPerItem = 2

			return &yaml.Node{
				Kind:    yaml.MappingNode,
				Content: make([]*yaml.Node, 0, len(a.ComponentItems)*nodesPerItem),
			}, true
		}

		return a.Node.Content[i+1], true
	}

	return nil, false
}

func (a *ComponentAppender) appendNewComponent() error {
	nodes, err := new(ComponentNodeBuilder).
		SetName(a.ComponentName).
		SetItems(a.ComponentItems).
		Build()
	if err != nil {
		return err
	}

	a.Node.Content = append(a.Node.Content, nodes...)

	return nil
}

func (a *ComponentAppender) appendExistComponent(itemNode *yaml.Node) error {
	err := a.validateItemNode(itemNode)
	if err != nil {
		return err
	}

	err = new(ComponentNodeBuilder).
		SetItemNode(itemNode).
		AppendItems(a.ComponentItems).
		Err()
	if err != nil {
		return err
	}

	return nil
}

var ErrComponentItemNameExists = errors.New("component item name exists")

func (a *ComponentAppender) validateItemNode(itemNode *yaml.Node) error {
	existsNames := make(map[string]struct{}, len(a.ComponentItems))
	for _, item := range a.ComponentItems {
		existsNames[item.Name] = struct{}{}
	}

	for i := 0; i < len(itemNode.Content); i += 2 {
		node := itemNode.Content[i]
		if node.Kind != yaml.ScalarNode {
			continue
		}

		_, exists := existsNames[node.Value]
		if exists {
			return fmt.Errorf("%w, name %s, line %d", ErrComponentItemNameExists, node.Value, node.Line)
		}
	}

	return nil
}

type ComponentNodeBuilder struct {
	namedNode *yaml.Node
	itemsNode *yaml.Node
	err       error
}

func (c *ComponentNodeBuilder) SetName(name string) *ComponentNodeBuilder {
	if c.err != nil {
		return c
	}

	c.namedNode = &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: name,
	}

	return c
}

func (c *ComponentNodeBuilder) SetItemNode(itemsNode *yaml.Node) *ComponentNodeBuilder {
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

	c.itemsNode = &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: make([]*yaml.Node, 0, len(items)),
	}
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
	namedNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: component.Name,
	}
	dec := yaml.NewDecoder(component.Content)

	itemNode, err := DecodeYamlNode(dec)
	if err != nil {
		return fmt.Errorf("%w, for %s, err: %w", ErrFailedDecodeFile, component.Name, err)
	}

	c.appendItems(namedNode, itemNode)

	return nil
}

func (c *ComponentNodeBuilder) appendItems(items ...*yaml.Node) {
	c.itemsNode.Content = append(c.itemsNode.Content, items...)
}

func (c *ComponentNodeBuilder) Build() ([]*yaml.Node, error) {
	return []*yaml.Node{c.namedNode, c.itemsNode}, c.err
}

func (c *ComponentNodeBuilder) Err() error {
	return c.err
}
