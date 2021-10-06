import {
	App,
} from "./App"

import "./reset.css"

if (document.getElementById("root").hasChildNodes()) {
	// For static-site generation (SSG) and server-side rendering (SSR)
	ReactDOM.hydrate(
		<App />,
		document.getElementById("root"),
	)
} else {
	// For client-side rendering (CSR)
	ReactDOM.render(
		<App />,
		document.getElementById("root"),
	)
}
