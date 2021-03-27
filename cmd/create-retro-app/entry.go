package create_retro_app

import (
	"fmt"
	"os"
)

func Run() {
	// Cover []string{"create-retro-app"} case:
	if len(os.Args) == 1 {
		fmt.Println(usage)
		os.Exit(0)
	}

	var cmd Command
	if arg := os.Args[1]; arg == "version" || arg == "--version" || arg == "-v" {
		fmt.Println(os.Getenv("retro_VERSION"))
		os.Exit(0)
	} else if arg == "help" || arg == "--help" || arg == "usage" || arg == "--usage" {
		fmt.Println(usage)
		os.Exit(0)
	} else {
		cmd = parseArguments(os.Args[1:]...)
	}
	cmd.CreateApp()
}
