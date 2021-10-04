// This is a post-installation script that resolves the canonical binary

const fsPromises = require("fs/promises")
const os = require("os")
const path = require("path")

const CANONICAL_BINARY_FILENAME = (function () {
	const package = require("./package.json")
	const binaryKey = Object.keys(package.bin)[0]
	return binaryKey + (process.platform === "win32" ? ".exe" : "")
})()

// Maps supported architecture keys to abstract binary filenames
const supported = {
	"platform=darwin arch=arm64 endianness=LE": "darwin-64",
	"platform=darwin arch=x64 endianness=LE": "darwin-64",
	"platform=linux arch=x64 endianness=LE": "linux-64",
	"platform=win32 arch=x64 endianness=LE": "windows-64.exe",
}

async function main() {
	const architectureKey =
		`platform=${process.platform} ` +
		`arch=${os.arch()} ` +
		`endianness=${os.endianness()}`
	const binaryFilename = supported[architectureKey]
	if (binaryFilename === undefined) {
		throw new Error(`postinstall.js: Architecture key \`${architectureKey}\` not supported. ` +
			`Create an issue at https://github.com/zaydek/retro.`)
	}
	const src = path.join(__dirname, "bin", binaryFilename)
	const dst = path.join(__dirname, "bin", CANONICAL_BINARY_FILENAME)
	await fsPromises.copyFile(src, dst)
	await fsPromises.chmod(dst, 0o755)
}

main()
