package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	VersionError = errors.New("cli: version error")
	UsageError   = errors.New("cli: usage error")
	HelpError    = errors.New("cli: help error")
)

type ErrorKind int

const (
	BadArgument ErrorKind = iota
	BadTemplateValue
	BadDirectoryValue
)

type CommandError struct {
	Kind        ErrorKind
	BadArgument string
	Err         error
}

func (e CommandError) Error() string {
	switch e.Kind {
	case BadArgument:
		return fmt.Sprintf("Unsupported argument '%s'.", e.BadArgument)
	case BadTemplateValue:
		return fmt.Sprintf("'--template' must be a 'js', 'javascript', 'ts', or 'typescript' (default 'javascript').")
	case BadDirectoryValue:
		cwd, _ := os.Getwd()
		return fmt.Sprintf("Use '.' to use the cwd '%s'.", cwd)
	}
	panic("Internal error")
}

func (e CommandError) Unwrap() error {
	return e.Err
}

func ParseCommand(args ...string) (Command, error) {
	cmd := Command{
		Template:  "javascript",
		Directory: ".",
	}
	var once sync.Once
	for _, arg := range args {
		err := CommandError{Kind: BadArgument, BadArgument: arg}
		if strings.HasPrefix(arg, "--template") {
			if len(arg) < len("--template=") {
				return Command{}, err
			}
			switch strings.ToLower(arg[len("--template="):]) {
			case "js":
			case "jsx":
			case "javascript":
				cmd.Template = "javascript"
			case "ts":
			case "tsx":
			case "typescript":
				cmd.Template = "typescript"
			default:
				return Command{}, err
			}
		} else if !strings.HasPrefix(arg, "--") {
			once.Do(func() {
				cmd.Directory = arg
			})
		} else {
			return Command{}, err
		}
	}
	return cmd, nil
}

func ParseCLIArguments() (interface{}, error) {
	if len(os.Args) < 2 {
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
	} else {
		cmd, err = ParseCommand(os.Args[1:]...)
	}
	return cmd, err
}
