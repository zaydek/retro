package retro

import (
	"os"
	"path/filepath"
	"strings"
)

func backtick(str string) string {
	return "`" + str + "`"
}

////////////////////////////////////////////////////////////////////////////////

// Removes `index.html` and or `.html`
func getCanonicalBrowserPath(url string) string {
	ret := url
	if strings.HasSuffix(url, "/index.html") {
		ret = ret[:len(ret)-len("index.html")]
	} else if strings.HasSuffix(url, "/index") {
		ret = ret[:len(ret)-len("index")]
	} else if strings.HasSuffix(url, ".html") {
		ret = ret[:len(ret)-len(".html")]
	}
	return ret
}

// Adds `index.html` and or `.html`
func getFilesystemPath(url string) string {
	ret := url
	if strings.HasSuffix(url, "/") {
		ret += "index.html"
	} else if strings.HasSuffix(url, "/index") {
		ret += ".html"
	} else if ext := filepath.Ext(url); ext == "" {
		ret += ".html"
	}
	return ret
}

////////////////////////////////////////////////////////////////////////////////

// Like `__dirname` in Node.js; gets the executable's directory
func getDirname() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	// Edge-case for local development: When running `go run main-*.go`, get the
	// current directory. `main-*` works as heuristic because the entry point
	// filenames are `main_create_retro_app.go` and `main_retro.go`.
	if strings.HasPrefix(filepath.Base(executable), "main_") {
		return os.Getwd()
	}
	// Follows symlinks; get `node_modules/.bin/@zaydek/bin/retro` not
	// `node_modules/.bin/retro`
	return filepath.EvalSymlinks(executable)
}
