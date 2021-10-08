// const sass = require("sass")
//
// try {
// 	sass.renderSync({
// 		file: "idea.scss",
// 		// callback(exception, result) {
// 		// 	// console.log("exception", exception)
// 		// 	if (exception?.message) {
// 		// 		console.error(exception.message)
// 		// 	}
// 		// 	console.log("result", result)
// 		// },
// 	})
// } catch (err) {
// 	console.log(err.formatted)
// }

const child_process = require("child_process")
const sass = require("sass")

let contents = ""
try {
	const result = child_process.execSync(`npx sass idea.scss --quiet`)
	contents = result.toString()
} catch { }

console.log(contents)
