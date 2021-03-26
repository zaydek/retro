package main

import (
	_ "embed"

	"fmt"

	"github.com/zaydek/retro/cmd/retro"
	"github.com/zaydek/retro/pkg/versions"
)

//go:embed versions.txt
var contents string

func main() {
	vs, err := versions.Parse(contents)
	if err != nil {
		panic(err)
	}
	fmt.Println(vs)
	return

	retro.Run()
}
