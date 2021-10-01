package retro

import (
	"fmt"
	"os"
)

var (
	// Global variables mirroring environmental variables
	NODE_ENV      = ""
	RETRO_CMD     = ""
	RETRO_WWW_DIR = ""
	RETRO_SRC_DIR = ""
	RETRO_OUT_DIR = ""
)

func setEnv(errPointer *error, envKey, defaultValue string) {
	if *errPointer != nil {
		return
	}
	envValue := os.Getenv(envKey)
	if envValue == "" {
		envValue = defaultValue
	}
	switch envKey {
	case "NODE_ENV":
		NODE_ENV = envValue
	case "RETRO_CMD":
		RETRO_CMD = envValue
	case "RETRO_WWW_DIR":
		RETRO_WWW_DIR = envValue
	case "RETRO_SRC_DIR":
		RETRO_SRC_DIR = envValue
	case "RETRO_OUT_DIR":
		RETRO_OUT_DIR = envValue
	}
	if err := os.Setenv(envKey, envValue); err != nil {
		*errPointer = fmt.Errorf("os.Setenv: %w", err)
	}
}

// Propagates environmental variables or sets default values
func setEnvsAndGlobalVariables(commandKind CommandKind) error {
	var err error
	switch commandKind {
	case KindDevCommand:
		setEnv(&err, "NODE_ENV", "development")
		setEnv(&err, "RETRO_CMD", "dev")
		setEnv(&err, "RETRO_WWW_DIR", "www")
		setEnv(&err, "RETRO_SRC_DIR", "src")
		setEnv(&err, "RETRO_OUT_DIR", "out")
	case KindBuildCommand:
		setEnv(&err, "NODE_ENV", "production")
		setEnv(&err, "RETRO_CMD", "build")
		setEnv(&err, "RETRO_WWW_DIR", "www")
		setEnv(&err, "RETRO_SRC_DIR", "src")
		setEnv(&err, "RETRO_OUT_DIR", "out")
	case KindServeCommand:
		setEnv(&err, "NODE_ENV", "production")
		setEnv(&err, "RETRO_CMD", "serve")
		setEnv(&err, "RETRO_WWW_DIR", "www")
		setEnv(&err, "RETRO_SRC_DIR", "src")
		setEnv(&err, "RETRO_OUT_DIR", "out")
	}
	return err
}
