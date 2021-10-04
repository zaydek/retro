import * as esbuild from "esbuild"

export interface BundleInfo {
	Metafile: esbuild.Metafile
	Warnings: esbuild.Message[]
	Errors: esbuild.Message[]
}

export interface BuildVendorAndClientDoneMessage {
	Kind: "build_done"
	Data: {
		Vendor: BundleInfo
		Client: BundleInfo
		ClientAppOnly: BundleInfo
	}
}

export interface RebuildClientDoneMessage {
	Kind: "rebuild_done"
	Data: {
		Client: BundleInfo
	}
}
