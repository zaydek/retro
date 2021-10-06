import {
	App,
} from "./App"

import "./reset.css"

if (document.getElementById("root").hasChildNodes()) {
	// For static-site generation (SSG) and server-side rendering (SSR)
	ReactDOM.hydrate(
		<React.StrictMode>
			<App />
		</React.StrictMode>,
		document.getElementById("root"),
	)
} else {
	// For client-side rendering (CSR)
	ReactDOM.render(
		<React.StrictMode>
			<App />
		</React.StrictMode>,
		document.getElementById("root"),
	)
}
