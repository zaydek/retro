if (typeof window === "undefined") {
	// For static-site generation (SSG) and server-side rendering (SSR)
	React = require("react")
	ReactDOM = require("react-dom")
	ReactDOMServer = require("react-dom/server")
} else {
	// For client-side rendering (CSR)
	window.require = moduleName => {
		switch (moduleName) {
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
}
