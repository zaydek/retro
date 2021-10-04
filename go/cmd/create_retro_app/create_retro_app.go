package create_retro_app

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/zaydek/retro/go/cmd/create_retro_app/cli"
	"github.com/zaydek/retro/go/cmd/format"
	"github.com/zaydek/retro/go/pkg/terminal"
)

// TODO: Can we deprecate this?
var cyan = func(str string) string { return format.Accent(str, terminal.Cyan) }

//go:embed static/*
var staticFS embed.FS

func (r App) CreateApp() error {
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
	}

	var copyPaths []string
	err := fs.WalkDir(staticFS, ".", func(root string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			rel, _ := filepath.Rel("static", root)
			copyPaths = append(copyPaths, rel)
		}
		return nil
	})
	if err != nil {
		return err
	}

	dirName := r.Command.Directory
	if r.Command.Directory == "." {
		wd, _ := os.Getwd()
		dirName = filepath.Base(wd)
	}

	var badCopyPaths []string
	for _, path := range copyPaths {
		path := filepath.Join(dirName, path)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			badCopyPaths = append(badCopyPaths, path)
		}
	}

	if len(badCopyPaths) > 0 {
		var badPathsStr string
		for badPathIndex, badPath := range badCopyPaths {
			var sep string
			if badPathIndex > 0 {
				sep = "\n"
			}
			badPathsStr += sep + "- " + badPath
		}
		fmt.Fprintln(
			os.Stderr,
			format.Error(
				fmt.Sprintf(
					"Refusing to overwrite files and or directories.\n\n"+
						badPathsStr,
				),
			),
		)
		os.Exit(1)
	}

	for _, copyPath := range copyPaths {
		path := filepath.Join(dirName, copyPath)
		if dir := filepath.Dir(path); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
		src, err := staticFS.Open(filepath.Join("static", copyPath))
		if err != nil {
			return err
		}
		dst, err := os.Create(path)
		if err != nil {
			return err
		}
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}
		src.Close()
		dst.Close()
	}

	pkg := fmt.Sprintf(
		`{
	"scripts": {
		"dev": "retro dev",
		"build": "retro build",
		"serve": "retro serve"
	},
	"dependencies": {
		"react": "%[2]s",
		"react-dom": "%[3]s"
	},
	"devDependencies": {
		"@zaydek/retro": "%[4]s",
		"esbuild": "%[1]s"
	}
}`,
		os.Getenv("ESBUILD_VERSION"),
		os.Getenv("REACT_VERSION"),
		os.Getenv("REACTDOM_VERSION"),
		os.Getenv("RETRO_VERSION"),
	)

	if err := os.WriteFile(filepath.Join(dirName, "package.json"), []byte(pkg+"\n"), 0644); err != nil {
		return err
	}

	if r.Command.Directory == "." {
		// TODO: Clean this up?
		fmt.Println(format.Tabs(createSuccessStr))
	} else {
		// TODO: Clean this up?
		fmt.Println(format.Tabs(fmt.Sprintf(createSuccessDirStr, dirName)))
	}

	return nil
}

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
