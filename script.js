const esbuild = require("esbuild")

async function main() {
	let formatted = []

	formatted.push(
		await esbuild.formatMessages(
			[
				{
					text: '"test" has already been declared',
					location: {
						file: "file.js",
						line: 2,
						column: 4,
						length: 4,
						lineText: 'let test = "second"',
					},
					// notes: [
					// 	{
					// 		text: '"test" was originally declared here',
					// 		location: {
					// 			file: "file.js",
					// 			line: 1,
					// 			column: 4,
					// 			length: 4,
					// 			lineText: 'let test = "first"',
					// 		},
					// 	},
					// ],
				},
			],
			{
				kind: "error",
				color: true,
				terminalWidth: 100,
			},
		),
	)

	// formatted.push(
	// 	await esbuild.formatMessages(
	// 		[
	// 			{
	// 				text: '"test" has already been declared',
	// 				location: {
	// 					file: "file.js",
	// 					line: 2,
	// 					column: 4,
	// 					length: 4,
	// 					lineText: 'let test = "second"',
	// 				},
	// 				notes: [
	// 					{
	// 						text: '"test" was originally declared here',
	// 						location: {
	// 							file: "file.js",
	// 							line: 1,
	// 							column: 4,
	// 							length: 4,
	// 							lineText: 'let test = "first"',
	// 						},
	// 					},
	// 				],
	// 			},
	// 		],
	// 		{
	// 			kind: "warning",
	// 			color: true,
	// 			terminalWidth: 100,
	// 		},
	// 	),
	// )

	console.log(formatted.join("").trimEnd())
}

main()
