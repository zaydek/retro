package create_retro_app

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/zaydek/retro/cmd/create_retro_app/cli"
	"github.com/zaydek/retro/cmd/create_retro_app/embeds"
	"github.com/zaydek/retro/cmd/retro/pretty"
	"github.com/zaydek/retro/cmd/shared"
	"github.com/zaydek/retro/pkg/terminal"
)

var (
	cyan    = func(str string) string { return pretty.Accent(str, terminal.Cyan) }
	magenta = func(str string) string { return pretty.Accent(str, terminal.Magenta) }
)

////////////////////////////////////////////////////////////////////////////////

func report(str string) {
	fmt.Fprintln(os.Stderr, pretty.Error(magenta(str)))
	os.Exit(1)
}

func (r Runner) CreateApp() {
	fsys := embeds.JavaScriptFS
	if r.Command.Template == "typescript" {
		fsys = embeds.TypeScriptFS
	}

	tmpl := embeds.JavaScriptPackageTemplate
	if r.Command.Template == "typescript" {
		tmpl = embeds.TypeScriptPackageTemplate
	}

	appName := r.Command.Directory
	if r.Command.Directory == "." {
		cwd, _ := os.Getwd()
		appName = filepath.Base(cwd)
	}

	if r.Command.Directory != "." {
		if info, err := os.Stat(r.Command.Directory); !os.IsNotExist(err) {
			var typ string
			if !info.IsDir() {
				typ = "file"
			} else {
				typ = "directory"
			}
			report(
				fmt.Sprintf("Aborted. "+
					"A %[1]s named %[3]s already exists. "+
					"Here’s what you can do:\n\n"+
					"- create-retro-app %[2]s\n\n"+
					"Or\n\n"+
					"- rm -r %[3]s && create-retro-app %[3]s",
					typ,
					increment(r.Command.Directory),
					r.Command.Directory,
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

	var paths []string
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
	for _, each := range paths {
		if _, err := os.Stat(each); !os.IsNotExist(err) {
			badPaths = append(badPaths, each)
		}
	}

	if len(badPaths) > 0 {
		var badPathsStr string
		for x, each := range badPaths {
			var sep string
			if x > 0 {
				sep = "\n"
			}
			badPathsStr += sep + "- " + terminal.Bold(each)
		}
		report(
			"Aborted. " +
				"These paths must be removed or renamed. " +
				"Use rm -r [paths] to remove them or mv [src] [dst] to rename them.\n\n" +
				badPathsStr,
		)
		os.Exit(1)
	}

	for _, each := range paths {
		if dir := filepath.Dir(each); dir != "." {
			if err := os.MkdirAll(dir, MODE_DIR); err != nil {
				panic(err)
			}
		}
		src, err := fsys.Open(each)
		if err != nil {
			panic(err)
		}
		dst, err := os.Create(each)
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(dst, src); err != nil {
			panic(err)
		}
		src.Close()
		dst.Close()
	}

	dot := embeds.PackageDot{
		APP_NAME:      appName,
		RETRO_VERSION: os.Getenv("RETRO_VERSION"),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, dot); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("package.json", buf.Bytes(), MODE_FILE); err != nil {
		panic(err)
	}

	if r.Command.Directory == "." {
		fmt.Println(fmt.Sprintf(successFormat, appName))
	} else {
		fmt.Println(fmt.Sprintf(successDirectoryFormat, appName))
	}
}

////////////////////////////////////////////////////////////////////////////////

type Runner struct {
	Command cli.Command
}

func Run() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}

	cmd, err := cli.ParseCLIArguments()
	switch err {
	case cli.VersionError:
		fmt.Println(shared.Package.Retro)
		return
	case cli.UsageError:
		fallthrough
	case cli.HelpError:
		fmt.Println(pretty.Inset(pretty.Spaces(cyan(usage))))
		return
	}

	runner := Runner{Command: cmd}
	runner.CreateApp()
}
