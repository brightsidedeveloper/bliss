package main

import (
	"fmt"
	"log"
	"master-gen/internal/generator"
)

func main() {

	err := (&generator.Generator{
		BlissPath:  "./api.bliss",
		ServerPath: "./app",
		ClientPaths: []string{
			"./clients/web/src/api",
		},
	}).Generate()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
}
