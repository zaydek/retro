import * as store from "../store"

import {
	getBrowserPathSSR,
} from "./helpers"

export const REPLACE_STATE = "REPLACE_STATE"
export const PUSH_STATE = "PUSH_STATE"

export const routerStore = store.createStore({
	type: PUSH_STATE,
	path: getBrowserPathSSR(),
	scrollTo: undefined,
})
