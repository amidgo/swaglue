package item

import (
	"errors"

	"github.com/amidgo/node"
)

type Source interface {
	Items() ([]Item, error)
}

type SliceSource []Item

func (s SliceSource) Items() ([]Item, error) {
	return s, nil
}

var (
	ErrGetItemsFromSource = errors.New("get items from source")
	ErrInvalidNode        = errors.New("invalid node")
)

type JoinSource struct {
	sources []Source
}

func NewJoinSource(sources ...Source) JoinSource {
	return JoinSource{sources: sources}
}

func (j JoinSource) Items() ([]Item, error) {
	totalRoutes := make([][]Item, 0, len(j.sources))
	totalSize := 0

	for _, source := range j.sources {
		routes, err := source.Items()
		if err != nil {
			return nil, errors.Join(ErrGetItemsFromSource, err)
		}

		totalSize += len(routes)
		totalRoutes = append(totalRoutes, routes)
	}

	joinedRoutes := make([]Item, totalSize)

	offset := 0
	for _, routes := range totalRoutes {
		copy(joinedRoutes[offset:len(routes)+offset], routes)
		offset += len(routes)
	}

	return joinedRoutes, nil
}

type NodeSource struct {
	nd       node.Node
	validate node.Validate
}

func NewNodeSource(nd node.Node) NodeSource {
	return NodeSource{
		nd:       nd,
		validate: node.NewKindValidate(node.Map, node.Empty),
	}
}

func (n NodeSource) Items() ([]Item, error) {
	err := n.validate.Validate(n.nd)
	if err != nil {
		return nil, errors.Join(ErrInvalidNode, err)
	}

	const nodesPerItem = 2

	items := make([]Item, 0, len(n.nd.Content())/nodesPerItem)

	iter := node.MakeMapNodeIterator(n.nd.Content())
	for iter.HasNext() {
		key, value := iter.Next()

		items = append(items, Item{Name: key.Value(), Content: value})
	}

	return items, nil
}

type RemoveDuplicatesSource struct {
	source Source
}

func NewRemoveDuplicatesSource(source Source) RemoveDuplicatesSource {
	return RemoveDuplicatesSource{
		source: source,
	}
}

func (r RemoveDuplicatesSource) Items() ([]Item, error) {
	items, err := r.source.Items()
	if err != nil {
		return nil, errors.Join(ErrGetItemsFromSource, err)
	}

	uniqueItems := make([]Item, 0, len(items))
	existsNames := make(map[string]struct{}, len(items))

	for _, item := range items {
		_, exists := existsNames[item.Name]
		if exists {
			continue
		}

		existsNames[item.Name] = struct{}{}

		uniqueItems = append(uniqueItems, item)
	}

	return uniqueItems, nil
}
