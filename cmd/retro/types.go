package retro

type CommandKind uint8

const (
	KindDevCommand CommandKind = iota
	KindBuildCommand
	KindServeCommand
)

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
