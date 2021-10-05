package main

import (
	_ "embed"
	"encoding/json"
	"os"
	"strings"

	"github.com/zaydek/retro/go/cmd/retro"
)

var (
	//go:embed package.json
	pkg string

	//go:embed version.txt
	RETRO_V_VERSION string // v0.0.0

	//go:embed version.txt
	RETRO_VERSION string // ^0.0.0
)

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func init() {
	RETRO_V_VERSION = strings.TrimRight(RETRO_VERSION, "\n")
	RETRO_VERSION = strings.Replace(strings.TrimRight(RETRO_VERSION, "\n"), "v", "^", 1)
}

func main() {
	var deps struct {
		DevDependencies struct {
			Esbuild  string `json:"esbuild"`
			React    string `json:"react"`
			ReactDOM string `json:"react-dom"`
			Retro    string
		} `json:"devDependencies"`
	}

	err := json.Unmarshal([]byte(pkg), &deps)
	must(err)

	err = os.Setenv("ESBUILD_VERSION", deps.DevDependencies.Esbuild)
	must(err)
	err = os.Setenv("REACT_VERSION", deps.DevDependencies.React)
	must(err)
	err = os.Setenv("REACTDOM_VERSION", deps.DevDependencies.ReactDOM)
	must(err)

	err = os.Setenv("RETRO_V_VERSION", RETRO_V_VERSION)
	must(err)
	err = os.Setenv("RETRO_VERSION", RETRO_VERSION)
	must(err)

	retro.Run()
}
