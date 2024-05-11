package app

import (
	"log"
	"os"

	dissolve "github.com/amidgo/swaglue/internal/dissolve"
	"github.com/amidgo/swaglue/internal/glue"
)

const (
	MinArgsCount = 2

	CommandGlue     = "glue"
	CommandDissolve = "dissolve"
)

func Run() {
	args := os.Args

	if len(os.Args) < MinArgsCount {
		log.Fatal("no args")
	}

	command := args[1]

	switch command {
	case CommandGlue:
		glue.Exec()
	case CommandDissolve:
		dissolve.Exec()
	default:
		glue.Exec()
	}
}
