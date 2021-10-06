import * as store from "../store"

import {
	getBrowserPathSSR
} from "./helpers"

import {
	PUSH_STATE,
	REPLACE_STATE,

	routerStore,
} from "./store"

// Syncs the window state to the router state
export function useSyncWindowToRouter() {
	const setState = store.useStateOnlySetState(routerStore)
	React.useEffect(() => {
		function handlePopState() {
			setState({
				type: REPLACE_STATE,
				path: getBrowserPathSSR(),
				scrollTo: undefined,
			})
		}
		window.addEventListener("popstate", handlePopState)
		return () => window.removeEventListener("popstate", handlePopState)
	}, [])
}

// Syncs the router state to the window state
export function useSyncRouterToWindow() {
	const state = store.useStateOnlyState(routerStore)
	const didMountRef = React.useRef(false)
	React.useEffect(() => {
		if (!didMountRef.current) {
			didMountRef.current = true
			return
		}
		if (state.path !== getBrowserPathSSR()) {
			if (state.type === REPLACE_STATE) {
				window.history.replaceState({}, "", state.path)
			} else if (state.type === PUSH_STATE) {
				window.history.pushState({}, "", state.path)
			}
		}
		if (state.scrollTo !== undefined) {
			window.scrollTo(0, state.scrollTo)
		}
	}, [state])
}
