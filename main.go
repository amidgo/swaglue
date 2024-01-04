package main

import (
	_ "embed"
	"log"

	"gopkg.in/yaml.v3"
)

//go:embed test.yaml
var testData []byte

func main() {
	var node yaml.Node
	yaml.Unmarshal(testData, &node)
	log.Println()
}
