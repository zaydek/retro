/**
 * @type { import("esbuild").Plugin }
 */
module.exports = {
	name: "sass",
	setup(build) {
		const sass = require("sass")

		build.onLoad({ filter: /\.scss$/ }, args => {
			const result = sass.renderSync({ file: args.path })
			return {
				contents: result.css.toString(),
				loader: "css",
			}
		})
	},
}
