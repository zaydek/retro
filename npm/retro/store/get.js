export function get(state, ...args) {
	const keys = args.slice(0, args.length - 1)
	const newValueOrUpdater = args[args.length - 1]

	// The focus reference for the current state
	let stateRef = state
	// The next state
	let nextState = { ...state }
	// The focus reference for the next state
	let nextStateFocusRef = nextState
	for (let keyIndex = 0; keyIndex < keys.length; keyIndex++) {
		const key = keys[keyIndex]
		const keyIsAtEnd = keyIndex + 1 === keys.length
		if (!keyIsAtEnd) {
			Object.assign(nextStateFocusRef, {
				// Allocate new references for arrays and objects
				[key]: Array.isArray(stateRef[key])
					? [...stateRef[key]]    // Array
					: { ...stateRef[key] }, // Object
			})
		} else {
			Object.assign(nextStateFocusRef, {
				// The deepest element needs to copy all properties
				...stateRef,                                   // Old properties
				[key]: typeof newValueOrUpdater === "function" // New property
					? newValueOrUpdater(nextStateFocusRef[key])
					: newValueOrUpdater,
			})
		}
		// Update the focus references
		stateRef = stateRef[key]
		nextStateFocusRef = nextStateFocusRef[key]
	}
	return nextState
}
