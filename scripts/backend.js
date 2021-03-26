const esbuild = require("esbuild")
const fs = require("fs")
const path = require("path")

////////////////////////////////////////////////////////////////////////////////

const stdout = res => console.log(JSON.stringify(res))
const stderr = console.error

const readline = (() => {
	async function* generator() {
		const read = require("readline").createInterface({ input: process.stdin })
		for await (const str of read) {
			yield str
		}
	}
	const generate = generator()
	return async () => (await generate.next()).value
})()

////////////////////////////////////////////////////////////////////////////////

const WWW_DIR = process.env["WWW_DIR"]
const SRC_DIR = process.env["SRC_DIR"]
const OUT_DIR = process.env["OUT_DIR"]

const env = process.env.NODE_ENV ?? "development"

const common = {
	color: true,
	define: {
		// __DEV__: JSON.stringify(env !== "production"),
		"process.env.NODE_ENV": JSON.stringify(env),
	},
	loader: {
		".js": "jsx",
	},
	logLevel: process.env["VERBOSE"] === "true" ? "warning" : "silent",
	minify: env === "production",
	sourcemap: true,
}

async function resolveUserConfig() {
	try {
		await fs.promises.stat("retro.config.js")
	} catch {
		return {}
	}
	return require(path.join(process.cwd(), "retro.config.js"))
}

let result = undefined

async function build(opt) {
	const buildRes = {
		warnings: [],
		errors: [],
	}

	const config = await resolveUserConfig()

	try {
		// out/vendor.js
		await esbuild.build({
			...common,
			bundle: true,
			entryPoints: ["scripts/shims/vendor.js"],
			outfile: path.join(OUT_DIR, "vendor.js"),
		})

		// out/bundle.js
		result = await esbuild.build({
			...config,
			...common,

			define: { ...config.define, ...common.define },
			loader: { ...config.loader, ...common.loader },

			bundle: true,
			entryPoints: [path.join(SRC_DIR, "index.js")], // TODO: Add support for ".jsx", ".ts", and ".tsx"
			outfile: path.join(OUT_DIR, "bundle.js"),

			external: ["react", "react-dom"], // Dedupe React APIs (because of vendor)
			inject: ["scripts/shims/require.js"], // Add support for vendor
			plugins: config?.plugins,

			// https://github.com/evanw/esbuild/issues/1063
			incremental: true,

			// watch: !opt.watch
			// 	? undefined
			// 	: {
			// 			onRebuild(error, result) {
			// 				require("child_process").execSync("osascript -e 'beep'")
			// 				require("child_process").execSync("osascript -e 'beep'")
			//
			// 				const rebuildRes = {
			// 					warnings: [],
			// 					errors: [],
			// 				}
			// 				// E.g. try
			// 				if (result !== null) {
			// 					if (result.warnings.length > 0) {
			// 						rebuildRes.warnings = result.warnings
			// 					}
			// 				}
			// 				// E.g. catch
			// 				if (error !== null) {
			// 					if (error.errors.length > 0) {
			// 						rebuildRes.errors = error.errors
			// 					}
			// 					if (error.warnings.length > 0) {
			// 						rebuildRes.warnings = error.warnings
			// 					}
			// 				}
			// 				stdout({
			// 					Kind: "rebuild-done",
			// 					Data: rebuildRes,
			// 				})
			// 			},
			// 	  },
		})
		if (result?.warnings?.length > 0) {
			buildRes.warnings = result.warnings
		}
	} catch (error) {
		if (error?.errors?.length > 0) {
			buildRes.errors = error.errors
		}
		if (error?.warnings?.length > 0) {
			buildRes.warnings = error.warnings
		}
	}

	require("child_process").execSync("osascript -e 'beep'")
	stdout({
		Kind: "build-done",
		Data: buildRes,
	})
}

async function rebuild() {
	if (result === null) return

	const buildRes = {
		warnings: [],
		errors: [],
	}

	try {
		const result2 = await result.rebuild()
		if (result2?.warnings?.length > 0) {
			buildRes.warnings = result2.warnings
		}
	} catch (error) {
		if (error?.errors?.length > 0) {
			buildRes.errors = error.errors
		}
		if (error?.warnings?.length > 0) {
			buildRes.warnings = error.warnings
		}
	}

	stdout({
		Kind: "rebuild-done",
		Data: buildRes,
	})
}

async function main() {
	esbuild.initialize()
	while (true) {
		const jsonstr = await readline()
		const msg = JSON.parse(jsonstr)
		try {
			switch (msg.Kind) {
				case "dev":
					await build({
						watch: true,
					})
					break
				case "build":
					await build({
						watch: false,
					})
					break
				case "rebuild":
					await rebuild()
					break
				case "done":
					if (result?.rebuild) {
						result.rebuild.dispose()
					}
					break
				default:
					throw new Error("Internal error")
			}
		} catch (error) {
			stderr(error.stack)
			process.exit(1)
		}
	}
}

main()
