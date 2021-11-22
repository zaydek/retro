import * as helpers from "./helpers"
import next from "./next"
import STORE_KEY from "./store-key"

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

	// Add the 'setState' to the store's subscriptions
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
			if (otherSetState !== setState) {
				if (helpers.isSelector(otherSelector)) {
					// Suppress useless rerenders
					const currSelected = helpers.querySelector(currState, otherSelector)
					const nextSelected = helpers.querySelector(nextState, otherSelector)
					if (currSelected !== nextSelected) {
						otherSetState(nextState)
					}
				} else {
					otherSetState(nextState)
				}
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

function useSelectorImpl(store, selector, { originator, flagIncludeState, flagIncludeSetState }) {
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

	const [state, setState] = React.useState(store.cachedState)
	const valueOrReference = helpers.querySelector(state, selector)

	// Add the 'setState' to the store's subscriptions
	React.useEffect(!flagIncludeState ? () => { /* No-op */ } : () => {
		store.subscriptions.set(setState, selector)
		return () => {
			store.subscriptions.delete(setState)
		}
	}, [])

	const setStore = React.useCallback(!flagIncludeSetState ? () => { /* No-op */ } : updater => {
		const currState = store.cachedState
		let nextState = next(currState, selector, updater)

		// Invalidate components
		setState(nextState)
		for (const [otherSetState, otherSelector] of store.subscriptions) {
			// Dedupe 'setState'
			if (otherSetState !== setState) {
				if (helpers.isSelector(otherSelector)) {
					// Suppress useless rerenders
					const currSelected = helpers.querySelector(currState, otherSelector)
					const nextSelected = helpers.querySelector(nextState, otherSelector)
					if (currSelected !== nextSelected) {
						otherSetState(nextState)
					}
				} else {
					otherSetState(nextState)
				}
			}
		}

		// Cache the current state
		store.cachedState = nextState
	}, [])

	return [
		valueOrReference,
		setStore,
	]
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

	// Add the 'setState' to the store's subscriptions
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
			if (otherSetState !== setState) {
				if (helpers.isSelector(otherSelector)) {
					// Suppress useless rerenders
					const currSelected = helpers.querySelector(currState, otherSelector)
					const nextSelected = helpers.querySelector(nextState, otherSelector)
					if (currSelected !== nextSelected) {
						otherSetState(nextState)
					}
				} else {
					otherSetState(nextState)
				}
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

export function useOnlyState(store) {
	return useStateImpl(store, {
		originator: "useOnlyState",
		flagIncludeState: true,
		flagIncludeSetState: false,
	})[0]
}

export function useOnlySetState(store) {
	return useStateImpl(store, {
		originator: "useOnlySetState",
		flagIncludeState: false,
		flagIncludeSetState: true,
	})[1]
}

export function useSelector(store, ...args) {
	const [selector, updater] = [
		args.slice(0, args.length - 1),
		args[args.length - 1],
	]
	return useSelectorImpl(store, selector, updater, {
		originator: "useSelector",
		flagIncludeState: true,
		flagIncludeSetState: true,
	})
}

export function useOnlySelector(store, ...args) {
	const [selector, updater] = [
		args.slice(0, args.length - 1),
		args[args.length - 1],
	]
	return useSelectorImpl(store, selector, updater, {
		originator: "useOnlySelector",
		flagIncludeState: true,
		flagIncludeSetState: false,
	})
}

export function useOnlySetSelector(store, ...args) {
	const [selector, updater] = [
		args.slice(0, args.length - 1),
		args[args.length - 1],
	]
	return useSelectorImpl(store, selector, updater, {
		originator: "useOnlySetSelector",
		flagIncludeState: false,
		flagIncludeSetState: true,
	})
}

export function useReducer(store, reducer) {
	return useReducerImpl(store, reducer, {
		originator: "useReducer",
		flagIncludeState: true,
		flagIncludeSetState: true,
	})
}

export function useOnlyDispatch(store, reducer) {
	return useReducerImpl(store, reducer, {
		originator: "useOnlyDispatch",
		flagIncludeState: false,
		flagIncludeSetState: true,
	})[1]
}
