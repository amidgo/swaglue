package components

import (
	"errors"
	"fmt"
	"strings"
)

var (
	InvalidInput = errors.New("invalid input")
)

type Component struct {
	Name string
	Path string
}

type Components []*Component

func ParseComponentsFromString(s string) ([]*Component, error) {
	components := strings.Split(s, ",")
	componentList := make([]*Component, 0)
	for _, cmnt := range components {
		cm := strings.Split(cmnt, "=")
		if len(cm) != 2 {
			return nil, fmt.Errorf("%w, parameter must follow pattern <name>=<dir path>", InvalidInput)
		}
		componentList = append(componentList, &Component{
			Name: cm[0],
			Path: cm[1],
		})
	}
	return componentList, nil
}
