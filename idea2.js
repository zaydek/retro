// const child_process = require("child_process")
//
// child_process.exec("npx sass idea.scss --quiet", (error, stdout, stderr) => {
// 	if (error) {
// 		console.error(`exec error: ${error}`)
// 		return;
// 	}
// 	if (stdout !== "") { console.log(`stdout: ${stdout}`) }
// 	if (stderr !== "") { console.error(`stderr: ${stderr}`) }
// })

const child_process = require("child_process")

let out = ""
try {
	const result = child_process.execSync("npx sass idea.scss --quiet")
	out = result.toString()
} catch { }
console.log(out)
