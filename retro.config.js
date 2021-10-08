// https://esbuild.github.io/api/#build-api

module.exports = {
	plugins: [
		require("./npm/retro/plugins/sass"),
		require("./npm/retro/plugins/sass-template-strings"),
	],
}
