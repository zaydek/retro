import * as esbuild from "esbuild"
import * as path from "path"
import readline from "./readline"

import {
	clientConfigFromUserConfig,
	vendorConfig,
} from "./configs"

let globalClientBundle: esbuild.BuildResult | esbuild.BuildIncremental | null = null

interface BundleInfo {
	Metafile: esbuild.Metafile
	Errors: esbuild.Message[]
	Warnings: esbuild.Message[]
}

async function buildVendorBundle(): Promise<BundleInfo> {
	const vendorInfo: BundleInfo = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}
	try {
		const bundle = await esbuild.build(vendorConfig)
		vendorInfo.Metafile = bundle.metafile
		if (bundle.errors.length > 0) { vendorInfo.Errors = bundle.errors }
		if (bundle.warnings.length > 0) { vendorInfo.Warnings = bundle.warnings }
	} catch (caught) {
		if (caught.errors.length > 0) { vendorInfo.Errors = caught.errors }
		if (caught.warnings.length > 0) { vendorInfo.Warnings = caught.warnings }
	}
	return vendorInfo
}

async function buildClientBundle(userConfig: esbuild.BuildOptions): Promise<BundleInfo> {
	const clientInfo: BundleInfo = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}
	try {
		globalClientBundle = await esbuild.build(clientConfigFromUserConfig(userConfig))
		clientInfo.Metafile = globalClientBundle.metafile
		if (globalClientBundle.errors.length > 0) { clientInfo.Errors = globalClientBundle.errors }
		if (globalClientBundle.warnings.length > 0) { clientInfo.Warnings = globalClientBundle.warnings }
	} catch (caught) {
		if (caught.errors.length > 0) { clientInfo.Errors = caught.errors }
		if (caught.warnings.length > 0) { clientInfo.Warnings = caught.warnings }
	}
	return clientInfo
}

async function buildVendorAndClientBundles(userConfig: esbuild.BuildOptions): Promise<{ vendorInfo: BundleInfo, clientInfo: BundleInfo }> {
	const vendorInfo = await buildVendorBundle()
	const clientInfo = await buildClientBundle(userConfig)
	return { vendorInfo, clientInfo }
}

async function rebuildClientBundle(userConfig: esbuild.BuildOptions): Promise<BundleInfo> {
	if (globalClientBundle === null) {
		return await buildClientBundle(userConfig)
	}
	const clientInfo: BundleInfo = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}
	try {
		await globalClientBundle.rebuild()
		clientInfo.Metafile = globalClientBundle.metafile
		if (globalClientBundle.errors.length > 0) { clientInfo.Errors = globalClientBundle.errors }
		if (globalClientBundle.warnings.length > 0) { clientInfo.Warnings = globalClientBundle.warnings }
	} catch (caught) {
		if (caught.errors.length > 0) { clientInfo.Errors = caught.errors }
		if (caught.warnings.length > 0) { clientInfo.Warnings = caught.warnings }
	}
	return clientInfo
}

async function main(): Promise<void> {
	let userConfig: esbuild.BuildOptions = {}
	try {
		userConfig = require(path.join(process.cwd(), "retro.config"))
	} catch { }

	esbuild.initialize({})
	while (true) {
		const action = await readline()
		switch (action) {
			case "build":
				const { vendorInfo, clientInfo } = await buildVendorAndClientBundles(userConfig)
				console.log(
					JSON.stringify({
						vendorInfo,
						clientInfo,
					}),
				)
				break
			case "rebuild": {
				const clientInfo = await rebuildClientBundle(userConfig)
				console.log(
					JSON.stringify({
						clientInfo,
					}),
				)
				break
			}
			default:
				throw new Error("Internal error")
		}
	}
}

main()
