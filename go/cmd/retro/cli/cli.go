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
	ErrVersion = errors.New("cli: version error")
	ErrUsage   = errors.New("cli: usage error")
	ErrHelp    = errors.New("cli: help error")
)

type ErrorKind int

const (
	BadCommandArgument ErrorKind = iota
	BadArgument
	BadPortValue
	BadSourcemapValue
	BadPortRange
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
	case BadPortValue:
		return "'--port' must be a number (default '8000')."
	case BadSourcemapValue:
		return "'--sourcemap' must be a 'true' or 'false' or empty (default 'true')."
	case BadPortRange:
		return fmt.Sprintf("'--port' must be between '1000' and '10000'; used '%d'.", e.BadPort)
	}
	panic("Internal error")
}

func (e CommandError) Unwrap() error {
	return e.Err
}

// Support _ separators
var portRegex = regexp.MustCompile(`^--port=([\d_]+)$`)

func ParseDevCommand(args ...string) (DevCommand, error) {
	command := DevCommand{
		Sourcemap: true,
		Port:      8000,
	}
	for _, arg := range args {
		err := CommandError{Kind: BadArgument, BadArgument: arg}
		if strings.HasPrefix(arg, "--port") {
			matches := portRegex.FindStringSubmatch(arg)
			if len(matches) == 2 {
				command.Port, _ = strconv.Atoi(strings.ReplaceAll(matches[1], "_", ""))
			} else {
				err.Kind = BadPortValue
				return DevCommand{}, err
			}
		} else if strings.HasPrefix(arg, "--sourcemap") {
			if arg == "--sourcemap" {
				command.Sourcemap = true
			} else if arg == "--sourcemap=true" || arg == "--sourcemap=false" {
				command.Sourcemap = arg == "--sourcemap=true"
			} else {
				err.Kind = BadSourcemapValue
				return DevCommand{}, err
			}
		} else {
			return DevCommand{}, err
		}
	}
	if command.Port < 1_000 || command.Port >= 10_000 {
		return DevCommand{}, CommandError{Kind: BadPortRange, BadPort: command.Port}
	}
	return command, nil
}

func ParseBuildCommand(args ...string) (BuildCommand, error) {
	command := BuildCommand{
		Sourcemap: true,
	}
	for _, arg := range args {
		err := CommandError{Kind: BadArgument, BadArgument: arg}
		if strings.HasPrefix(arg, "--sourcemap") {
			if arg == "--sourcemap" {
				command.Sourcemap = true
			} else if arg == "--sourcemap=true" || arg == "--sourcemap=false" {
				command.Sourcemap = arg == "--sourcemap=true"
			} else {
				err.Kind = BadSourcemapValue
				return BuildCommand{}, err
			}
		} else {
			return BuildCommand{}, err
		}
	}
	return command, nil
}

func ParseServeCommand(args ...string) (ServeCommand, error) {
	command := ServeCommand{
		Port: 8000,
	}
	for _, arg := range args {
		err := CommandError{Kind: BadArgument, BadArgument: arg}
		if strings.HasPrefix(arg, "--port") {
			matches := portRegex.FindStringSubmatch(arg)
			if len(matches) == 2 {
				command.Port, _ = strconv.Atoi(strings.ReplaceAll(matches[1], "_", ""))
			} else {
				err.Kind = BadPortValue
				return ServeCommand{}, err
			}
		} else {
			return ServeCommand{}, err
		}
	}
	if command.Port < 1_000 || command.Port >= 10_000 {
		return ServeCommand{}, CommandError{Kind: BadPortRange, BadPort: command.Port}
	}
	return command, nil
}

func ParseCLIArguments() (interface{}, error) {
	if len(os.Args) < 2 {
		return nil, ErrUsage
	}

	var (
		command interface{}
		err     error
	)

	if cmdArg := os.Args[1]; cmdArg == "version" || cmdArg == "--version" || cmdArg == "-v" {
		return nil, ErrVersion
	} else if cmdArg == "usage" || cmdArg == "--usage" {
		return nil, ErrUsage
	} else if cmdArg == "help" || cmdArg == "--help" {
		return nil, ErrHelp
	} else if cmdArg == "dev" {
		command, err = ParseDevCommand(os.Args[2:]...)
	} else if cmdArg == "build" {
		command, err = ParseBuildCommand(os.Args[2:]...)
	} else if cmdArg == "serve" {
		command, err = ParseServeCommand(os.Args[2:]...)
	} else {
		err = CommandError{Kind: BadCommandArgument, BadCmdArgument: cmdArg}
	}
	return command, err
}
