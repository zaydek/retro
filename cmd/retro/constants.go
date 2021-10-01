package retro

const (
	// Permission bits for writing files
	permBitsFile = 0644

	// Permission bits for writing directories
	permBitsDirectory = 0755
)

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
    <div id="retro_root"></div>
    <script src="/vendor.js"></script>
    <script src="/client.js"></script>
  </body>
</html>` + "\n"

	serverSentEventsStub = `<script type="module">
	// Retro sends server-sent events (SSE) to ` + "`/__retro_dev__`" + ` to force-reload
	// tabs. Because server-sent events propagate to ~6 tabs max (depending on the
	// browser), use localStorage as a mechanism to force-reload other tabs.
	//
	// https://github.com/evanw/esbuild/issues/802#issuecomment-778852594
	// https://github.com/evanw/esbuild/issues/802#issuecomment-803297488
	const dev = new EventSource("/__retro_dev__")
	dev.addEventListener("reload", () => {
		localStorage.setItem("/__retro_dev__", "" + Date.now())
		window.location.reload()
	})
	dev.addEventListener("error", e => {
		try {
			console.error(JSON.parse(e.data))
		} catch { }
	})
	window.addEventListener("storage", e => {
		if (e.key === "/__retro_dev__") {
			window.location.reload()
		}
	})
</script>` + "\n"

	indexJS = `import App from "./App"

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
