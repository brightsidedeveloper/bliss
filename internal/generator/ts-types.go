package generator

import (
	"master-gen/internal/parser"
	"master-gen/internal/writer"
)

func genTsTypes(bliss parser.Bliss, path string) error {

	code := ``

	for _, op := range bliss.Operations {
		// 		if op.QueryParams != nil {
		// 			code += `
		// export interface ` + op.Method + op.Name + `Query {`
		// 			for k, t := range op.QueryParams {
		// 				code += `
		// 	` + k + `: ` + t + `;`
		// 			}
		// 			code += `
		// }
		// `
		// 		}
		// 		if op.Body != nil {
		// 			code += `
		// export interface ` + op.Method + op.Name + `Body {`
		// 			for k, t := range op.Body {
		// 				code += `
		// 	` + k + `: ` + t + `;`
		// 			}
		// 			code += `
		// }
		// `
		// 		}

		// 		if op.Res != nil {
		// 			code += `
		// export interface ` + op.Method + op.Name + `Res {`
		// 			for k, t := range op.Res {
		// 				code += `
		// 	` + k + `: ` + convertGoTypeToTs(t) + `;`
		// 			}
		// 			code += `
		// }
		// 			`
		// 		} else {
		code += `
export interface ` + op.Query + `Res {}`
	}

	return writer.WriteFile(path+"/types.ts", code)
}

func convertGoTypeToTs(t string) string {
	switch t {
	case "int":
		return "number"
	case "string":
		return "string"
	case "bool":
		return "boolean"
	default:
		return "any"
	}
}
