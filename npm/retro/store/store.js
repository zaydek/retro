import {
	isFunction,
	isSelector,
	isStore,
	querySelector,
} from "./helpers"
import { STORE_KEY } from "./store-key"

const ERR_BAD_STORE = originator => `${originator}: Bad store; expected \`createStore({ ... })\`.`
const ERR_BAD_REDUCER = originator => `${originator}: Bad reducer; expected \`function reducer(state, action) { ... }\`.`
const ERR_BAD_SELECTOR = (originator, selector) => `${originator}: Bad selector; want \`["foo", "bar", ...]\` got ${JSON.stringify(selector)}.`

export function createStore(initialStateOrInitializer) {
	const initializerIsFunction = isFunction(initialStateOrInitializer)

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
		if (!isStore(store)) {
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
				if (isSelector(otherSelector)) {
					// Suppress useless rerenders
					const currSelected = querySelector(currState, otherSelector)
					const nextSelected = querySelector(nextState, otherSelector)
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

function useReducerImpl(store, reducer, { originator, flagIncludeState, flagIncludeSetState }) {
	React.useMemo(() => {
		if (!isStore(store)) {
			throw new Error(ERR_BAD_STORE(originator))
		}
		if (!flagIncludeState) {
			if (!isFunction(reducer)) {
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
				if (isSelector(otherSelector)) {
					// Suppress useless rerenders
					const currSelected = querySelector(currState, otherSelector)
					const nextSelected = querySelector(nextState, otherSelector)
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

function useSelectorImpl(store, selector, { originator }) {
	React.useMemo(() => {
		if (!isStore(store)) {
			throw new Error(ERR_BAD_STORE(originator))
		}
		if (!isSelector(selector)) {
			throw new Error(ERR_BAD_SELECTOR(originator, selector))
		}
		let focus = store.cachedState
		for (const id of selector) {
			if (!(id in focus)) {
				throw new Error(ERR_BAD_SELECTOR(originator, selector))
			}
			focus = focus[id]
		}
	}, [])

	const [state, setState] = React.useState(store.cachedState)
	const valueOrReference = querySelector(state, selector)

	const memoSelector = React.useMemo(() => {
		return selector
	}, [selector])

	// Add the 'setState' to the store's subscriptions
	React.useEffect(() => {
		store.subscriptions.set(setState, memoSelector)
		return () => {
			store.subscriptions.delete(setState)
		}
	}, [memoSelector])

	return valueOrReference
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

export function useSelector(store, selector) {
	return useSelectorImpl(store, selector, {
		originator: "useSelector",
	})
}
