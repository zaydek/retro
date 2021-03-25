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
	path := filepath.Join(WWW_DIR, "index.html")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), MODE_DIR); err != nil {
			return err
		}
		err := ioutil.WriteFile(path,
			[]byte(`<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Hello, world!</title>
		<link rel="stylesheet" href="/bundle.css" />
	</head>
	<body>
		<div id="root"></div>
		<script src="/vendor.js"></script>
		script src="/bundle.js"></script>
	</body>
</html>
`), MODE_FILE)
		if err != nil {
			return err
		}
	}
	html, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	//////////////////////////////////////////////////////////////////////////////

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

	//////////////////////////////////////////////////////////////////////////////

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

	//////////////////////////////////////////////////////////////////////////////

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

	//////////////////////////////////////////////////////////////////////////////

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

	return nil
}
