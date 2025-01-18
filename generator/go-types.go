package generator

import (
	"master-gen/parser"
	"master-gen/util"
	"master-gen/writer"
)

func TypesFileContent(ops parser.AhaJSON) error {

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
	` + util.Capitalize(k) + ` ` + t + " `" + `json:"` + k + `"` + "`"
		}
		goCode += `
}
	`
		if op.QueryParams != nil {
			goCode += `
type ` + op.Name + `Query struct {`
			for k, t := range op.QueryParams {
				goCode += `
	` + util.Capitalize(k) + ` ` + t + " `" + `json:"` + k + `"` + "`"
			}
			goCode += `
}

	`
		}
	}

	writer.WriteFile("./app/genesis/types/types.go", goCode)

	return nil
}
