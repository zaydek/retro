window.require = modName => {
	switch (modName) {
		case "react":
			return window["React"]
		case "react-dom":
			return window["ReactDOM"]
		case "react-dom/server":
			return window["ReactDOMServer"]
		default:
			throw new Error("Internal error")
	}
}
