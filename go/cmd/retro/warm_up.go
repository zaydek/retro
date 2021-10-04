package retro

import (
	"os"
	"path/filepath"

	"github.com/zaydek/retro/go/cmd/retro/unix"
)

func warmUp(commandKind CommandKind) error {
	if err := setEnvAndGlobalVariables(commandKind); err != nil { // Takes precedence
		return err
	}
	if err := guardEntryPoints(); err != nil {
		return err
	}
	if err := os.RemoveAll(RETRO_OUT_DIR); err != nil {
		return err
	}
	if err := os.MkdirAll(RETRO_OUT_DIR, 0755); err != nil {
		return err
	}
	target := filepath.Join(RETRO_OUT_DIR, RETRO_WWW_DIR)
	if err := unix.CopyRecursively(RETRO_WWW_DIR, target, []string{filepath.Join(RETRO_WWW_DIR, "index.html")}); err != nil {
		return err
	}
	return nil
}
