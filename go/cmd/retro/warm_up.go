package retro

import (
	"os"
	"path/filepath"

	"github.com/zaydek/retro/go/cmd/perm"
	"github.com/zaydek/retro/go/pkg/fsUtils"
)

func warmUp(commandKind CommandKind) error {
	if err := setEnv(commandKind); err != nil { // Takes precedence
		return decorate(&err, "setEnv")
	}
	if err := guardEntryPoints(); err != nil {
		return decorate(&err, "guardEntryPoints")
	}
	if err := os.RemoveAll(RETRO_OUT_DIR); err != nil {
		return decorate(&err, "os.RemoveAll")
	}
	if err := os.MkdirAll(RETRO_OUT_DIR, perm.BitsDirectory); err != nil {
		return decorate(&err, "os.MkdirAll")
	}
	target := filepath.Join(RETRO_OUT_DIR, RETRO_WWW_DIR)
	if err := fsUtils.CopyRecursively(RETRO_WWW_DIR, target, []string{filepath.Join(RETRO_WWW_DIR, "index.html")}); err != nil {
		return decorate(&err, "copyDirectory")
	}
	return nil
}
