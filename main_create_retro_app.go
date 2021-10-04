package main

import (
	_ "embed"
	"encoding/json"
	"os"
	"strings"

	"github.com/zaydek/retro/go/cmd/create_retro_app"
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

	bstr, err := os.ReadFile("package.json")
	must(err)
	err = json.Unmarshal(bstr, &deps)
	must(err)

	err = os.Setenv("ESBUILD_VERSION", deps.DevDependencies.Esbuild)
	must(err)
	err = os.Setenv("REACT_VERSION", deps.DevDependencies.React)
	must(err)
	err = os.Setenv("REACTDOM_VERSION", deps.DevDependencies.ReactDOM)
	must(err)
	err = os.Setenv("RETRO_VERSION", RETRO_VERSION)
	must(err)

	create_retro_app.Run()
}
