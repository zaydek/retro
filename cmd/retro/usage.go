package retro

import "github.com/zaydek/retro/pkg/terminal"

var usage = `
` + terminal.Bold("retro dev") + `

	Start the dev server

		--port=...       Use port (default '8000')
		--sourcemap=...  Add source maps (default 'true')

` + terminal.Bold("retro build") + `

	Build the production-ready build

		--sourcemap=...  Add source maps (default 'true')

` + terminal.Bold("retro serve") + `

	Serve the production-ready build

		--port=...       Use port (default '8000')

` + terminal.Bold("Repository") + `

	` + terminal.Underline("https://github.com/zaydek/retro") + `
	` + terminal.Underline("https://github.com/evanw/esbuild") + `
`
