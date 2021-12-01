// Inspired by https://simpler-state.js.org
export function createStore(initialStateOrInitializer) {
	return {
		/*
		 * Properties
		 */
		__currentState: typeof initialStateOrInitializer === "function"
			? initialStateOrInitializer()
			: initialStateOrInitializer,
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
			// Cache the next state eagerly
			this.__currentState = nextState
			for (const [setState, selector] of this.__subscriptions) {
				if (selector === undefined) {
					// Force rerender
					setState(next)
				} else {
					// Suppress useless rerenders
					const selected1 = selector(currentState)
					const selected2 = selector(nextState)
					if (selected1 !== selected2) {
						setState(next)
					}
				}
			}
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
