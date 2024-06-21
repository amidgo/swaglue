package route

import "io"

type Route struct {
	Name    string
	Methods []Method
}

type Method struct {
	Method  string
	Content io.Reader
}
