package app

import (
	"log"
	"os"

	"github.com/amidgo/swaglue/internal/components"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/parser"
)

func Run() {
	config := parseConfigFromFlags()
	head, err := head.ParseHeadFromFile(config.HeadFile)
	if err != nil {
		log.Fatalf("failed parse head from file, %s", err)
	}
	pathsFiles := config.ParsePaths()
	err = head.SetPaths(pathsFiles)
	if err != nil {
		log.Fatalf("failed set head paths, %s", err)
	}

	routes := config.ParseRoutes()
	err = head.AppendRoutes(routes)
	if err != nil {
		log.Fatalf("failed set append routes, %s", err)
	}

	components, err := components.ParseComponentsFromString(config.ComponentsData)
	if err != nil {
		log.Fatalf("failed parse components from string, %s", err)
	}
	config.PrintComponents(components)
	var componentParser *parser.SwaggerComponentParser
	for _, cmpnt := range components {
		componentParser = parser.NewSwaggerComponentParser(cmpnt.Path)
		err := componentParser.Parse()
		if err != nil {
			log.Fatalf("failed parse component, %s", err)
		}
		if config.Debug {
			log.Printf("component %s, component path %s, parser files: %v", cmpnt.Name, cmpnt.Path, componentParser.Files())
		}
		err = head.AppendComponent(cmpnt.Name, componentParser.Files())
		if err != nil {
			log.Fatalf("failed append component to head, %s", err)
		}
	}
	file, err := os.Create(config.Output)
	if err != nil {
		log.Fatalf("failed create output file %s, %s", config.Output, err)
	}
	err = head.SaveTo(file)
	if err != nil {
		log.Fatalf("failed save head to output file %s, %s", config.Output, err)
	}
}
