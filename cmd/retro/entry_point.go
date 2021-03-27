package retro

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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

func guardHTMLEntryPoint() error {
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
		<link rel="stylesheet" href="/index.css" />
	</head>
	<body>
		<div id="root"></div>
		<script src="/react.js"></script>
		script src="/bundle.js"></script>
	</body>
</html>
` /* EOF */), MODE_FILE)
		if err != nil {
			return err
		}
	}

	bstr, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	//////////////////////////////////////////////////////////////////////////////

	contents := string(bstr)
	if !strings.Contains(contents, `<link rel="stylesheet" href="/index.css" />`) {
		return newHTMLError(fmt.Sprintf(`Add %s somewhere to %s.`, terminal.Magenta(`'<link rel="stylesheet" href="/index.css" />'`), terminal.Magenta("'<head>'")) + `

For example:

` + terminal.Dimf(`// %s`, path) + `
<!DOCTYPE html>
	<head lang="en">
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		` + terminal.Green(`<link rel="stylesheet" href="/index.css" />`) + `
		` + terminal.Dim("...") + `
	</head>
	<body>
		` + terminal.Dim("...") + `
	</body>
</html>`)
	}

	//////////////////////////////////////////////////////////////////////////////

	if !strings.Contains(contents, `<div id="root"></div>`) {
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
		` + terminal.Dim("...") + `
		` + terminal.Green(`<div id="root"></div>`) + `
	</body>
</html>`)
	}

	//////////////////////////////////////////////////////////////////////////////

	if !strings.Contains(contents, `<script src="/react.js"></script>`) {
		return newHTMLError(fmt.Sprintf(`Add %s somewhere to %s.`, terminal.Magenta(`'<script src="/react.js"></script>'`), terminal.Magenta("'<body>'")) + `

For example:

` + terminal.Dimf(`// %s`, path) + `
<!DOCTYPE html>
	<head lang="en">
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		` + terminal.Dim("...") + `
	</head>
	<body>
		` + terminal.Dim("...") + `
		` + terminal.Green(`<script src="/react.js"></script>`) + `
	</body>
</html>`)
	}

	//////////////////////////////////////////////////////////////////////////////

	if !strings.Contains(contents, `<script src="/index.js"></script>`) {
		return newHTMLError(fmt.Sprintf(`Add %s somewhere to %s.`, terminal.Magenta(`'<script src="/index.js"></script>'`), terminal.Magenta("'<body>'")) + `

For example:

` + terminal.Dimf(`// %s`, path) + `
<!DOCTYPE html>
	<head lang="en">
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		` + terminal.Dim("...") + `
	</head>
	<body>
		` + terminal.Dim("...") + `
		` + terminal.Green(`<script src="/index.js"></script>`) + `
	</body>
</html>`)
	}

	return nil
}

func copyHTMLEntryPoint(react_js, index_js, index_css string) error {
	bstr, err := ioutil.ReadFile(filepath.Join(WWW_DIR, "index.html"))
	if err != nil {
		return err
	}

	// Swap cache busted paths
	contents := string(bstr)
	contents = strings.Replace(
		contents,
		`<script src="/react.js"></script>`,
		fmt.Sprintf(`<script src="/%s"></script>`,
			react_js,
		),
		1,
	)

	contents = strings.Replace(
		contents,
		`<script src="/index.js"></script>`,
		fmt.Sprintf(`<script src="/%s"></script>`,
			index_js,
		),
		1,
	)

	contents = strings.Replace(
		contents,
		`<link rel="stylesheet" href="/index.css" />`,
		fmt.Sprintf(`<link rel="stylesheet" href="/%s" />`,
			index_css,
		),
		1,
	)

	// Add the dev stub
	if os.Getenv("ENV") == "development" {
		contents = strings.Replace(contents, "</html>", fmt.Sprintf("\t%s\n</html>", devStub), 1)
	}

	if err := ioutil.WriteFile(filepath.Join(OUT_DIR, "index.html"), []byte(contents), MODE_FILE); err != nil {
		return err
	}
	return nil
}
