package main

import (
	_ "embed"

	"os"

	"github.com/zaydek/retro/cmd/create_retro_app"
)

//go:embed version.txt
var RETRO_VERSION string

func main() {
	os.Setenv("RETRO_VERSION", RETRO_VERSION)
	create_retro_app.Run()
}
