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
	VersionError = errors.New("cli: version error")
	UsageError   = errors.New("cli: usage error")
	HelpError    = errors.New("cli: help error")
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
	case BadPortValue:
		return fmt.Sprintf("'--port' must be a number (default '8000').")
	case BadSourcemapValue:
		return fmt.Sprintf("'--sourcemap' must be a 'true' or 'false' or empty (default 'true').")
	case BadCommandArgument:
		return fmt.Sprintf("Unsupported command argument '%s'.", e.BadCmdArgument)
	case BadArgument:
		return fmt.Sprintf("Unsupported argument '%s'.", e.BadArgument)
	case BadPortRange:
		return fmt.Sprintf("'--port' must be between '1000' and '10_000' (you may use '_' as a separator); used '%d'.", e.BadPort)
	}
	panic("Internal error")
}

func (e CommandError) Unwrap() error {
	return e.Err
}

// Support '_' as a separator
var portRegex = regexp.MustCompile(`^--port=([\d_]+)$`)

func ParseDevCommand(args ...string) (DevCommand, error) {
	cmd := DevCommand{
		Sourcemap: true,
		Port:      8000,
	}
	for _, arg := range args {
		err := CommandError{Kind: BadArgument, BadArgument: arg}
		if strings.HasPrefix(arg, "--port") {
			matches := portRegex.FindStringSubmatch(arg)
			if len(matches) == 2 {
				cmd.Port, _ = strconv.Atoi(strings.ReplaceAll(matches[1], "_", ""))
			} else {
				err.Kind = BadPortValue
				return DevCommand{}, err
			}
		} else if strings.HasPrefix(arg, "--sourcemap") {
			if arg == "--sourcemap" {
				cmd.Sourcemap = true
			} else if arg == "--sourcemap=true" || arg == "--sourcemap=false" {
				cmd.Sourcemap = arg == "--sourcemap=true"
			} else {
				err.Kind = BadSourcemapValue
				return DevCommand{}, err
			}
		} else {
			return DevCommand{}, err
		}
	}
	if cmd.Port < 1_000 || cmd.Port >= 10_000 {
		return DevCommand{}, CommandError{Kind: BadPortRange, BadPort: cmd.Port}
	}
	return cmd, nil
}

func ParseBuildCommand(args ...string) (BuildCommand, error) {
	cmd := BuildCommand{
		Sourcemap: true,
	}
	for _, arg := range args {
		err := CommandError{Kind: BadArgument, BadArgument: arg}
		if strings.HasPrefix(arg, "--sourcemap") {
			if arg == "--sourcemap" {
				cmd.Sourcemap = true
			} else if arg == "--sourcemap=true" || arg == "--sourcemap=false" {
				cmd.Sourcemap = arg == "--sourcemap=true"
			} else {
				err.Kind = BadSourcemapValue
				return BuildCommand{}, err
			}
		} else {
			return BuildCommand{}, err
		}
	}
	return cmd, nil
}

func ParseServeCommand(args ...string) (ServeCommand, error) {
	cmd := ServeCommand{
		Port: 8000,
	}
	for _, arg := range args {
		err := CommandError{Kind: BadArgument, BadArgument: arg}
		if strings.HasPrefix(arg, "--port") {
			matches := portRegex.FindStringSubmatch(arg)
			if len(matches) == 2 {
				cmd.Port, _ = strconv.Atoi(strings.ReplaceAll(matches[1], "_", ""))
			} else {
				err.Kind = BadPortValue
				return ServeCommand{}, err
			}
		} else {
			return ServeCommand{}, err
		}
	}
	if cmd.Port < 1_000 || cmd.Port >= 10_000 {
		return ServeCommand{}, CommandError{Kind: BadPortRange, BadPort: cmd.Port}
	}
	return cmd, nil
}

func ParseCLIArguments() (interface{}, error) {
	if len(os.Args) == 1 {
		return nil, UsageError
	}

	var (
		cmd interface{}
		err error
	)

	if cmdArg := os.Args[1]; cmdArg == "version" || cmdArg == "--version" || cmdArg == "-v" {
		return nil, VersionError
	} else if cmdArg == "usage" || cmdArg == "--usage" {
		return nil, UsageError
	} else if cmdArg == "help" || cmdArg == "--help" {
		return nil, HelpError
	} else if cmdArg == "dev" {
		cmd, err = ParseDevCommand(os.Args[2:]...)
	} else if cmdArg == "build" {
		cmd, err = ParseBuildCommand(os.Args[2:]...)
	} else if cmdArg == "serve" {
		cmd, err = ParseServeCommand(os.Args[2:]...)
	} else {
		err = CommandError{Kind: BadCommandArgument, BadCmdArgument: cmdArg}
	}
	return cmd, err
}
