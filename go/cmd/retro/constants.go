package retro

const (
	// The HTML entry point
	indexHTML = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Hello, world!</title>
    <link rel="stylesheet" href="/client.css" />
  </head>
  <body>
    <div id="root"></div>
    <script src="/vendor.js"></script>
    <script src="/client.js"></script>
  </body>
</html>` + "\n"

	// Server-sent events (SSE) stub for the dev command
	serverSentEventsStub = `<script type="module">const __dev__=new EventSource("/__dev__");__dev__.addEventListener("reload",()=>{localStorage.setItem("__dev__",""+Date.now()),window.location.reload()}),__dev__.addEventListener("error",e=>{try{console.error(JSON.parse(e.data))}catch{}}),window.addEventListener("storage",e=>{e.key==="__dev__"&&window.location.reload()});</script>`

	// The JS entry point
	indexJS = `import App from "./App"

import "./reset.css"

if (document.getElementById("root").hasChildNodes()) {
	// For static-side generation (SSG) and server-side rendering (SSR)
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
` + "\n"

	// The JS app entry point
	appJS = `import "./App.css"

export default function App() {
  return (
    <div className="App">
      <h1>Hello, world!</h1>
    </div>
  )
}` + "\n"
)
