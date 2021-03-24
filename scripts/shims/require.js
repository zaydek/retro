window.require = modName => {
	switch (modName) {
		case "react":
			return window["React"]
		case "react-dom":
			return window["ReactDOM"]
		default:
			throw new Error("Internal error")
	}
}
