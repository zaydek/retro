// https://esbuild.github.io/api/#build-api

const sass = {
	name: "scss",
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

module.exports = {
	target: ["es2017"],
	plugins: [
		sass,
	],
}
