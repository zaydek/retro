package retro

const (
	// The HTML entry point
	htmlEntryPoint = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Hello, world!</title>
		<link rel="stylesheet" href="/client.css" />
	</head>
	<body>
		<div id="root"></div>
		<script src="/vendor.js" type="module"></script>
		<script src="/client.js" type="module"></script>
	</body>
</html>`

	// Server-sent events (SSE) for the dev command
	htmlServerSentEvents = `<script type="module">const dev=new EventSource("/__dev__");dev.addEventListener("reload",()=>{localStorage.setItem("__dev__",""+Date.now()),window.location.reload()}),dev.addEventListener("error",e=>{try{console.error(JSON.parse(e.data))}catch{}}),window.addEventListener("storage",e=>{e.key==="__dev__"&&window.location.reload()});</script>`

	// The JavaScript entry point
	jsEntryPoint = `import "./reset.css"

import { App } from "./App"

ReactDOM.render(
	<React.StrictMode>
		<App />
	</React.StrictMode>,
	document.getElementById("root"),
)`

	// The JavaScript app entry point
	appJSEntryPoint = `import "./App.css"

export default function App() {
	return (
		<div className="App">
			<h1>Hello, world!</h1>
		</div>
	)
}`
)
