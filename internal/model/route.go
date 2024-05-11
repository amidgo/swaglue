package model

import "io"

type Route struct {
	Name    string
	Methods []*RouteMethod
}

type RouteMethod struct {
	Method  string
	Content io.Reader
}
