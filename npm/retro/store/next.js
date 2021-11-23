function nextImpl(state, selector, updater) {
	// The focus reference for the current state
	let stateRef = state
	// The next state
	let nextState = { ...state }
	// The focus reference for the next state
	let nextStateRef = nextState
	for (let keyIndex = 0; keyIndex < selector.length; keyIndex++) {
		const key = selector[keyIndex]
		const keyIsAtEnd = keyIndex + 1 === selector.length
		if (!keyIsAtEnd) {
			Object.assign(nextStateRef, {
				// Allocate new references for arrays and objects
				[key]: Array.isArray(stateRef[key])
					? [...stateRef[key]]    // Array
					: { ...stateRef[key] }, // Object
			})
		} else {
			Object.assign(nextStateRef, {
				// The deepest element needs to copy all properties
				...stateRef,                         // Old properties
				[key]: typeof updater === "function" // New property
					? updater(nextStateRef[key])
					: updater,
			})
		}
		// Update the focus references
		stateRef = stateRef[key]
		nextStateRef = nextStateRef[key]
	}
	return nextState
}

export default function next(state, ...args) {
	const [selector, updater] = [
		args.slice(0, args.length - 1),
		args[args.length - 1],
	]
	return nextImpl(state, selector, updater)
}
