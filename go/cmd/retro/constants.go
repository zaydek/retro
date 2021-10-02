package retro

const (
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

	// serverSentEventsStub = `<script type="module">
	// 		// Retro sends server-sent events (SSE) to ` + "`" + `/__dev__` + "`" + ` to force-reload
	// 		// tabs. Because server-sent events propagate to ~6 tabs max (depending on
	// 		// the browser), use localStorage as a mechanism to force-reload other
	// 		// tabs.
	// 		//
	// 		// https://github.com/evanw/esbuild/issues/802#issuecomment-778852594
	// 		// https://github.com/evanw/esbuild/issues/802#issuecomment-803297488
	// 		const __dev__ = new EventSource("/__dev__")
	// 		__dev__.addEventListener("reload", () => {
	// 			localStorage.setItem("__dev__", "" + Date.now())
	// 			window.location.reload()
	// 		})
	// 		__dev__.addEventListener("error", e => {
	// 			try {
	// 				console.error(JSON.parse(e.data))
	// 			} catch { }
	// 		})
	// 		window.addEventListener("storage", e => {
	// 			if (e.key === "__dev__") {
	// 				window.location.reload()
	// 			}
	// 		})
	// 	</script>`

	serverSentEventsStub = `<script type="module">const __dev__=new EventSource("/__dev__");__dev__.addEventListener("reload",()=>{localStorage.setItem("__dev__",""+Date.now()),window.location.reload()}),__dev__.addEventListener("error",e=>{try{console.error(JSON.parse(e.data))}catch{}}),window.addEventListener("storage",e=>{e.key==="__dev__"&&window.location.reload()});</script>`

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

	appJS = `import "./App.css"

export default function App() {
  return (
    <div className="App">
      <h1>Hello, world!</h1>
    </div>
  )
}` + "\n"
)
