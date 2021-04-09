package create_retro_app

import "github.com/zaydek/retro/pkg/terminal"

var usage = `
` + terminal.Bold("create-retro-app [app]") + `

	Create a new app at directory app

		--template=...  'starter', 'sass', or 'mdx' (default 'starter')

` + terminal.Bold("Repositories") + `

	` + terminal.Underline("https://github.com/zaydek/retro") + `
	` + terminal.Underline("https://github.com/evanw/esbuild") + `
`

var successFmt = terminal.Cyan("Success!") + `

  npm:

	  1. npm
	  2. npm run dev

  yarn:

	  1. yarn
	  2. yarn dev

Happy hacking!`

var successDirFmt = terminal.Cyan("Success!") + `

  npm:

	  1. cd %[1]s
	  2. npm i
	  3. npm run dev

  yarn:

	  1. cd %[1]s
	  2. yarn
	  3. yarn dev

Happy hacking!`
