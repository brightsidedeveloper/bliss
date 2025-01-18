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
		if op.QueryParams != nil {
			code += ` ` + op.Method + op.Name + `Query,`
		}
		if op.Body != nil {
			code += ` ` + op.Method + op.Name + `Body,`
		}
		code += ` ` + op.Method + op.Name + `Res,`
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
		goCode += `
	` + op.Name + `(params: ` + op.Method + op.Name + `Query) {
	return get<` + op.Method + op.Name + `Res>('` + op.Endpoint + `', params as unknown as Record<string, unknown>)
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
		goCode += `
	` + op.Name + `(body: ` + op.Method + op.Name + `Body) {
	return post<` + op.Method + op.Name + `Res>('` + op.Endpoint + `', body)
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
		goCode += `
	` + op.Name + `(body: ` + op.Method + op.Name + `Body) {
	return put<` + op.Method + op.Name + `Res>('` + op.Endpoint + `', body)
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
		goCode += `
	` + op.Name + `(body: ` + op.Method + op.Name + `Body) {
	return patch<` + op.Method + op.Name + `Res>('` + op.Endpoint + `', body)
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
		goCode += `
	` + op.Name + `(body: ` + op.Method + op.Name + `Body) {
	return del<` + op.Method + op.Name + `Res>('` + op.Endpoint + `', body as unknown as Record<string, unknown>)
},`

	}

	return goCode
}
