package head

import (
	"errors"
	"fmt"

	"github.com/amidgo/swaglue/internal/model"
	"gopkg.in/yaml.v3"
)

var (
	NoComponentsTag = errors.New("'components' tag not found")
)

func (h Head) AppendComponent(componentName string, componentItems []*model.Component) error {
	componentTag := h.SearchTag("components")
	if componentTag == nil {
		return NoComponentsTag
	}
	if len(componentTag.Content) == 0 {
		componentTag.Kind = yaml.MappingNode
		componentTag.Tag = ""
	}
	nodes, err := new(ComponentNodeBuilder).
		SetName(componentName).
		SetItems(componentItems).
		Build()
	if err != nil {
		return err
	}
	componentTag.Content = append(componentTag.Content, nodes...)
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

func (c *ComponentNodeBuilder) SetItems(items []*model.Component) *ComponentNodeBuilder {
	if c.err != nil {
		return c
	}
	c.itemsNode = &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: make([]*yaml.Node, 0, len(items)),
	}

	for _, component := range items {
		namedNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: component.Name,
		}
		decoder := Decoder{yaml.NewDecoder(component.Content)}
		itemNode := decoder.Node()
		if itemNode == nil {
			c.err = fmt.Errorf("%w, for %s", FailedDecodeFile, component.Name)
			return c
		}
		c.itemsNode.Content = append(c.itemsNode.Content, namedNode, itemNode)
	}
	return c
}

func (c *ComponentNodeBuilder) Build() ([]*yaml.Node, error) {
	return []*yaml.Node{c.namedNode, c.itemsNode}, c.err
}
