// Retro sends server-sent events (SSE) to ` + "`" + `/__dev__` + "`" + ` to force-reload
// tabs. Because server-sent events propagate to ~6 tabs max (depending on
// the browser), use localStorage as a mechanism to force-reload other
// tabs.
//
// https://github.com/evanw/esbuild/issues/802#issuecomment-778852594
// https://github.com/evanw/esbuild/issues/802#issuecomment-803297488
const dev = new EventSource("/__dev__")
dev.addEventListener("reload", () => {
	localStorage.setItem("/__dev__", "" + Date.now())
	window.location.reload()
})
dev.addEventListener("error", e => {
	try {
		console.error(JSON.parse(e.data))
	} catch { }
})
window.addEventListener("storage", e => {
	if (e.key === "/__dev__") {
		window.location.reload()
	}
})
