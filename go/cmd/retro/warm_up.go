package retro

import (
	"os"
	"path/filepath"

	"github.com/zaydek/retro/go/cmd/perm"
	"github.com/zaydek/retro/go/pkg/sys"
)

func warmUp(commandKind CommandKind) error {
	// Set environmental variables and global variables
	if err := setEnv(commandKind); err != nil { // Takes precedence
		return err
	}
	// Guard HTML and JS entry points
	if err := guardEntryPoints(); err != nil {
		return err
	}
	// Remove `RETRO_OUT_DIR` directory
	if err := os.RemoveAll(RETRO_OUT_DIR); err != nil {
		return err
	}
	// Create `RETRO_OUT_DIR` directgory
	if err := os.MkdirAll(RETRO_OUT_DIR, perm.BitsDirectory); err != nil {
		return err
	}
	// Copy `RETRO_WWW_DIR` directory recursively
	target := filepath.Join(RETRO_OUT_DIR, RETRO_WWW_DIR)
	if err := sys.CopyRecursively(RETRO_WWW_DIR, target, []string{filepath.Join(RETRO_WWW_DIR, "index.html")}); err != nil {
		return err
	}
	return nil
}
