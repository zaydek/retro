/**
 * @type { import("esbuild").Plugin }
 */
module.exports = {
	name: "mdx",
	setup(build) {
		const fs = require("fs")
		const mdx = require("@mdx-js/mdx")
		build.onLoad({ filter: /\.mdx?$/ }, async args => {
			const contents = await fs.promises.readFile(args.path, "utf8")
			const reactContents = await mdx(contents)
			return {
				contents: `
					import { mdx } from "@mdx-js/react"
					${reactContents}
				`,
				loader: "jsx",
			}
		})
	},
}
