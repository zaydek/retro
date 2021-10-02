export default function App() {
	const [state, setState] = React.useState(800

	return (
		<div>
			<span>{state}</span>
			<button onClick={e => setState(state - 1)}>-</button>
			<button onClick={e => setState(state + 1)}>+</button>
		</div>
	)
}
