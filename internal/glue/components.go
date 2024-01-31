package glue

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidInput = errors.New("invalid input")

type component struct {
	Name string
	Path string
}

func parseComponentsFromString(s string) ([]*component, error) {
	if s == "" {
		return make([]*component, 0), nil
	}

	const validComponentSplitCount = 2

	components := strings.Split(s, ",")
	componentList := make([]*component, 0, len(components)/validComponentSplitCount)

	for _, cmnt := range components {
		cm := strings.Split(cmnt, "=")
		if len(cm) != validComponentSplitCount {
			return nil, fmt.Errorf("%w, parameter must follow pattern <name>=<dir path>", ErrInvalidInput)
		}

		componentList = append(componentList, &component{
			Name: cm[0],
			Path: cm[1],
		})
	}

	return componentList, nil
}
