// Inspired by https://simpler-state.js.org
export function createStore(initialState, initializer) {
	return {
		/*
		 * Properties
		 */
		__currentState: typeof initializer === "function"
			? initializer(initialState)
			: initialState,
		__subscriptions: new Map(),
		/*
		 * Methods
		 */
		// Get the current state
		get(selector) {
			return typeof selector === "function"
				? selector(this.__currentState)
				: this.__currentState
		},
		// Broadcast the next state to subscribed components
		set(updater) {
			const currentState = this.__currentState
			let nextState = typeof updater === "function"
				? updater(currentState)
				: updater
			for (const [setState, selector] of this.__subscriptions) {
				if (selector === undefined) {
					// Rerender (no selectors to scope)
					setState(nextState)
				} else {
					// Suppress useless rerenders
					const selected1 = selector(currentState)
					const selected2 = selector(nextState)
					if (selected1 !== selected2) {
						setState(nextState)
					}
				}
			}
			this.__currentState = nextState
		},
		// Subscribe a component
		use(selector) {
			// Subscribe and unsubscribe components based on lifecycle events
			const [_, setState] = React.useState(this.__currentState)
			React.useEffect(() => {
				this.__subscriptions.set(setState, selector)
				return () => {
					this.__subscriptions.delete(setState)
				}
			}, [selector])
			return typeof selector === "function"
				? selector(this.__currentState)
				: this.__currentState
		},
	}
}
