package create_retro_app

import "github.com/zaydek/retro/go/pkg/terminal"

var usage = `
 ` + terminal.Bold("create-retro-app [dir]") + `

   Create a Retro app at directory ` + terminal.Bold("[dir]") + `

 ` + terminal.Bold("Repositories") + `

   ` + terminal.Underline("https://github.com/zaydek/retro") + `
   ` + terminal.Underline("https://github.com/evanw/esbuild") + `
 `

var successStr = terminal.Cyan("Success!") + `

   npm:

     1. npm
     2. npm run dev

   yarn:

     1. yarn
     2. yarn dev

 Happy hacking!`

var successDirStr = terminal.Cyan("Success!") + `

   npm:

     1. cd %[1]s
     2. npm i
     3. npm run dev

   yarn:

     1. cd %[1]s
     2. yarn
     3. yarn dev

 Happy hacking!`
