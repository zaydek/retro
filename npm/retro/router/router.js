import * as store from "../store"

import {
	useSyncWindowToRouter,
	useSyncRouterToWindow,
} from "./hooks"

import {
	PUSH_STATE,
	REPLACE_STATE,

	routerStore,
} from "./store"

import {
	useLayoutEffectSSR,
} from "../use-layout-effect-ssr"

export function Link({ path, scrollTo, children, ...props }) {
	const setState = store.useStateOnlySetState(routerStore)

	const isLocallyScoped = React.useMemo(() => {
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
		// TODO: Right now we are setting `REPLACE_STATE` and `PUSH_STATE` events
		// when a user clicks on a link that points to a route that redirects. This
		// happens because we aren't checking for route presence here.
		e.preventDefault()
		setState({
			type: PUSH_STATE,
			path,
			scrollTo,
		})
	}

	if (isLocallyScoped) {
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
		setState({
			type: REPLACE_STATE,
			path,
			scrollTo,
		})
	}, [path, scrollTo])
	return null
}

export function Route({ path: _, children }) {
	return children
}

function getSurroundingChildrenAndRouteMap(children) {
	const above = []    // Components above the route
	const below = []    // Components below the route
	const routeMap = {} // Maps paths to route components

	function isRouteAndDistinct(component) {
		return component?.type === Route &&              // <Route>
			typeof component?.props?.path === "string" &&  // <Route path="/hello">
			routeMap[component?.props?.path] === undefined // <Route path="/world">
	}

	const components = [children].flat()
	let isAbove = true
	for (let componentIndex = 0; componentIndex < components.length; componentIndex++) {
		const component = components[componentIndex]
		if (isRouteAndDistinct(component)) {
			routeMap[component.props.path] = component
			isAbove = false
			continue
		}
		if (isAbove) {
			above.push(component)
		} else {
			below.push(component)
		}
	}

	return [above, below, routeMap]
}

export function Router({ children }) {
	const state = store.useStateOnlyState(routerStore)

	useSyncWindowToRouter()
	useSyncRouterToWindow()

	const [above, below, routeMap] = React.useMemo(() => {
		return getSurroundingChildrenAndRouteMap(children)
	}, [children])

	return (
		<>
			{above}
			{routeMap[state.path] ??
				routeMap["/404"]}
			{below}
		</>
	)
}
