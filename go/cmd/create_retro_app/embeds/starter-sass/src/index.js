import App from "./App"

import "./reset.css"

if (document.getElementById("retro_root").hasChildNodes()) {
	// For static-side generation (SSG) and server-side rendering (SSR)
	ReactDOM.hydrate(
		<React.StrictMode>
			<App />
		</React.StrictMode>,
		document.getElementById("retro_root"),
	)
} else {
	// For client-side rendering (CSR)
	ReactDOM.render(
		<React.StrictMode>
			<App />
		</React.StrictMode>,
		document.getElementById("retro_root"),
	)
}
