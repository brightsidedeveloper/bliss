package generator

import (
	"fmt"
	"master-gen/internal/parser"
	"path"
)

type Generator struct {
	BlissPath  string
	ServerPath string
	WebPath    string
}

func (g *Generator) Generate() error {
	bliss, err := parser.GetBliss(g.BlissPath)
	if err != nil {
		return err
	}
	if err := createServerIfNotExists(g.ServerPath); err != nil {
		return err
	}
	// Go
	if err := genSqlc(); err != nil {
		return fmt.Errorf("genSqlc: %w", err)
	}
	if err := genHandlers(g, bliss, path.Join(g.ServerPath, "genesis")); err != nil {
		return err
	}
	if err := genMountRoutes(bliss, path.Join(g.ServerPath, "genesis")); err != nil {
		return err
	}
	// TS
	if g.WebPath != "" {
		if err := createWebIfNotExists(g.WebPath); err != nil {
			return err
		}

		if err := genRequest(g.WebPath + "/src/api"); err != nil {
			return err
		}

		params := getParams(bliss)
		responses := getResponses(bliss)

		if err := genTsTypes(params, responses, g.WebPath+"/src/api"); err != nil {
			return err
		}
		if err := genBlissClient(bliss, params, responses, g.WebPath+"/src/api"); err != nil {
			return err
		}
	}
	return nil
}

type Params struct {
	interfaceType string // "" or "interface ExampleQuery = { Wow: string }"
	args          string // "" or "params: ExampleQuery"
	use           string // "" or ", params"
}

type Response struct {
	interfaceType string // "" or "interface ExampleRes {}"
	use           string // "" or "<ExampleRes>"
}

type ParamStrings map[string]Params
type ResponseStrings map[string]Response

func getParams(bliss parser.Bliss) ParamStrings {
	params := make(ParamStrings)
	for _, op := range bliss.Operations {
		name := op.Query
		if op.Handler != "" {
			name = op.Handler
		}
		params[name] = Params{
			interfaceType: "",
			args:          "",
			use:           "",
		}
	}
	return params

}

func getResponses(bliss parser.Bliss) ResponseStrings {
	responses := make(ResponseStrings)
	for _, op := range bliss.Operations {
		name := op.Query
		if op.Handler != "" {
			name = op.Handler
		}
		responses[name] = Response{
			interfaceType: "",
			use:           "",
		}
	}
	return responses
}
