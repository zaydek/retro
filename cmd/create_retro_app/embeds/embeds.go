package embeds

import (
	"embed"
	"io/fs"
)

//go:embed javascript/*
var jsFS embed.FS

//go:embed typescript/*
var tsFS embed.FS

var (
	JavaScriptFS fs.FS
	TypeScriptFS fs.FS
)

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func init() {
	var err error
	JavaScriptFS, err = fs.Sub(jsFS, "javascript")
	must(err)
	TypeScriptFS, err = fs.Sub(tsFS, "typescript")
	must(err)
}
