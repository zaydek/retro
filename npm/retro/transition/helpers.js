export function buildStyles(dir, { flagIncludeTranslateZ }) {
	const {
		duration: _1, // No-op
		ease: _2,     // No-op
		delay: _3,    // No-op
		x,
		y,
		scale,
		...styles
	} = dir

	let transformStr = ""
	if (x !== undefined) {
		if (transformStr !== "") transformStr += " "
		transformStr += `translateX(${typeof x === "string" ? x : `${x / 16}rem`})`
	}
	if (y !== undefined) {
		if (transformStr !== "") transformStr += " "
		transformStr += `translateY(${typeof y === "string" ? y : `${y / 16}rem`})`
	}
	if (scale !== undefined) {
		if (transformStr !== "") transformStr += " "
		transformStr += `scale(${scale})`
	}
	if (flagIncludeTranslateZ) {
		if (transformStr !== "") transformStr += " "
		transformStr += `translateZ(0)`
	}
	if (transformStr !== "") {
		styles.transform = transformStr
	}

	return styles
}

function convertToKebabCase(str) {
	return str.replace(/([A-Z])/g, (_, $1, $1Index) =>
		($1Index === 0 ? "" : "-") + $1.toLowerCase())
}

function getCSSPropertyForShorthand(shorthand) {
	switch (shorthand) {
		case "duration": // No-op time-related properties
		case "ease":     // No-op time-related properties
			return false
		case "scale":
		case "x":
		case "y":
			return "transform"
	}
	return convertToKebabCase(shorthand)
}

export function getCSSPropertiesForShorthands(shorthands) {
	return shorthands
		.map(shorthand => getCSSPropertyForShorthand(shorthand))
		.filter(Boolean)
}

export function safeEase(ease) {
	if (typeof ease === "string") { return ease }
	return `cubic-bezier(${ease.join(", ")})`
}
