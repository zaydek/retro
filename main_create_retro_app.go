package main

import (
	_ "embed"
	"strings"

	"os"

	"github.com/zaydek/retro/cmd/create_retro_app"
)

//go:embed version.txt
var RETRO_VERSION string

func init() {
	os.Setenv("RETRO_VERSION", strings.TrimSpace(RETRO_VERSION))
}

func main() {
	create_retro_app.Run()
}
