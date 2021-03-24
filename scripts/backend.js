const esbuild = require("esbuild")
const fs = require("fs")
const path = require("path")

const WWW_DIR = process.env["WWW_DIR"]
const SRC_DIR = process.env["SRC_DIR"]
const OUT_DIR = process.env["OUT_DIR"]

const env = process.env.NODE_ENV ?? "development"

const common = {
	color: true,
	define: {
		__DEV__: JSON.stringify(env !== "production"),
		"process.env.NODE_ENV": JSON.stringify(env),
	},
	loader: {
		".js": "jsx",
	},
	minify: env === "production",
	sourcemap: true,
}

async function resolvePlugins() {
	try {
		await fs.promises.stat("plugins.js")
	} catch {
		return []
	}
	return require(path.join(process.cwd(), "plugins.js"))
}

let result = null

async function build() {
	// TODO: Refactor to Go
	const sources = (await fs.promises.readdir(OUT_DIR)).map(source => path.join(OUT_DIR, source))
	for (const v of sources) {
		// prettier-ignore
		if (!v.endsWith(".css") && !v.endsWith(".css.map") &&
				!v.endsWith(".js") && !v.endsWith(".js.map")) {
			continue
		}
		await fs.promises.unlink(v)
	}

	const response = {
		warnings: [],
		errors: [],
	}

	// TODO: Refactor to Go
	// prettier-ignore
	await fs.promises.copyFile(
		path.join(WWW_DIR, "index.html"),
		path.join(OUT_DIR, "index.html"),
	)

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
			...common,
			bundle: true,
			entryPoints: [path.join(SRC_DIR, "index.js")], // TODO: Add support for ".jsx", ".ts", and ".tsx"
			outfile: path.join(OUT_DIR, "bundle.js"),

			external: ["react", "react-dom"], // Dedupe React APIs (because of vendor)
			inject: ["scripts/shims/require.js"], // Add support for vendor
			plugins: await resolvePlugins(),

			incremental: true,
			watch: true,
		})
		if (result?.warnings?.length > 0) {
			response.warnings = result.warnings
		}
	} catch (error) {
		if (error?.errors?.length > 0) {
			response.errors = error.errors
		}
		if (error?.warnings?.length > 0) {
			response.warnings = error.warnings
		}
	}

	// TODO: Refactor to Go
	try {
		await fs.promises.stat(path.join(OUT_DIR, "bundle.css"))
	} catch {
		// Create bundle.css so <link ... href="/bundle.css"> does not error
		await fs.promises.writeFile(path.join(OUT_DIR, "bundle.css"), "")
	}

	stdout({
		Kind: "build-done",
		Data: response,
	})
}

async function rebuild() {
	if (result === null) return

	const response = {
		warnings: [],
		errors: [],
	}

	try {
		const result2 = await result.rebuild()
		if (result2?.warnings?.length > 0) {
			response.warnings = result2.warnings
		}
	} catch (error) {
		if (error?.errors?.length > 0) {
			response.errors = error.errors
		}
		if (error?.warnings?.length > 0) {
			response.warnings = error.warnings
		}
	}

	stdout({
		Kind: "rebuild-done",
		Data: response,
	})
}

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

async function main() {
	esbuild.build({})

	while (true) {
		const jsonstr = await readline()
		const msg = JSON.parse(jsonstr)
		try {
			switch (msg.Kind) {
				case "build":
					await build()
					break
				case "rebuild":
					await rebuild()
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
