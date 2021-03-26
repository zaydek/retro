import "./App.scss"

export default function App() {
	return (
		<div className="App">
			<h1>Hello {JSON.stringify(__DEV__)}</h1>
		</div>
	)
}
