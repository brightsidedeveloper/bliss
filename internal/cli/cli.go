package cli

import (
	"log"
	"master-gen/internal/generator"

	"github.com/charmbracelet/huh"
)

func generate() {
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
}

func Run() {
	var cmd string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Action").
				Options(
					huh.NewOption("Generate", "generate"),
				).
				Value(&cmd),
		),
	)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	switch cmd {
	case "generate":
		generate()
	default:
		log.Fatalf("unknown command %s", cmd)
	}
}
