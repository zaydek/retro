import * as fs from "fs"
import * as os from "os"
import * as path from "path"

const CANON_BIN_NAME = Object.keys(require("./package.json").bin)[0]!

const supported: Record<string, string> = {
	"darwin arm64 LE": "darwin-64",
	"darwin x64 LE": "darwin-64",
	"linux x64 LE": "linux-64",
	"win32 x64 LE": "windows-64.exe",
}

function main(): void {
	const name = supported[`${process.platform} ${os.arch()} ${os.endianness()}`]!

	// https://github.com/evanw/esbuild/blob/9522ed1a4b4e2c9cabf427592c8a2ecaeecbcb74/npm/esbuild-windows-64/bin/esbuild
	if (process.platform === "win32") {
		const target = path.join(__dirname, "bin", name)
		fs.writeFileSync(
			target,
			`#!/usr/bin/env node

const exe = require.resolve("./windows-64.exe")
const child_process = require("child_process")
child_process.spawnSync(exe, process.argv.slice(2), { stdio: "inherit" })
`,
		)
		fs.chmodSync(target, 0o755)
		return
	}

	const tar1 = path.join(__dirname, "bin", name)
	const tar2 = path.join(__dirname, "bin", CANON_BIN_NAME)
	fs.copyFileSync(tar1, tar2)
	fs.chmodSync(tar2, 0o755)
}

main()
