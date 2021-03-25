package retro

import (
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
