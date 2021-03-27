import * as React from "react"

import "./App.scss"

export default function App() {
	const [state, setState] = React.useState(0)

	return (
		<>
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
			</div>
			{state}
			<br />
			<button onClick={() => setState(state - 1)}>-</button>
			<button onClick={() => setState(state + 1)}>+</button>
		</>
	)
}
