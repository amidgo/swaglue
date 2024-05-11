package model

import "io"

type Item struct {
	Name    string
	Content io.Reader
}
