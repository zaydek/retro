package retro

import (
	"os"
	"path/filepath"

	"github.com/zaydek/retro/cmd/retro/cli"
)

const (
	KindDevCommand CommandKind = iota
	KindBuildCommand
	KindServeCommand
)

type Runner struct {
	Command interface{}
}

type CommandKind uint8

func (r Runner) preflight() (copyHTML func(string, string, string) error, err error) {
	cmd := r.getCommandKind()

	// Set env vars
	switch cmd {

	// TODO: Add support for source map, etc.
	case KindDevCommand:
		os.Setenv("CMD", "dev")
		os.Setenv("ENV", "development")
		os.Setenv("NODE_ENV", "development")
		os.Setenv("WWW_DIR", WWW_DIR)
		os.Setenv("SRC_DIR", SRC_DIR)
		os.Setenv("OUT_DIR", OUT_DIR)

	// TODO: Add support for source map, etc.
	case KindBuildCommand:
		os.Setenv("CMD", "build")
		os.Setenv("ENV", "production")
		os.Setenv("NODE_ENV", "production")
		os.Setenv("WWW_DIR", WWW_DIR)
		os.Setenv("SRC_DIR", SRC_DIR)
		os.Setenv("OUT_DIR", OUT_DIR)

	// TODO: Add support for source map, etc.
	case KindServeCommand:
		os.Setenv("CMD", "serve")
		os.Setenv("ENV", "production")
		os.Setenv("NODE_ENV", "production")
		os.Setenv("WWW_DIR", WWW_DIR)
		os.Setenv("SRC_DIR", SRC_DIR)
		os.Setenv("OUT_DIR", OUT_DIR)
		return
	}

	if err := guardHTMLEntryPoint(); err != nil { // Takes precedence
		return nil, err
	}
	if err := os.RemoveAll(OUT_DIR); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(OUT_DIR, MODE_DIR); err != nil {
		return nil, err
	}
	if err := cpdir(WWW_DIR, filepath.Join(OUT_DIR, WWW_DIR), []string{filepath.Join(WWW_DIR, "index.html")}); err != nil {
		return nil, err
	}

	copyHTML = copyHTMLEntryPoint
	return
}

func (r Runner) getCommandKind() (out CommandKind) {
	switch r.Command.(type) {
	// % retro dev
	case cli.DevCommand:
		return KindDevCommand
	// % retro build
	case cli.BuildCommand:
		return KindBuildCommand
	// % retro serve
	case cli.ServeCommand:
		return KindServeCommand
	}
	return
}

func (r Runner) getPort() (out int) {
	if cmd := r.getCommandKind(); cmd == KindDevCommand {
		return r.Command.(cli.DevCommand).Port
	} else if cmd == KindServeCommand {
		return r.Command.(cli.ServeCommand).Port
	}
	return
}
