import * as store from "../store"

import { getCurrentPathSSR } from "./helpers"
import { useLayoutEffectSSR } from "../use-layout-effect-ssr"
import { useSyncWindowToRouter, useSyncRouterToWindow } from "./hooks"

export const actions = {
	REPLACE_STATE: "REPLACE_STATE",
	PUSH_STATE: "PUSH_STATE",
}

export const routerStore = store.createStore({
	type: actions.PUSH_STATE,
	path: getCurrentPathSSR(),
	scrollTo: 0,
})

function evaluateScrollTo(scrollTo) {
	if (typeof scrollTo === "function") {
		return scrollTo()
	}
	return scrollTo ?? 0
}

export function Link({ path, scrollTo, children, ...props }) {
	const setState = store.useStateOnlySetState(routerStore)

	const flagIsLocal = React.useMemo(() => {
		return path.startsWith("/") || (
			!path.startsWith("https://") &&
			!path.startsWith("http://") &&
			!path.startsWith("www.")
		)
	}, [path])

	function handleClick(e) {
		// https://github.com/molefrog/wouter
		if (e.button > 0 || e.shiftKey || e.ctrlKey || e.altKey || e.metaKey) {
			return
		}
		e.preventDefault()
		setState(current => ({
			...current,
			type: actions.PUSH_STATE,
			path,
			scrollTo: evaluateScrollTo(scrollTo),
		}))
	}

	if (flagIsLocal) {
		return (
			<a href={path} onClick={handleClick} {...props}>
				{children}
			</a>
		)
	} else {
		return (
			<a href={path} target="_blank" rel="noreferrer noopener" {...props}>
				{children}
			</a>
		)
	}
}

export function Redirect({ path, scrollTo }) {
	const setState = store.useStateOnlySetState(routerStore)
	useLayoutEffectSSR(() => {
		setState(current => ({
			...current,
			type: actions.REPLACE_STATE,
			path,
			scrollTo,
		}))
	}, [path, scrollTo])
	return null
}

export function Route({ path: _, children }) {
	return children
}

export function Router({ children }) {
	const path = store.useSelector(routerStore, ["path"])

	useSyncWindowToRouter()
	useSyncRouterToWindow()

	const [Surrounding, RenderRoute] = React.useMemo(() => {
		const routeMap = {} // Maps paths to route components
		const above = []
		const below = []

		let flagIsAbove = false
		for (const component of [...children].flat()) {
			if (component?.type === Route && typeof component?.props?.path === "string") {
				routeMap[component.props.path] = component
				flagIsAbove = true
				continue
			}
			if (!flagIsAbove) {
				above.push(component)
			} else {
				below.push(component)
			}
		}

		// Higher-order component for rendering the above and below components
		// around the render route
		function Surrounding({ children }) {
			return <>{above}{children}{below}</>
		}

		// The render route or the '/404' route or 'null'
		function RenderRoute() {
			return routeMap[path] ?? routeMap["/404"] ?? null
		}

		return [Surrounding, RenderRoute]
	}, [children, path])

	return (
		<Surrounding>
			<RenderRoute />
		</Surrounding>
	)
}
