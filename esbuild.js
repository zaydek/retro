const esbuild = require("esbuild")
const path = require("path")

////////////////////////////////////////////////////////////////////////////////

const nodeReadline = require("readline")

const readline = (() => {
	async function* createReadlineGenerator() {
		const nodeReadlineInterface = nodeReadline.createInterface({ input: process.stdin })
		for await (const line of nodeReadlineInterface) {
			yield line
		}
	}
	const generate = createReadlineGenerator()
	return async () => {
		const result = await generate.next()
		return result.value
	}
})()


////////////////////////////////////////////////////////////////////////////////

let globalClientResult = null

async function buildClient(config) {
	const client = {
		Metafile: null,
		Errors: [],
		Warnings: [],
	}

	try {
		globalClientResult = await esbuild.build({
			...config,
			bundle: true,
			define: {
				...config.define,
				"process.env.NODE_ENV": JSON.stringify("development"),   // TODO: Change to an environmental variable
				"process.env.RETRO_CMD": JSON.stringify("dev"),          // TODO: Change to an environmental variable
				"process.env.RETRO_WWW_DIR": JSON.stringify("www"),      // TODO: Change to an environmental variable
				"process.env.RETRO_SRC_DIR": JSON.stringify("examples"), // TODO: Change to an environmental variable
				"process.env.RETRO_OUT_DIR": JSON.stringify("out"),      // TODO: Change to an environmental variable
			},
			entryPoints: {
				...config.entryPoints,
				"client": path.join("examples", "index.js"),
			},
			external: [
				// Only React APIs are shimmed
				"react",
				"react-dom",
				"react-dom/server",
			],
			inject: [
				// Only React APIs are shimmed
				path.join("scripts", "require.js"),
			],
			loader: {
				...config.loader,
				".js": "jsx",
			},
			logLevel: "silent",
			metafile: true,
			minify: process.env.NODE_ENV === "production", // TODO: Change to an environmental variable
			outdir: "out",                                 // TODO: Change to an environmental variable
			sourcemap: true,
		})
		if (globalClientResult.errors.length > 0) { client.Errors = globalClientResult.errors }
		if (globalClientResult.warnings.length > 0) { client.Warnings = globalClientResult.warnings }
		client.Metafile = globalClientResult.metafile
	} catch (caught) {
		if (caught.errors.length > 0) { client.Errors = caught.errors }
		if (caught.warnings.length > 0) { client.Warnings = caught.warnings }
	}

	return client
}

async function main() {
	let config = {}
	try {
		config = require("./retro.config")
	} catch { }

	esbuild.initialize({})
	while (true) {
		const action = await readline()
		switch (action) {
			case "build": {
				const client = await buildClient(config)
				console.log(
					JSON.stringify({
						client,
					}),
				)
				break
			}
			// case "rebuild": {
			// 	const client = await rebuildClient()
			// 	console.log(
			// 		JSON.stringify({
			// 			client,
			// 		}),
			// 	)
			// 	break
			// }
			default:
				throw new Error("Internal error")
		}
	}
}

main()
