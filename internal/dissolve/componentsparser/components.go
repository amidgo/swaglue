package parser

import (
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/model"
)

type ComponentsParser struct {
	headNode      *head.Head
	componentName string
	components    []*model.Item
}

func NewComponentsParser(headNode *head.Head, componentName string) *ComponentsParser {
	return &ComponentsParser{
		headNode:      headNode,
		componentName: componentName,
		components:    make([]*model.Item, 0),
	}
}

func (c *ComponentsParser) Parse() error {
	return nil
}

func (c *ComponentsParser) Components() []*model.Item {
	return c.components
}
