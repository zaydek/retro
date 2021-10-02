import esbuild from "esbuild"
import fsPromises from "fs/promises"
import path from "path"

import {
	NODE_ENV,
	RETRO_CMD,
	RETRO_WWW_DIR,
	RETRO_SRC_DIR,
	RETRO_OUT_DIR,
} from "./env"

// Resolves `retro.config.js`
export async function resolveUserConfiguration(): Promise<esbuild.BuildOptions> {
	try {
		await fsPromises.stat("retro.config.js")
	} catch {
		return {}
	}
	return require(path.join(process.cwd(), "retro.config.js"))
}

// Base configuration for vendor and client bundles
export const baseConfiguration: esbuild.BuildOptions = {
	// Always bundle
	bundle: true,

	// Propagate environmental variables
	define: {
		// React and React DOM environmental variables
		"process.env.NODE_ENV": JSON.stringify(NODE_ENV),

		// Retro environmental variables
		"process.env.RETRO_CMD": JSON.stringify(RETRO_CMD),
		"process.env.RETRO_WWW_DIR": JSON.stringify(RETRO_WWW_DIR),
		"process.env.RETRO_SRC_DIR": JSON.stringify(RETRO_SRC_DIR),
		"process.env.RETRO_OUT_DIR": JSON.stringify(RETRO_OUT_DIR),
	},

	// Load JavaScript as JavaScript React
	loader: {
		".js": "jsx",
	},

	// Don't log because warnings and errors are handled programmatically
	logLevel: "silent",

	// Includes the generated hashed filenames
	metafile: true,

	// Minify for production
	minify: NODE_ENV === "production",

	// Add sourcemaps
	sourcemap: true,
}

// Builds the client configuration from the base and user configurations
export const buildClientConfiguration = (userConfiguration: esbuild.BuildOptions): esbuild.BuildOptions => ({
	...baseConfiguration,
	...userConfiguration,

	// Global variables
	define: {
		...baseConfiguration.define,
		...userConfiguration.define,
	},

	// Dedupe React APIs from `bundle.js`; React APIs are bundled in `vendor.js`

	// Dedupe vendor APIs; vendor APIs are bundled in `vendor.js`
	external: [
		"react",
		"react-dom",
		"react-dom/server",
	],

	// Enable incremental compilation for development
	incremental: NODE_ENV === "development",

	// Vendor API shims
	inject: [path.join(__dirname, "require.js")],

	loader: {
		...baseConfiguration.loader,
		...userConfiguration.loader,
	},
})
