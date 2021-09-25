package main

import (
	_ "embed"
	"os"
	"strings"

	"github.com/zaydek/retro/cmd/create_retro_app"
)

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

//go:embed version.txt
var RETRO_VERSION string

func init() {
	must(os.Setenv("RETRO_VERSION", strings.TrimSpace(RETRO_VERSION)))
}

func main() {
	create_retro_app.Run()
}
