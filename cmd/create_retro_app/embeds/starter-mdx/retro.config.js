// https://esbuild.github.io/api/#build-api

const mdx = {
	name: "mdx",
	setup(build) {
		const fs = require("fs")
		const mdx = require("@mdx-js/mdx")

		build.onLoad({ filter: /\.mdx$/ }, async args => {
			const text = await fs.promises.readFile(args.path, "utf8")
			const contents = await mdx(text)
			return {
				contents: `
					import { mdx } from "@mdx-js/react"
					${contents}
				`,
				loader: "jsx",
			}
		})
	},
}

module.exports = {
	target: ["es2017"],
	plugins: [
		mdx,
	],
}
