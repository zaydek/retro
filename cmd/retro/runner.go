package retro

import (
	"os"
	"path/filepath"

	"github.com/zaydek/retro/cmd/retro/cli"
)

const (
	DevCommand CmdKind = iota
	BuildCommand
	ServeCommand
)

type Runner struct {
	Command interface{}
}

type CmdKind uint8

func (r Runner) preflight() (copyHTML func(string, string, string) error, err error) {
	cmd := r.getCommandKind()

	// Set env vars
	switch cmd {
	case DevCommand:
		os.Setenv("CMD", "dev")
		os.Setenv("ENV", "development")
		os.Setenv("NODE_ENV", "development")
		os.Setenv("WWW_DIR", WWW_DIR)
		os.Setenv("SRC_DIR", SRC_DIR)
		os.Setenv("OUT_DIR", OUT_DIR)

	case BuildCommand:
		os.Setenv("CMD", "build")
		os.Setenv("ENV", "production")
		os.Setenv("NODE_ENV", "production")
		os.Setenv("WWW_DIR", WWW_DIR)
		os.Setenv("SRC_DIR", SRC_DIR)
		os.Setenv("OUT_DIR", OUT_DIR)

	case ServeCommand:
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
	if err := cpdir(WWW_DIR, filepath.Join(OUT_DIR, WWW_DIR), []string{"index.html"}); err != nil {
		return nil, err
	}
	copyHTML = copyHTMLEntryPoint

	return
}

func (r Runner) getCommandKind() (out CmdKind) {
	switch r.Command.(type) {
	// % retro dev
	case cli.DevCommand:
		return DevCommand
	// % retro build
	case cli.BuildCommand:
		return BuildCommand
	// % retro serve
	case cli.ServeCommand:
		return ServeCommand
	}
	return
}

func (r Runner) getPort() (out int) {
	if cmd := r.getCommandKind(); cmd == DevCommand {
		return r.Command.(cli.DevCommand).Port
	} else if cmd == ServeCommand {
		return r.Command.(cli.ServeCommand).Port
	}
	return
}
