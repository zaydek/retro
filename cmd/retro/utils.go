package retro

import (
	"path/filepath"
	"strings"
)

func getBrowserPath(url string) string {
	out := url
	if strings.HasSuffix(url, "/index.html") {
		out = out[:len(out)-len("index.html")] // Keep "/"
	} else if strings.HasSuffix(url, "/index") {
		out = out[:len(out)-len("index")] // Keep "/"
	} else if strings.HasSuffix(url, ".html") {
		out = out[:len(out)-len(".html")]
	}
	return out
}

func getFSPath(url string) string {
	out := url
	if strings.HasSuffix(url, "/") {
		out += "index.html"
	} else if strings.HasSuffix(url, "/index") {
		out += ".html"
	} else if ext := filepath.Ext(url); ext == "" {
		out += ".html"
	}
	return out
}
