import esbuild from "esbuild"
import path from "path"

import {
	NODE_ENV,
	RETRO_CMD,
	RETRO_OUT_DIR,
	RETRO_SRC_DIR,
	RETRO_WWW_DIR,
} from "./env"

export const vendorConfig: esbuild.BuildOptions = {
	bundle: true,
	entryNames: NODE_ENV !== "production"
		? undefined
		: "[dir]/[name]__[hash]",
	entryPoints: {
		"vendor": path.join(__dirname, "vendor.js"),
	},
	logLevel: "silent",
	metafile: true,
	minify: NODE_ENV === "production",
	outdir: RETRO_OUT_DIR,
	sourcemap: true,
}

export const clientConfigFromUserConfig = (userConfig: esbuild.BuildOptions): esbuild.BuildOptions => ({
	...userConfig,
	bundle: true,
	define: {
		...userConfig.define,
		"process.env.NODE_ENV": JSON.stringify(NODE_ENV),
		"process.env.RETRO_CMD": JSON.stringify(RETRO_CMD),
		"process.env.RETRO_WWW_DIR": JSON.stringify(RETRO_WWW_DIR),
		"process.env.RETRO_SRC_DIR": JSON.stringify(RETRO_SRC_DIR),
		"process.env.RETRO_OUT_DIR": JSON.stringify(RETRO_OUT_DIR),
	},
	entryNames: NODE_ENV !== "production"
		? undefined
		: "[dir]/[name]__[hash]",
	entryPoints: {
		...userConfig.entryPoints,
		"client": path.join(RETRO_SRC_DIR, "index.js"),
	},
	external: [
		// Only React APIs are shimmed
		"react",
		"react-dom",
		"react-dom/server",
	],
	incremental: RETRO_CMD === "dev",
	inject: [
		// Only React APIs are shimmed
		path.join(__dirname, "require.js"),
	],
	loader: {
		...userConfig.loader,
		".js": "jsx",
	},
	logLevel: "silent",
	metafile: true,
	minify: NODE_ENV === "production",
	outdir: RETRO_OUT_DIR,
	sourcemap: true,
})

export const clientAppConfigFromUserConfig = (userConfig: esbuild.BuildOptions): esbuild.BuildOptions => ({
	...clientConfigFromUserConfig(userConfig),
	entryNames: undefined, // No-op
	entryPoints: {
		...userConfig.entryPoints,
		"App.js": path.join(RETRO_SRC_DIR, "App.js"),
	},
	outdir: path.join(RETRO_OUT_DIR, ".retro"),
	platform: "node",
})
