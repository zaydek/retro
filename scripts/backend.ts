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
		// For React and React DOM
		"process.env.NODE_ENV": JSON.stringify(ENV),

		// For Retro
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

let vendorResult: esbuild.BuildResult | null = null
let bundleResult: esbuild.BuildResult | esbuild.BuildIncremental | null = null

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
		vendorResult = await esbuild.build({
			...common,

			// Add support for target
			target: config.target,

			bundle: true,
			entryNames: ENV !== "production" ? undefined : "[dir]/[name]__[hash]",
			entryPoints: {
				"vendor": path.join(__dirname, "react.js"),
			},
			metafile: true,
			outdir: OUT_DIR,
		})
		buildRes.Metafile.Vendor = vendorResult.metafile!

		bundleResult = await esbuild.build({
			...config,
			...common,

			define: { ...config.define, ...common.define },
			loader: { ...config.loader, ...common.loader },

			bundle: true,
			entryNames: ENV !== "production" ? undefined : "[dir]/[name]__[hash]",
			entryPoints: {
				"bundle": path.join(SRC_DIR, "index.js"),
			},
			metafile: true,
			outdir: OUT_DIR,

			external: ["react", "react-dom"], // Dedupe React APIs
			inject: [path.join(__dirname, "shims/require.js")], // Add React APIs
			plugins: config?.plugins,

			incremental: ENV === "development",
		})
		buildRes.Metafile.Bundle = bundleResult.metafile!

		if (bundleResult.warnings.length > 0) {
			buildRes.Warnings = bundleResult.warnings
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
	if (bundleResult?.rebuild === undefined) {
		return await build()
	}

	const rebuildRes: BackendResponse = {
		Metafile: {
			Vendor: null,
			Bundle: null,
		},
		Warnings: [],
		Errors: [],
	}

	try {
		const result2 = await bundleResult.rebuild()
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
