package cli

import (
	"log"
	"master-gen/internal/generator"
	"master-gen/internal/parser"
	"master-gen/internal/writer"
	"os"

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

func initBliss() {
	writer.WriteFile("./config.bliss", `{
  "blissPath": "./api.bliss",
  "serverPath": "./app",
  "webPath": "./clients/web"
}
`)
	writer.WriteFile("./api.bliss", `{
  "operations": [
    {
      "name": "Example",
      "endpoint": "/api/v1/example",
      "method": "Get",
      "queryParams": {
        "example": "string"
      },
      "type": "row",
      "query": "SELECT {name}$ FROM public.example e WHERE e.example = ${example}",
      "response": {
        "name": "string"
      }
    }
  ]
}
`)
}

func Run() {
	var cmd string
	var opt huh.Option[string]

	var hasConfig bool
	if _, err := os.Stat("./config.bliss"); !os.IsNotExist(err) {
		hasConfig = true
	}

	if hasConfig {
		opt = huh.NewOption("Generate", "generate")
	} else {
		opt = huh.NewOption("Initialize", "init")
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Action").
				Options(
					opt,
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
	case "init":
		initBliss()
	default:
		log.Fatalf("unknown command %s", cmd)
	}
}
