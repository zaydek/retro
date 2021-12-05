import "./App.css"

import { ease, Transition } from "../../npm/retro/transition" // -> "@zaydek/retro/transition"

export default function App() {
	const [flagIsOpen, setFlagIsOpen] = React.useState(false)

	return (
		<>
			<button onClick={() => setFlagIsOpen(!flagIsOpen)}>
				Press me
			</button>
			<div className="center min-h-screen">
				<Transition
					on={flagIsOpen}
					from={{
						boxShadow: `0 0 1px hsla(0, 0%, 0%, 0.25),
							0 0 transparent,
							0 0 transparent`,
						opacity: 0,
						y: -20,
						scale: 0.75,
						duration: 0.5,
					}}
					to={{
						boxShadow: `0 0 1px hsla(0, 0%, 0%, 0.25),
							0 8px 8px hsla(0, 0%, 0%, 0.1),
							0 2px 8px hsla(0, 0%, 0%, 0.1)`,
						opacity: 1,
						y: 0,
						scale: 1,
						duration: 0.25,
					}}
					ease={ease.outQuart}
					flagDisableMountTransition
					flagUnmountFrom
					flagIncludeTranslateZ
				>
					<div
						className="modal center"
						onClick={e => setFlagIsOpen(!flagIsOpen)}
					>
						<h1>Hello, world!</h1>
					</div>
				</Transition>
			</div>
		</>
	)
}
