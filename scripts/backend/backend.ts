import * as esbuild from "esbuild"
import * as path from "path"
import readline from "./readline"

import {
	clientAppConfigFromUserConfig,
	clientConfigFromUserConfig,
	vendorConfig,
} from "./configs"

import {
	RETRO_CMD,
} from "./env"

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

async function buildClientBundle(userConfig: esbuild.BuildOptions): Promise<{ clientInfo: BundleInfo, clientAppInfo: BundleInfo }> {
	const clientInfo: BundleInfo = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}
	const clientAppInfo: BundleInfo = {
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
	if (RETRO_CMD !== "build") {
		return { clientInfo, clientAppInfo }
	}
	try {
		const bundle = await esbuild.build(clientAppConfigFromUserConfig(userConfig))
		clientAppInfo.Metafile = bundle.metafile
		if (bundle.errors.length > 0) { clientAppInfo.Errors = bundle.errors }
		if (bundle.warnings.length > 0) { clientAppInfo.Warnings = bundle.warnings }
	} catch (caught) {
		if (caught.errors.length > 0) { clientAppInfo.Errors = caught.errors }
		if (caught.warnings.length > 0) { clientAppInfo.Warnings = caught.warnings }
	}
	return { clientInfo, clientAppInfo }
}

async function buildVendorAndClientBundles(userConfig: esbuild.BuildOptions): Promise<{ vendorInfo: BundleInfo, clientInfo: BundleInfo, clientAppInfo: BundleInfo }> {
	const vendorInfo = await buildVendorBundle()
	const { clientInfo, clientAppInfo } = await buildClientBundle(userConfig)
	return { vendorInfo, clientInfo, clientAppInfo }
}

async function rebuildClientBundle(userConfig: esbuild.BuildOptions): Promise<BundleInfo> {
	if (globalClientBundle === null) {
		return (await buildClientBundle(userConfig)).clientInfo
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
				const { vendorInfo, clientInfo, clientAppInfo } = await buildVendorAndClientBundles(userConfig)
				console.log(
					JSON.stringify(
						{
							vendorInfo,
							clientInfo,
							clientAppInfo,
						},
					),
				)
				break
			case "rebuild": {
				const clientInfo = await rebuildClientBundle(userConfig)
				console.log(
					JSON.stringify(
						{
							clientInfo,
						},
					),
				)
				break
			}
			default:
				throw new Error("Internal error")
		}
	}
}

main()
