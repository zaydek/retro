package retro

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zaydek/retro/go/cmd/perm"
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
	if err := os.MkdirAll(RETRO_OUT_DIR, perm.BitsDirectory); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}
	// Copy `www` to `out/www`. Note that `www/index.html` is excluded because the
	// vendor and client are transformed by `transformAndCopyIndexHTMLEntryPoint`.
	target := filepath.Join(RETRO_OUT_DIR, RETRO_WWW_DIR)
	if err := copyDirectory(RETRO_WWW_DIR, target, []string{filepath.Join(RETRO_WWW_DIR, "index.html")}); err != nil {
		return fmt.Errorf("copyDirectory: %w", err)
	}
	return nil
}
