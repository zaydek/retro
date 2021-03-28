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

interface BackendResponse {
	Metafile: {
		Vendor: esbuild.Metafile | null
		Bundle: esbuild.Metafile | null
	}
	Errors: esbuild.Message[]
	Warnings: esbuild.Message[]
}

////////////////////////////////////////////////////////////////////////////////

function InternalError<T>(param: T): T {
	throw new Error("Internal error")
	return param
}

const CMD = process.env["CMD"] ?? InternalError("")
const ENV = process.env["ENV"] ?? InternalError("")
const WWW_DIR = process.env["WWW_DIR"] ?? InternalError("")
const SRC_DIR = process.env["SRC_DIR"] ?? InternalError("")
const OUT_DIR = process.env["OUT_DIR"] ?? InternalError("")

const common: esbuild.BuildOptions = {
	color: true,

	// Propagate env vars
	define: {
		// React, React DOM
		"process.env.NODE_ENV": JSON.stringify(ENV),

		// Retro
		"process.env.CMD": JSON.stringify(CMD),
		"process.env.ENV": JSON.stringify(ENV),
		"process.env.WWW_DIR": JSON.stringify(WWW_DIR),
		"process.env.SRC_DIR": JSON.stringify(SRC_DIR),
		"process.env.OUT_DIR": JSON.stringify(OUT_DIR),
	},
	loader: {
		".js": "jsx",
	},
	logLevel: "silent",
	minify: ENV === "production",
	sourcemap: true, // TODO
}

async function resolveConfig(): Promise<esbuild.BuildOptions> {
	try {
		await fs.promises.stat("retro.config.js")
	} catch {
		return {}
	}
	return require(path.join(process.cwd(), "retro.config.js"))
}

let reactResult: esbuild.BuildResult | null = null
let indexResult: esbuild.BuildResult | esbuild.BuildIncremental | null = null

async function build(): Promise<BackendResponse> {
	const buildRes: BackendResponse = {
		Metafile: {
			Vendor: null,
			Bundle: null,
		},
		Warnings: [],
		Errors: [],
	}

	const config = await resolveConfig()

	try {
		// React, React DOM
		reactResult = await esbuild.build({
			...common,

			bundle: true,
			entryNames: ENV !== "production" ? undefined : "[dir]/[name]__[hash]",
			entryPoints: [path.join(__dirname, "react.js")],
			metafile: true,
			outdir: OUT_DIR,
		})

		// Attach metafile
		buildRes.Metafile.Vendor = reactResult.metafile!

		// User code
		indexResult = await esbuild.build({
			...config,
			...common,

			define: { ...config.define, ...common.define },
			loader: { ...config.loader, ...common.loader },

			bundle: true,
			entryNames: ENV !== "production" ? undefined : "[dir]/[name]__[hash]",
			entryPoints: [path.join(SRC_DIR, "index.js")],
			metafile: true,
			outdir: OUT_DIR,

			external: ["react", "react-dom"], // Dedupe React APIs (because of vendor)
			inject: [path.join(__dirname, "shims/require.js")], // Add support for vendor
			plugins: config?.plugins,

			incremental: ENV === "development",
		})

		// Attach metafile
		buildRes.Metafile.Bundle = indexResult.metafile!

		if (indexResult.warnings.length > 0) {
			buildRes.Warnings = indexResult.warnings
		}
	} catch (caught) {
		if (caught.errors.length > 0) {
			buildRes.Errors = caught.errors
		}
		if (caught.warnings.length > 0) {
			buildRes.Warnings = caught.warnings
		}
	}

	return buildRes
}

async function rebuild(): Promise<BackendResponse> {
	if (indexResult?.rebuild === undefined) throw new Error("Internal error")

	const rebuildRes: BackendResponse = {
		Metafile: {
			Vendor: null,
			Bundle: null,
		},
		Warnings: [],
		Errors: [],
	}

	try {
		const result2 = await indexResult.rebuild()
		if (result2.warnings.length > 0) {
			rebuildRes.Warnings = result2.warnings
		}
	} catch (caught) {
		if (caught.errors.length > 0) {
			rebuildRes.Errors = caught.errors
		}
		if (caught.warnings.length > 0) {
			rebuildRes.Warnings = caught.warnings
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
