import * as store from "../../npm/retro/store" // -> "@zaydek/retro/store"

export const actions = {
	LOAD: "LOAD",
	TOGGLE_TODO: "TOGGLE_TODO",
	CHANGE_TODO: "CHANGE_TODO",
	COMMIT_TODO: "COMMIT_TODO",
	TOGGLE_TODO_BY_ID: "TOGGLE_TODO_BY_ID",
	CHANGE_TODO_BY_ID: "CHANGE_TODO_BY_ID",
	DELETE_TODO_BY_ID: "DELETE_TODO_BY_ID",
}

export const todosInitialState = {
	form: {
		checked: false,
		value: ""
	},
	todos: [
		// {
		// 	id: string
		// 	checked: boolean
		// 	value: string
		// },
	],
}

export const todosStore = store.createStore(todosInitialState)

function shortID() {
	return Math.random().toString(36).slice(2, 6)
}

export function todosReducer(state, { type, data }) {
	if (type === actions.TOGGLE_TODO) {
		return {
			...state,
			form: {
				...state.form,
				checked: data.checked,
			},
		}
	} else if (type === actions.CHANGE_TODO) {
		return {
			...state,
			form: {
				...state.form,
				value: data.value,
			},
		}
	} else if (type === actions.COMMIT_TODO) {
		if (state.form.value === "") {
			return state
		}
		return {
			...state,
			form: {
				...todosInitialState.form,
			},
			todos: [
				{
					id: shortID(),
					...state.form,
				},
				...state.todos,
			],
		}
	} else if (type === actions.TOGGLE_TODO_BY_ID) {
		const todoIndex = data.todoIndex
		return {
			...state,
			todos: [
				...state.todos.slice(0, todoIndex),
				{
					...state.todos[todoIndex],
					checked: data.checked,
				},
				...state.todos.slice(todoIndex + 1),
			],
		}
	} else if (type === actions.CHANGE_TODO_BY_ID) {
		const todoIndex = data.todoIndex
		return {
			...state,
			todos: [
				...state.todos.slice(0, todoIndex),
				{
					...state.todos[todoIndex],
					value: data.value,
				},
				...state.todos.slice(todoIndex + 1),
			],
		}
	} else if (type === actions.DELETE_TODO_BY_ID) {
		const todoIndex = data.todoIndex
		return {
			...state,
			todos: [
				...state.todos.slice(0, todoIndex),
				...state.todos.slice(todoIndex + 1),
			],
		}
	} else {
		throw new Error("Internal error")
	}
}
