import * as esbuild from "esbuild"
import * as path from "path"
import * as t from "./types"
import readline from "./readline"

import {
	baseConfiguration,
	buildClientConfiguration,
	resolveUserConfiguration,
} from "./configuration"

import {
	NODE_ENV,
	RETRO_CMD,
	RETRO_OUT_DIR,
	RETRO_SRC_DIR,
} from "./env"

function stdout(message:
	| t.BuildVendorAndClientDoneMessage
	| t.RebuildClientDoneMessage
): void {
	console.log(JSON.stringify(message))
}

// retro.config.js
let globalUserConfiguration: esbuild.BuildOptions | null = null

// react, react-dom, react-dom/server
let globalVendorEntryPoint: esbuild.BuildResult | null = null

// src/index.js
let globalClientEntryPoint: esbuild.BuildResult | esbuild.BuildIncremental | null = null

////////////////////////////////////////////////////////////////////////////////

async function buildVendorBundle(): Promise<t.BundleInfo> {
	const vendor: t.BundleInfo = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}

	try {
		globalVendorEntryPoint = await esbuild.build({
			...baseConfiguration,
			entryNames: NODE_ENV !== "production"
				? undefined
				: "[dir]/[name]__[hash]",
			entryPoints: {
				"vendor": path.join(__dirname, "vendor.js"),
			},
			outdir: RETRO_OUT_DIR,
		})
		if (globalVendorEntryPoint.errors.length > 0) { vendor.Errors = globalVendorEntryPoint.errors }
		if (globalVendorEntryPoint.warnings.length > 0) { vendor.Warnings = globalVendorEntryPoint.warnings }
		vendor.Metafile = globalVendorEntryPoint.metafile
	} catch (caught) {
		if (caught.errors.length > 0) { vendor.Errors = caught.errors }
		if (caught.warnings.length > 0) { vendor.Warnings = caught.warnings }
	}

	return vendor
}

async function buildClientBundle(): Promise<[t.BundleInfo, t.BundleInfo]> {
	const client: t.BundleInfo = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}

	const clientAppOnly: t.BundleInfo = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}

	try {
		globalClientEntryPoint = await esbuild.build({
			...buildClientConfiguration(globalUserConfiguration),
			entryNames: NODE_ENV !== "production"
				? undefined
				: "[dir]/[name]__[hash]",
			entryPoints: {
				"client": path.join(RETRO_SRC_DIR, "index.js"),
			},
			outdir: RETRO_OUT_DIR,
		})
		if (globalClientEntryPoint.errors.length > 0) { client.Errors = globalClientEntryPoint.errors }
		if (globalClientEntryPoint.warnings.length > 0) { client.Warnings = globalClientEntryPoint.warnings }
		client.Metafile = globalClientEntryPoint.metafile
	} catch (caught) {
		if (caught.errors.length > 0) { client.Errors = caught.errors }
		if (caught.warnings.length > 0) { client.Warnings = caught.warnings }
	}

	if (RETRO_CMD === "build") {
		try {
			const clientAppEntryPoint = await esbuild.build({
				...buildClientConfiguration(globalUserConfiguration),
				entryPoints: [path.join(RETRO_SRC_DIR, "App.js")],
				outdir: path.join(RETRO_OUT_DIR, ".retro"),
				platform: "node",
			})
			if (clientAppEntryPoint.errors.length > 0) { clientAppOnly.Errors = clientAppEntryPoint.errors }
			if (clientAppEntryPoint.warnings.length > 0) { clientAppOnly.Warnings = clientAppEntryPoint.warnings }
			clientAppOnly.Metafile = clientAppEntryPoint.metafile
		} catch (caught) {
			if (caught.errors.length > 0) { clientAppOnly.Errors = caught.errors }
			if (caught.warnings.length > 0) { clientAppOnly.Warnings = caught.warnings }
		}
	}

	return [client, clientAppOnly]
}

async function buildVendorAndClientBundles(): Promise<{ vendor: t.BundleInfo, client: t.BundleInfo, clientAppOnly: t.BundleInfo }> {
	const vendor = await buildVendorBundle()
	const [client, clientAppOnly] = await buildClientBundle()
	return { vendor, client, clientAppOnly }
}

async function rebuildClientBundle(): Promise<t.BundleInfo> {
	if (globalClientEntryPoint === null || globalClientEntryPoint.rebuild === undefined) {
		return await buildClientBundle()[0]
	}

	const client: t.BundleInfo = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}

	try {
		const clientResult = await globalClientEntryPoint.rebuild()
		if (clientResult.errors.length > 0) { client.Errors = clientResult.errors }
		if (clientResult.warnings.length > 0) { client.Warnings = clientResult.warnings }
		client.Metafile = clientResult.metafile
	} catch (caught) {
		if (caught.errors.length > 0) { client.Errors = caught.errors }
		if (caught.warnings.length > 0) { client.Warnings = caught.warnings }
	}

	return client
}

function sleep(durationMs: number): Promise<void> {
	return new Promise(resolve => setTimeout(resolve, durationMs))
}

// This becomes a Node.js IPC process, from Go to JavaScript. Messages are sent
// as plaintext strings (actions) and received as JSON-encoded payloads.
//
// stdout messages that aren't encoded should be logged regardless because
// plugins can implement logging. stderr messages are exceptions and should
// terminate the Node.js runtime.
async function main(): Promise<void> {
	esbuild.initialize({})
	globalUserConfiguration = await resolveUserConfiguration()

	while (true) {
		const action = await readline()
		switch (action) {
			case "build":
				const { vendor, client, clientAppOnly } = await buildVendorAndClientBundles()
				stdout({
					Kind: "build_done",
					Data: {
						Vendor: vendor,
						Client: client,
						ClientAppOnly: clientAppOnly,
					},
				})
				break
			case "rebuild": {
				const client = await rebuildClientBundle()
				stdout({
					Kind: "rebuild_done",
					Data: {
						Client: client,
					},
				})
				break
			}
			default:
				throw new Error("Internal error")
		}
		// Add a micro-delay to prevent high CPU usage
		await sleep(10)
	}
}

main()
