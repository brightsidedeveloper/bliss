package main

import (
	"log"
	"master-gen/generator"
	"master-gen/parser"
)

func main() {
	aha, err := parser.ParseAhaJSON("./app/aha.json")
	if err != nil {
		log.Fatal(err)
	}

	g := &generator.Generator{
		GenesisPath: "./app/genesis",
		Bliss:       aha,
	}

	if err := g.Generate(); err != nil {
		log.Fatal(err)
	}

}
