import * as fs from "fs"
import * as os from "os"
import * as path from "path"

const CANONICAL_NAME = Object.keys(require("./package.json").bin)[0]!

const supported: { [key: string]: string } = {
	"darwin arm64 LE": "darwin-64",
	"darwin x64 LE": "darwin-64",
	"linux x64 LE": "linux-64",
	"win32 x64 LE": "windows-64.exe",
}

async function main(): Promise<void> {
	const name = supported[`${process.platform} ${os.arch()} ${os.endianness()}`]!
	const p1 = path.join(__dirname, "bin", name)
	const p2 = path.join(__dirname, "bin", CANONICAL_NAME)
	await fs.promises.copyFile(p1, p2)
	await fs.promises.chmod(p2, 0o755)
}

main()
