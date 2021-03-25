package cli

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	VersionError = errors.New("version")
	UsageError   = errors.New("usage")
)

type ErrorKind int

const (
	BadCommandArgument ErrorKind = iota
	BadArgument
	BadPort
)

type CommandError struct {
	Kind ErrorKind

	BadCmdArgument string
	BadArgument    string
	BadPort        int

	Err error
}

func (e CommandError) Error() string {
	switch e.Kind {
	case BadCommandArgument:
		return fmt.Sprintf("Unsupported command argument '%s'.", e.BadCmdArgument)
	case BadArgument:
		return fmt.Sprintf("Unsupported argument '%s'.", e.BadArgument)
	case BadPort:
		return fmt.Sprintf("'--port' must be between '1000' and '10000'; used '%d'.", e.BadPort)
	}
	panic("Internal error")
}

func (e CommandError) Unwrap() error {
	return e.Err
}

var portRegex = regexp.MustCompile(`^--port=(\d+)$`)

func ParseDevCommand(args ...string) (DevCommand, error) {
	cmd := DevCommand{
		Sourcemap: true,
		Port:      8000,
	}
	for _, arg := range args {
		cmdErr := CommandError{Kind: BadArgument, BadArgument: arg}
		if strings.HasPrefix(arg, "--port") {
			matches := portRegex.FindStringSubmatch(arg)
			if len(matches) == 2 {
				cmd.Port, _ = strconv.Atoi(matches[1])
			} else {
				return DevCommand{}, cmdErr
			}
		} else if strings.HasPrefix(arg, "--sourcemap") {
			if arg == "--sourcemap" {
				cmd.Sourcemap = true
			} else if arg == "--sourcemap=true" || arg == "--sourcemap=false" {
				cmd.Sourcemap = arg == "--sourcemap=true"
			} else {
				return DevCommand{}, cmdErr
			}
		} else {
			return DevCommand{}, cmdErr
		}
	}
	if cmd.Port < 1_000 || cmd.Port >= 10_000 {
		return DevCommand{}, CommandError{Kind: BadPort, BadPort: cmd.Port}
	}
	return cmd, nil
}

func ParseExportCommand(args ...string) (BuildCommand, error) {
	cmd := BuildCommand{
		Sourcemap: true,
	}
	for _, arg := range args {
		cmdErr := CommandError{Kind: BadArgument, BadArgument: arg}
		if strings.HasPrefix(arg, "--sourcemap") {
			if arg == "--sourcemap" {
				cmd.Sourcemap = true
			} else if arg == "--sourcemap=true" || arg == "--sourcemap=false" {
				cmd.Sourcemap = arg == "--sourcemap=true"
			} else {
				return BuildCommand{}, cmdErr
			}
		} else {
			return BuildCommand{}, cmdErr
		}
	}
	return cmd, nil
}

func ParseServeCommand(args ...string) (ServeCommand, error) {
	cmd := ServeCommand{
		Port: 8000,
	}
	for _, arg := range args {
		cmdErr := CommandError{Kind: BadArgument, BadArgument: arg}
		if strings.HasPrefix(arg, "--port") {
			matches := portRegex.FindStringSubmatch(arg)
			if len(matches) == 2 {
				cmd.Port, _ = strconv.Atoi(matches[1])
			} else {
				return ServeCommand{}, cmdErr
			}
		} else {
			return ServeCommand{}, cmdErr
		}
	}
	if cmd.Port < 1_000 || cmd.Port >= 10_000 {
		return ServeCommand{}, CommandError{Kind: BadPort, BadPort: cmd.Port}
	}
	return cmd, nil
}

func ParseCLIArguments() (interface{}, error) {
	if len(os.Args) == 1 {
		return nil, UsageError
	}

	var cmd interface{}
	var cmdErr error

	if cmdArg := os.Args[1]; cmdArg == "version" || cmdArg == "--version" || cmdArg == "--v" {
		return nil, VersionError
	} else if cmdArg == "usage" || cmdArg == "--usage" || cmdArg == "help" || cmdArg == "--help" {
		return nil, UsageError
	} else if cmdArg == "dev" {
		cmd, cmdErr = ParseDevCommand(os.Args[2:]...)
	} else if cmdArg == "export" {
		cmd, cmdErr = ParseExportCommand(os.Args[2:]...)
	} else if cmdArg == "serve" {
		cmd, cmdErr = ParseServeCommand(os.Args[2:]...)
	} else {
		cmdErr = CommandError{Kind: BadCommandArgument, BadCmdArgument: cmdArg}
	}
	return cmd, cmdErr
}
