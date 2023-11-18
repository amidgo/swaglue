package model

import "io"

type SwaggerComponentItem struct {
	Name    string
	Content io.Reader
}
