package retro

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/zaydek/retro/pkg/terminal"
)

type HTMLError struct {
	err error
}

func newHTMLError(str string) HTMLError {
	return HTMLError{err: errors.New(str)}
}

func (t HTMLError) Error() string {
	return t.err.Error()
}

func guards() error {
	// Read 'www/index.html'
	path := filepath.Join(WWW_DIR, "index.html")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), MODE_DIR); err != nil {
			return err
		}
		err := ioutil.WriteFile(path,
			[]byte(`<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Hello, world!</title>
	</head>
	<body></body>
</html>`), MODE_FILE)
		if err != nil {
			return err
		}
	}
	html, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Guard '<link rel="stylesheet" href="/bundle.css" />'
	if !bytes.Contains(html, []byte(`<link rel="stylesheet" href="/bundle.css" />`)) {
		return newHTMLError(fmt.Sprintf(`Add %s somewhere to %s.`, terminal.Magenta(`'<link rel="stylesheet" href="/bundle.css" />'`), terminal.Magenta("'<head>'")) + `

For example:

` + terminal.Dimf(`// %s`, path) + `
<!DOCTYPE html>
	<head lang="en">
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		` + terminal.Magenta(`<link rel="stylesheet" href="/bundle.css" />`) + `
		` + terminal.Dim("...") + `
	</head>
	<body>
		` + terminal.Dim("...") + `
	</body>
</html>`)
	}

	// Guard '<div id="root"></div>'
	if !bytes.Contains(html, []byte(`<div id="root"></div>`)) {
		return newHTMLError(fmt.Sprintf(`Add %s somewhere to %s.`, terminal.Magenta(`'<div id="root"></div>'`), terminal.Magenta("'<body>'")) + `

For example:

` + terminal.Dimf(`// %s`, path) + `
<!DOCTYPE html>
	<head lang="en">
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		` + terminal.Dim("...") + `
	</head>
	<body>
		` + terminal.Magenta(`<div id="root"></div>`) + `
		` + terminal.Dim("...") + `
	</body>
</html>`)
	}

	// Guard '<script src="/vendor.js"></script>'
	if !bytes.Contains(html, []byte(`<script src="/vendor.js"></script>`)) {
		return newHTMLError(fmt.Sprintf(`Add %s somewhere to %s.`, terminal.Magenta(`'<script src="/vendor.js"></script>'`), terminal.Magenta("'<body>'")) + `

For example:

` + terminal.Dimf(`// %s`, path) + `
<!DOCTYPE html>
	<head lang="en">
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		` + terminal.Dim("...") + `
	</head>
	<body>
		` + terminal.Magenta(`<script src="/vendor.js"></script>`) + `
		` + terminal.Dim("...") + `
	</body>
</html>`)
	}

	// Guard '<script src="/bundle.js"></script>'
	if !bytes.Contains(html, []byte(`<script src="/bundle.js"></script>`)) {
		return newHTMLError(fmt.Sprintf(`Add %s somewhere to %s.`, terminal.Magenta(`'<script src="/bundle.js"></script>'`), terminal.Magenta("'<body>'")) + `

For example:

` + terminal.Dimf(`// %s`, path) + `
<!DOCTYPE html>
	<head lang="en">
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		` + terminal.Dim("...") + `
	</head>
	<body>
		` + terminal.Magenta(`<script src="/bundle.js"></script>`) + `
		` + terminal.Dim("...") + `
	</body>
</html>`)
	}

	// Remove 'out'
	rmdirs := []string{OUT_DIR}
	for _, rmdir := range rmdirs {
		if err := os.RemoveAll(rmdir); err != nil {
			return err
		}
	}

	// Create 'www', 'src/pages', 'out'
	mkdirs := []string{WWW_DIR, SRC_DIR, OUT_DIR}
	for _, mkdir := range mkdirs {
		if err := os.MkdirAll(mkdir, MODE_DIR); err != nil {
			return err
		}
	}

	// Copy 'www' to 'out'
	excludes := []string{path}
	if err := cpdir(WWW_DIR, OUT_DIR, excludes); err != nil {
		return err
	}

	return nil
}
