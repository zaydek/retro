import {
	isFunction,
	isStore,
	isValidSelector,
	query,
} from "./helpers"

import {
	STORE_KEY,
} from "./store-key"

export function createStore(initialStateOrInitializer) {
	const initializerIsFunction = isFunction(initialStateOrInitializer)

	let initialState = initialStateOrInitializer
	if (initializerIsFunction) {
		initialState = initialStateOrInitializer()
	}

	return {
		// Reference for checking whether a store is a store
		key: STORE_KEY,
		// Subscriptions for all setters (`setState`)
		subscriptions: new Map(),
		// Initial state
		initialState: initialState,
		// Initializer
		initializer: !initializerIsFunction
			? () => initialState
			: initialStateOrInitializer,
		// Cached state
		cachedState: initialState,
	}
}

function useStateImpl(store, { flagIncludeState, flagIncludeSetState }) {
	React.useMemo(() => {
		if (!isStore(store)) {
			throw new Error("useState: First argument is not a store. " +
				"Use `createStore` to create a store.")
		}
	}, [])

	const [state, setState] = React.useState(store.cachedState)

	// Add the `setState` to the store's subscriptions
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
			// Dedupe `setState`
			if (otherSetState !== setState) {
				if (isValidSelector(otherSelector)) {
					// Suppress useless rerenders
					const currSelected = query(currState, otherSelector)
					const nextSelected = query(nextState, otherSelector)
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

function useReducerImpl(store, reducer, { flagIncludeState, flagIncludeSetState }) {
	React.useMemo(() => {
		if (!isStore(store)) {
			throw new Error("useReducer: First argument is not a store. " +
				"Use `createStore` to create a store.")
		}
		if (!flagIncludeState) {
			if (!isFunction(reducer)) {
				throw new Error("useReducer: Second argument is not a reducer.")
			}
		}
	}, [])

	const [state, setState] = React.useState(store.cachedState)

	// Add the `setState` to the store's subscriptions
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
			// Dedupe `setState`
			if (otherSetState !== setState) {
				if (isValidSelector(otherSelector)) {
					// Suppress useless rerenders
					const currSelected = query(currState, otherSelector)
					const nextSelected = query(nextState, otherSelector)
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

function useSelectorImpl(store, selector) {
	React.useMemo(() => {
		if (!isStore(store)) {
			throw new Error("useSelector: First argument is not a store. " +
				"Use `createStore` to create a store.")
		}
		if (!isValidSelector(selector)) {
			throw new Error("useSelector: Second argument is not a selector.")
		}
		let focus = store.cachedState
		for (const id of selector) {
			if (!(id in focus)) {
				throw new Error(`useSelector: Selector path \`[${selector.join(", ")}]\` is unreachable.`)
			}
			focus = focus[id]
		}
	}, [])

	const [state, setState] = React.useState(store.cachedState)
	const valueOrReference = query(state, selector)

	// Add the `setState` to the store's subscriptions
	React.useEffect(() => {
		store.subscriptions.set(setState, selector)
		return () => {
			store.subscriptions.delete(setState)
		}
	}, [selector])

	return valueOrReference
}

export function useState(store) {
	return useStateImpl(store, {
		flagIncludeState: true,
		flagIncludeSetState: true,
	})
}

export function useStateOnlyState(store) {
	return useStateImpl(store, {
		flagIncludeState: true,
		flagIncludeSetState: false,
	})[0]
}

export function useStateOnlySetState(store) {
	return useStateImpl(store, {
		flagIncludeState: false,
		flagIncludeSetState: true,
	})[1]
}

export function useReducer(store, reducer) {
	return useReducerImpl(store, reducer, {
		flagIncludeState: true,
		flagIncludeSetState: true,
	})
}

export function useReducerOnlyState(store, reducer) {
	return useReducerImpl(store, reducer, {
		flagIncludeState: true,
		flagIncludeSetState: false,
	})[0]
}

export function useReducerOnlyDispatch(store, reducer) {
	return useReducerImpl(store, reducer, {
		flagIncludeState: false,
		flagIncludeSetState: true,
	})[1]
}

export function useSelector(store, selector) {
	return useSelectorImpl(store, selector)
}
