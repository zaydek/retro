export function getCurrentPathSSR() {
	let path = "/"
	if (typeof process.env.LOCATION !== "undefined") {
		path = process.env.LOCATION
	}
	if (path.endsWith(".html")) {
		path = path.slice(0, -".html".length)
		if (path.endsWith("/index")) {
			path = path.slice(0, -"/index".length)
		}
	}
	return path
}
