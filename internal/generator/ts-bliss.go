package generator

import (
	"fmt"
	"master-gen/internal/parser"
	"master-gen/internal/util"
	"master-gen/internal/writer"
)

func genBlissClient(bliss parser.Bliss, path string) error {

	code := genRequestImport(map[string]bool{"Post": true})
	code += genTypeImports(bliss)
	code += `

export default class Bliss {`

	for _, op := range bliss.Operations {

		name := op.Query
		if op.Handler != "" {
			name = op.Handler
		}

		code += fmt.Sprintf(`
	static %s = () => {
		return post<%s>("%s")
	}`, name, name+"Res", op.Endpoint)
	}

	code += `
}
	`

	return writer.WriteFile(path+"/bliss.ts", code)
}

func genRequestImport(methods map[string]bool) string {
	code := ""
	if len(methods) == 0 {
		return code
	}
	code += `import { `
	for method := range methods {
		requestMethod := util.Lower(method)
		if requestMethod == "delete" {
			requestMethod = "del"
		}
		code += requestMethod + `, `
	}
	code += ` } from './request'`
	return code
}

func genTypeImports(bliss parser.Bliss) string {
	if len(bliss.Operations) == 0 {
		return ""
	}
	code := `
import { `
	for _, op := range bliss.Operations {
		name := op.Query
		if op.Handler != "" {
			name = op.Handler
		}
		if false {
			// Params
		}
		code += ` ` + name + `Res,`
	}
	code += ` } from './types'`
	return code
}
