package embeds

import (
	"embed"
	"io/fs"
	"text/template"
)

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

var (
	//go:embed package.json
	pkg string

	//go:embed package.sass.json
	sassPkg string

	//go:embed package.mdx.json
	mdxPkg string

	StarterPackage = template.Must(template.New("package.json").Parse(pkg))
	SassPackage    = template.Must(template.New("package.json").Parse(sassPkg))
	MDXPackage     = template.Must(template.New("package.json").Parse(mdxPkg))
)

var (
	//go:embed starter/*
	starterFS embed.FS

	//go:embed starter-sass/*
	sassFS embed.FS

	//go:embed starter-mdx/*
	mdxFS embed.FS

	StarterFS fs.FS
	SassFS    fs.FS
	MDXFS     fs.FS
)

func init() {
	var err error
	StarterFS, err = fs.Sub(starterFS, "starter")
	must(err)
	SassFS, err = fs.Sub(sassFS, "starter-sass")
	must(err)
	MDXFS, err = fs.Sub(mdxFS, "starter-mdx")
	must(err)
}
