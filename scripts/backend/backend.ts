import * as esbuild from "esbuild"
import * as path from "path"
import * as t from "./types"
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

////////////////////////////////////////////////////////////////////////////////

async function buildVendorBundle(): Promise<t.BundleInfo> {
	const vendorInfo: t.BundleInfo = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}

	try {
		const bundle = await esbuild.build(vendorConfig)
		if (bundle.errors.length > 0) { vendorInfo.Errors = bundle.errors }
		if (bundle.warnings.length > 0) { vendorInfo.Warnings = bundle.warnings }
		vendorInfo.Metafile = bundle.metafile
	} catch (caught) {
		if (caught.errors.length > 0) { vendorInfo.Errors = caught.errors }
		if (caught.warnings.length > 0) { vendorInfo.Warnings = caught.warnings }
	}

	return vendorInfo
}

////////////////////////////////////////////////////////////////////////////////

async function buildClientBundle(config: esbuild.BuildOptions): Promise<{ clientInfo: t.BundleInfo, clientAppInfo: t.BundleInfo }> {
	const clientInfo: t.BundleInfo = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}

	const clientAppInfo: t.BundleInfo = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}

	try {
		globalClientBundle = await esbuild.build(clientConfigFromUserConfig(config))
		if (globalClientBundle.errors.length > 0) { clientInfo.Errors = globalClientBundle.errors }
		if (globalClientBundle.warnings.length > 0) { clientInfo.Warnings = globalClientBundle.warnings }
		clientInfo.Metafile = globalClientBundle.metafile
	} catch (caught) {
		if (caught.errors.length > 0) { clientInfo.Errors = caught.errors }
		if (caught.warnings.length > 0) { clientInfo.Warnings = caught.warnings }
	}

	if (RETRO_CMD !== "build") {
		return { clientInfo, clientAppInfo }
	}

	try {
		const bundle = await esbuild.build(clientAppConfigFromUserConfig(config))
		if (bundle.errors.length > 0) { clientAppInfo.Errors = bundle.errors }
		if (bundle.warnings.length > 0) { clientAppInfo.Warnings = bundle.warnings }
		clientAppInfo.Metafile = bundle.metafile
	} catch (caught) {
		if (caught.errors.length > 0) { clientAppInfo.Errors = caught.errors }
		if (caught.warnings.length > 0) { clientAppInfo.Warnings = caught.warnings }
	}

	return { clientInfo, clientAppInfo }
}

////////////////////////////////////////////////////////////////////////////////

async function buildVendorAndClientBundles(config): Promise<{ vendorInfo: t.BundleInfo, clientInfo: t.BundleInfo, clientAppInfo: t.BundleInfo }> {
	const vendorInfo = await buildVendorBundle()
	const { clientInfo, clientAppInfo } = await buildClientBundle(config)
	return { vendorInfo, clientInfo, clientAppInfo }
}

////////////////////////////////////////////////////////////////////////////////

async function rebuildClientBundle(config: esbuild.BuildOptions): Promise<t.BundleInfo> {
	if (globalClientBundle === null) {
		return (await buildClientBundle(config)).clientInfo
	}

	const clientInfo: t.BundleInfo = {
		Metafile: null,
		Warnings: [],
		Errors: [],
	}

	try {
		await globalClientBundle.rebuild()
		if (globalClientBundle.errors.length > 0) { clientInfo.Errors = globalClientBundle.errors }
		if (globalClientBundle.warnings.length > 0) { clientInfo.Warnings = globalClientBundle.warnings }
		clientInfo.Metafile = globalClientBundle.metafile
	} catch (caught) {
		if (caught.errors.length > 0) { clientInfo.Errors = caught.errors }
		if (caught.warnings.length > 0) { clientInfo.Warnings = caught.warnings }
	}

	return clientInfo
}

////////////////////////////////////////////////////////////////////////////////

async function main(): Promise<void> {
	let config: esbuild.BuildOptions = {}
	try {
		config = require(path.join(process.cwd(), "retro.config"))
	} catch { }

	esbuild.initialize({})
	while (true) {
		const action = await readline()
		switch (action) {
			case "build":
				const { vendorInfo, clientInfo, clientAppInfo } = await buildVendorAndClientBundles(config)
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
				const clientInfo = await rebuildClientBundle(config)
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
