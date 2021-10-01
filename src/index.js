import "env"

import App from "./App"

import "./index.css"

if (document.getElementById("root").hasChildNodes()) {
	ReactDOM.hydrate(
		<React.StrictMode>
			<App />
		</React.StrictMode>,
		document.getElementById("root"),
	)
} else {
	ReactDOM.render(
		<React.StrictMode>
			<App />
		</React.StrictMode>,
		document.getElementById("root"),
	)
}
