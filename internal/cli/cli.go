package cli

import (
	"fmt"
	"log"
	"master-gen/internal/docker"
	"master-gen/internal/generator"
	"master-gen/internal/parser"
	"master-gen/internal/server"
	"master-gen/internal/writer"
	"os"
	"os/exec"

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
  "webPath": "./web"
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
      "query": "GetExample",
      "response": {
        "name": "string"
      }
    }
  ]
}
`)
	writer.WriteFile("schemas.sql", `CREATE TABLE IF NOT EXISTS public.example (
  id   SERIAL PRIMARY KEY,
  name text      NOT NULL
);`)
	writer.WriteFile("api.schema.bliss", `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "operations": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "description": "The name of the operation."
          },
          "endpoint": {
            "type": "string",
            "description": "The path of the endpoint."
          },
          "method": {
            "type": "string",
            "description": "The HTTP method of the endpoint."
          },
          "query": {
            "type": "string",
            "description": "The query to execute."
          },
          "queryParams": {
            "type": "object",
            "patternProperties": {
              "\\w+": {
                "type": "string",
                "enum": [
                  "string",
                  "int",
                  "bool"
                ]
              }
            },
            "description": "The query parameters."
          },
          "handler": {
            "type": "string",
            "description": "The handler to execute."
          },
          "response": {
            "type": "object",
            "properties": {
              "test": {
                "type": "string"
              }
            },
            "description": "The response json."
          }
        },
        "anyOf": [
          {
            "required": [
              "query"
            ]
          },
          {
            "required": [
              "handler"
            ]
          }
        ],
        "required": [
          "name",
          "endpoint",
          "method"
        ]
      }
    }
  }
}
`)
	writer.WriteFile("queries.sql", `-- name: GetExample :one
SELECT * FROM public.example
WHERE name = $1 LIMIT 1;`)
	writer.WriteFile("sqlc.yaml", `version: "2"
sql:
  - engine: "postgresql"
    queries: "queries.sql"
    schema: "schemas.sql"
    gen:
      go:
        package: "queries"
        out: "./app/genesis/queries"
        sql_package: "pgx/v5"
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

	if !hasConfig {
		initBliss()
		fmt.Println("Initialized bliss config")
		return
	}

	opts = append(opts, huh.NewOption("Generate", "generate"))
	config, err := parser.GetConfig("./config.bliss")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(config.ServerPath + "/docker-compose.yaml"); !os.IsNotExist(err) {

		if server.IsServerRunning(config.ServerPath) {
			opts = append(opts, huh.NewOption("Stop Server", "stopServer"))
		} else {
			opts = append(opts, huh.NewOption("Start Server", "startServer"))
		}

		if docker.CheckDockerDB() {
			opts = append(opts, huh.NewOption("DB Down", "docker-compose down"))
		} else {
			opts = append(opts, huh.NewOption("DB Up", "docker-compose up -d"))
		}
		dockerComposePath = config.ServerPath
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
	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	switch cmd {
	case "generate":
		generate()
		fmt.Println("Generated code")
	case "startServer":
		err := server.StartServer(config.ServerPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Started Server")
	case "stopServer":
		err := server.StopServer(config.ServerPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Stopped Server")
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
