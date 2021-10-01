package retro

type CommandKind uint8

const (
	KindDevCommand CommandKind = iota
	KindBuildCommand
	KindServeCommand
)

// TODO: Why do we need a runner abstraction? Can we not use a CLI command
// directly? I think the idea was to provide an interface that can resolve the
// CLI type at runtime because Go's type-system is more strict.

// func (r Runner) getCommandKind() (out CommandKind) {
// 	switch r.Command.(type) {
// 	// % retro dev
// 	case cli.DevCommand:
// 		return KindDevCommand
// 	// % retro build
// 	case cli.BuildCommand:
// 		return KindBuildCommand
// 	// % retro serve
// 	case cli.ServeCommand:
// 		return KindServeCommand
// 	}
// 	return
// }
//
// func (r Runner) getPort() (out int) {
// 	if cmd := r.getCommandKind(); cmd == KindDevCommand {
// 		return r.Command.(cli.DevCommand).Port
// 	} else if cmd == KindServeCommand {
// 		return r.Command.(cli.ServeCommand).Port
// 	}
// 	return
// }
