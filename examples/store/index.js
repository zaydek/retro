import * as store from "../../npm/retro/store" // -> "@zaydek/retro/store"

import {
	actions,
	todosInitialState,
	todosReducer,
	todosStore,
} from "./reducer"

const LOCALSTORAGE_KEY = "todos-app"
const LOCALSTORAGE_DEBOUNCE_MS = 100

function TodoForm() {
	const form = store.useSelector(todosStore, ["form"])
	const dispatch = store.useReducerOnlyDispatch(todosStore, todosReducer)
	return (
		<form
			onSubmit={e => {
				e.preventDefault()
				dispatch({
					type: actions.COMMIT_TODO,
				})
			}}
		>
			<input
				type="checkbox"
				checked={form.checked}
				onChange={e => {
					dispatch({
						type: actions.TOGGLE_TODO,
						data: {
							checked: e.target.checked,
						},
					})
				}}
			/>
			<input
				type="text"
				value={form.value}
				onChange={e => {
					dispatch({
						type: actions.CHANGE_TODO,
						data: {
							value: e.target.value,
						},
					})
				}}
			/>
			<button type="submit">
				+
			</button>
		</form>
	)
}

function Todos() {
	const todos = store.useSelector(todosStore, ["todos"])
	return (
		todos.map((todo, todoIndex) => (
			<Todo
				key={todo.id}
				todoIndex={todoIndex}
			/>
		))
	)
}

function Todo({ todoIndex }) {
	const todo = store.useSelector(todosStore, ["todos", todoIndex])
	const dispatch = store.useReducerOnlyDispatch(todosStore, todosReducer)
	return (
		<div id={todo.id}>
			<input
				type="checkbox"
				checked={todo.checked}
				onChange={e => {
					dispatch({
						type: actions.TOGGLE_TODO_BY_ID,
						data: {
							todoIndex,
							checked: e.target.checked,
						},
					})
				}}
			/>
			<input
				type="text"
				value={todo.value}
				onChange={e => {
					dispatch({
						type: actions.CHANGE_TODO_BY_ID,
						data: {
							todoIndex,
							value: e.target.value,
						},
					})
				}}
			/>
			<button
				onClick={e => {
					dispatch({
						type: actions.DELETE_TODO_BY_ID,
						data: {
							todoIndex,
						},
					})
				}}
			>
				-
			</button>
		</div>
	)
}

function DEBUG_Todos() {
	const state = store.useStateOnlyState(todosStore)
	return (
		<pre style={{ fontSize: 14 }}>
			{JSON.stringify(state, null, 2)}
		</pre>
	)
}

export default function App() {
	const [state, setState] = store.useState(todosStore)
	const [loaded, setLoaded] = React.useState(false)

	React.useEffect(() => {
		let restoredState = todosInitialState
		const jsonState = localStorage.getItem(LOCALSTORAGE_KEY)
		if (jsonState !== null) {
			try {
				restoredState = JSON.parse(jsonState)
			} catch { }
		}
		setState(restoredState)
		setLoaded(true)
	}, [])

	React.useEffect(() => {
		if (!loaded) {
			return
		}
		const id = setTimeout(() => {
			localStorage.setItem(LOCALSTORAGE_KEY, JSON.stringify(state))
		}, LOCALSTORAGE_DEBOUNCE_MS)
		return () => {
			clearTimeout(id)
		}
	}, [loaded, state])

	return (
		<>
			{!loaded
				? "Loading..."
				: (
					<>
						<TodoForm />
						<Todos />
						<DEBUG_Todos />
					</>
				)
			}
		</>
	)
}
