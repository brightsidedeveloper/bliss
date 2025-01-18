package main

import (
	"fmt"
	"log"
	"master-gen/generator"
)

func main() {

	err := (&generator.Generator{
		BlissPath:   "./api.bliss",
		GenesisPath: "./app/genesis",
		ApiPaths: []string{
			"./clients/web/src/api",
		},
	}).Generate()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
}
