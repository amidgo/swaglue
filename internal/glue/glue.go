package glue

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/amidgo/swaglue/internal/fileformats"
	"github.com/amidgo/swaglue/internal/glue/gluer"
	"github.com/amidgo/swaglue/internal/glue/head"
	"github.com/amidgo/swaglue/internal/glue/parser"
	"github.com/amidgo/swaglue/pkg/logger"
)

const PathsPrefix = "#/paths/"

func Exec() {
	config := parseGlueConfigFromFlags()

	fileFormat := fileFormat(config.TargetFileFormat)

	head, err := head.ParseHeadFromFile(config.HeadFile)
	if err != nil {
		log.Fatalf("failed parse head from file, %s", err)
	}

	gluerLogger := gluerLogger(&config)

	gluerContainer := gluer.NewContainer()

	if config.Paths != "" {
		pathsGluer := gluer.NewPathsGluer(
			gluerLogger,
			parser.NewSwaggerPathsParser(config.Paths, PathsPrefix, fileFormat),
			head,
		)
		gluerContainer.AddGluer(pathsGluer)
	}

	if config.Routes != "" {
		routesGluer := gluer.NewRoutesGluer(
			gluerLogger,
			parser.NewRouteParser(config.Routes),
			head,
		)
		gluerContainer.AddGluer(routesGluer)
	}

	if config.Tags != "" {
		tagsGluer := gluer.NewTagsGluer(
			gluerLogger,
			parser.NewSwaggerComponentParser(config.Tags, fileFormat),
			head,
		)
		gluerContainer.AddGluer(tagsGluer)
	}

	componentsData, err := parseComponentsFromString(config.ComponentsData)
	if err != nil {
		log.Fatalf("failed parse components, %s", err)
	}

	for _, cmpt := range componentsData {
		componentGluer := gluer.NewComponentsGluer(
			cmpt.Name,
			gluerLogger,
			parser.NewSwaggerComponentParser(cmpt.Path, fileFormat),
			head,
		)
		gluerContainer.AddGluer(componentGluer)
	}

	glue(&gluerContainer)

	save(head, config.Output)
}

type glueConfig struct {
	HeadFile, ComponentsData, Paths, Routes, Tags, TargetFileFormat, Output string
	Debug                                                                   bool
}

func parseGlueConfigFromFlags() glueConfig {
	var cnf glueConfig

	flag.StringVar(&cnf.HeadFile, "head", "", "head swagger yaml, should be in .yaml format")
	flag.StringVar(
		&cnf.ComponentsData,
		"components",
		"",
		"components with name and directory path, example --components=<name>=<path>,<name>=<path>",
	)
	flag.StringVar(&cnf.Paths, "paths", "", "path to paths directory")
	flag.StringVar(&cnf.Routes, "routes", "", "path to routes directory")
	flag.StringVar(&cnf.Tags, "tags", "", "path to tags directory")
	flag.StringVar(&cnf.Output, "output", "", "output file")
	flag.StringVar(&cnf.TargetFileFormat, "format", fileformats.YamlFileFormat, "scan file formats")
	flag.BoolVar(&cnf.Debug, "_debug", false, "")

	flag.Parse()

	return cnf
}

func fileFormat(format string) parser.FileFormat {
	fileFormat, err := fileformats.Detect(format)
	if err != nil {
		log.Fatalf("failed detect fileformat, %s", err)
	}

	return fileFormat
}

func gluerLogger(config *glueConfig) logger.Logger {
	var slogLevel slog.Level
	if config.Debug {
		slogLevel = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel})

	return &logger.SlogWrapper{Logger: slog.New(handler)}
}

func glue(gluerContainer gluer.Gluer) {
	err := gluerContainer.Glue()
	if err != nil {
		log.Fatalf("failed gluer container glue: %s", err)
	}
}

func save(head *head.Head, output string) {
	file, err := os.Create(output)
	if err != nil {
		log.Fatalf("failed create output file %s, %s", output, err)
	}

	err = head.SaveTo(file)
	if err != nil {
		log.Fatalf("failed save head to output file %s, %s", output, err)
	}
}
