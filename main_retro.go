package main

import (
	_ "embed"
	"strings"

	"os"

	"github.com/zaydek/retro/cmd/retro"
)

//go:embed version.txt
var RETRO_VERSION string

func init() {
	os.Setenv("RETRO_VERSION", strings.TrimSpace(RETRO_VERSION))
}

func main() {
	retro.Run()
}
