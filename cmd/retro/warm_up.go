package retro

import (
	"fmt"
	"os"
	"path/filepath"
)

func warmUp(commandKind CommandKind) error {
	// Set environmental variables and global variables
	if err := setEnvsAndGlobalVariables(commandKind); err != nil { // Takes precedence
		return fmt.Errorf("setEnvsAndGlobalVariables: %w", err)
	}
	// Guard entry points
	if err := guardEntryPoints(); err != nil {
		return fmt.Errorf("guardEntryPoints: %w", err)
	}
	// Remove previous builds
	if err := os.RemoveAll(RETRO_OUT_DIR); err != nil {
		return fmt.Errorf("os.RemoveAll: %w", err)
	}
	// Make the out directory so static assets can be copied over
	if err := os.MkdirAll(RETRO_OUT_DIR, permBitsDirectory); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}
	// Copy `www` to `out/www`. Note that `www/index.html` is excluded because the
	// vendor and client script tags depend on `NODE_ENV`.
	target := filepath.Join(RETRO_OUT_DIR, RETRO_WWW_DIR)
	if err := copyDirectory(RETRO_WWW_DIR, target, []string{filepath.Join(RETRO_WWW_DIR, "index.html")}); err != nil {
		return fmt.Errorf("copyDirectory: %w", err)
	}
	return nil
}

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
