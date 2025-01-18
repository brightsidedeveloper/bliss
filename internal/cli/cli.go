package cli

import (
	"bytes"
	"fmt"
	"log"
	"master-gen/internal/generator"
	"master-gen/internal/parser"
	"master-gen/internal/writer"
	"os"
	"os/exec"
	"strings"

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
	var opts []huh.Option[string]

	var hasConfig bool
	var dockerComposePath string
	if _, err := os.Stat("./config.bliss"); !os.IsNotExist(err) {
		hasConfig = true
	}

	if hasConfig {
		opts = append(opts, huh.NewOption("Generate", "generate"))
		config, err := parser.GetConfig("./config.bliss")
		if err != nil {
			log.Fatal(err)
		}
		if _, err := os.Stat(config.ServerPath + "/docker-compose.yaml"); !os.IsNotExist(err) {
			if checkDockerDB() {
				opts = append(opts, huh.NewOption("DB Down", "docker-compose down"))
			} else {
				opts = append(opts, huh.NewOption("DB Up", "docker-compose up -d"))
			}
			dockerComposePath = config.ServerPath
		}
	} else {
		opts = append(opts, huh.NewOption("Initialize", "init"))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Action").
				Options(
					opts...,
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
		fmt.Println("Generated code")
	case "init":
		initBliss()
		fmt.Println("Initialized bliss config")
	case "docker-compose up -d":
		cmd := exec.Command("docker-compose", "up", "-d")
		cmd.Dir = dockerComposePath
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Started Database")
	case "docker-compose down":
		cmd := exec.Command("docker-compose", "down")
		cmd.Dir = dockerComposePath
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Stopped Database")
	default:
		log.Fatalf("unknown command %s", cmd)
	}
}

func checkDockerDB() bool {
	cmd := exec.Command("docker", "ps")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing docker ps: %v\n", err)
		return false
	}

	// Read the output of `docker ps`
	output := out.String()
	lines := strings.Split(output, "\n")

	// Check if any line contains the container name "db"
	containerFound := false
	for _, line := range lines {
		if strings.Contains(line, "db") {
			containerFound = true
			break
		}
	}
	return containerFound
}
