const sass = {
	name: "scss",
	setup(build) {
		const path = require("path")
		const sass = require("sass")

		build.onResolve({ filter: /\.scss$/ }, args => ({
			path: args.path,
			namespace: "scss-ns",
		}))

		build.onLoad({ filter: /.*/, namespace: "scss-ns" }, async args => {
			// NOTE: esbuild does not yet support CSS sourcemaps. Tracked by
			// https://github.com/evanw/esbuild/issues/519.
			const result = sass.renderSync({
				file: path.join(process.env.SRC_DIR, args.path),
			})
			return {
				contents: result.css.toString(),
				loader: "css",
			}
		})
	},
}

module.exports = {
	define: {
		__DEV__: JSON.stringify(process.env["NODE_ENV"] !== "production"),
	},
	plugins: [sass],
}
