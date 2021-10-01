package retro

import (
	"path/filepath"
	"strings"
)

func backtick(str string) string {
	return "`" + str + "`"
}

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
