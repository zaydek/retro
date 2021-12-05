import * as router from "../npm/retro/router" // -> @zaydek/retro/router

import PluginMDXApp from "./plugin-mdx"
import PluginSassApp from "./plugin-sass"
import RouterApp, { PATHS } from "./router"
import StoreApp from "./store"
import TransitionApp from "./transition"
import { DocumentTitle } from "../npm/retro/document-title" // -> "@zaydek/retro/title"

export default function App() {
	return (
		<>

			<ul>
				<li>
					<router.Link path="/router">
						Open the <code>/router</code> app
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
				<li>
					<router.Link path="/plugin-sass">
						Open the <code>/plugin-sass</code> app
					</router.Link>
				</li>
				<li>
					<router.Link path="/plugin-sass-template-strings">
						Open the <code>/plugin-sass-template-strings</code> app
					</router.Link>
				</li>
				<li>
					<router.Link path="/plugin-mdx">
						Open the <code>/plugin-mdx</code> app
					</router.Link>
				</li>
			</ul>

			<router.Router>

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

				<router.Route path="/store">
					<StoreApp />
				</router.Route>

				<router.Route path="/transition">
					<TransitionApp />
				</router.Route>

				<router.Route path="/plugin-sass">
					<PluginSassApp />
				</router.Route>

				<router.Route path="/plugin-mdx">
					<PluginMDXApp />
				</router.Route>

			</router.Router>

		</>
	)
}
