import * as esbuild from "esbuild"
import * as fs from "fs"
import * as path from "path"

////////////////////////////////////////////////////////////////////////////////

const stdout = (msg: Message) => console.log(JSON.stringify(msg))
const stderr = console.error

const readline = ((): (() => Promise<string>) => {
	async function* generator(): AsyncGenerator<string> {
		const read = require("readline").createInterface({ input: process.stdin })
		for await (const str of read) {
			yield str
		}
	}
	const generate = generator()
	return async () => (await generate.next()).value
})()

////////////////////////////////////////////////////////////////////////////////

interface Message {
	Kind: string
	Data: any
}

interface BuildResponse {
	errors: esbuild.Message[]
	warnings: esbuild.Message[]
}

////////////////////////////////////////////////////////////////////////////////

function InternalError<T>(param: T): T {
	throw new Error("Internal error")
	return param
}

const WWW_DIR = process.env["WWW_DIR"] ?? InternalError("")
const SRC_DIR = process.env["SRC_DIR"] ?? InternalError("")
const OUT_DIR = process.env["OUT_DIR"] ?? InternalError("")

const ENV = process.env["NODE_ENV"] ?? InternalError("")

const common: esbuild.BuildOptions = {
	color: true,
	define: {
		"process.env.NODE_ENV": JSON.stringify(ENV),
	},
	loader: {
		".js": "jsx",
	},
	logLevel: "silent",
	minify: ENV === "production",
	sourcemap: true,
}

async function resolveUserConfig(): Promise<esbuild.BuildOptions> {
	try {
		await fs.promises.stat("retro.config.js")
	} catch {
		return {}
	}
	return require(path.join(process.cwd(), "retro.config.js"))
}

let result: esbuild.BuildResult | esbuild.BuildIncremental | null = null

async function build(): Promise<BuildResponse> {
	const buildRes: BuildResponse = {
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

			// TODO: Add support for ".jsx", ".ts", and ".tsx"
			bundle: true,
			entryPoints: [path.join(SRC_DIR, "index.js")],
			outfile: path.join(OUT_DIR, "bundle.js"),

			external: ["react", "react-dom"], // Dedupe React APIs (because of vendor)
			inject: ["scripts/shims/require.js"], // Add support for vendor
			plugins: config?.plugins,

			incremental: true,
		})
		if (result.warnings.length > 0) {
			buildRes.warnings = result.warnings
		}
	} catch (error) {
		if (error.errors.length > 0) {
			buildRes.errors = error.errors
		}
		if (error.warnings.length > 0) {
			buildRes.warnings = error.warnings
		}
	}

	return buildRes
}

async function rebuild(): Promise<BuildResponse> {
	if (result?.rebuild === undefined) throw new Error("Internal error")

	const rebuildRes: BuildResponse = {
		warnings: [],
		errors: [],
	}

	try {
		const result2 = await result.rebuild()
		if (result2.warnings.length > 0) {
			rebuildRes.warnings = result2.warnings
		}
	} catch (error) {
		if (error.errors.length > 0) {
			rebuildRes.errors = error.errors
		}
		if (error.warnings.length > 0) {
			rebuildRes.warnings = error.warnings
		}
	}

	return rebuildRes
}

async function main() {
	esbuild.initialize({})

	while (true) {
		const jsonstr = await readline()
		const msg = JSON.parse(jsonstr)
		try {
			switch (msg.Kind) {
				case "build":
					const buildRes = await build()
					stdout({ Kind: "build-done", Data: buildRes })
					break
				case "rebuild":
					const rebuildRes = await rebuild()
					stdout({ Kind: "rebuild-done", Data: rebuildRes })
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
