package app

import (
	"flag"
	"io"
	"log"

	"github.com/amidgo/swaglue/internal/components"
	"github.com/amidgo/swaglue/internal/parser"
)

type Config struct {
	HeadFile, ComponentsData, Paths, Output string
	Debug                                   bool
}

func (c Config) ParsePaths() map[string]io.Reader {
	pathParser := parser.NewSwaggerPathParser(c.Paths, "#/paths/")
	err := pathParser.Parse()
	if err != nil {
		log.Fatalf("parse paths, err: %s", err)
	}
	if c.Debug {
		log.Printf("paths: %v", pathParser.Files())
	}
	return pathParser.Files()
}

func (c Config) PrintComponents(components []*components.Component) {
	if c.Debug {
		log.Printf("components: %v, raw:%s", components, c.ComponentsData)
	}
}

func parseConfigFromFlags() Config {
	var headFile, componentsData, paths, output string
	var debug bool
	flag.StringVar(&headFile, "head", "", "head swagger yaml, should be in .yaml format")
	flag.StringVar(&componentsData, "components", "", "components with name and directory path, example --components=<name>=<path>,<name>=<path>")
	flag.StringVar(&paths, "paths", "", "path to paths directory")
	flag.StringVar(&output, "output", "", "output file")
	flag.BoolVar(&debug, "_debug", false, "")
	flag.Parse()
	return Config{
		HeadFile:       headFile,
		ComponentsData: componentsData,
		Paths:          paths,
		Output:         output,
		Debug:          debug,
	}
}
