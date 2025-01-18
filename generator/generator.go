package generator

import (
	"master-gen/parser"
)

type Generator struct {
	GenesisPath string
	Bliss       parser.AhaJSON
}

func (g *Generator) Generate() error {
	// Go
	if err := genGoTypes(g.Bliss, g.GenesisPath); err != nil {
		return err
	}
	if err := genHandlers(g.Bliss, g.GenesisPath); err != nil {
		return err
	}
	if err := genMountRoutes(g.Bliss, g.GenesisPath); err != nil {
		return err
	}
	// TS

	return nil
}
