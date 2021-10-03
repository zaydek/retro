package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var versions = struct {
	Retro string

	Esbuild string
	React   string
	Sass    string
}{
	Retro: strings.Replace(os.Getenv("RETRO_VERSION"), "v", "^", 1),

	Esbuild: "^0.13.3",
	React:   "^17.0.2",
	Sass:    "^1.32.8",
}

func main() {
	bstr, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bstr))
}
