package create_retro_app

import "github.com/zaydek/retro/pkg/terminal"

var usage = `
` + terminal.Bold("create-retro-app [app]") + `

	Create a new app at directory app

		--template=...  'javascript' or 'typescript' (default 'javascript')

` + terminal.Bold("Repositories") + `

	` + terminal.Underline("https://github.com/zaydek/retro") + `
	` + terminal.Underline("https://github.com/evanw/esbuild") + `
`

var successFormat = terminal.Cyan("Success!") + `

npm:

	1. npm
	2. npm run start

yarn:

	1. yarn
	2. yarn start

Happy hacking!`

var successDirectoryFormat = terminal.Cyan("Success!") + `

npm:

	1. cd %[1]s
	2. npm
	3. npm run start

yarn:

	1. cd %[1]s
	2. yarn
	3. yarn start

Happy hacking!`
