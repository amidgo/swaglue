package app

import (
	"flag"
	"log"
	"os"

	"github.com/amidgo/swaglue/internal/components"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/internal/parser"
)

func Run() {
	var headFile, componentsData, paths, output string
	var debug bool
	flag.StringVar(&headFile, "head", "", "head swagger yaml, should be in .yaml format")
	flag.StringVar(&componentsData, "components", "", "components with name and directory path, example --components=<name>=<path>,<name>=<path>")
	flag.StringVar(&paths, "paths", "", "path to paths directory")
	flag.StringVar(&output, "output", "", "output file")
	flag.BoolVar(&debug, "_debug", false, "")
	flag.Parse()
	head, err := head.ParseHeadFromFile(headFile)
	if err != nil {
		log.Fatal(err)
	}
	pathsFiles := ParsePaths(paths, debug)
	err = head.SetPaths(pathsFiles)
	if err != nil {
		log.Fatal(err)
	}
	components := make(components.Components, 0)
	err = components.ParseFromString(componentsData)
	if err != nil {
		log.Fatal(err)
	}
	if debug {
		log.Printf("components: %v, raw:%s", components, componentsData)
	}
	var componentParser *parser.SwaggerComponentParser
	for _, cmpnt := range components {
		componentParser = parser.NewSwaggerComponentParser(cmpnt.Path)
		err := componentParser.Parse()
		if err != nil {
			log.Fatal(err)
		}
		if debug {
			log.Printf("component %s, component path %s, parser files: %v", cmpnt.Name, cmpnt.Path, componentParser.Files())
		}
		err = head.AppendComponent(cmpnt.Name, componentParser.Files())
		if err != nil {
			log.Fatal(err)
		}
	}
	file, err := os.Create(output)
	if err != nil {
		log.Fatalf("failed create output file, %s", err)
	}
	err = head.SaveTo(file)
	if err != nil {
		log.Fatal(err)
	}
}
