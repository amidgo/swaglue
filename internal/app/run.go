package app

import (
	"log"
	"os"

	"github.com/amidgo/swaglue/internal/glue"
	"github.com/amidgo/swaglue/internal/solv"
)

const (
	MinArgsCount = 2

	CommandGlue = "glue"
	CommandSolv = "solv"
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
	case CommandSolv:
		solv.Exec()
	default:
		glue.Exec()
	}
}
