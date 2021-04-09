package main

import (
	_ "embed"
	"strings"

	"os"

	"github.com/zaydek/retro/cmd/retro"
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
	retro.Run()
}
