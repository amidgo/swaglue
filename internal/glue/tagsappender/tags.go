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
	decoder node.Decoder
	head    *head.Head
}

func New(head *head.Head, decoder node.Decoder) *HeadTagsAppender {
	return &HeadTagsAppender{head: head, decoder: decoder}
}

func (h *HeadTagsAppender) AppendTags(tags []*model.Item) error {
	index, ok := h.head.SearchRootTag(tagsTag)
	if !ok {
		return ErrNoTagsTag
	}

	tagsNode := h.head.Content()[index]

	if len(tagsNode.Content()) == 0 {
		tagsNode = node.MakeArrayNode()
		h.head.Content()[index] = tagsNode
	}

	var tagsExistsNames TagsExistsNames

	err := tagsExistsNames.ScanNode(tagsNode)
	if err != nil {
		return fmt.Errorf("failed scan tags node, err: %w", err)
	}

	appender := TagsAppender{
		Decoder:         h.decoder,
		Node:            tagsNode,
		TagsExistsNames: tagsExistsNames,
	}

	err = appender.AppendTags(tags)
	if err != nil {
		return fmt.Errorf("failed append tags to node, %w", err)
	}

	return nil
}

type TagsAppender struct {
	Decoder         node.Decoder
	Node            node.Node
	TagsExistsNames TagsExistsNames
}

var (
	ErrTagNameExists      = errors.New("tag name exists")
	ErrInvalidContentName = errors.New("tag content name not equal filename")
)

func (n *TagsAppender) AppendTags(tags []*model.Item) error {
	for _, tag := range tags {
		if n.TagsExistsNames.TagNameExists(tag.Name) {
			return fmt.Errorf("%w, tag name: %s", ErrTagNameExists, tag.Name)
		}

		node, err := head.DecodeNodeFrom(tag.Content, n.Decoder)
		if err != nil {
			return fmt.Errorf("failed decode item content, err: %w", err)
		}

		name, err := extractTagNodeName(node)
		if err != nil {
			return fmt.Errorf("failed extract tag node name, err: %w", err)
		}

		if tag.Name != name {
			return fmt.Errorf("%w, file name: %s", ErrInvalidContentName, tag.Name)
		}

		n.Node.AppendNode(node)
	}

	return nil
}
