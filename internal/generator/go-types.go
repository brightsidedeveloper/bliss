package generator

import (
	"master-gen/internal/parser"
	"master-gen/internal/util"
	"master-gen/internal/writer"
)

func genGoTypes(ops parser.Bliss, path string) error {

	goCode := `package types

	`
	for _, op := range ops.Operations {
		theType := "Row"
		if op.Handler != "" {
			theType = "Res"
		}
		goCode += `
type ` + op.Name + theType + ` struct {`
		for k, t := range op.Res {

			goCode += `
	` + util.ConvertToCamelCase(k) + ` ` + t + " `" + `json:"` + k + `"` + "`"
		}
		goCode += `
}
	`
		if op.QueryParams != nil {
			goCode += `
type ` + op.Name + `Query struct {`
			for k, t := range op.QueryParams {
				goCode += `
	` + util.ConvertToCamelCase(k) + ` ` + t + " `" + `json:"` + k + `"` + "`"
			}
			goCode += `
}

	`
		}

		if op.Body != nil {
			goCode += `
type ` + op.Name + `Body struct {`
			for k, t := range op.Body {
				goCode += `
	` + util.ConvertToCamelCase(k) + ` ` + t + " `" + `json:"` + k + `"` + "`"
			}
			goCode += `
}
		`
		}
	}

	writer.WriteFile(path+"/types/types.go", goCode)

	return nil
}
