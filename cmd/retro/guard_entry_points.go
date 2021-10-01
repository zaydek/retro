package retro

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/zaydek/retro/cmd/format"
	"github.com/zaydek/retro/pkg/terminal"
)

// type HTMLError struct {
// 	err error
// }
//
// func newHTMLError(str string) HTMLError {
// 	return HTMLError{err: errors.New(str)}
// }
//
// func (t HTMLError) Error() string {
// 	return t.err.Error()
// }

const (
	indexHTML = `<!DOCTYPE html>
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
    <script src="/bundle.js"></script>
  </body>
</html>` + "\n"

	indexJS = `import App from "./App"

import "./index.css"

if (document.getElementById("root").hasChildNodes()) {
  ReactDOM.hydrate(
    <React.StrictMode>
      <App />
    </React.StrictMode>,
    document.getElementById("root"),
  )
} else {
  ReactDOM.render(
    <React.StrictMode>
      <App />
    </React.StrictMode>,
    document.getElementById("root"),
  )
}` + "\n"

	appJS = `import "./App.css"

export default function App() {
  return (
    <div className="App">
      <h1>Hello, world!</h1>
    </div>
  )
}` + "\n"
)

////////////////////////////////////////////////////////////////////////////////

func copyDefaultIndexHTMLEntryPoint() error {
	filename := filepath.Join(RETRO_WWW_DIR, "index.html")
	if err := os.MkdirAll(filepath.Dir(filename), permBitsDirectory); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}
	if err := ioutil.WriteFile(filename, []byte(indexHTML), permBitsFile); err != nil {
		return fmt.Errorf("ioutil.WriteFile: %w", err)
	}
	return nil
}

func copyDefaultIndexJSEntryPoint() error {
	filename := filepath.Join(RETRO_SRC_DIR, "index.js")
	if err := os.MkdirAll(filepath.Dir(filename), permBitsDirectory); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}
	if err := ioutil.WriteFile(filename, []byte(indexJS), permBitsFile); err != nil {
		return fmt.Errorf("ioutil.WriteFile: %w", err)
	}
	return nil
}

func copyDefaultAppJSEntryPoint() error {
	filename := filepath.Join(RETRO_SRC_DIR, "app.js")
	if err := os.MkdirAll(filepath.Dir(filename), permBitsDirectory); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}
	if err := ioutil.WriteFile(filename, []byte(appJS), permBitsFile); err != nil {
		return fmt.Errorf("ioutil.WriteFile: %w", err)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func guardIndexHTMLEntryPoint() error {
	// Guard `index.html`
	filename := filepath.Join(RETRO_WWW_DIR, "index.html")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := copyDefaultIndexHTMLEntryPoint(); err != nil {
			return fmt.Errorf("copyDefaultIndexHTMLEntryPoint: %w", err)
		}
	}

	// Read contents of `index.html`
	byteStr, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	contents := string(byteStr)

	// <link rel="stylesheet" href="/bundle.css" />
	if !strings.Contains(contents, `<link rel="stylesheet" href="/bundle.css" />`) {
		fmt.Fprintln(
			os.Stderr,
			format.Error(
				fmt.Sprintf("Add %s somewhere to %s", `Add `+terminal.Magenta(backtick(`<link rel="stylesheet" href="/bundle.css" />`)), terminal.Magenta(backtick(`<head>`)))+`.

For example:

`+terminal.Dimf("// %s", filename)+`
<!DOCTYPE html>
  <head lang="en">
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    `+terminal.Green(`<link rel="stylesheet" href="/bundle.css" />`)+`
    `+terminal.Dim("...")+`
  </head>
  <body>
    `+terminal.Dim("...")+`
  </body>
</html>`,
			),
		)
		os.Exit(1)
	}

	// <div id="root"></div>
	if !strings.Contains(contents, `<div id="root"></div>`) {
		fmt.Fprintln(
			os.Stderr,
			format.Error(
				fmt.Sprintf("Add %s somewhere to %s", `Add `+terminal.Magenta(backtick(`<div id="root"></div>`)), terminal.Magenta(backtick(`<body>`)))+`.

For example:

`+terminal.Dimf("// %s", filename)+`
<!DOCTYPE html>
  <head lang="en">
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    `+terminal.Dim("...")+`
  </head>
  <body>
    `+terminal.Green(`<div id="root"></div>`)+`
    `+terminal.Dim("...")+`
  </body>
</html>`,
			),
		)
		os.Exit(1)
	}

	// <script src="/vendor.js"></script>
	if !strings.Contains(contents, `<script src="/vendor.js"></script>`) {
		fmt.Fprintln(
			os.Stderr,
			format.Error(
				fmt.Sprintf("Add %s somewhere to %s", `Add `+terminal.Magenta(backtick(`<script src="/vendor.js"></script>`)), terminal.Magenta(backtick(`<body>`)))+`.

For example:

`+terminal.Dimf("// %s", filename)+`
<!DOCTYPE html>
  <head lang="en">
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    `+terminal.Dim("...")+`
  </head>
  <body>
    <div id="root"></div>
    `+terminal.Green(`<script src="/vendor.js"></script>`)+`
    `+terminal.Dim("...")+`
  </body>
</html>`,
			),
		)
		os.Exit(1)
	}

	// <script src="/client.js"></script>
	if !strings.Contains(contents, `<script src="/client.js"></script>`) {
		fmt.Fprintln(
			os.Stderr,
			format.Error(
				fmt.Sprintf("Add %s somewhere to %s", `Add `+terminal.Magenta(backtick(`<script src="/client.js"></script>`)), terminal.Magenta(backtick(`<body>`)))+`.

For example:

`+terminal.Dimf("// %s", filename)+`
<!DOCTYPE html>
  <head lang="en">
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    `+terminal.Dim("...")+`
  </head>
  <body>
    <div id="root"></div>
    <script src="/vendor.js"></script>
    `+terminal.Green(`<script src="/client.js"></script>`)+`
    `+terminal.Dim("...")+`
  </body>
</html>`,
			),
		)
		os.Exit(1)
	}

	return nil
}

func guardIndexJSEntryPoint() error {
	filename := filepath.Join(RETRO_WWW_DIR, "index.js")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := copyDefaultIndexJSEntryPoint(); err != nil {
			return fmt.Errorf("copyDefaultIndexJSEntryPoint: %w", err)
		}
	}
	return nil
}

func guardAppJSEntryPoint() error {
	filename := filepath.Join(RETRO_WWW_DIR, "app.js")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := copyDefaultAppJSEntryPoint(); err != nil {
			return fmt.Errorf("copyDefaultAppJSEntryPoint: %w", err)
		}
	}
	return nil
}

// Guards entry points:
//
// - www/index.html
// - src/index.{jsx?|tsx?}
// - src/app.{jsx?|tsx?}
//
func guardEntryPoints() error {
	if err := guardIndexHTMLEntryPoint(); err != nil {
		return fmt.Errorf("guardIndexHTMLEntryPoint: %w")
	}
	if err := guardIndexJSEntryPoint(); err != nil {
		return fmt.Errorf("guardIndexJSEntryPoint: %w")
	}
	if err := guardAppJSEntryPoint(); err != nil {
		return fmt.Errorf("guardAppJSEntryPoint: %w")
	}
	return nil
}

// func copyHTMLEntryPoint(vendorDotJS, bundleDotJS, bundleDotCSS string) error {
// 	bstr, err := os.ReadFile(filepath.Join(RETRO_WWW_DIR, "index.html"))
// 	if err != nil {
// 		return err
// 	}
//
// 	// Swap cache busted paths
// 	contents := string(bstr)
// 	contents = strings.Replace(
// 		contents,
// 		`<script src="/vendor.js"></script>`,
// 		fmt.Sprintf(`<script src="/%s"></script>`,
// 			vendorDotJS,
// 		),
// 		1,
// 	)
//
// 	contents = strings.Replace(
// 		contents,
// 		`<script src="/bundle.js"></script>`,
// 		fmt.Sprintf(`<script src="/%s"></script>`,
// 			bundleDotJS,
// 		),
// 		1,
// 	)
//
// 	contents = strings.Replace(
// 		contents,
// 		`<link rel="stylesheet" href="/bundle.css" />`,
// 		fmt.Sprintf(`<link rel="stylesheet" href="/%s" />`,
// 			bundleDotCSS,
// 		),
// 		1,
// 	)
//
// 	if err := ioutil.WriteFile(filepath.Join(RETRO_OUT_DIR, "index.html"), []byte(contents), permBitsFile); err != nil {
// 		return err
// 	}
// 	return nil
// }
