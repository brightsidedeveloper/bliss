package cli

import (
	"log"
	"master-gen/internal/generator"
	"master-gen/internal/parser"

	"github.com/charmbracelet/huh"
)

func generate() {
	config, err := parser.GetConfig("./config.bliss")
	if err != nil {
		log.Fatal(err)
	}
	err = (&generator.Generator{
		BlissPath:  config.BlissPath,
		ServerPath: config.ServerPath,
		WebPath:    config.WebPath,
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
