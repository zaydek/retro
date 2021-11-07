export function getCurrentPathSSR() {
	let path = "/"
	if (typeof __location__ !== "undefined") {
		path = __location__
	}
	if (path.endsWith(".html")) {
		path = path.slice(0, -".html".length)
		if (path.endsWith("/index")) {
			path = path.slice(0, -"/index".length)
		}
	}
	return path
}
