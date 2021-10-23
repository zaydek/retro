package retro

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zaydek/retro/go/pkg/terminal"
)

type EntryPointError struct {
	err error
}

func newEntryPointError(str string) EntryPointError {
	return EntryPointError{err: errors.New(str)}
}

func (e EntryPointError) Error() string {
	return e.err.Error()
}

func copyDefaultIndexHTMLEntryPoint() error {
	filename := filepath.Join(RETRO_WWW_DIR, "index.html")
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filename, []byte(htmlEntryPoint+"\n"), 0644); err != nil {
		return err
	}
	return nil
}

func copyDefaultIndexJSEntryPoint() error {
	filename := filepath.Join(RETRO_SRC_DIR, "index.js")
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filename, []byte(jsEntryPoint+"\n"), 0644); err != nil {
		return err
	}
	return nil
}

func copyDefaultAppJSEntryPoint() error {
	filename := filepath.Join(RETRO_SRC_DIR, "app.js")
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filename, []byte(appJSEntryPoint+"\n"), 0644); err != nil {
		return err
	}
	return nil
}

// Guards for the presence of 'www/index.js' and:
//
// - <link rel="stylesheet" href="/client.css" />
// - <div id="root"></div>
// - <script src="/vendor.js"></script>
// - <script src="/client.js"></script>
//
func guardHTMLEntryPoint() error {
	// Guard 'www/index.html'
	filename := filepath.Join(RETRO_WWW_DIR, "index.html")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := copyDefaultIndexHTMLEntryPoint(); err != nil {
			return err
		}
	}

	// www/index.html
	bstr, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	contents := string(bstr)

	// <link rel="stylesheet" href="/client.css" />
	if !strings.Contains(contents, `<link rel="stylesheet" href="/client.css" />`) {
		return newEntryPointError(
			fmt.Sprintf("Add %s somewhere to %s", `Add `+terminal.Magenta(quote(`<link rel="stylesheet" href="/client.css" />`)), terminal.Magenta(quote(`<head>`))) + `.

For example:

` + terminal.Dimf("// %s", filename) + `
<!DOCTYPE html>
	<head lang="en">
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		` + terminal.Green(`<link rel="stylesheet" href="/client.css" />`) + `
		` + terminal.Dim("...") + `
	</head>
	<body>
		` + terminal.Dim("...") + `
	</body>
</html>`,
		)
	}

	// <div id="root"></div>
	if !strings.Contains(contents, `<div id="root"></div>`) {
		return newEntryPointError(
			fmt.Sprintf("Add %s somewhere to %s", `Add `+terminal.Magenta(quote(`<div id="root"></div>`)), terminal.Magenta(quote(`<body>`))) + `.

For example:

` + terminal.Dimf("// %s", filename) + `
<!DOCTYPE html>
	<head lang="en">
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		` + terminal.Dim("...") + `
	</head>
	<body>
		` + terminal.Green(`<div id="root"></div>`) + `
		` + terminal.Dim("...") + `
	</body>
</html>`,
		)
	}

	// <script src="/vendor.js"></script>
	if !strings.Contains(contents, `<script src="/vendor.js"></script>`) {
		return newEntryPointError(
			fmt.Sprintf("Add %s somewhere to %s", `Add `+terminal.Magenta(quote(`<script src="/vendor.js"></script>`)), terminal.Magenta(quote(`<body>`))) + `.

For example:

` + terminal.Dimf("// %s", filename) + `
<!DOCTYPE html>
	<head lang="en">
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		` + terminal.Dim("...") + `
	</head>
	<body>
		<div id="root"></div>
		` + terminal.Green(`<script src="/vendor.js"></script>`) + `
		` + terminal.Dim("...") + `
	</body>
</html>`,
		)
	}

	// <script src="/client.js"></script>
	if !strings.Contains(contents, `<script src="/client.js"></script>`) {
		return newEntryPointError(
			fmt.Sprintf("Add %s somewhere to %s", `Add `+terminal.Magenta(quote(`<script src="/client.js"></script>`)), terminal.Magenta(quote(`<body>`))) + `.

For example:

` + terminal.Dimf("// %s", filename) + `
<!DOCTYPE html>
	<head lang="en">
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		` + terminal.Dim("...") + `
	</head>
	<body>
		<div id="root"></div>
		<script src="/vendor.js"></script>
		` + terminal.Green(`<script src="/client.js"></script>`) + `
		` + terminal.Dim("...") + `
	</body>
</html>`,
		)
	}

	return nil
}

// Guards for the presence of 'src/index.js'
func guardJSEntryPoint() error {
	filename := filepath.Join(RETRO_SRC_DIR, "index.js")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := copyDefaultIndexJSEntryPoint(); err != nil {
			return err
		}
	}
	return nil
}

// Guards for the presence of 'src/App.js'
func guardAppJSEntryPoint() error {
	filename := filepath.Join(RETRO_SRC_DIR, "app.js")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := copyDefaultAppJSEntryPoint(); err != nil {
			return err
		}
	}
	return nil
}

// Guards entry points:
//
// - www/index.html
// - src/index.js
// - src/App.js
//
func guardEntryPoints() error {
	if err := guardHTMLEntryPoint(); err != nil {
		return err
	}
	if err := guardJSEntryPoint(); err != nil {
		return err
	}
	if err := guardAppJSEntryPoint(); err != nil {
		return err
	}
	return nil
}

type entryPoints struct {
	clientCSS string // The bundled CSS filename
	vendorJS  string // The bundled vendor JS filename
	clientJS  string // The bundled client JS filename
}

func copyIndexHTMLEntryPoint(entries entryPoints) error {
	// www/index.html
	srcPath := filepath.Join(RETRO_WWW_DIR, "index.html")
	bstr, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	// <link rel="stylesheet" href="/client.css" />
	contents := string(bstr)
	contents = strings.Replace(
		contents,
		`<link rel="stylesheet" href="/client.css" />`,
		fmt.Sprintf(`<link rel="stylesheet" href="/%s" />`, entries.clientCSS),
		1,
	)
	// <script src="/vendor.js"></script>
	contents = strings.Replace(
		contents,
		`<script src="/vendor.js"></script>`,
		fmt.Sprintf(`<script src="/%s"></script>`, entries.vendorJS),
		1,
	)
	// <script src="/client.js"></script>
	contents = strings.Replace(
		contents,
		`<script src="/client.js"></script>`,
		fmt.Sprintf(`<script src="/%s"></script>`, entries.clientJS),
		1,
	)
	// out/index.html
	dstPath := filepath.Join(RETRO_OUT_DIR, "index.html")
	if err := os.WriteFile(dstPath, []byte(contents), 0644); err != nil {
		return err
	}
	return nil
}
