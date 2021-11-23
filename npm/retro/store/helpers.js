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

export function toPath(selector) {
	return `[${selector.map(key => JSON.stringify(key)).join("][")}]`
}

export function querySelector(state, selector) {
	let stateRef = state
	for (const key of selector) {
		stateRef = stateRef[key]
	}
	return stateRef
}
