package retro

import "github.com/zaydek/retro/go/cmd/retro/cli"

type CommandKind string

var (
	KindDevCommand   CommandKind = "dev"
	KindBuildCommand CommandKind = "build"
	KindServeCommand CommandKind = "serve"
)

type App struct {
	// One of:
	//
	// - cli.DevCommand
	// - cli.BuildCommand
	// - cli.ServeCommand
	//
	Command interface{}
}

// Gets the app's command kind; one of dev, build, or serve
func (a *App) getCommandKind() CommandKind {
	var zeroValue CommandKind
	switch a.Command.(type) {
	case cli.DevCommand:
		return KindDevCommand
	case cli.BuildCommand:
		return KindBuildCommand
	case cli.ServeCommand:
		return KindServeCommand
	}
	return zeroValue
}

// Gets the app's port number
func (a *App) getPort() int {
	var zeroValue int
	if commandKind := a.getCommandKind(); commandKind == KindDevCommand {
		return a.Command.(cli.DevCommand).Port
	} else if commandKind == KindServeCommand {
		return a.Command.(cli.ServeCommand).Port
	}
	return zeroValue
}
