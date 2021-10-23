const sassTplStrGlobalRegex = /(?!\/\/\s*)sass\.global`((?:\r?\n.*?)+)`/gm
const sassTplStrInlineRegex = /(?!\/\/\s*)sass(?!\.global)`((?:\r?\n.*?)+)`/gm

// Scans for global and inline Sass template strings
function scanMatches(contents) {
	const globals = []
	const inlines = []
	let match = null
	while ((match = sassTplStrGlobalRegex.exec(contents))) {
		globals.push(match[1])
	}
	while ((match = sassTplStrInlineRegex.exec(contents))) {
		inlines.push(match[1])
	}
	return [globals, inlines]
}

/**
 * @type { import("esbuild").Plugin }
 */
module.exports = {
	name: "sass-template-strings",
	setup(build) {
		const fs = require("fs")
		const sass = require("sass")

		const importers = new Set()

		build.onResolve({ filter: /^sass-template-strings$/ }, args => {
			importers.add(args.importer)
			return {
				path: args.path,
				namespace: "sass-template-strings-ns",
			}
		})

		build.onLoad({ filter: /.*/, namespace: "sass-template-strings-ns" }, async args => {
			const allGlobals = []
			const allInlines = []

			// Aggregate global and inline matches
			for (const importer of importers) {
				const buffer = await fs.promises.readFile(importer)
				const [globals, inlines] = scanMatches(buffer.toString())
				allGlobals.push(...globals)
				allInlines.push(...inlines)
			}

			// Build globals string
			let globalsStr = ""
			for (const global of allGlobals) {
				globalsStr += `
					${global}
				`
			}

			// Build inlines string
			let inlinesStr = ""
			for (const inline of allInlines) {
				inlinesStr += `
					@at-root {
						${inline}
					}
				`
			}

			// Render Sass
			const result = sass.renderSync({
				data: `
					${globalsStr}
					${inlinesStr}
				`,
			})
			const css = result.css.toString()

			// Stub 'sass.global' and 'sass' functions
			return {
				contents: `
					import "data:text/css,${encodeURI(css)}"
					function sass() { /* No-op */ }
					Object.assign(sass, {
						global() { /* No-op */ }
					})
					export default sass
				`,
				loader: "js",
			}
		})
	},
}
