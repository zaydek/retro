package retro

import "github.com/zaydek/retro/go/cmd/retro/cli"

type CommandKind string

var (
	KindDevCommand   CommandKind = "dev"
	KindBuildCommand CommandKind = "build"
	KindServeCommand CommandKind = "serve"
)

// Abstraction on top of the returned CLI command
type App struct {
	// One of `*cli.DevCommand`, `*cli.BuildCommand`, or `*cli.ServeCommand`. Note
	// that pointers are used because port numbers may be incremented.
	Command interface{}
}

// Gets the app's command kind; one of dev, build, or serve
func (a *App) getCommandKind() CommandKind {
	var zeroValue CommandKind
	switch a.Command.(type) {
	case *cli.DevCommand:
		return KindDevCommand
	case *cli.BuildCommand:
		return KindBuildCommand
	case *cli.ServeCommand:
		return KindServeCommand
	}
	return zeroValue
}

// Gets the app's port number
func (a *App) getPort() int {
	var zeroValue int
	if cmd := a.getCommandKind(); cmd == KindDevCommand {
		return a.Command.(*cli.DevCommand).Port
	} else if cmd == KindServeCommand {
		return a.Command.(*cli.ServeCommand).Port
	}
	return zeroValue
}

// Sets the app's port number
func (a *App) setPort(port int) {
	if cmd := a.getCommandKind(); cmd == KindDevCommand {
		a.Command.(*cli.DevCommand).Port = port
	} else if cmd == KindServeCommand {
		a.Command.(*cli.ServeCommand).Port = port
	}
}
