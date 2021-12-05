import * as router from "../../npm/retro/router" // -> "@zaydek/retro/router"

export const PATHS = {
	A: "/router/a",
	B: "/router/b",
	C: "/router/c",
	D: "/router/d",
	INLINE: "/router/inline",
	FOUR_ZERO_FOUR: "/router/404",
	FOUR_ZERO_FOUR_ZERO_FOUR: "/router/40404",
}

export const URLS = {
	GOOGLE: "https://google.com",
}

export default function App() {
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
