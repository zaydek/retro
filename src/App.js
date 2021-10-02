export default function App() {
	const [state, setState] = React.useState(0)

	return (
		<div>
			<span>Hello, world!{state}</span>
			<button onClick={e => setState(state - 1)}>-</button>
			<button onClick={e => setState(state + 1)}>+</button>
		</div>
	)
}
