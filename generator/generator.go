package generator

import "master-gen/parser"

type Generator struct {
	BlissPath   string
	GenesisPath string
	ApiPaths    []string
}

func (g *Generator) Generate() error {
	bliss, err := parser.GetBliss(g.BlissPath)
	if err != nil {
		return err
	}

	// Go
	if err := genGoTypes(bliss, g.GenesisPath); err != nil {
		return err
	}
	if err := genHandlers(bliss, g.GenesisPath); err != nil {
		return err
	}
	if err := genMountRoutes(bliss, g.GenesisPath); err != nil {
		return err
	}
	// TS

	return nil
}
