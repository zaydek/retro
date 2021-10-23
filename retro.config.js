// https://esbuild.github.io/api/#build-api

module.exports = {
	plugins: [
		require("./npm/retro/plugins/mdx"),
		require("./npm/retro/plugins/sass-template-strings"),
		require("./npm/retro/plugins/sass"),
	],
}
