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
					CMD: {process.env.CMD}<br />
					ENV: {process.env.ENV}<br />
					WWW_DIR: {process.env.WWW_DIR}<br />
					SRC_DIR: {process.env.SRC_DIR}<br />
					OUT_DIR: {process.env.OUT_DIR}<br />
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
