import * as store from "../store"

import {
	getCurrentPathSSR,
} from "./helpers"

import {
	actions,
	routerStore,
} from "./router"

export function useSyncWindowToRouter() {
	const setState = store.useStateOnlySetState(routerStore)
	React.useEffect(() => {
		function handlePopState() {
			setState(current => ({
				...current,
				type: actions.REPLACE_STATE,
				path: getCurrentPathSSR(),
				scrollTo: undefined,
			}))
		}
		window.addEventListener("popstate", handlePopState)
		return () => window.removeEventListener("popstate", handlePopState)
	}, [])
}

export function useSyncRouterToWindow() {
	const type = store.useSelector(routerStore, ["type"])
	const path = store.useSelector(routerStore, ["path"])
	const scrollTo = store.useSelector(routerStore, ["scrollTo"])
	const didMountRef = React.useRef(false)
	React.useEffect(() => {
		if (!didMountRef.current) {
			didMountRef.current = true
			return
		}
		if (path !== getCurrentPathSSR()) {
			if (type === actions.REPLACE_STATE) {
				// TODO: Add support for push or replacing relative URLs. For example:
				//
				//   window.history.pushState({}, "", path.startsWith("/")
				//     ? path
				//     : window.location.pathname + "/" + path
				//   )
				//
				window.history.replaceState({}, "", path)
			} else if (type === actions.PUSH_STATE) {
				// TODO: Add support for push or replacing relative URLs. For example:
				//
				//   window.history.pushState({}, "", path.startsWith("/")
				//     ? path
				//     : window.location.pathname + "/" + path
				//   )
				//
				window.history.pushState({}, "", path)
			}
		}
		if (scrollTo !== undefined) {
			window.scrollTo(0, scrollTo)
		}
	}, [
		type,
		path,
		scrollTo,
	])
}
