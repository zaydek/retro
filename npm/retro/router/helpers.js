// Gets the browser path or the home path
export function getBrowserPathSSR() {
	let path = "/"
	if (typeof window !== "undefined") {
		path = window.location.pathname
	}
	if (path.endsWith(".html")) {
		path = path.slice(0, -".html".length)
		if (path.endsWith("/index")) {
			path = path.slice(0, -"/index".length)
		}
	}
	return path
}
