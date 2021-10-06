import * as router from "../../npm/retro/router" // -> "@zaydek/retro/router"

import {
	PATHS,
	URLS,
} from "./paths"

import {
	Title,
} from "../../title" // -> "@zaydek/retro/title"

function Links() {
	return (
		<ul>
			<li>
				<router.Link path={PATHS.A}>
					Click to open the <code>{PATHS.A}</code> page
				</router.Link>
			</li>
			<li>
				<router.Link path={PATHS.B}>
					Click to open the <code>{PATHS.B}</code> page
				</router.Link>
			</li>
			<li>
				<router.Link path={PATHS.C}>
					Click to open the <code>{PATHS.C}</code> page
				</router.Link>
			</li>
			<li>
				<router.Link path={PATHS.D}>
					Click to open the <code>{PATHS.D}</code> page (redirects to <code>{PATHS.A}</code>)
				</router.Link>
			</li>
			<li>
				<router.Link path={PATHS.INLINE}>
					Click to open the <code>{PATHS.INLINE}</code> page
				</router.Link>
			</li>
			<li>
				<router.Link path={PATHS.FOUR_ZERO_FOUR}>
					Click to open the <code>{PATHS.FOUR_ZERO_FOUR}</code> page
				</router.Link>
			</li>
			<li>
				<router.Link path={PATHS.FOUR_ZERO_FOUR_ZERO_FOUR}>
					Click to open the <code>{PATHS.FOUR_ZERO_FOUR_ZERO_FOUR}</code> page
				</router.Link>
			</li>
			<li>
				<router.Link path={URLS.GOOGLE}>
					Click to open <code>{URLS.GOOGLE}</code>
				</router.Link>
			</li>
		</ul>
	)
}

function ComponentA() {
	return <div>You are currently on the <code>{PATHS.A}</code> page</div>
}

function ComponentB() {
	return <div>You are currently on the <code>{PATHS.B}</code> page</div>
}

function ComponentC() {
	return <div>You are currently on the <code>{PATHS.C}</code> page</div>
}

export function App() {
	return (
		<>

			<Links />

			<router.Router>

				<router.Route path={PATHS.A}>
					<Title title={`Welcome to the ${PATHS.A} page`}>
						<ComponentA />
					</Title>
				</router.Route>

				<router.Route path={PATHS.B}>
					<Title title={`Welcome to the ${PATHS.B} page`}>
						<ComponentB />
					</Title>
				</router.Route>

				<router.Route path={PATHS.C}>
					<Title title={`Welcome to the ${PATHS.C} page`}>
						<ComponentC />
					</Title>
				</router.Route>

				<router.Route path={PATHS.D}>
					<router.Redirect path={PATHS.A} />
				</router.Route>

				<router.Route path={PATHS.INLINE}>
					<Title title={`Welcome to the ${PATHS.INLINE} page`}>
						<div>You are currently on the <code>{PATHS.INLINE}</code> page</div>
					</Title>
				</router.Route>

				<router.Route path={PATHS.FOUR_ZERO_FOUR}>
					<Title title={`Welcome to the ${PATHS.FOUR_ZERO_FOUR} page`}>
						<div>You are currently on the <code>{PATHS.FOUR_ZERO_FOUR}</code> page</div>
					</Title>
				</router.Route>

			</router.Router>

		</>
	)
}
