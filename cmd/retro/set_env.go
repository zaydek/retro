package retro

import (
	"fmt"
	"os"
)

var (
	NODE_ENV      = ""
	RETRO_CMD     = ""
	RETRO_WWW_DIR = ""
	RETRO_SRC_DIR = ""
	RETRO_OUT_DIR = ""
)

// Propagates environmental variables or sets default values
func setEnvsAndGlobalVariables(commandKind CommandKind) error {
	var err error
	setEnv := func(envKey, defaultValue string) {
		if err != nil {
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
		err = os.Setenv(envKey, envValue)
		if err != nil {
			err = fmt.Errorf("os.Setenv: %w", err)
		}
	}
	switch commandKind {
	case KindDevCommand:
		setEnv("NODE_ENV", "development")
		setEnv("RETRO_CMD", "dev")
		setEnv("RETRO_WWW_DIR", "www")
		setEnv("RETRO_SRC_DIR", "src")
		setEnv("RETRO_OUT_DIR", "out")
	case KindBuildCommand:
		setEnv("NODE_ENV", "production")
		setEnv("RETRO_CMD", "build")
		setEnv("RETRO_WWW_DIR", "www")
		setEnv("RETRO_SRC_DIR", "src")
		setEnv("RETRO_OUT_DIR", "out")
	case KindServeCommand:
		setEnv("NODE_ENV", "production")
		setEnv("RETRO_CMD", "serve")
		setEnv("RETRO_WWW_DIR", "www")
		setEnv("RETRO_SRC_DIR", "src")
		setEnv("RETRO_OUT_DIR", "out")
	}
	return err
}
