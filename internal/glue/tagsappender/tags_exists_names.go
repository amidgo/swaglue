package tagsappender

import (
	"errors"

	"github.com/amidgo/node"
)

var (
	ErrExtractTagNodeName = errors.New("extract tag node name")
)

type TagsExistsNames struct {
	names map[string]struct{}
}

func (e *TagsExistsNames) TagNameExists(tagName string) bool {
	_, ok := e.names[tagName]

	return ok
}

func (e *TagsExistsNames) ScanNode(tagsNode node.Node) error {
	e.names = make(map[string]struct{}, len(tagsNode.Content()))

	for _, tag := range tagsNode.Content() {
		tagName, err := extractTagNodeName(tag)
		if err != nil {
			return errors.Join(ErrExtractTagNodeName, err)
		}

		e.names[tagName] = struct{}{}
	}

	return nil
}

var (
	ErrNoTagName      = errors.New("no 'name' tag")
	ErrInvalidTagYaml = errors.New("invalid tag yaml")
)

func extractTagNodeName(tagNode node.Node) (string, error) {
	const nameKey = "name"

	if len(tagNode.Content())%2 != 0 {
		return "", ErrInvalidTagYaml
	}

	for i := 0; i < len(tagNode.Content()); i += 2 {
		if tagNode.Content()[i].Value() != nameKey {
			continue
		}

		return tagNode.Content()[i+1].Value(), nil
	}

	return "", ErrNoTagName
}
