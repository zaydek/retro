package main

import (
	_ "embed"
	"os"
	"strings"

	"github.com/zaydek/retro/go/cmd/create_retro_app"
)

//go:embed version.txt
var RETRO_VERSION string

func main() {
	if err := os.Setenv("RETRO_VERSION", strings.TrimSpace(RETRO_VERSION)); err != nil {
		panic(err)
	}
	create_retro_app.Run()
}
