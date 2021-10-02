export default function App() {
	const [state, setState] = React.useState("Hell")

	return (
		<div>
			<span>Hello, world! Hahaha {state}</span>
			<button onClick={e => setState(state - 1)}>-</button>
			<button onClick={e => setState(state + 1)}>+</button>
		</div>
	)
}
