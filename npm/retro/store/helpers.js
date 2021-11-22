import STORE_KEY from "./STORE_KEY"

export function isFunction(arg) {
	return typeof arg === "function"
}

export function isSelector(arg) {
	return arg !== undefined && Array.isArray(arg) && arg.length > 0
}

export function isStore(arg) {
	return arg?.$$key === STORE_KEY
}

export function querySelector(state, selector) {
	let focus = state
	for (const key of selector) {
		focus = focus[key]
	}
	return focus
}
