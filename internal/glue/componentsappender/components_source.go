package componentsappender

import (
	"errors"

	"github.com/amidgo/node"
	"github.com/amidgo/swaglue/internal/item"
)

var (
	ErrIterationSource = errors.New("iterate on component node")
	ErrInvalidNode     = errors.New("invalid node")
)

type ComponentIterationStep struct {
	name   string
	source item.Source
}

func NewComponentIterationStep(name string, source item.Source) ComponentIterationStep {
	return ComponentIterationStep{
		name:   name,
		source: source,
	}
}

func (c ComponentIterationStep) KeyValue(key, value node.Node) (resKey, resValue node.Node, err error) {
	if !node.StringEquals(key, c.name) {
		return key, value, nil
	}

	componentNode, err := componentNodeSource{
		source: item.NewRemoveDuplicatesSource(
			item.NewJoinSource(
				item.NewNodeSource(value),
				c.source,
			),
		),
	}.ComponentNode()
	if err != nil {
		return nil, nil, err
	}

	return key, componentNode, nil
}

type componentNodeSource struct {
	source item.Source
}

func (c componentNodeSource) ComponentNode() (node.MapNode, error) {
	items, err := c.source.Items()
	if err != nil {
		return node.MakeMapNode(), err
	}

	const nodesPerItem = 2

	mapNode := node.MakeMapNode(make([]node.Node, 0, len(items)*nodesPerItem)...)

	for _, item := range items {
		mapNode = node.MapAppend(mapNode, node.MakeStringNode(item.Name), item.Content)
	}

	return mapNode, nil
}

type IterationStep struct {
	componentsSteps []ComponentIterationStep
}

func NewIterationStep(componentsSteps ...ComponentIterationStep) IterationStep {
	return IterationStep{componentsSteps: componentsSteps}
}

func (i IterationStep) iterationSteps() []node.IterationStep {
	steps := make([]node.IterationStep, 0, len(i.componentsSteps))

	for _, step := range i.componentsSteps {
		steps = append(steps, step)
	}

	return steps
}

func (i IterationStep) names() []string {
	names := make([]string, 0, len(i.componentsSteps))

	for _, step := range i.componentsSteps {
		names = append(names, step.name)
	}

	return names
}

func (i IterationStep) KeyValue(key, value node.Node) (resKey, resValue node.Node, err error) {
	if !node.StringEquals(key, componentsTag) {
		return key, value, nil
	}

	gen := node.NewIterationMapSource(
		newfilledWithNamesComponentNode(i.names(), value),
		node.NewJoinIterationStep(i.iterationSteps()...),
	)

	componentNode, err := gen.MapNode()
	if err != nil {
		return nil, nil, errors.Join(ErrIterationSource, err)
	}

	return key, componentNode, nil
}

type filledWithNamesComponentNode struct {
	names         []string
	componentNode node.Node
	validate      node.Validate
}

func newfilledWithNamesComponentNode(names []string, componentNode node.Node) filledWithNamesComponentNode {
	return filledWithNamesComponentNode{
		names:         names,
		componentNode: componentNode,
		validate:      node.NewKindValidate(node.Map, node.Empty),
	}
}

func (f filledWithNamesComponentNode) MapNode() (node.MapNode, error) {
	err := f.validate.Validate(f.componentNode)
	if err != nil {
		return node.MakeMapNode(), errors.Join(ErrInvalidNode, err)
	}

	names := f.makeNames()

	iter := node.MakeMapNodeIterator(f.componentNode.Content())
	for iter.HasNext() {
		key, _ := iter.Next()

		_, isTargetName := names[key.Value()]
		if !isTargetName {
			continue
		}

		names[key.Value()] = true
	}

	mapNode := node.MakeMapNode(f.componentNode.Content()...)

	for _, name := range f.names {
		existsInNode := names[name]
		if !existsInNode {
			mapNode = node.MapAppend(mapNode, node.MakeStringNode(name), node.MakeMapNode())
		}
	}

	return mapNode, nil
}

func (f filledWithNamesComponentNode) makeNames() map[string]bool {
	names := make(map[string]bool, len(f.names))

	for _, name := range f.names {
		names[name] = false
	}

	return names
}
