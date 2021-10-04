export const NODE_ENV = (function () {
	const env = process.env["NODE_ENV"]
	if (env === "") {
		throw new Error(`process.env["NODE_ENV"] === ""`)
	}
	return env
})()

export const RETRO_CMD = (function () {
	const env = process.env["RETRO_CMD"]
	if (env === "") {
		throw new Error(`process.env["RETRO_CMD"] === ""`)
	}
	return env
})()

export const RETRO_WWW_DIR = (function () {
	const env = process.env["RETRO_WWW_DIR"]
	if (env === "") {
		throw new Error(`process.env["RETRO_WWW_DIR"] === ""`)
	}
	return env
})()

export const RETRO_SRC_DIR = (function () {
	const env = process.env["RETRO_SRC_DIR"]
	if (env === "") {
		throw new Error(`process.env["RETRO_SRC_DIR"] === ""`)
	}
	return env
})()

export const RETRO_OUT_DIR = (function () {
	const env = process.env["RETRO_OUT_DIR"]
	if (env === "") {
		throw new Error(`process.env["RETRO_OUT_DIR"] === ""`)
	}
	return env
})()
