/**
 * @type { import("esbuild").Plugin }
 */
module.exports = {
	name: "sass",
	setup(build) {
		const sass = require("sass")
		build.onLoad({ filter: /\.scss$/ }, args => {
			const result = sass.renderSync({
				file: args.path,

				// Suppress debug and warn messages
				//
				//   export interface Logger {
				//     /** This method is called when Sass emits a debug message due to a @debug rule. */
				//     debug?(message: string, options: { deprecation: boolean; span?: SourceSpan; stack?: string }): void;
				//
				//     /** This method is called when Sass emits a debug message due to a @warn rule. */
				//     warn?(message: string, options: { span: SourceSpan }): void;
				//   }
				//
				logger: {
					debug() { /* No-op */ },
					warn() { /* No-op */ },
				}
			})
			return {
				contents: result.css.toString(),
				loader: "css",
			}
		})
	},
}
