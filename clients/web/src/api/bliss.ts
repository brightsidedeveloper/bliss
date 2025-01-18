import { get,  } from './request'
import {  GetExampleQuery, GetExampleRes, } from './types'

export default class Bliss {
	static get = {
	Example(params: GetExampleQuery) {
	return get<GetExampleRes>('/api/v1/example', params as unknown as Record<string, unknown>)
},
	}
}
	