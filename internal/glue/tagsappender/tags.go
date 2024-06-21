package tagsappender

import (
	"errors"
	"fmt"

	"github.com/amidgo/node"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/model"
)

const tagsTag = "tags"

var ErrNoTagsTag = errors.New("'tags' tag not found")

type HeadTagsAppender struct {
	decoder node.DecoderFrom
	head    *head.Head
}

func New(head *head.Head, decoder node.DecoderFrom) *HeadTagsAppender {
	return &HeadTagsAppender{head: head, decoder: decoder}
}

func (h *HeadTagsAppender) AppendTags(tags []model.Item) error {
	index := node.MapSearchByStringKey(h.head.Node(), tagsTag)
	if index == -1 {
		return ErrNoTagsTag
	}

	tagsNode := h.head.Node().Content()[index]

	var tagsExistsNames TagsExistsNames

	err := tagsExistsNames.ScanNode(tagsNode)
	if err != nil {
		return fmt.Errorf("scan tags node, %w", err)
	}

	tagsAppend := TagsAppend{
		decoder:         h.decoder,
		arrayNode:       node.MakeArrayNodeWithContent(tagsNode.Content()...),
		tagsExistsNames: tagsExistsNames,
	}

	nd, err := tagsAppend.Node(tags)
	if err != nil {
		return fmt.Errorf("append tags to node, %w", err)
	}

	h.head.Node().Content()[index] = nd

	return nil
}

type TagsAppend struct {
	decoder         node.DecoderFrom
	arrayNode       node.ArrayNode
	tagsExistsNames TagsExistsNames
}

var (
	ErrTagNameExists      = errors.New("tag name exists")
	ErrInvalidContentName = errors.New("tag content name not equal filename")
)

func (n *TagsAppend) Node(tags []model.Item) (node.Node, error) {
	for _, tag := range tags {
		if n.tagsExistsNames.TagNameExists(tag.Name) {
			return nil, fmt.Errorf("%w, tag name: %s", ErrTagNameExists, tag.Name)
		}

		nd, err := n.decoder.DecodeFrom(tag.Content)
		if err != nil {
			return nil, fmt.Errorf("decode item content, %w", err)
		}

		name, err := extractTagNodeName(nd)
		if err != nil {
			return nil, fmt.Errorf("extract tag node name, %w", err)
		}

		if tag.Name != name {
			return nil, fmt.Errorf("%w, file name: %s", ErrInvalidContentName, tag.Name)
		}

		n.arrayNode = node.ArrayAppend(n.arrayNode, nd)
	}

	return n.arrayNode, nil
}
