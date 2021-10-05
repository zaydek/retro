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

//go:embed static/*
var staticFS embed.FS

type App struct {
	Command cli.CreateCommand
}

func (r App) CreateApp() error {
	if r.Command.Directory != "." {
		if _, err := os.Stat(r.Command.Directory); !os.IsNotExist(err) {
			errStr := fmt.Sprintf("Refusing to overwrite directory `%s`.", r.Command.Directory)
			fmt.Fprintln(os.Stderr, errStr)
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

	dir := r.Command.Directory
	if r.Command.Directory == "." {
		wd, _ := os.Getwd()
		dir = filepath.Base(wd)
	}

	var badCopyPaths []string
	for _, path := range copyPaths {
		path := filepath.Join(dir, path)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			badCopyPaths = append(badCopyPaths, path)
		}
	}

	if len(badCopyPaths) > 0 {
		var badCopyPathsStr string
		for _, badCopyPath := range badCopyPaths {
			if badCopyPathsStr != "" {
				badCopyPathsStr += "\n"
			}
			badCopyPathsStr += "- " + badCopyPath
		}
		errStr := format.Stderr(fmt.Sprintf("Refusing to overwrite files and or directories.\n\n%s", badCopyPathsStr))
		fmt.Fprintln(os.Stderr, errStr)
		os.Exit(1)
	}

	for _, copyPath := range copyPaths {
		path := filepath.Join(dir, copyPath)
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

	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg+"\n"), 0644); err != nil {
		return err
	}

	if r.Command.Directory == "." {
		fmt.Println(terminal.Cyanf("Success! %s", terminal.Dimf("(%s)", os.Getenv("RETRO_V_VERSION"))) + `

 npm:
   npm i
   npm run dev

 yarn:
   yarn
   yarn dev

`)
	} else {
		fmt.Println(fmt.Sprintf(terminal.Cyanf("Success! %s", terminal.Dimf("(%s)", os.Getenv("RETRO_V_VERSION")))+`

 npm:
   cd %[1]s
   npm i
   npm run dev

 yarn:
   cd %[1]s
   yarn
   yarn dev

Happy hacking!`, dir))
	}

	return nil
}

func Run() {
	// Non-command errors
	command, err := cli.ParseCLIArguments()
	switch err {
	case cli.ErrVersion:
		fmt.Println(os.Getenv("RETRO_V_VERSION"))
		return
	case cli.ErrUsage:
		fallthrough
	case cli.ErrHelp:
		fmt.Println(format.Stdout(usage))
		return
	}

	// Command errors
	switch err.(type) {
	case cli.CommandError:
		fmt.Fprintln(os.Stderr, format.Stderr(err))
		os.Exit(1)
	default:
		must(err)
	}

	app := &App{Command: command}
	must(app.CreateApp())
}
