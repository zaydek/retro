package create_retro_app

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/zaydek/retro/cmd/create_retro_app/cli"
	"github.com/zaydek/retro/cmd/create_retro_app/embeds"
	"github.com/zaydek/retro/cmd/deps"
	"github.com/zaydek/retro/cmd/format"
	"github.com/zaydek/retro/pkg/terminal"
)

var cyan = func(str string) string { return format.Accent(str, terminal.Cyan) }

////////////////////////////////////////////////////////////////////////////////

func (r Runner) getFSAndPkg() (fs.FS, *template.Template) {
	switch r.Command.Template {
	case "starter":
		return embeds.StarterFS, embeds.StarterPackage
	case "sass":
		return embeds.SassFS, embeds.SassPackage
	case "mdx":
		return embeds.MDXFS, embeds.MDXPackage
	}
	panic("Internal error")
}

func (r Runner) CreateApp() {
	fsys, pkg := r.getFSAndPkg()

	appName := r.Command.Directory
	if r.Command.Directory == "." {
		cwd, _ := os.Getwd()
		appName = filepath.Base(cwd)
	}

	if r.Command.Directory != "." {
		if _, err := os.Stat(r.Command.Directory); !os.IsNotExist(err) {
			fmt.Fprintln(
				os.Stderr,
				format.Error(
					fmt.Sprintf(
						"Aborted. Refusing to overwrite directory '%s'.",
						r.Command.Directory,
					),
				),
			)
			os.Exit(1)
		}
		if err := os.MkdirAll(r.Command.Directory, MODE_DIR); err != nil {
			panic(err)
		}
		if err := os.Chdir(r.Command.Directory); err != nil {
			panic(err)
		}
		defer os.Chdir("..")
	}

	// Add package.json
	paths := []string{"package.json"}
	err := fs.WalkDir(fsys, ".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !dirEntry.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	var badPaths []string
	for _, v := range paths {
		if _, err := os.Stat(v); !os.IsNotExist(err) {
			badPaths = append(badPaths, v)
		}
	}

	if len(badPaths) > 0 {
		var badPathsStr string
		for x, v := range badPaths {
			var sep string
			if x > 0 {
				sep = "\n"
			}
			badPathsStr += sep + "- " + v
		}
		fmt.Fprintln(
			os.Stderr,
			format.Error(
				fmt.Sprintf(
					"Aborted. Refusing to overwrite paths. Use 'rm -r [...paths]' to remove them or 'mv [src] [dst]' to rename them.\n\n"+
						badPathsStr,
				),
			),
		)
		os.Exit(1)
	}

	// Remove package.json
	paths = paths[1:]
	for _, v := range paths {
		if dir := filepath.Dir(v); dir != "." {
			if err := os.MkdirAll(dir, MODE_DIR); err != nil {
				panic(err)
			}
		}
		src, err := fsys.Open(v)
		if err != nil {
			panic(err)
		}
		dst, err := os.Create(v)
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(dst, src); err != nil {
			panic(err)
		}
		src.Close()
		dst.Close()
	}

	var buf bytes.Buffer
	deps.Deps.RetroVersion = os.Getenv("RETRO_VERSION") // Add '@zaydek/retro'
	if err := pkg.Execute(&buf, deps.Deps); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("package.json", buf.Bytes(), MODE_FILE); err != nil {
		panic(err)
	}

	if r.Command.Directory == "." {
		fmt.Println(format.TabsToSpaces(successFmt))
	} else {
		fmt.Println(format.TabsToSpaces(fmt.Sprintf(successDirFmt, appName)))
	}
}

////////////////////////////////////////////////////////////////////////////////

type Runner struct {
	Command cli.Command
}

func Run() {
	cmd, err := cli.ParseCLIArguments()
	switch err {
	case cli.ErrVersion:
		fmt.Println(os.Getenv("RETRO_VERSION"))
		return
	case cli.ErrUsage:
		fallthrough
	case cli.ErrHelp:
		fmt.Println(format.SpaceInset(format.TabsToSpaces(cyan(usage))))
		return
	}

	switch err.(type) {
	case cli.CommandError:
		fmt.Fprintln(os.Stderr, format.Error(err))
		os.Exit(1)
	default:
		if err != nil {
			panic(err)
		}
	}

	runner := Runner{Command: cmd}
	runner.CreateApp()
}
