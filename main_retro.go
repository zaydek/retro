package main

import (
	_ "embed"
	"encoding/json"
	"os"
	"strings"

	"github.com/zaydek/retro/go/cmd/retro"
)

//go:embed version.txt
var RETRO_VERSION string

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func init() {
	RETRO_VERSION = strings.TrimRight(RETRO_VERSION, "\n")
}

func main() {
	var deps struct {
		Esbuild  string
		React    string
		ReactDOM string
		Retro    string
	}

	bstr, err := os.ReadFile("package.json")
	must(err)
	err = json.Unmarshal(bstr, &deps)
	must(err)

	err = os.Setenv("ESBUILD_VERSION", deps.Esbuild)
	must(err)
	err = os.Setenv("REACT_VERSION", deps.React)
	must(err)
	err = os.Setenv("REACTDOM_VERSION", deps.ReactDOM)
	must(err)
	err = os.Setenv("RETRO_VERSION", RETRO_VERSION)
	must(err)

	retro.Run()
}
