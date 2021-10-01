import * as esbuild from "esbuild"

// Bundle metadata and structured warnings and errors
export interface BundleMetadata {
	Metafile: esbuild.Metafile
	Warnings: esbuild.Message[]
	Errors: esbuild.Message[]
}

// Message for completed build vendor and client events
export interface BuildVendorAndClientDoneMessage {
	Kind: "build_done"
	Data: {
		Vendor: BundleMetadata
		Client: BundleMetadata
	}
}

// Message for completed rebuild client events
export interface RebuildClientDoneMessage {
	Kind: "rebuild_done"
	Data: {
		Client: BundleMetadata
	}
}
