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
	RETRO_OUT_DIR,
	RETRO_SRC_DIR,
} from "./env"

function stdout(message:
	| t.BuildVendorAndClientDoneMessage
	| t.RebuildClientDoneMessage
): void {
	console.log(JSON.stringify(message))
}

// Describes `retro.config.js`
let globalUserConfiguration: esbuild.BuildOptions | null = null

// Describes the bundled vendor esbuild result
let globalVendorBuildResult: esbuild.BuildResult | null = null

// Describes the bundled client esbuild result
let globalClientBuildResult: esbuild.BuildResult | esbuild.BuildIncremental | null = null

// Builds the vendor bundle (e.g. React) and sets the global vendor variable
async function buildVendorBundle(): Promise<t.BundleMetadata> {
	const vendor: t.BundleMetadata = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}

	try {
		globalVendorBuildResult = await esbuild.build({
			...baseConfiguration,
			entryNames: NODE_ENV !== "production"
				? undefined
				: "[dir]/[name]__[hash]",
			entryPoints: {
				"vendor": path.join(__dirname, "vendor.js"),
			},
			outdir: RETRO_OUT_DIR,
		})
		if (globalVendorBuildResult.errors.length > 0) { vendor.Errors = globalVendorBuildResult.errors }
		if (globalVendorBuildResult.warnings.length > 0) { vendor.Warnings = globalVendorBuildResult.warnings }
		vendor.Metafile = globalVendorBuildResult.metafile
	} catch (caught) {
		if (caught.errors.length > 0) { vendor.Errors = caught.errors }
		if (caught.warnings.length > 0) { vendor.Warnings = caught.warnings }
	}

	return vendor
}

// Builds the client bundle (e.g. Retro) and sets the global client variable
async function buildClientBundle(): Promise<t.BundleMetadata> {
	const client: t.BundleMetadata = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}

	try {
		globalClientBuildResult = await esbuild.build({
			...buildClientConfiguration(globalUserConfiguration),
			entryNames: NODE_ENV !== "production"
				? undefined
				: "[dir]/[name]__[hash]",
			entryPoints: {
				"client": path.join(RETRO_SRC_DIR, "index.js"),
			},
			outdir: RETRO_OUT_DIR,
		})
		if (globalClientBuildResult.errors.length > 0) { client.Errors = globalClientBuildResult.errors }
		if (globalClientBuildResult.warnings.length > 0) { client.Warnings = globalClientBuildResult.warnings }
		client.Metafile = globalClientBuildResult.metafile
	} catch (caught) {
		if (caught.errors.length > 0) { client.Errors = caught.errors }
		if (caught.warnings.length > 0) { client.Warnings = caught.warnings }
	}

	try {
		await esbuild.build({
			...buildClientConfiguration(globalUserConfiguration),
			// entryNames: NODE_ENV !== "production"
			// 	? undefined
			// 	: "[dir]/[name]__[hash]",
			// entryPoints: {
			// 	"client": path.join(RETRO_SRC_DIR, "index.js"),
			// },
			entryPoints: [path.join(RETRO_SRC_DIR, "App.js")],
			outdir: path.join(RETRO_OUT_DIR, ".retro"),
			platform: "node",
		})
		// if (globalClientBuildResult.errors.length > 0) { client.Errors = globalClientBuildResult.errors }
		// if (globalClientBuildResult.warnings.length > 0) { client.Warnings = globalClientBuildResult.warnings }
		// client.Metafile = globalClientBuildResult.metafile
	} catch (caught) {
		// if (caught.errors.length > 0) { client.Errors = caught.errors }
		// if (caught.warnings.length > 0) { client.Warnings = caught.warnings }
	}

	return client
}

// Builds the vendor and client bundles
async function buildVendorAndClientBundles(): Promise<[t.BundleMetadata, t.BundleMetadata]> {
	const vendor = await buildVendorBundle()
	const client = await buildClientBundle()
	return [vendor, client]
}

// Builds or rebuild the client bundle
async function rebuildClientBundle(): Promise<t.BundleMetadata> {
	if (globalClientBuildResult === null || globalClientBuildResult.rebuild === undefined) {
		return await buildClientBundle()
	}

	const client: t.BundleMetadata = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}

	try {
		const clientResult = await globalClientBuildResult.rebuild()
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
				const [vendor, client] = await buildVendorAndClientBundles()
				stdout({
					Kind: "build_done",
					Data: {
						Vendor: vendor,
						Client: client,
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
