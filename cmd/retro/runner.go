package retro

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zaydek/retro/cmd/retro/cli"
)

// const (
// 	KindDevCommand CommandKind = iota
// 	KindBuildCommand
// 	KindServeCommand
// )
//
// type CommandKind uint8

type Runner struct {
	Command interface{}
}

func (r Runner) warmUp() (copyHTML func(string, string, string) error, err error) {
	if err := setEnvsAndGlobalVariables(r.getCommandKind()); err != nil {
		return nil, fmt.Errorf("setEnvsAndGlobalVariables: %w", err)
	}

	if err := guardEntryPoints(); err != nil {
		return nil, fmt.Errorf("guardEntryPoints: %w", err)
	}

	if err := os.RemoveAll(RETRO_OUT_DIR); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(RETRO_OUT_DIR, permBitsDirectory); err != nil {
		return nil, err
	}
	if err := copyRecursively(RETRO_WWW_DIR, filepath.Join(RETRO_OUT_DIR, RETRO_WWW_DIR), []string{filepath.Join(RETRO_WWW_DIR, "index.html")}); err != nil {
		return nil, err
	}

	// TODO: The return statement here is weird
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
