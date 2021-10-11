import {
	buildStyles,
	getCSSPropertiesForShorthands,
	safeEase,
} from "./helpers"

import {
	useLayoutEffectSSR,
} from "../use-layout-effect-ssr"

const DEFAULTS = {
	duration: 0.3,
	ease: "ease-out",
	delay: 0,
}

export function Transition({
	on,
	from,
	to,
	duration,
	ease,
	delay,
	children,
	flagDisableMountTransition,
	flagUnmountFrom,
	flagUnmountTo,
	flagIncludeTranslateZ,
}) {
	const [computedChildren, setComputedChildren] = React.useState(() => on && children)

	const memoFrom = React.useMemo(() => {
		return from
	}, [from])
	const memoTo = React.useMemo(() => {
		return to
	}, [to])

	const memoKeys = React.useMemo(() => {
		return [
			...new Set([
				...getCSSPropertiesForShorthands(Object.keys(memoFrom)),
				...getCSSPropertiesForShorthands(Object.keys(memoTo)),
			])
		]
	}, [
		memoFrom,
		memoTo,
	])

	const timeoutIDRef = React.useRef(0)

	// Tracks whether to suppress a rerender when on=false once
	const didMountRef = React.useRef(false)
	useLayoutEffectSSR(
		React.useCallback(() => {
			const setDir = dir => {
				setComputedChildren(
					React.cloneElement(
						children,
						{
							style: {
								...children.props.style,
								...buildStyles(dir, { flagIncludeTranslateZ }),
								willChange: memoKeys.join(", "),
								transition: memoKeys.map(property =>
									`
										${property}
										${dir.duration ?? duration ?? DEFAULTS.duration}s
										${safeEase(dir.ease ?? ease ?? DEFAULTS.ease)}
										${dir.delay ?? delay ?? DEFAULTS.delay}s
									`
										.trim()
										.replace(/\n\t+/g, " "),
								).join(", "),
							},
						},
					),
				)
			}

			if (timeoutIDRef.current !== 0) {
				clearTimeout(timeoutIDRef.current)
				timeoutIDRef.current = 0
			}

			let dir1 = null
			let dir2 = null
			let flagUmountOnDir2 = false
			if (!on) {
				// Backwards
				dir1 = memoTo
				dir2 = memoFrom
				flagUmountOnDir2 = flagUnmountFrom ?? false
			} else {
				// Forwards
				dir1 = memoFrom
				dir2 = memoTo
				flagUmountOnDir2 = flagUnmountTo ?? false
			}

			if (flagUmountOnDir2 && !didMountRef.current) {
				setComputedChildren(null)
				didMountRef.current = true
				return
			} else if (flagDisableMountTransition && !didMountRef.current) {
				setDir(dir2)
				didMountRef.current = true
				return
			} else {
				didMountRef.current = true
				// Passthrough
			}

			setDir(dir1)
			setTimeout(() => {
				setDir(dir2)
				if (flagUmountOnDir2) {
					timeoutIDRef.current = setTimeout(() => {
						setComputedChildren(null)
						timeoutIDRef.current = 0
					}, ((dir2?.delay ?? 0) + (dir2?.duration ?? 0)) * 1_000)
				}
			}, 0)

			return () => {
				if (timeoutIDRef.current !== 0) {
					clearTimeout(timeoutIDRef.current)
					timeoutIDRef.current = 0
				}
			}
		}, [
			on,                         // Value
			memoFrom,                   // Reference
			memoTo,                     // Reference
			duration,                   // Value
			ease,                       // Value
			delay,                      // Value
			children,                   // Reference
			memoKeys,                   // Reference
			flagDisableMountTransition, // Value
			flagUnmountFrom,            // Value
			flagUnmountTo,              // Value
			flagIncludeTranslateZ,      // Value
		]),
		[,
			on, // Value
		],
	)

	return computedChildren
}
