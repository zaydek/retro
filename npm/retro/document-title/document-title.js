import { useLayoutEffectSSR } from "../use-layout-effect-ssr"

export function LayoutDocumentTitle({ title, children }) {
	useLayoutEffectSSR(() => {
		document.title = title
	}, [title])
	return children
}

export function DocumentTitle({ title, children }) {
	React.useEffect(() => {
		document.title = title
	}, [title])
	return children
}
