package model

import "io"

type Component struct {
	Name    string
	Content io.Reader
}
