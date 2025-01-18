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

	if err := generator.TypesFileContent(aha); err != nil {
		log.Fatal(err)
	}

	if err := generator.Handlers(aha); err != nil {
		log.Fatal(err)
	}

	if err := generator.MountRoutesFileContent(aha); err != nil {
		log.Fatal(err)
	}

}
