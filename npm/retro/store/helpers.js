import {
	STORE_KEY,
} from "./store-key"

export function isFunction(arg) {
	return typeof arg === "function"
}

export function isSelector(arg) {
	return arg !== undefined &&
		Array.isArray(arg) &&
		arg.length > 0
}

export function isStore(arg) {
	return arg?.$$key === STORE_KEY
}

export function querySelector(state, selector) {
	let focus = state
	for (const id of selector) {
		focus = focus[id]
	}
	return focus
}
