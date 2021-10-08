package retro

import "github.com/zaydek/retro/go/pkg/terminal"

var usage = `
 ` + terminal.Bold("retro dev") + `

   Start the development server

     --port=...  Use port number (default ` + terminal.Cyan("8000") + `)

 ` + terminal.Bold("retro build") + `

   Build the production-ready build

 ` + terminal.Bold("retro serve") + `

   Serve the production-ready build

     --port=...  Use port number (default ` + terminal.Cyan("8000") + `)

 ` + terminal.Bold("Repositories") + `

   ` + terminal.Underline("https://github.com/zaydek/retro") + `
   ` + terminal.Underline("https://github.com/evanw/esbuild") + `
 `
