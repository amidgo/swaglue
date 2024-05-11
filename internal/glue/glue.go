package glue

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/amidgo/node"
	"github.com/amidgo/node/json"
	"github.com/amidgo/node/yaml"
	"github.com/amidgo/swaglue/internal/fileformats"
	"github.com/amidgo/swaglue/internal/glue/componentsappender"
	"github.com/amidgo/swaglue/internal/glue/parser"
	"github.com/amidgo/swaglue/internal/glue/pathssetter"
	"github.com/amidgo/swaglue/internal/glue/routesappender"
	"github.com/amidgo/swaglue/internal/glue/tagsappender"
	"github.com/amidgo/swaglue/internal/gluer"
	"github.com/amidgo/swaglue/internal/head"
	"github.com/amidgo/swaglue/pkg/logger"
)

const PathsPrefix = "#/paths/"

//nolint:funlen // it is an entrypoint function
func Exec() {
	config := parseGlueConfigFromFlags()

	fileFormat := config.fileFormat()
	decoder := config.decoder()
	encoder := config.encoder()
	logger := config.logger()

	gluerContainer := gluer.NewContainer()

	head, err := head.ParseHeadFromFile(config.HeadFile, decoder)
	if err != nil {
		log.Fatalf("parse head from file, %s", err)
	}

	if config.Paths != "" {
		gluerContainer.AddGluer(
			gluer.NewPathsGluer(
				logger,
				parser.NewSwaggerPathsParser(config.Paths, PathsPrefix, fileFormat),
				pathssetter.New(head, decoder),
			),
		)
	}

	if config.Routes != "" {
		gluerContainer.AddGluer(
			gluer.NewRoutesGluer(
				logger,
				parser.NewRouteParser(config.Routes),
				routesappender.New(head, decoder),
			),
		)
	}

	if config.Tags != "" {
		gluerContainer.AddGluer(
			gluer.NewTagsGluer(
				logger,
				parser.NewSwaggerComponentParser(config.Tags, fileFormat),
				tagsappender.New(head, decoder),
			),
		)
	}

	componentsData, err := parseComponentsFromString(config.ComponentsData)
	if err != nil {
		log.Fatalf("parse components, %s", err)
	}

	for _, cmpt := range componentsData {
		gluerContainer.AddGluer(
			gluer.NewComponentsGluer(
				cmpt.Name,
				logger,
				parser.NewSwaggerComponentParser(cmpt.Path, fileFormat),
				componentsappender.New(head, decoder),
			),
		)
	}

	glue(&gluerContainer)

	save(head, config.Output, encoder)
}

type glueConfig struct {
	HeadFile, ComponentsData, Paths, Routes, Tags, TargetFileFormat, Output string
	YamlIndent                                                              int
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
	flag.IntVar(&cnf.YamlIndent, "yamlindent", head.DefaultYamlIndent, "yaml indent")
	flag.BoolVar(&cnf.Debug, "debug", false, "use for debug logging")

	flag.Parse()

	return cnf
}

func (cnf *glueConfig) fileFormat() parser.FileFormat {
	fileFormat, err := fileformats.Detect(cnf.TargetFileFormat)
	if err != nil {
		log.Fatalf("detect fileformat, %s", err)
	}

	return fileFormat
}

func (cnf *glueConfig) logger() logger.Logger {
	var slogLevel slog.Level
	if cnf.Debug {
		slogLevel = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel})

	return &logger.SlogWrapper{Logger: slog.New(handler)}
}

func (cnf *glueConfig) decoder() node.Decoder {
	switch cnf.TargetFileFormat {
	case fileformats.JSONFileFormat:
		return &json.Decoder{}
	case fileformats.YamlFileFormat:
		return &yaml.Decoder{}
	default:
		log.Fatalf("unsupported fileformat, %s", cnf.TargetFileFormat)
	}

	return nil
}

func (cnf *glueConfig) encoder() node.Encoder {
	switch cnf.TargetFileFormat {
	case fileformats.JSONFileFormat:
		return &json.Encoder{}
	case fileformats.YamlFileFormat:
		return &yaml.Encoder{Indent: cnf.YamlIndent}
	default:
		log.Fatalf("unsupported fileformat, %s", cnf.TargetFileFormat)
	}

	return nil
}

func glue(gluerContainer gluer.Gluer) {
	err := gluerContainer.Glue()
	if err != nil {
		log.Fatalf("gluer container glue: %s", err)
	}
}

func save(head *head.Head, output string, encoder node.Encoder) {
	file, err := os.Create(output)
	if err != nil {
		log.Fatalf("create output file %s, %s", output, err)
	}

	err = head.SaveTo(file, encoder)
	if err != nil {
		log.Fatalf("save head to output file %s, %s", output, err)
	}
}
