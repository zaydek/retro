import { useLayoutEffectSSR } from "../use-layout-effect-ssr"

export function LayoutTitle({ title, children }) {
	useLayoutEffectSSR(() => {
		document.title = title
	}, [title])
	return children
}

export function Title({ title, children }) {
	React.useEffect(() => {
		document.title = title
	}, [title])
	return children
}
