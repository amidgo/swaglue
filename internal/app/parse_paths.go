package app

import (
	"io"
	"log"

	"github.com/amidgo/swaglue/internal/parser"
)

func ParsePaths(paths string, debug bool) map[string]io.Reader {
	pathParser := parser.NewSwaggerPathParser(paths, "#/paths/")
	err := pathParser.Parse()
	if err != nil {
		log.Fatalf("parse paths, err: %s", err)
	}
	if debug {
		log.Printf("paths: %v", pathParser.Files())
	}
	return pathParser.Files()
}
