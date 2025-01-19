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
		if err := genTsTypes(bliss, g.WebPath+"/src/api"); err != nil {
			return err
		}
		if err := genBlissClient(bliss, g.WebPath+"/src/api"); err != nil {
			return err
		}
	}
	return nil
}
