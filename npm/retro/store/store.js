import * as helpers from "./helpers"
import STORE_KEY from "./STORE_KEY"

const ERR_BAD_STORE = originator => `${originator}: bad store; use 'createStore({ ... })'`
const ERR_BAD_REDUCER = originator => `${originator}: bad store reducer; use 'function reducer(state, action) { ... }'`
const ERR_BAD_SELECTOR = (originator, selector) => `${originator}: bad store selector; ${JSON.stringify(selector)}`

export function createStore(initialStateOrInitializer) {
	const initializerIsFunction = helpers.isFunction(initialStateOrInitializer)
	let initialState = initialStateOrInitializer
	if (initializerIsFunction) {
		initialState = initialStateOrInitializer()
	}
	return {
		$$key: STORE_KEY,
		// Component subscriptions
		subscriptions: new Map(),
		// Initial state
		initialState: initialState,
		// Initial state initializer
		initializer: !initializerIsFunction
			? () => initialState
			: initialStateOrInitializer,
		// Current cached state
		cachedState: initialState,
	}
}

function useStateImpl(store, { originator, flagIncludeState, flagIncludeSetState }) {
	React.useMemo(() => {
		if (!helpers.isStore(store)) {
			throw new Error(ERR_BAD_STORE(originator))
		}
	}, [])

	const [state, setState] = React.useState(store.cachedState)

	// Add 'setState' to the store's subscriptions
	React.useEffect(!flagIncludeState ? () => { /* No-op */ } : () => {
		store.subscriptions.set(setState, undefined /* selector=undefined */)
		return () => {
			store.subscriptions.delete(setState)
		}
	}, [])

	const setStore = React.useCallback(!flagIncludeSetState ? () => { /* No-op */ } : updater => {
		const currState = store.cachedState
		let nextState = updater
		if (typeof updater === "function") {
			nextState = updater(currState)
		}

		// Invalidate components
		setState(nextState)
		for (const [otherSetState, otherSelector] of store.subscriptions) {
			// Dedupe 'setState'
			if (setState !== otherSetState) {
				// Suppress useless rerenders
				if (otherSelector !== undefined) {
					const curr = helpers.querySelector(currState, otherSelector)
					const next = helpers.querySelector(nextState, otherSelector)
					if (curr === next) {
						continue
					}
				}
				otherSetState(nextState)
			}
		}

		// Cache the current state
		store.cachedState = nextState
	}, [])

	return [
		state,
		setStore,
	]
}

function useSelectorImpl(store, selector, { originator }) {
	React.useMemo(() => {
		if (!helpers.isStore(store)) {
			throw new Error(ERR_BAD_STORE(originator))
		}
		// Selector guards
		if (!helpers.isSelector(selector)) {
			throw new Error(ERR_BAD_SELECTOR(originator, selector))
		}
		let cachedRef = store.cachedState
		for (const key of selector) {
			if (!(key in cachedRef)) {
				throw new Error(ERR_BAD_SELECTOR(originator, selector))
			}
			cachedRef = cachedRef[key]
		}
	}, [])

	const [_, setState] = React.useState(store.cachedState)
	const valueOrReference = helpers.querySelector(store.cachedState, selector)

	const memoSelector = React.useMemo(() => {
		return selector
	}, [helpers.toPath(selector)])

	// Add 'setState' to the store's subscriptions
	React.useEffect(() => {
		store.subscriptions.set(setState, memoSelector)
		return () => {
			store.subscriptions.delete(setState)
		}
	}, [memoSelector])

	return valueOrReference
}

function useReducerImpl(store, reducer, { originator, flagIncludeState, flagIncludeSetState }) {
	React.useMemo(() => {
		if (!helpers.isStore(store)) {
			throw new Error(ERR_BAD_STORE(originator))
		}
		// Reducer guards
		if (!flagIncludeState) {
			if (!helpers.isFunction(reducer)) {
				throw new Error(ERR_BAD_REDUCER(originator))
			}
		}
	}, [])

	const [state, setState] = React.useState(store.cachedState)

	// Add 'setState' to the store's subscriptions
	React.useEffect(!flagIncludeState ? () => { /* No-op */ } : () => {
		store.subscriptions.set(setState, undefined /* selector=undefined */)
		return () => {
			store.subscriptions.delete(setState)
		}
	}, [])

	const dispatch = React.useCallback(!flagIncludeSetState ? () => { /* No-op */ } : action => {
		const currState = store.cachedState
		const nextState = reducer(currState, action)

		// Invalidate components
		setState(nextState)
		for (const [otherSetState, otherSelector] of store.subscriptions) {
			// Dedupe 'setState'
			if (setState !== otherSetState) {
				// Suppress useless rerenders
				if (otherSelector !== undefined) {
					const curr = helpers.querySelector(currState, otherSelector)
					const next = helpers.querySelector(nextState, otherSelector)
					if (curr === next) {
						continue
					}
				}
				otherSetState(nextState)
			}
		}

		// Cache the current state
		store.cachedState = nextState
	}, [])

	return [
		state,
		dispatch,
	]
}

export function useState(store) {
	return useStateImpl(store, {
		originator: "useState",
		flagIncludeState: true,
		flagIncludeSetState: true,
	})
}

export function useStateOnlyState(store) {
	return useStateImpl(store, {
		originator: "useStateOnlyState",
		flagIncludeState: true,
		flagIncludeSetState: false,
	})[0]
}

export function useStateOnlySetState(store) {
	return useStateImpl(store, {
		originator: "useStateOnlySetState",
		flagIncludeState: false,
		flagIncludeSetState: true,
	})[1]
}

export function useSelector(store, ...selector) {
	return useSelectorImpl(store, selector, {
		originator: "useSelector",
	})
}

export function useReducer(store, reducer) {
	return useReducerImpl(store, reducer, {
		originator: "useReducer",
		flagIncludeState: true,
		flagIncludeSetState: true,
	})
}

export function useReducerOnlyState(store, reducer) {
	return useReducerImpl(store, reducer, {
		originator: "useReducerOnlyState",
		flagIncludeState: true,
		flagIncludeSetState: false,
	})[0]
}

export function useReducerOnlyDispatch(store, reducer) {
	return useReducerImpl(store, reducer, {
		originator: "useReducerOnlyDispatch",
		flagIncludeState: false,
		flagIncludeSetState: true,
	})[1]
}

// Based on https://github.com/pelotom/use-methods/blob/master/src/index.ts
export function useCallbacks(store, methods) {
	const dispatch = useReducerOnlyDispatch(store, (state, action) => {
		return methods(state)[action.type](...action.payload)
	})

	const memoCallbacks = React.useMemo(() => {
		const callbacks = {}
		for (const key of Object.keys(methods(undefined))) {
			callbacks[key] = (...payload) => dispatch({
				type: key,
				payload,
			})
		}
		return callbacks
	}, [])

	return memoCallbacks
}
