package embeds

import (
	"embed"
	"fmt"
	"io/fs"
	"text/template"
)

func must(err error) {
	if err == nil {
		return
	}
	panic(fmt.Errorf("must: %w", err))
}

var (
	//go:embed package.json
	pkg string

	//go:embed package.sass.json
	sassPkg string

	StarterPackage = template.Must(template.New("package.json").Parse(pkg))
	SassPackage    = template.Must(template.New("package.json").Parse(sassPkg))
)

var (
	//go:embed starter/*
	starterFS embed.FS

	//go:embed starter-sass/*
	sassFS embed.FS

	StarterFS fs.FS
	SassFS    fs.FS
)

func init() {
	var err error
	StarterFS, err = fs.Sub(starterFS, "starter")
	must(err)
	SassFS, err = fs.Sub(sassFS, "starter-sass")
	must(err)
}
