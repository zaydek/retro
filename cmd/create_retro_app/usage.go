package create_retro_app

import "github.com/zaydek/retro/pkg/terminal"

var usage = `
` + terminal.Bold("create-retro-app dir") + `

	Create a new app at directory dir

		--template=...  'javascript' or 'typescript' (default 'javascript')

` + terminal.Bold("Repository") + `

	` + terminal.Underline("https://github.com/zaydek/retro") + `
	` + terminal.Underline("https://github.com/evanw/esbuild") + `
`

var successFormat = `Success!

npm:

	1. npm
	2. npm run start

yarn:

	1. yarn
	2. yarn start

Happy hacking!`

var successDirectoryFormat = `Success!

npm:

	1. cd %[1]s
	2. npm
	3. npm run start

yarn:

	1. cd %[1]s
	2. yarn
	3. yarn start

Happy hacking!`
