import "./App.scss"

export default function App() {
	return (
		<div className="App">
			<pre>
				{/* prettier-ignore */}
				<code>
					CMD: {CMD}<br />
					ENV: {ENV}<br />
					WWW_DIR: {WWW_DIR}<br />
					SRC_DIR: {SRC_DIR}<br />
					OUT_DIR: {OUT_DIR}<br />

				</code>
			</pre>
			Hello, world!
		</div>
	)
}
