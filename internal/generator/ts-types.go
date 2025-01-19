package generator

import (
	"master-gen/internal/parser"
	"master-gen/internal/writer"
)

func genTsTypes(bliss parser.Bliss, path string) error {

	code := ``

	for _, op := range bliss.Operations {
		name := op.Query
		if op.Handler != "" {
			name = op.Handler
		}
		code += `
export interface ` + name + `Res {}`
	}

	return writer.WriteFile(path+"/types.ts", code)
}
