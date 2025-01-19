package generator

import (
	"master-gen/internal/parser"
	"master-gen/internal/util"
	"master-gen/internal/writer"
)

func genBlissClient(bliss parser.Bliss, path string) error {
	methods := getMethodsUsed(bliss)
	code := genRequestImport(methods)
	code += genTypeImports(bliss)
	code += `

export default class Bliss {`

	for method := range methods {

		requestMethod := util.Lower(method)
		if requestMethod == "delete" {
			requestMethod = "del"
		}

		code += `
	static ` + util.Lower(method) + ` = {`

		switch method {
		case "Get":
			code += genGets(bliss)
		case "Post":
			code += genPosts(bliss)
		case "Put":
			code += genPuts(bliss)
		case "Patch":
			code += genPatches(bliss)
		case "Delete":
			code += genDeletes(bliss)
		}

		code += `
	}`
	}

	code += `
}
	`

	return writer.WriteFile(path+"/bliss.ts", code)
}

func getMethodsUsed(bliss parser.Bliss) map[string]bool {
	methods := make(map[string]bool)
	for _, op := range bliss.Operations {
		methods[op.Method] = true
	}
	return methods
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
		// if op.QueryParams != nil {
		// 	code += ` ` + op.Query + `Query,`
		// }
		// if op.Body != nil {
		// 	code += ` ` + op.Query + `Body,`
		// }
		code += ` ` + op.Query + `Res,`
	}
	code += ` } from './types'`
	return code
}

func genGets(bliss parser.Bliss) string {
	goCode := ""

	for _, op := range bliss.Operations {
		if op.Method != "Get" {
			continue
		}
		args := ""
		params := ""
		// if op.QueryParams != nil {
		// 	args = `params: ` + op.Query + `Query`
		// 	params = `, params as unknown as Record<string, unknown>`
		// }
		goCode += `
	` + op.Query + `(` + args + `) {
	return get<` + op.Query + `Res>('` + op.Endpoint + `'` + params + `)
},`

	}

	return goCode
}

func genPosts(bliss parser.Bliss) string {
	goCode := ""

	for _, op := range bliss.Operations {
		if op.Method != "Post" {
			continue
		}

		args := ""
		body := ""
		// if op.Body != nil {
		// 	args = `body: ` + op.Query + `Body`
		// 	body = `, body as unknown as Record<string, unknown>`
		// }

		goCode += `
	` + op.Query + `(` + args + `) {
	return post<` + op.Query + `Res>('` + op.Endpoint + `'` + body + `)
},`

	}

	return goCode
}

func genPuts(bliss parser.Bliss) string {
	goCode := ""

	for _, op := range bliss.Operations {
		if op.Method != "Put" {
			continue
		}
		args := ""
		body := ""
		// if op.Body != nil {
		// 	args = `body: ` + op.Query + `Body`
		// 	body = `, body as unknown as Record<string, unknown>`
		// }

		goCode += `
	` + op.Query + `(` + args + `) {
	return put<` + op.Query + `Res>('` + op.Endpoint + `'` + body + `)
},`

	}

	return goCode
}

func genPatches(bliss parser.Bliss) string {
	goCode := ""

	for _, op := range bliss.Operations {
		if op.Method != "Patch" {
			continue
		}
		args := ""
		body := ""
		// if op.Body != nil {
		// 	args = `body: ` + op.Query + `Body`
		// 	body = `, body as unknown as Record<string, unknown>`
		// }

		goCode += `
	` + op.Query + `(` + args + `) {
	return patch<` + op.Query + `Res>('` + op.Endpoint + `'` + body + `)
},`

	}

	return goCode
}

func genDeletes(bliss parser.Bliss) string {
	goCode := ""

	for _, op := range bliss.Operations {
		if op.Method != "Delete" {
			continue
		}
		args := ""
		body := ""
		// if op.Body != nil {
		// 	args = `body: ` + op.Query + `Body`
		// 	body = `, body as unknown as Record<string, unknown>`
		// }

		goCode += `
	` + op.Query + `(` + args + `) {
	return del<` + op.Query + `Res>('` + op.Endpoint + `'` + body + `)
},`

	}

	return goCode
}
