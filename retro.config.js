// https://esbuild.github.io/api/#build-api

/**
 * @type { import("esbuild").BuildOptions }
 */
module.exports = {
	plugins: [
		require("./npm/retro/plugins/mdx"),
		require("./npm/retro/plugins/sass"),
	],
	target: ["es2017"],
}
