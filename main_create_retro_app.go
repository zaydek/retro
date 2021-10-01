package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/zaydek/retro/cmd/create_retro_app"
)

//go:embed version.txt
var RETRO_VERSION string

func main() {
	if err := os.Setenv("RETRO_VERSION", strings.TrimSpace(RETRO_VERSION)); err != nil {
		panic(fmt.Errorf("os.Setenv: %w", err))
	}
	create_retro_app.Run()
}
