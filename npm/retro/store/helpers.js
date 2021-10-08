import {
	STORE_KEY,
} from "./store-key"

export function isFunction(v) {
	return typeof v === "function"
}

export function isSelector(v) {
	return v !== undefined && Array.isArray(v) && v.length > 0
}

export function isStore(v) {
	return v?.$$key === STORE_KEY
}

export function querySelector(state, selector) {
	let selected = state
	for (const id of selector) {
		selected = selected[id]
	}
	return selected
}
