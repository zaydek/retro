import * as router from "../npm/retro/router" // -> @zaydek/retro/router

import {
	App as RouterApp,
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
				<router.Route path="/router">
					<RouterApp />
				</router.Route>
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
