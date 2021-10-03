package create_retro_app

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/zaydek/retro/go/cmd/create_retro_app/cli"
	"github.com/zaydek/retro/go/cmd/create_retro_app/embeds"
	"github.com/zaydek/retro/go/cmd/deps"
	"github.com/zaydek/retro/go/cmd/format"
	"github.com/zaydek/retro/go/cmd/perm"
	"github.com/zaydek/retro/go/pkg/terminal"
)

// TODO: Can we deprecate this?
var cyan = func(str string) string { return format.Accent(str, terminal.Cyan) }

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

////////////////////////////////////////////////////////////////////////////////

func (r App) mustGetFSAndPKG() (fs.FS, *template.Template) {
	switch r.Command.Template {
	case "starter":
		return embeds.StarterFS, embeds.StarterPackage
	case "sass":
		return embeds.SassFS, embeds.SassPackage
	}
	panic("Internal error")
}

func (r App) CreateApp() error {
	fsys, pkg := r.mustGetFSAndPKG()

	appName := r.Command.Directory
	if r.Command.Directory == "." {
		wd, _ := os.Getwd()
		appName = filepath.Base(wd)
	}

	if r.Command.Directory != "." {
		if _, err := os.Stat(r.Command.Directory); !os.IsNotExist(err) {
			fmt.Fprintln(
				os.Stderr,
				format.Error(
					fmt.Sprintf(
						"Refusing to overwrite directory `%s`.",
						r.Command.Directory,
					),
				),
			)
			os.Exit(1)
		}
		if err := os.MkdirAll(r.Command.Directory, perm.BitsDirectory); err != nil {
			return err
		}
		if err := os.Chdir(r.Command.Directory); err != nil {
			return err
		}
		defer os.Chdir("..")
	}

	// Add package.json
	paths := []string{"package.json"}
	err := fs.WalkDir(fsys, ".", func(root string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			paths = append(paths, root)
		}
		return nil
	})
	if err != nil {
		return err
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
					"Refusing to overwrite paths. Use `rm -r [...paths]` to remove them or `mv [src] [dst]` to rename them.\n\n"+
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
			if err := os.MkdirAll(dir, perm.BitsDirectory); err != nil {
				return err
			}
		}
		src, err := fsys.Open(v)
		if err != nil {
			return err
		}
		dst, err := os.Create(v)
		if err != nil {
			return err
		}
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}
		src.Close()
		dst.Close()
	}

	var buf bytes.Buffer
	deps.Deps.RetroVersion = os.Getenv("RETRO_VERSION") // Add @zaydek/retro
	if err := pkg.Execute(&buf, deps.Deps); err != nil {
		return err
	}

	if err := os.WriteFile("package.json", buf.Bytes(), perm.BitsFile); err != nil {
		return err
	}

	if r.Command.Directory == "." {
		fmt.Println(format.Tabs(successFmt))
	} else {
		fmt.Println(format.Tabs(fmt.Sprintf(successDirFmt, appName)))
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

type App struct {
	Command cli.CreateCommand
}

func Run() {
	// Non-command errors
	command, err := cli.ParseCLIArguments()
	switch err {
	case cli.ErrVersion:
		fmt.Println(os.Getenv("RETRO_VERSION"))
		return
	case cli.ErrUsage:
		fallthrough
	case cli.ErrHelp:
		fmt.Println(format.Pad(format.Tabs(cyan(usage))))
		return
	}

	// Command errors
	switch err.(type) {
	case cli.CommandError:
		fmt.Fprintln(os.Stderr, format.Error(err))
		os.Exit(1)
	default:
		must(err)
	}

	app := &App{Command: command}
	must(app.CreateApp())
}
