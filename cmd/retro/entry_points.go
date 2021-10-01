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

// TODO: In theory we can also access default values from
// `create_retro_app/embeds`. However, this is more self-contained.
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

// TODO: In theory we can also access default values from
// `create_retro_app/embeds`. However, this is more self-contained.
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

// TODO: In theory we can also access default values from
// `create_retro_app/embeds`. However, this is more self-contained.
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

// Guards for the presence of `www/index.js` and:
//
// - <link rel="stylesheet" href="/client.css" />
// - <div id="retro_root"></div>
// - <script src="/vendor.js"></script>
// - <script src="/client.js"></script>
//
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

	// <link rel="stylesheet" href="/client.css" />
	if !strings.Contains(contents, `<link rel="stylesheet" href="/client.css" />`) {
		fmt.Fprintln(
			os.Stderr,
			format.Error(
				fmt.Sprintf("Add %s somewhere to %s", `Add `+terminal.Magenta(backtick(`<link rel="stylesheet" href="/client.css" />`)), terminal.Magenta(backtick(`<head>`)))+`.

For example:

`+terminal.Dimf("// %s", filename)+`
<!DOCTYPE html>
  <head lang="en">
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    `+terminal.Green(`<link rel="stylesheet" href="/client.css" />`)+`
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

	// <div id="retro_root"></div>
	if !strings.Contains(contents, `<div id="rootretro_"></div>`) {
		fmt.Fprintln(
			os.Stderr,
			format.Error(
				fmt.Sprintf("Add %s somewhere to %s", `Add `+terminal.Magenta(backtick(`<div id="rootretro_"></div>`)), terminal.Magenta(backtick(`<body>`)))+`.

For example:

`+terminal.Dimf("// %s", filename)+`
<!DOCTYPE html>
  <head lang="en">
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    `+terminal.Dim("...")+`
  </head>
  <body>
    `+terminal.Green(`<div id="rootretro_"></div>`)+`
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
    <div id="retro_root"></div>
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
    <div id="retro_root"></div>
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

// Guards for the presence of `src/index.js`
func guardIndexJSEntryPoint() error {
	filename := filepath.Join(RETRO_SRC_DIR, "index.js")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := copyDefaultIndexJSEntryPoint(); err != nil {
			return fmt.Errorf("copyDefaultIndexJSEntryPoint: %w", err)
		}
	}
	return nil
}

// Guards for the presence of `src/App.js`
func guardAppJSEntryPoint() error {
	filename := filepath.Join(RETRO_SRC_DIR, "app.js")
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
// - src/index.js
// - src/App.js
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

func transformAndCopyIndexHTMLEntryPoint(vendorJSFilename, clientJSFilename, clientCSSFilename string) error {
	filename := filepath.Join(RETRO_WWW_DIR, "index.html")
	byteStr, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}
	// <link rel="stylesheet" href="/client.css" />
	contents := string(byteStr)
	contents = strings.Replace(
		contents,
		`<link rel="stylesheet" href="/client.css" />`,
		fmt.Sprintf(`<link rel="stylesheet" href="/%s" />`, clientCSSFilename),
		1,
	)
	// <script src="/vendor.js"></script>
	contents = strings.Replace(
		contents,
		`<script src="/vendor.js"></script>`,
		fmt.Sprintf(`<script src="/%s"></script>`, vendorJSFilename),
		1,
	)
	// <script src="/client.js"></script>
	contents = strings.Replace(
		contents,
		`<script src="/client.js"></script>`,
		fmt.Sprintf(`<script src="/%s"></script>`, clientJSFilename),
		1,
	)
	// Copy the transformed `www/index.html` to `out/www/index.html`
	target := filepath.Join(RETRO_OUT_DIR, "index.html")
	if err := os.MkdirAll(filepath.Dir(target), permBitsDirectory); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}
	if err := ioutil.WriteFile(target, []byte(contents), permBitsFile); err != nil {
		return fmt.Errorf("ioutil.WriteFile: %w", err)
	}
	return nil
}
