import * as router from "../npm/retro/router" // -> @zaydek/retro/router

import {
	App as RouterApp,
	PATHS,
} from "./router"

import {
	App as SassApp,
} from "./sass"

import {
	App as SassTemplateStringsApp,
} from "./sass-template-strings"

import {
	App as StoreApp,
} from "./store"

import {
	DocumentTitle,
} from "../npm/retro/document-title" // -> "@zaydek/retro/title"

import {
	App as TransitionApp,
} from "./transition"

export function App() {
	return (
		<>

			<ul>
				<li>
					<router.Link path="/router">
						Open the <code>/router</code> app
					</router.Link>
				</li>
				<li>
					<router.Link path="/sass">
						Open the <code>/sass</code> app
					</router.Link>
				</li>
				<li>
					<router.Link path="/sass-template-strings">
						Open the <code>/sass-template-strings</code> app
					</router.Link>
				</li>
				<li>
					<router.Link path="/store">
						Open the <code>/store</code> app
					</router.Link>
				</li>
				<li>
					<router.Link path="/transition">
						Open the <code>/transition</code> app
					</router.Link>
				</li>
			</ul>

			<router.Router>

				<router.Route path="/">
					Hello, world!
				</router.Route>

				<router.Route path="/router">
					<RouterApp />
				</router.Route>

				{/********************************************************************/}

				<router.Route path={PATHS.A}>
					<DocumentTitle title={`Welcome to the ${PATHS.A} page`}>
						<div>You are currently on the <code>{PATHS.A}</code> page</div>
					</DocumentTitle>
				</router.Route>

				<router.Route path={PATHS.B}>
					<DocumentTitle title={`Welcome to the ${PATHS.B} page`}>
						<div>You are currently on the <code>{PATHS.B}</code> page</div>
					</DocumentTitle>
				</router.Route>

				<router.Route path={PATHS.C}>
					<DocumentTitle title={`Welcome to the ${PATHS.C} page`}>
						<div>You are currently on the <code>{PATHS.C}</code> page</div>
					</DocumentTitle>
				</router.Route>

				<router.Route path={PATHS.D}>
					<router.Redirect path={PATHS.A} />
				</router.Route>

				<router.Route path={PATHS.FOUR_ZERO_FOUR}>
					<DocumentTitle title={`Welcome to the ${PATHS.FOUR_ZERO_FOUR} page`}>
						<div>You are currently on the <code>{PATHS.FOUR_ZERO_FOUR}</code> page</div>
					</DocumentTitle>
				</router.Route>

				{/********************************************************************/}

				<router.Route path="/sass">
					<SassApp />
				</router.Route>

				<router.Route path="/sass-template-strings">
					<SassTemplateStringsApp />
				</router.Route>

				<router.Route path="/store">
					<StoreApp />
				</router.Route>

				<router.Route path="/transition">
					<TransitionApp />
				</router.Route>

			</router.Router>

		</>
	)
}
