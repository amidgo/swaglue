package head

import (
	"errors"
	"fmt"

	"github.com/amidgo/swaglue/internal/glue/model"
	"gopkg.in/yaml.v3"
)

const tagsTag = "tags"

var ErrNoTagsTag = errors.New("'tags' tag not found")

func (h *Head) AppendTags(tags []*model.Item) error {
	tagsNode := h.SearchTag(tagsTag)
	if tagsNode == nil {
		return ErrNoTagsTag
	}

	if len(tagsNode.Content) == 0 {
		tagsNode.Kind = yaml.MappingNode
		tagsNode.Tag = ""
	}

	var tagsExistsNames TagsExistsNames

	err := tagsExistsNames.ScanNode(tagsNode)
	if err != nil {
		return fmt.Errorf("failed scan tags node, err: %w", err)
	}

	appender := TagsAppender{
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
	Node            *yaml.Node
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

		node, err := DecodeYamlNode(yaml.NewDecoder(tag.Content))
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

		n.Node.Content = append(n.Node.Content, node)
	}

	return nil
}
