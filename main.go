package main

import (
	"log"
	"master-gen/generator"
	"master-gen/parser"
)

func main() {
	bliss, err := parser.ParseBliss("./app/aha.json")
	if err != nil {
		log.Fatal(err)
	}

	g := &generator.Generator{
		GenesisPath: "./app/genesis",
		Bliss:       bliss,
	}

	if err := g.Generate(); err != nil {
		log.Fatal(err)
	}
}
