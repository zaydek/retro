require("esbuild")

// https://esbuild.github.io/plugins/#using-plugins
const envPlugin = {
	name: "env",
	setup(build) {
		// Intercept import paths called "env" so esbuild doesn't attempt
		// to map them to a file system location. Tag them with the "env-ns"
		// namespace to reserve them for this plugin.
		build.onResolve({ filter: /^env$/ }, args => ({
			path: args.path,
			namespace: "env-ns",
		}))

		// Load paths tagged with the "env-ns" namespace and behave as if
		// they point to a JSON file containing the environment variables.
		build.onLoad({ filter: /.*/, namespace: "env-ns" }, () => ({
			contents: JSON.stringify(process.env),
			loader: "json",
		}))
	},
}

const scssPlugin = {
	name: "scss",
	setup(build) {
		const path = require("path")
		const sass = require("sass")

		build.onResolve({ filter: /\.scss$/ }, args => ({
			path: args.path,
			namespace: "scss-ns",
		}))

		// NOTE: esbuild does not support sourcemaps for CSS. Tracked by
		// https://github.com/evanw/esbuild/issues/519.
		build.onLoad({ filter: /.*/, namespace: "scss-ns" }, async args => {
			const result = sass.renderSync({ file: path.join("src", args.path) })
			return {
				contents: result.css.toString(),
				loader: "css",
				// watchDirs: result.stats.includedFiles, // TODO
				// watchFiles: result.stats.entry, // TODO
			}
		})
	},
}

// TODO: Change to a configuration-based object
module.exports = [
	envPlugin,
	scssPlugin,
	// ...
]
