package retro

// 	// Read www/index.html
// 	index_html := filepath.Join(runtime.Dirs.WwwDir, "index.html")
// 	if _, err := os.Stat(index_html); os.IsNotExist(err) {
// 		if err := os.MkdirAll(filepath.Dir(index_html), MODE_DIR); err != nil {
// 			return Runtime{}, err
// 		}
// 		err := ioutil.WriteFile(index_html,
// 			[]byte(`<!DOCTYPE html>
// <html lang="en">
// 	<head>
// 		<meta charset="UTF-8" />
// 		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
// 		<title>Hello, world!</title>
// 	</head>
// 	<body></body>
// </html>
// `), MODE_FILE)
// 		if err != nil {
// 			return Runtime{}, err
// 		}
// 	}
// 	tmpl, err := ioutil.ReadFile(index_html)
// 	if err != nil {
// 		return Runtime{}, err
// 	}
//
// 	// Check %head%
// 	if !bytes.Contains(tmpl, []byte("%head%")) {
// 		return Runtime{}, newTemplateError(index_html + `: Add '%head%' somewhere to '<head>'.
//
// For example:
//
// ` + terminal.Dimf(`// %s`, index_html) + `
// <!DOCTYPE html>
// 	<head lang="en">
// 		<meta charset="utf-8" />
// 		<meta name="viewport" content="width=device-width, initial-scale=1" />
// 		` + terminal.Magenta("%head%") + `
// 		` + terminal.Dim("...") + `
// 	</head>
// 	<body>
// 		` + terminal.Dim("...") + `
// 	</body>
// </html>
// `)
// 	}
//
// 	// Check %body%
// 	if !bytes.Contains(tmpl, []byte("%body%")) {
// 		return Runtime{}, newTemplateError(index_html + `: Add '%body%' somewhere to '<body>'.
//
// For example:
//
// ` + terminal.Dimf(`// %s`, index_html) + `
// <!DOCTYPE html>
// 	<head lang="en">
// 		<meta charset="utf-8" />
// 		<meta name="viewport" content="width=device-width, initial-scale=1" />
// 		` + terminal.Dim("...") + `
// 	</head>
// 	<body>
// 		` + terminal.Magenta("%body%") + `
// 		` + terminal.Dim("...") + `
// 	</body>
// </html>
// `)
// 	}
//
// 	runtime.Template = string(tmpl)
//
// 	// Remove __cache__, __export__
// 	rmdirs := []string{runtime.Dirs.CacheDir, runtime.Dirs.ExportDir}
// 	for _, rmdir := range rmdirs {
// 		if err := os.RemoveAll(rmdir); err != nil {
// 			return Runtime{}, err
// 		}
// 	}
//
// 	// Create www, src/pages, __cache__, __export__
// 	mkdirs := []string{runtime.Dirs.WwwDir, runtime.Dirs.SrcPagesDir, runtime.Dirs.CacheDir, runtime.Dirs.ExportDir}
// 	for _, mkdir := range mkdirs {
// 		if err := os.MkdirAll(mkdir, MODE_DIR); err != nil {
// 			return Runtime{}, err
// 		}
// 	}
//
// 	// Copy www to __export__
// 	excludes := []string{index_html}
// 	if err := copyDir(runtime.Dirs.WwwDir, runtime.Dirs.ExportDir, excludes); err != nil {
// 		return Runtime{}, err
// 	}
