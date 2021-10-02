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

func setEnvImpl(errPointer *error, envKey, defaultValue string) {
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
func setEnv(commandKind CommandKind) error {
	var err error
	switch commandKind {
	case KindDevCommand:
		setEnvImpl(&err, "NODE_ENV", "development")
		setEnvImpl(&err, "RETRO_CMD", "dev")
		setEnvImpl(&err, "RETRO_WWW_DIR", "www")
		setEnvImpl(&err, "RETRO_SRC_DIR", "src")
		setEnvImpl(&err, "RETRO_OUT_DIR", "out")
	case KindBuildCommand:
		setEnvImpl(&err, "NODE_ENV", "production")
		setEnvImpl(&err, "RETRO_CMD", "build")
		setEnvImpl(&err, "RETRO_WWW_DIR", "www")
		setEnvImpl(&err, "RETRO_SRC_DIR", "src")
		setEnvImpl(&err, "RETRO_OUT_DIR", "out")
	case KindServeCommand:
		setEnvImpl(&err, "NODE_ENV", "production")
		setEnvImpl(&err, "RETRO_CMD", "serve")
		setEnvImpl(&err, "RETRO_WWW_DIR", "www")
		setEnvImpl(&err, "RETRO_SRC_DIR", "src")
		setEnvImpl(&err, "RETRO_OUT_DIR", "out")
	}
	return err
}
