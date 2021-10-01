package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/zaydek/retro/cmd/retro"
)

//go:embed version.txt
var RETRO_VERSION string

func main() {
	if err := os.Setenv("RETRO_VERSION", strings.TrimSpace(RETRO_VERSION)); err != nil {
		panic(fmt.Errorf("os.Setenv: %w", err))
	}
	retro.Run()
}
